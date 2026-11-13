# 📚 Personal Reading Library — Go Backend

> A lightweight, secure, and synchronization-focused backend powering a cross-platform, local-first digital reading application.

## 1. Project Overview

The Personal Reading Library Backend is a RESTful API built with Go.

It provides the cloud infrastructure required by the Flutter application to synchronize books, annotations, bookmarks, collections, reading progress, and device information across multiple devices.

The backend is designed around a Local-First architecture.

The Flutter application maintains its own local database and remains fully functional without an internet connection. The backend acts as the centralized synchronization service and cloud metadata manager.

The backend is responsible for:

- User authentication and authorization.
- Book metadata management.
- Cloud file storage integration.
- Device registration.
- Annotation synchronization.
- Bookmark synchronization.
- Reading progress synchronization.
- Collection and tag management.
- Incremental synchronization.
- Conflict detection and resolution.

The backend should remain lightweight, stateless where practical, and inexpensive to deploy.

## 2. System Architecture

```text
                  FLUTTER APPLICATIONS
             Android | Windows | macOS
                          |
             +------------+------------+
             |                         |
        Local SQLite               Local Files
             |                         |
             +------------+------------+
                          |
                      REST API
                          |
                  GO BACKEND SERVER
                          |
          +---------------+----------------+
          |               |                |
      PostgreSQL      Cloudflare R2     Auth System
      Metadata        Book Files        JWT / Tokens
          |
    Synchronization
       Metadata
```

### Architectural Responsibilities

**Flutter Application**
- Maintains local books and annotations.
- Handles offline reading.
- Queues local modifications.
- Synchronizes changes with the backend.

**Go Backend**
- Authenticates users.
- Validates incoming data.
- Manages cloud metadata.
- Coordinates synchronization.
- Generates secure file upload/download URLs.
- Enforces ownership and authorization.

**PostgreSQL**
- Stores user accounts.
- Stores book metadata.
- Stores annotations and bookmarks.
- Stores reading progress.
- Tracks registered devices.
- Maintains synchronization history.

**Cloudflare R2**
- Stores PDF files.
- Stores EPUB files.
- Stores book covers and other associated assets.
- Provides object storage through an S3-compatible API.

The backend must never store large book files directly inside PostgreSQL.

## 3. Technology Stack

| Component | Technology |
|---|---|
| Language | Go |
| HTTP Framework | Chi or Gin |
| Database | PostgreSQL |
| Database Driver | pgx |
| Database Queries | sqlc |
| Migrations | golang-migrate |
| Authentication | JWT |
| Password Hashing | Argon2id or bcrypt |
| Object Storage | Cloudflare R2 |
| Object Storage SDK | AWS SDK for Go v2 |
| Configuration | Environment variables |
| Logging | slog or Zerolog |
| Containerization | Docker |
| API Documentation | OpenAPI |
| Testing | Go testing package |

### Technology Choices

**Go**

Go provides a lightweight runtime, efficient concurrency, simple deployment, and excellent support for building HTTP services.

**PostgreSQL**

PostgreSQL is the primary persistent database for all cloud-side application data.

**sqlc**

SQLC generates type-safe Go code from SQL queries, keeping database access explicit and predictable.

**Cloudflare R2**

R2 provides object storage for books without requiring the backend to proxy every large file transfer.

## 4. Backend Responsibilities

### 4.1 Authentication

The backend must support user authentication.

Features:
- User registration.
- User login.
- Access token generation.
- Refresh token rotation.
- Logout.
- Token revocation.
- Current user information.

Authentication should use short-lived access tokens and securely managed refresh tokens.

Passwords must be hashed using a suitable password-hashing algorithm. Plaintext passwords must never be stored.

### 4.2 User Management

Users should be able to:
- View their profile.
- Update basic profile information.
- Change their password.
- Manage registered devices.
- Revoke device access.
- Manage account sessions.

### 4.3 Book Management

The backend stores metadata about the books belonging to each user.

