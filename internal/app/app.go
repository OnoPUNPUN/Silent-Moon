package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/silentmoon/api/internal/config"
	"github.com/silentmoon/api/internal/handler"
	annhandler "github.com/silentmoon/api/internal/handler/annotation"
	authhandler "github.com/silentmoon/api/internal/handler/auth"
	bmhandler "github.com/silentmoon/api/internal/handler/bookmark"
	bookhandler "github.com/silentmoon/api/internal/handler/book"
	collectionhandler "github.com/silentmoon/api/internal/handler/collection"
	devicehandler "github.com/silentmoon/api/internal/handler/device"
	progressohandler "github.com/silentmoon/api/internal/handler/progress"
	storagehandler "github.com/silentmoon/api/internal/handler/storage"
	synchandler "github.com/silentmoon/api/internal/handler/sync"
	taghandler "github.com/silentmoon/api/internal/handler/tag"
	"github.com/silentmoon/api/internal/middleware"
	"github.com/silentmoon/api/internal/platform/database"
	"github.com/silentmoon/api/internal/platform/security"
	"github.com/silentmoon/api/internal/platform/server"
	postgresrepo "github.com/silentmoon/api/internal/repository/postgres"
	sqlcdb "github.com/silentmoon/api/internal/repository/postgres/db"
	annservice "github.com/silentmoon/api/internal/service/annotation"
	authservice "github.com/silentmoon/api/internal/service/auth"
	bmservice "github.com/silentmoon/api/internal/service/bookmark"
	bookservice "github.com/silentmoon/api/internal/service/book"
	collectionservice "github.com/silentmoon/api/internal/service/collection"
	deviceservice "github.com/silentmoon/api/internal/service/device"
	progressservice "github.com/silentmoon/api/internal/service/progress"
	storageservice "github.com/silentmoon/api/internal/service/storage"
	syncservice "github.com/silentmoon/api/internal/service/sync"
	tagservice "github.com/silentmoon/api/internal/service/tag"
	"github.com/silentmoon/api/internal/storage"
	"github.com/silentmoon/api/internal/storage/r2"
	"github.com/silentmoon/api/migrations"
)

type App struct {
	server *server.Server
	dbConn *database.DB
	logger *slog.Logger
}

