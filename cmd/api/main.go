package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
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
	taghandler "github.com/silentmoon/api/internal/handler/tag"
	"github.com/silentmoon/api/internal/middleware"
	"github.com/silentmoon/api/internal/platform/database"
	"github.com/silentmoon/api/internal/platform/logger"
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
	tagservice "github.com/silentmoon/api/internal/service/tag"
	"github.com/silentmoon/api/migrations"
)

func main() {
	runMigrate := flag.Bool("migrate", false, "run database migrations and exit")
	flag.Parse()

	cfg := config.Load()
	appLogger := logger.New(cfg.AppEnv)

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
		defer dbConn.Close()

		// Run migrations if connected
		if err := database.RunMigrations(migrations.FS, ".", cfg.DatabaseURL); err != nil {
			appLogger.Error("failed to run database migrations", slog.String("error", err.Error()))
			if *runMigrate {
				os.Exit(1)
			}
		}
	}

	if *runMigrate {
		appLogger.Info("migrations completed, exiting")
		return
	}

	// Security & Token Manager
	tokenManager := security.NewTokenManager(
		cfg.JWTAccessSecret,
		cfg.JWTAccessTTL,
		cfg.JWTRefreshTTL,
	)

	srv := server.New(cfg, appLogger)
	r := srv.Router()

	// Health check handlers
	healthHandler := handler.NewHealthHandler(dbConn)
	r.Get("/health", healthHandler.Check)

	// API v1 group
	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/health", healthHandler.Check)

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

			// Services
			authSvc := authservice.NewService(userRepo, deviceRepo, tokenRepo, tokenManager)
			deviceSvc := deviceservice.NewService(deviceRepo, tokenRepo)
			bookSvc := bookservice.NewService(bookRepo)
			collSvc := collectionservice.NewService(collRepo, bookRepo)
			tagSvc := tagservice.NewService(tagRepo, bookRepo)
			annSvc := annservice.NewService(annRepo, bookRepo)
			bmSvc := bmservice.NewService(bmRepo, bookRepo)
			progSvc := progressservice.NewService(progRepo, bookRepo)

			// Handlers
			authHdl := authhandler.NewHandler(authSvc)
			deviceHdl := devicehandler.NewHandler(deviceSvc)
			bookHdl := bookhandler.NewHandler(bookSvc)
			collHdl := collectionhandler.NewHandler(collSvc)
			tagHdl := taghandler.NewHandler(tagSvc)
			annHdl := annhandler.NewHandler(annSvc)
			bmHdl := bmhandler.NewHandler(bmSvc)
			progHdl := progressohandler.NewHandler(progSvc)

			// Public authentication routes
			api.Route("/auth", func(auth chi.Router) {
				auth.Post("/register", authHdl.Register)
				auth.Post("/login", authHdl.Login)
				auth.Post("/refresh", authHdl.Refresh)
				auth.Post("/logout", authHdl.Logout)

				// Protected auth routes
				auth.Group(func(protected chi.Router) {
					protected.Use(middleware.Authenticate(tokenManager))
					protected.Get("/me", authHdl.Me)
				})
			})

			// Protected routes
			api.Group(func(protected chi.Router) {
				protected.Use(middleware.Authenticate(tokenManager))

				// Devices
				protected.Route("/devices", func(dev chi.Router) {
					dev.Get("/", deviceHdl.List)
					dev.Post("/", deviceHdl.Register)
					dev.Patch("/{id}", deviceHdl.Update)
					dev.Delete("/{id}", deviceHdl.Revoke)
				})

				// Books
				protected.Route("/books", func(bk chi.Router) {
					bk.Get("/", bookHdl.List)
					bk.Post("/", bookHdl.Create)
					bk.Get("/{id}", bookHdl.Get)
					bk.Patch("/{id}", bookHdl.Update)
					bk.Delete("/{id}", bookHdl.Delete)

					// Book Annotations & Bookmarks & Progress
					bk.Get("/{id}/annotations", annHdl.ListByBook)
					bk.Get("/{id}/bookmarks", bmHdl.ListByBook)
					bk.Get("/{id}/progress", progHdl.Get)
					bk.Put("/{id}/progress", progHdl.Update)

					// Book Tags
					bk.Get("/{id}/tags", tagHdl.ListTagsForBook)
					bk.Post("/{id}/tags/{tagId}", tagHdl.AddTagToBook)
					bk.Delete("/{id}/tags/{tagId}", tagHdl.RemoveTagFromBook)
				})

				// Annotations
				protected.Route("/annotations", func(ann chi.Router) {
					ann.Post("/", annHdl.Create)
					ann.Get("/{id}", annHdl.Get)
					ann.Patch("/{id}", annHdl.Update)
					ann.Delete("/{id}", annHdl.Delete)
				})

				// Bookmarks
				protected.Route("/bookmarks", func(bm chi.Router) {
					bm.Post("/", bmHdl.Create)
					bm.Get("/{id}", bmHdl.Get)
					bm.Patch("/{id}", bmHdl.Update)
					bm.Delete("/{id}", bmHdl.Delete)
				})

				// Collections
				protected.Route("/collections", func(col chi.Router) {
					col.Get("/", collHdl.List)
					col.Post("/", collHdl.Create)
					col.Get("/{id}", collHdl.Get)
					col.Patch("/{id}", collHdl.Update)
					col.Delete("/{id}", collHdl.Delete)

					// Collection Books
					col.Get("/{id}/books", collHdl.ListBooks)
					col.Post("/{id}/books/{bookId}", collHdl.AddBook)
					col.Delete("/{id}/books/{bookId}", collHdl.RemoveBook)
				})

				// Tags
				protected.Route("/tags", func(tg chi.Router) {
					tg.Get("/", tagHdl.List)
					tg.Post("/", tagHdl.Create)
					tg.Delete("/{id}", tagHdl.Delete)
				})
			})
		}
	})

	if err := srv.Start(); err != nil {
		appLogger.Error("server terminated with error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