Features:
- Create book records.
- Retrieve book details.
- List user books.
- Update book metadata.
- Delete books.
- Track cloud file availability.
- Manage file versions.
- Generate secure upload and download URLs.

A book record should reference an object storage key instead of storing the actual file contents.

### 4.4 Cloud File Management

Book files are stored in Cloudflare R2.

The backend should use presigned URLs for direct client uploads and downloads.

#### Upload Workflow

1. Client requests an upload session.
2. Backend verifies authentication.
3. Backend validates the requested file metadata.
4. Backend generates an upload URL.
5. Client uploads the file directly to R2.
6. Client notifies the backend that the upload is complete.
7. Backend verifies the uploaded object.
8. Backend updates the book record.

#### Download Workflow

1. Client requests a download URL.
2. Backend verifies book ownership.
3. Backend generates a short-lived download URL.
4. Client downloads the book directly from R2.

This approach prevents the Go server from becoming a bottleneck for large file transfers.

### 4.5 Annotation Management

Annotations are a central part of the application.

The backend must support:

- Creating annotations.
- Updating annotations.
- Retrieving annotations.
- Listing annotations by book.
- Filtering annotations by type.
- Deleting annotations.
- Synchronizing annotations across devices.

Supported annotation types:

| Type | Description |
|---|---|
| highlight | Highlights selected text |
| underline | Underlines selected text |
| strikethrough | Marks selected text |
| note | Standalone or passage-associated note |

Annotations should store location information as JSONB because PDF and EPUB use different positioning systems.

### 4.6 Bookmark Management

Features:
- Create bookmarks.
- Update bookmark labels.
- Retrieve bookmarks by book.
- Delete bookmarks.
- Synchronize bookmarks across devices.

### 4.7 Reading Progress

The backend should maintain synchronized reading progress.

Features:
- Save reading position.
- Update reading percentage.
- Retrieve the latest synchronized position.
- Track device-specific reading activity.

Reading positions must support both PDF page locations and EPUB CFI locations.

### 4.8 Collections and Tags

Users should be able to organize books into collections and assign custom tags.

Features:
- Create collections.
- Update collections.
- Delete collections.
- Add books to collections.
- Remove books from collections.
- Create and manage tags.
- Assign and remove tags from books.

### 4.9 Device Management

Each installation should have a unique device identifier.

Features:
- Register devices.
- Retrieve registered devices.
- Update device information.
- Track last synchronization time.
- Revoke devices.
- Track device-specific synchronization cursors.

Device registration is important for synchronization, security, and conflict detection.

## 5. Domain Models

The following are conceptual Go domain models. Database-specific representations should remain in the persistence layer.

### 5.1 User

Represents an authenticated application user.