func New(cfg *config.Config, appLogger *slog.Logger, migrateOnly bool) (*App, error) {
	appLogger.Info("initializing personal reading library api",
		slog.String("env", cfg.AppEnv),
		slog.String("port", cfg.HTTPPort),
	)

	// Attempt database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var dbConn *database.DB
	var err error

	dbConn, err = database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		appLogger.Warn("database connection could not be established; starting with degraded state",
			slog.String("error", err.Error()),
		)
	} else {
		appLogger.Info("database connection established successfully")

		// Run database migrations
		if err := database.RunMigrations(migrations.FS, ".", cfg.DatabaseURL); err != nil {
			appLogger.Error("failed to run database migrations", slog.String("error", err.Error()))
			if migrateOnly {
				return nil, fmt.Errorf("migration failure: %w", err)
			}
		}
	}

	if migrateOnly {
		appLogger.Info("migrations completed")
		return &App{dbConn: dbConn, logger: appLogger}, nil
	}

	// Storage client initialization
	var storageClient storage.Client
	if cfg.R2AccessKeyID != "" && cfg.R2SecretKey != "" {
		r2Cli, err := r2.New(r2.Config{
			AccountID: cfg.R2AccountID,
			AccessKey: cfg.R2AccessKeyID,
			SecretKey: cfg.R2SecretKey,
			Bucket:    cfg.R2BucketName,
			Endpoint:  cfg.R2Endpoint,
		})
		if err != nil {
			appLogger.Warn("failed to initialize R2 storage client; using noop fallback", slog.String("error", err.Error()))
			storageClient = storage.NewNoopClient()
		} else {
			appLogger.Info("Cloudflare R2 storage client configured successfully", slog.String("bucket", cfg.R2BucketName))
			storageClient = r2Cli
		}
	} else {
		appLogger.Info("R2 credentials not provided; using noop storage fallback")
		storageClient = storage.NewNoopClient()
	}

	// Security & Token Manager
	tokenManager := security.NewTokenManager(
		cfg.JWTAccessSecret,
		cfg.JWTAccessTTL,
		cfg.JWTRefreshTTL,
	)
	authMiddleware := middleware.Authenticate(tokenManager)

	srv := server.New(cfg, appLogger)
	r := srv.Router()

	// Health check handler
	healthHdl := handler.NewHealthHandler(dbConn)
	healthHdl.RegisterRoutes(r)

	// API v1 group & domain handlers
	r.Route("/api/v1", func(api chi.Router) {
		healthHdl.RegisterRoutes(api)

		if dbConn != nil {
			queries := sqlcdb.New(dbConn.Pool)

			// Repositories
			userRepo := postgresrepo.NewUserRepository(queries)
			deviceRepo := postgresrepo.NewDeviceRepository(queries)
			tokenRepo := postgresrepo.NewRefreshTokenRepository(queries)
			bookRepo := postgresrepo.NewBookRepository(queries)
			collRepo := postgresrepo.NewCollectionRepository(queries)
			tagRepo := postgresrepo.NewTagRepository(queries)
			annRepo := postgresrepo.NewAnnotationRepository(queries)
			bmRepo := postgresrepo.NewBookmarkRepository(queries)
			progRepo := postgresrepo.NewReadingProgressRepository(queries)
			syncRepo := postgresrepo.NewSyncRepository(queries)

			// Services
			authSvc := authservice.NewService(userRepo, deviceRepo, tokenRepo, tokenManager)
			deviceSvc := deviceservice.NewService(deviceRepo, tokenRepo)
			bookSvc := bookservice.NewService(bookRepo)
			collSvc := collectionservice.NewService(collRepo, bookRepo)
			tagSvc := tagservice.NewService(tagRepo, bookRepo)
			annSvc := annservice.NewService(annRepo, bookRepo)
			bmSvc := bmservice.NewService(bmRepo, bookRepo)
			progSvc := progressservice.NewService(progRepo, bookRepo)
			storageSvc := storageservice.NewService(storageClient, bookRepo)
			syncSvc := syncservice.NewService(syncRepo, deviceRepo, annRepo, bmRepo, progRepo, bookRepo, collRepo)

			// Handlers
			authHdl := authhandler.NewHandler(authSvc)
			deviceHdl := devicehandler.NewHandler(deviceSvc)
			bookHdl := bookhandler.NewHandler(bookSvc)
			storageHdl := storagehandler.NewHandler(storageSvc)
			annHdl := annhandler.NewHandler(annSvc)
			bmHdl := bmhandler.NewHandler(bmSvc)
			progHdl := progressohandler.NewHandler(progSvc)
			collHdl := collectionhandler.NewHandler(collSvc)
			tagHdl := taghandler.NewHandler(tagSvc)
			syncHdl := synchandler.NewHandler(syncSvc)

			// Delegate route registration to each domain handler
			authHdl.RegisterRoutes(api, authMiddleware)
			deviceHdl.RegisterRoutes(api, authMiddleware)
			bookHdl.RegisterRoutes(api, authMiddleware)
			storageHdl.RegisterRoutes(api, authMiddleware)
			annHdl.RegisterRoutes(api, authMiddleware)
			bmHdl.RegisterRoutes(api, authMiddleware)
			progHdl.RegisterRoutes(api, authMiddleware)
			collHdl.RegisterRoutes(api, authMiddleware)
			tagHdl.RegisterRoutes(api, authMiddleware)
			syncHdl.RegisterRoutes(api, authMiddleware)
		}
	})

	return &App{
		server: srv,
		dbConn: dbConn,
		logger: appLogger,
	}, nil
}

func (a *App) Start() error {
	return a.server.Start()
}

func (a *App) Close() {
	if a.dbConn != nil {
		a.dbConn.Close()
	}
}