```go
type User struct {
    ID           uuid.UUID
    Email        string
    PasswordHash string
    DisplayName  string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

Responsibilities:
- Identify the account owner.
- Own books and annotations.
- Authenticate API requests.
- Control access to private resources.

### 5.2 Book

Represents a book in the cloud library.

```go
type Book struct {
    ID          uuid.UUID
    UserID      uuid.UUID
    Title       string
    Author      *string
    Description *string
    Format      string
    StorageKey  string
    CoverKey    *string
    FileSize    int64
    Checksum    string
    Status      string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   *time.Time
}
```

Responsibilities:
- Store book metadata.
- Reference the original object in R2.
- Track file availability.
- Associate the book with its owner.

### 5.3 Annotation

Represents a highlight, note, or text modification marker.

```go
type Annotation struct {
    ID           uuid.UUID
    UserID       uuid.UUID
    BookID       uuid.UUID
    DeviceID     uuid.UUID
    Type         string
    SelectedText *string
    Color        *string
    Note         *string
    Location     json.RawMessage
    Version      int64
    CreatedAt    time.Time
    UpdatedAt    time.Time
    DeletedAt    *time.Time
}
```

Responsibilities:
- Preserve highlighted passages.
- Store personal notes.
- Maintain format-specific location information.
- Track modifications.
- Support multi-device synchronization.

### 5.4 Bookmark

Represents a saved location inside a book.

```go
type Bookmark struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    BookID    uuid.UUID
    DeviceID  uuid.UUID
    Title     *string
    Location  json.RawMessage
    Version   int64
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
}
```

### 5.5 ReadingProgress

Represents synchronized reading progress.

```go
type ReadingProgress struct {
    ID                 uuid.UUID
    UserID             uuid.UUID
    BookID             uuid.UUID
    DeviceID           uuid.UUID
    Location           json.RawMessage
    ProgressPercentage float64
    LastReadAt         time.Time
    UpdatedAt          time.Time
}
```

### 5.6 Collection

Represents a user-created book collection.

```go
type Collection struct {
    ID          uuid.UUID
    UserID      uuid.UUID
    Name        string
    Description *string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   *time.Time
}
```

### 5.7 Tag

Represents a custom label.

```go
type Tag struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    Name      string
    CreatedAt time.Time
}
```

### 5.8 Device

Represents a registered application installation.

```go
type Device struct {
    ID           uuid.UUID
    UserID       uuid.UUID
    Name         string
    Platform     string
    RefreshTokenHash *string
    LastSyncedAt *time.Time
    CreatedAt    time.Time
    RevokedAt    *time.Time
}
```

### 5.9 SyncChange

Represents a change that can be delivered to synchronized devices.

```go
type SyncChange struct {
    ID         int64
    UserID     uuid.UUID
    DeviceID   uuid.UUID
    EntityType string
    EntityID   uuid.UUID
    Operation  string
    CreatedAt  time.Time
}
```

Responsibilities:
- Track synchronized changes.
- Provide incremental synchronization.
- Notify other devices about modifications.
- Preserve deletion information through tombstones.

## 6. PostgreSQL Database Structure

The following schema describes the logical database structure.

### Entity Relationships

```text
users
  |
  +---- books
  |       |
  |       +---- annotations
  |       |
  |       +---- bookmarks
  |       |
  |       +---- reading_progress
  |       |
  |       +---- book_collections
  |       |
  |       +---- book_tags
  |
  +---- collections
  |
  +---- tags
  |
  +---- devices
  |
  +---- sync_changes
  |
  +---- refresh_tokens
```

### 6.1 users

Stores registered application accounts.

| Column | Type | Description |
|---|---|---|
| id | UUID | Primary key |
| email | VARCHAR | Unique email |
| password_hash | TEXT | Hashed password |
| display_name | VARCHAR | User display name |
| created_at | TIMESTAMPTZ | Creation time |
| updated_at | TIMESTAMPTZ | Last update |

### 6.2 books

Stores cloud library metadata.

| Column | Type | Description |
|---|---|---|
| id | UUID | Primary key |
| user_id | UUID | Owner |
| title | TEXT | Book title |
| author | TEXT | Author |
| description | TEXT | Description |
| format | VARCHAR | PDF or EPUB |
| storage_key | TEXT | R2 object key |
| cover_key | TEXT | R2 cover object key |
| file_size | BIGINT | Original file size |
| checksum | VARCHAR | File integrity hash |
| status | VARCHAR | Upload state |
| created_at | TIMESTAMPTZ | Creation time |
| updated_at | TIMESTAMPTZ | Last update |
| deleted_at | TIMESTAMPTZ | Soft deletion |

Suggested upload states:
- pending
- uploaded
- failed

### 6.3 annotations

Stores highlights and notes.

| Column | Type | Description |
|---|---|---|
| id | UUID | Primary key |
| user_id | UUID | Owner |
| book_id | UUID | Related book |
| device_id | UUID | Originating device |
| type | VARCHAR | Annotation type |
| selected_text | TEXT | Selected passage |
| color | VARCHAR | Highlight color |
| note | TEXT | Personal note |
| location | JSONB | Format-specific position |
| version | BIGINT | Entity version |
| created_at | TIMESTAMPTZ | Creation time |
| updated_at | TIMESTAMPTZ | Last update |
| deleted_at | TIMESTAMPTZ | Soft deletion |

### 6.4 bookmarks

| Column | Type | Description |
|---|---|---|
| id | UUID | Primary key |
| user_id | UUID | Owner |
| book_id | UUID | Related book |
| device_id | UUID | Originating device |
| title | TEXT | Optional label |
| location | JSONB | Reading position |
| version | BIGINT | Entity version |
| created_at | TIMESTAMPTZ | Creation time |
| updated_at | TIMESTAMPTZ | Last update |
| deleted_at | TIMESTAMPTZ | Soft deletion |

### 6.5 reading_progress

| Column | Type | Description |
|---|---|---|
| id | UUID | Primary key |
| user_id | UUID | Owner |
| book_id | UUID | Related book |
| device_id | UUID | Device |
| location | JSONB | Current position |
| progress_percentage | NUMERIC | Completion percentage |
| last_read_at | TIMESTAMPTZ | Last reading activity |
| updated_at | TIMESTAMPTZ | Last update |

Unique constraint: `(book_id, device_id)`.

### 6.6 collections

| Column | Type | Description |
|---|---|---|
| id | UUID | Primary key |
| user_id | UUID | Owner |
| name | VARCHAR | Collection name |
| description | TEXT | Description |
| created_at | TIMESTAMPTZ | Creation time |
| updated_at | TIMESTAMPTZ | Last update |
| deleted_at | TIMESTAMPTZ | Soft deletion |

### 6.7 book_collections

Many-to-many relationship between books and collections.

| Column | Type | Description |
|---|---|---|
| book_id | UUID | Book |
| collection_id | UUID | Collection |
| created_at | TIMESTAMPTZ | Assignment time |

Composite primary key: `(book_id, collection_id)`.

### 6.8 tags

| Column | Type | Description |
|---|---|---|
| id | UUID | Primary key |
| user_id | UUID | Owner |
| name | VARCHAR | Tag name |
| created_at | TIMESTAMPTZ | Creation time |

### 6.9 book_tags

Many-to-many relationship between books and tags.

| Column | Type | Description |
|---|---|---|
| book_id | UUID | Book |
| tag_id | UUID | Tag |

Composite primary key: `(book_id, tag_id)`.

### 6.10 devices

Stores registered devices.

| Column | Type | Description |
|---|---|---|
| id | UUID | Primary key |
| user_id | UUID | Owner |
| name | VARCHAR | Device name |
| platform | VARCHAR | Android, Windows, macOS |
| last_synced_at | TIMESTAMPTZ | Last synchronization |
| created_at | TIMESTAMPTZ | Registration time |
| revoked_at | TIMESTAMPTZ | Revocation time |

### 6.11 refresh_tokens

Stores hashed refresh token records.

| Column | Type | Description |
|---|---|---|
| id | UUID | Primary key |
| user_id | UUID | Owner |
| device_id | UUID | Associated device |
| token_hash | TEXT | Hashed token |
| expires_at | TIMESTAMPTZ | Expiration |
| revoked_at | TIMESTAMPTZ | Revocation |
| created_at | TIMESTAMPTZ | Creation time |

### 6.12 sync_changes

Append-only synchronization change log.

| Column | Type | Description |
|---|---|---|
| id | BIGINT | Monotonic sequence |
| user_id | UUID | Owner |
| device_id | UUID | Originating device |
| entity_type | VARCHAR | Entity category |
| entity_id | UUID | Entity identifier |
| operation | VARCHAR | upsert or delete |
| created_at | TIMESTAMPTZ | Change time |

The `id` column is used as a synchronization cursor.

The log should retain deletion events until all relevant synchronization clients have had an opportunity to receive them, or until a safe retention policy permits cleanup.

## 7. Annotation Location Structure

Location information is format-specific and must be stored as JSONB.

### PDF Location Example

```json
{
  "format": "pdf",
  "page": 42,
  "selectedText": "An important passage",
  "rectangles": [
    {
      "x": 120.5,
      "y": 340.2,
      "width": 240.0,
      "height": 18.5
    }
  ],
  "pageWidth": 595.0,
  "pageHeight": 842.0
}
```

Coordinates should use a documented normalized coordinate system so rendering differences between platforms do not invalidate annotations.

### EPUB Location Example

```json
{
  "format": "epub",
  "cfi": "epubcfi(/6/8!/4/2/1:120)",
  "selectedText": "An important passage",
  "prefix": "The preceding text",
  "suffix": "The following text"
}
```

The EPUB CFI identifies the original passage, while the selected text and surrounding context provide fallback anchoring.

## 8. Synchronization System

Synchronization is one of the most important backend responsibilities.

The backend must support incremental synchronization instead of requiring clients to download the entire library every time.

### 8.1 Synchronization Workflow

```text
Flutter Application
       |
       | Local changes
       v
Local SQLite Database
       |
       | Pending operations
       v
Go Synchronization API
       |
       +------ Validate changes
       |
       +------ Persist entities
       |
       +------ Append sync_changes
       |
       v
PostgreSQL
       |
       | Return accepted changes
       v
Flutter Application
       |
       | Request changes after cursor
       v
Apply changes to SQLite
       |
       v
Update local sync cursor
```

### 8.2 Push Changes

The client sends pending local changes to the backend.

Example request:

```json
{
  "deviceId": "device-uuid",
  "operations": [
    {
      "operationId": "operation-uuid",
      "entityType": "annotation",
      "entityId": "annotation-uuid",
      "operation": "upsert",
      "baseVersion": 1,
      "data": {
        "bookId": "book-uuid",
        "type": "highlight",
        "selectedText": "Important passage",
        "color": "#FFE082",
        "note": "Remember this concept"
      }
    }
  ]
}
```

The backend should:
- Authenticate the user.
- Validate device ownership.
- Validate entity ownership.
- Validate the operation.
- Detect duplicate operation IDs.
- Check entity versions.
- Persist accepted changes.
- Append synchronization events.
- Return accepted or conflicting operations.

Each operation should be idempotent. Retrying the same operation must not create duplicate records.

### 8.3 Pull Changes

The client requests changes after its last known cursor.

Example:

```http
GET /api/v1/sync/changes?cursor=1250&limit=100
```

Example response:

```json
{
  "changes": [
    {
      "sequence": 1251,
      "entityType": "annotation",
      "entityId": "annotation-uuid",
      "operation": "upsert",
      "data": {}
    }
  ],
  "nextCursor": 1251,
  "hasMore": false
}
```

The client should apply changes transactionally to SQLite and update its cursor only after successful application.

### 8.4 Conflict Resolution

The initial implementation should use version-based conflict detection.

Every synchronized entity should maintain a version number.

When updating an entity, the client sends the version it originally observed.

If the server version has changed, the backend reports a conflict rather than silently overwriting data.

Initial conflict policy:
- Annotations: entity-level conflict detection with a documented last-write-wins resolution strategy.
- Bookmarks: entity-level conflict detection.
- Reading progress: latest-position policy, with per-device positions retained.
- Book metadata: version-based updates.
- Deletions: tombstones.

The client must be able to recover from a conflict by receiving the authoritative server version and applying the documented resolution.

## 9. REST API Design

All endpoints should use the `/api/v1` prefix.

### Authentication

| Method | Endpoint | Purpose |
|---|---|---|
| POST | /auth/register | Register account |
| POST | /auth/login | Login |
| POST | /auth/refresh | Refresh tokens |
| POST | /auth/logout | Revoke session |
| GET | /auth/me | Current user |

### Books

| Method | Endpoint | Purpose |
|---|---|---|
| GET | /books | List books |
| POST | /books | Create book metadata |
| GET | /books/{id} | Get book |
| PATCH | /books/{id} | Update book |
| DELETE | /books/{id} | Delete book |

### Book Storage

| Method | Endpoint | Purpose |
|---|---|---|
| POST | /books/{id}/upload-url | Generate upload URL |
| POST | /books/{id}/upload-complete | Confirm upload |
| GET | /books/{id}/download-url | Generate download URL |

### Annotations

| Method | Endpoint | Purpose |
|---|---|---|
| GET | /books/{id}/annotations | List book annotations |
| POST | /annotations | Create annotation |
| PATCH | /annotations/{id} | Update annotation |
| DELETE | /annotations/{id} | Delete annotation |

### Bookmarks

| Method | Endpoint | Purpose |
|---|---|---|
| GET | /books/{id}/bookmarks | List bookmarks |
| POST | /bookmarks | Create bookmark |
| PATCH | /bookmarks/{id} | Update bookmark |
| DELETE | /bookmarks/{id} | Delete bookmark |

### Reading Progress

| Method | Endpoint | Purpose |
|---|---|---|
| GET | /books/{id}/progress | Retrieve progress |
| PUT | /books/{id}/progress | Update progress |

### Collections

| Method | Endpoint | Purpose |
|---|---|---|
| GET | /collections | List collections |
| POST | /collections | Create collection |
| PATCH | /collections/{id} | Update collection |
| DELETE | /collections/{id} | Delete collection |
| POST | /collections/{id}/books/{bookId} | Add book |
| DELETE | /collections/{id}/books/{bookId} | Remove book |

### Tags

| Method | Endpoint | Purpose |
|---|---|---|
| GET | /tags | List tags |
| POST | /tags | Create tag |
| DELETE | /tags/{id} | Delete tag |
| POST | /books/{id}/tags/{tagId} | Assign tag |
| DELETE | /books/{id}/tags/{tagId} | Remove tag |

### Devices

| Method | Endpoint | Purpose |
|---|---|---|
| GET | /devices | List devices |
| POST | /devices | Register device |
| PATCH | /devices/{id} | Update device |
| DELETE | /devices/{id} | Revoke device |

### Synchronization

| Method | Endpoint | Purpose |
|---|---|---|
| POST | /sync/push | Upload local changes |
| GET | /sync/changes | Retrieve remote changes |
| GET | /sync/status | Get synchronization status |

All endpoints must enforce user ownership and authorization.

## 10. Go Project Structure

Recommended project layout:

```text
backend/
│
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── config/
│   │
│   ├── domain/
│   │   ├── user/
│   │   ├── book/
│   │   ├── annotation/
│   │   ├── bookmark/
│   │   ├── progress/
│   │   ├── collection/
│   │   ├── tag/
│   │   ├── device/
│   │   └── sync/
│   │
│   ├── handler/
│   │   ├── auth/
│   │   ├── book/
│   │   ├── annotation/
│   │   ├── bookmark/
│   │   ├── progress/
│   │   ├── collection/
│   │   ├── device/
│   │   └── sync/
│   │
│   ├── service/
│   │   ├── auth/
│   │   ├── book/
│   │   ├── annotation/
│   │   ├── storage/
│   │   └── sync/
│   │
│   ├── repository/
│   │   └── postgres/
│   │
│   ├── middleware/
│   ├── storage/
│   │   └── r2/
│   │
│   └── platform/
│       ├── database/
│       ├── logger/
│       └── server/
│
├── migrations/
│
├── sql/
│   ├── queries/
│   └── schema/
│
├── api/
│   └── openapi.yaml
│
├── configs/
│
├── Dockerfile
├── docker-compose.yml
├── sqlc.yaml
├── go.mod
├── go.sum
└── README.md
```

### Layer Responsibilities

**Domain**
- Core entities.
- Repository interfaces.
- Business rules.
- Domain errors.

**Handler**
- HTTP request handling.
- Input validation.
- Authentication context.
- HTTP response formatting.

**Service**
- Application use cases.
- Business logic orchestration.
- Transaction management.
- Synchronization rules.

**Repository**
- Database access.
- SQL queries.
- Persistence implementation.

**Storage**
- R2 integration.
- Presigned URL generation.
- File metadata verification.

**Middleware**
- Authentication.
- Authorization.
- Request logging.
- Rate limiting.
- Recovery and error handling.

## 11. Security Requirements

Security must be implemented from the beginning.

### Authentication
- Password hashing using Argon2id or bcrypt.
- Short-lived access tokens.
- Refresh token rotation.
- Refresh token hashes stored in PostgreSQL.
- Device-associated sessions.
- Session revocation.

### Authorization
- Users can only access their own books.
- Users can only modify their own annotations.
- Users cannot access another user's storage objects.
- Device registration must be validated against account ownership.
- Every database query involving user resources must enforce ownership.

### File Security
- Validate file extensions and actual file formats.
- Enforce upload size limits.
- Generate unpredictable storage keys.
- Use short-lived presigned URLs.
- Keep R2 credentials exclusively on the backend.
- Verify uploaded objects before marking uploads complete.
- Never trust client-supplied file paths.

### API Security
- HTTPS only in production.
- Request validation.
- Rate limiting.
- CORS configuration.
- Secure error responses.
- Structured logging without passwords or tokens.
- Database transactions for multi-step operations.

## 12. Error Response Format

All API errors should follow a consistent structure.

Example:

```json
{
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "The requested book was not found."
  }
}
```

Suggested error codes:

| Code | Meaning |
|---|---|
| VALIDATION_ERROR | Invalid request |
| UNAUTHORIZED | Authentication required |
| FORBIDDEN | Insufficient permissions |
| RESOURCE_NOT_FOUND | Resource does not exist |
| CONFLICT | Version conflict |
| DUPLICATE_OPERATION | Operation already processed |
| STORAGE_ERROR | File storage failure |
| INTERNAL_ERROR | Unexpected server error |

## 13. Database and Synchronization Considerations

Important implementation rules:

1. Use UUIDs for externally referenced entities.
2. Use foreign keys to enforce relational integrity.
3. Use transactions for multi-entity operations.
4. Add indexes for frequently queried columns.
5. Index `(user_id, book_id)` where appropriate.
6. Index synchronization sequence IDs.
7. Use JSONB for format-specific locations, not for the entire relational data model.
8. Use soft deletion for synchronized entities.
9. Never expose internal storage credentials.
10. Keep synchronization operations idempotent.

### Important Sync Design Decision

The `sync_changes` table is an append-only change log.

Every accepted synchronized modification should create a new change event inside the same database transaction as the entity update.

The backend should provide a stable, ordered cursor to clients.

A device must not advance its cursor until it has successfully applied all changes in the corresponding batch.

## 14. Deployment Architecture

The initial deployment should prioritize low operating costs.

| Component | Provider |
|---|---|
| Go API | Self-hosted homelab or low-cost container hosting |
| Database | Neon PostgreSQL |
| File Storage | Cloudflare R2 |
| Remote Access | HTTPS reverse proxy or Cloudflare Tunnel |
| Containerization | Docker |

The Go API can run as a single container.

PostgreSQL can initially be hosted through Neon, avoiding the need to maintain a public database server.

Cloudflare R2 handles book storage independently of the API server.

### Environment Variables

Example:

```env
APP_ENV=development
HTTP_PORT=8080

DATABASE_URL=postgresql://user:password@localhost:5432/reading_library

JWT_ACCESS_SECRET=replace_me
JWT_ACCESS_TTL=15m

R2_ACCOUNT_ID=your_account_id
R2_ACCESS_KEY_ID=your_access_key
R2_SECRET_ACCESS_KEY=your_secret_key
R2_BUCKET_NAME=reading-library

R2_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com

ALLOWED_ORIGINS=http://localhost:3000
```

Production secrets must be managed securely and must never be committed to version control.

### Docker

The backend should provide:
- A multi-stage Dockerfile.
- A small production runtime image.
- Docker Compose for local development.
- Database migration commands.
- Health checks.
- Persistent configuration through environment variables.

## 15. Testing Strategy

### Unit Tests
- Domain business rules.
- Authentication logic.
- Annotation validation.
- Conflict detection.
- Synchronization logic.

### Integration Tests
- PostgreSQL repository operations.
- Database transactions.
- Authentication middleware.
- R2 upload and download workflows.
- Synchronization endpoints.

### API Tests
- Registration and login.
- Book management.
- Annotation CRUD.
- Bookmark CRUD.
- Device registration.
- Incremental synchronization.
- Unauthorized access prevention.

### Critical Synchronization Tests

The following scenarios must be tested:

1. Device A creates an annotation.
2. Device B retrieves the annotation.
3. Device B modifies the annotation.
4. Device A retrieves the modification.
5. Device A modifies an outdated version.
6. The backend detects the conflict.
7. Device A deletes the annotation.
8. Device B receives the deletion tombstone.
9. A retry of the same operation does not create duplicates.
10. A device reconnects after several days offline and successfully synchronizes its pending changes.

## 16. Development Roadmap

### Phase 1 — Backend Foundation
- [ ] Initialize Go project.
- [ ] Configure HTTP server.
- [ ] Configure PostgreSQL.
- [ ] Set up migrations.
- [ ] Integrate sqlc.
- [ ] Configure logging.
- [ ] Add middleware.
- [ ] Create health endpoint.

### Phase 2 — Authentication and Devices
- [ ] User registration.
- [ ] User login.
- [ ] JWT authentication.
- [ ] Refresh token rotation.
- [ ] Logout.
- [ ] Device registration.
- [ ] Device revocation.

### Phase 3 — Library Management
- [ ] Book CRUD.
- [ ] Collection CRUD.
- [ ] Book-collection relationships.
- [ ] Tag management.
- [ ] Book-tag relationships.

### Phase 4 — Annotation System
- [ ] Annotation CRUD.
- [ ] Bookmark CRUD.
- [ ] Reading progress.
- [ ] Format-specific location validation.
- [ ] Soft deletion.

### Phase 5 — Object Storage
- [ ] R2 integration.
- [ ] Presigned upload URLs.
- [ ] Presigned download URLs.
- [ ] Upload completion verification.
- [ ] File integrity checks.
- [ ] Storage cleanup.

### Phase 6 — Synchronization Engine
- [ ] Change log.
- [ ] Push synchronization.
- [ ] Pull synchronization.
- [ ] Cursor management.
- [ ] Idempotency.
- [ ] Version-based conflict detection.
- [ ] Deletion tombstones.
- [ ] Multi-device testing.

### Phase 7 — Production Readiness
- [ ] API documentation.
- [ ] Integration tests.
- [ ] Security review.
- [ ] Docker production setup.
- [ ] Database backup strategy.
- [ ] Deployment.
- [ ] Monitoring and logging.

## 17. Out of Scope for Initial Release

The following features are not part of the initial backend implementation:

- AI-powered book analysis.
- Social networking.
- Public book libraries.
- Collaborative annotations.
- Real-time collaborative reading.
- Book purchasing.
- Recommendation engines.
- Payment processing.
- Subscription management.

The backend should focus on reliable synchronization, secure file management, and user-owned reading data.

## 18. Definition of Success

The backend's primary success criterion is reliable multi-device synchronization.

A successful implementation must allow:

1. A user to register and authenticate.
2. A user to upload a PDF or EPUB.
3. Multiple devices to access the same cloud library.
4. Devices to independently create annotations.
5. Annotation changes to synchronize without duplication.
6. Reading progress to move between devices.
7. Offline modifications to synchronize after reconnection.
8. Conflicts to be detected and handled predictably.
9. Deleted records to propagate to offline devices.
10. Users to access only their own data.

---

**The backend is not intended to be a complicated enterprise platform. Its purpose is to provide a reliable synchronization and storage foundation for a personal, local-first reading ecosystem.**