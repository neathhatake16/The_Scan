# Doc Scanner

A full-stack mobile document scanning application. Point your camera at any document, and the app auto-detects the edges, corrects perspective, enhances contrast, and stores the result as a clean PDF — all linked to your account.

---

## Tech Stack

| Layer            | Technology                                         |
| ---------------- | -------------------------------------------------- |
| Mobile           | Flutter 3.x · Dart 3.x · BLoC · Clean Architecture |
| API              | Go 1.22 · Gin · GORM · JWT                         |
| Image processing | Python 3.12 · FastAPI · OpenCV · ReportLab         |
| Database         | MySQL 8.4                                          |
| Infrastructure   | Docker · Docker Compose · Nginx                    |

---

## Project Structure

```
doc-scanner/
├── backend/                     # Go API (Gin + GORM)
│   ├── cmd/server/main.go       # Entry point — wires all layers
│   ├── internal/
│   │   ├── config/              # Environment configuration
│   │   ├── database/            # GORM connection + AutoMigrate
│   │   ├── models/              # GORM models + request/response DTOs
│   │   ├── repository/          # Data access layer (interfaces + impls)
│   │   ├── services/            # Business logic
│   │   ├── handlers/            # HTTP handlers (thin — bind → service → respond)
│   │   ├── middleware/          # JWT auth, request logger
│   │   ├── response/            # Unified JSON envelope helpers
│   │   ├── apperrors/           # Typed HTTP-aware error types
│   │   └── logger/              # Structured Zap logger
│   ├── Dockerfile               # Multi-stage build (builder + alpine runtime)
│   └── go.mod
│
├── python-scanner/              # OpenCV PDF microservice (FastAPI)
│   ├── main.py                  # /scan endpoint
│   ├── requirements.txt
│   └── Dockerfile
│
├── flutter_app/                 # Flutter mobile app
│   ├── lib/
│   │   ├── main.dart            # App root + dependency injection
│   │   ├── app_shell.dart       # Bottom-nav shell
│   │   ├── core/
│   │   │   ├── constants/       # API endpoints, app constants
│   │   │   ├── errors/          # Failures + Exceptions
│   │   │   ├── network/         # Dio client + JWT interceptor + auto-refresh
│   │   │   ├── storage/         # Token persistence (SharedPreferences)
│   │   │   └── utils/           # Either<L,R> type
│   │   ├── features/
│   │   │   ├── auth/            # Login · Register
│   │   │   ├── documents/       # List · Download · Rename · Delete
│   │   │   ├── scan/            # Camera / gallery → upload
│   │   │   └── profile/         # User info · storage stats · logout
│   │   └── shared/
│   │       ├── theme/           # AppTheme (Material 3)
│   │       └── widgets/         # AppButton, AppTextField, ErrorView
│   └── pubspec.yaml
│
├── docker/
│   └── nginx/nginx.conf         # Reverse proxy + rate limiting
├── docker-compose.yml           # Full stack (mysql + go-api + scanner + nginx)
└── .env.example                 # Environment variable template
```

Each feature in `flutter_app/lib/features/` follows the same four-layer pattern:

```
feature/
├── data/
│   ├── datasources/    # Raw HTTP calls (throws Exceptions)
│   ├── models/         # JSON-serialisable models (extend domain entities)
│   └── repositories/   # Catch exceptions → return Either<Failure, T>
├── domain/
│   ├── entities/       # Pure Dart classes — no Flutter, no JSON
│   ├── repositories/   # Abstract contracts
│   └── usecases/       # One public method per use case
└── presentation/
    ├── bloc/           # Events · States · BLoC
    └── pages/          # Stateful widgets — read BLoC, never call HTTP directly
```

---

## Architecture: Go Backend

```
HTTP Request
     │
     ▼
  Middleware          (JWT auth, request logger, CORS)
     │
     ▼
  Handler             (bind request → call service → call response helper)
     │
     ▼
  Service             (business logic, orchestration, no HTTP types)
     │
     ▼
  Repository          (interface) ← implemented by GORM concrete type
     │
     ▼
  MySQL (GORM)
```

**Key principles:**

- Handlers are thin — they never touch `*gorm.DB` directly
- Services depend on repository interfaces, making them fully testable with mocks
- All errors flow as typed `*apperrors.AppError`; the `response` package converts them to the correct HTTP status
- The `response.Envelope` wrapper ensures every API response has the same `{ success, data, error }` shape

---

## API Reference

### Auth

| Method | Path             | Auth | Description                        |
| ------ | ---------------- | ---- | ---------------------------------- |
| POST   | `/auth/register` | —    | Create account, returns token pair |
| POST   | `/auth/login`    | —    | Login, returns token pair          |
| POST   | `/auth/refresh`  | —    | Rotate refresh token               |
| POST   | `/auth/logout`   | —    | Revoke refresh token               |

### User

| Method | Path                        | Auth | Description                   |
| ------ | --------------------------- | ---- | ----------------------------- |
| GET    | `/users/me`                 | ✓    | Get profile                   |
| PATCH  | `/users/me`                 | ✓    | Update full_name / avatar_url |
| POST   | `/users/me/change-password` | ✓    | Change password               |
| GET    | `/users/me/storage`         | ✓    | Storage usage stats           |

### Documents

| Method | Path                         | Auth | Description                  |
| ------ | ---------------------------- | ---- | ---------------------------- |
| POST   | `/scan?title=...`            | ✓    | Upload image → get PDF saved |
| GET    | `/documents?page=1&limit=20` | ✓    | Paginated document list      |
| GET    | `/documents/:id/download`    | ✓    | Download PDF file            |
| PATCH  | `/documents/:id`             | ✓    | Rename document              |
| DELETE | `/documents/:id`             | ✓    | Delete document + file       |

**Response envelope** (all endpoints):

```json
{
  "success": true,
  "data": { ... }
}
```

```json
{
  "success": false,
  "error": "human-readable message"
}
```

---

## Quick Start

### Prerequisites

- Docker ≥ 24 and Docker Compose v2
- Flutter SDK ≥ 3.3 (for mobile)
- An Android emulator or physical device

### 1. Clone and configure

```bash
git clone https://github.com/yourname/doc-scanner.git
cd doc-scanner
cp .env.example .env
```

Edit `.env` — at minimum set:

```env
DB_PASSWORD=your_strong_password
JWT_SECRET=$(openssl rand -hex 32)
```

### 2. Start the backend stack

```bash
docker compose up --build -d
```

This starts four services in order:

1. **mysql** — waits until healthy
2. **python-scanner** — OpenCV PDF service on port 8001
3. **go-api** — REST API on port 8080 (auto-migrates the database on first start)
4. **nginx** — reverse proxy on port 80

Check everything is running:

```bash
docker compose ps
curl http://localhost/health
# {"status":"ok","version":"2.0.0"}
```

View logs:

```bash
docker compose logs -f go-api
docker compose logs -f python-scanner
```

### 3. Run the Flutter app

```bash
cd flutter_app
flutter pub get
```

Open `lib/core/constants/api_constants.dart` and set the base URL:

```dart
// Android emulator → host machine localhost
defaultValue: 'http://10.0.2.2:80',

// Physical device → your machine's local IP
defaultValue: 'http://192.168.1.x:80',
```

Then run:

```bash
flutter run
```

---

## Development

### Backend only (no Docker)

```bash
cd backend

# Start only MySQL in Docker
docker compose up mysql -d

# Run Go server locally
APP_ENV=development \
DB_HOST=localhost \
DB_PASSWORD=secret \
JWT_SECRET=dev-secret-at-least-32-characters-long \
go run ./cmd/server
```

### Python scanner only

```bash
cd python-scanner
pip install -r requirements.txt
uvicorn main:app --host 0.0.0.0 --port 8001 --reload
```

### Rebuild a single service

```bash
docker compose up --build go-api -d
docker compose up --build python-scanner -d
```

---

## Environment Variables

| Variable                      | Default                      | Description                                            |
| ----------------------------- | ---------------------------- | ------------------------------------------------------ |
| `APP_ENV`                     | `development`                | `development` or `production`                          |
| `SERVER_ADDR`                 | `:8080`                      | Go server listen address                               |
| `DB_HOST`                     | `mysql`                      | MySQL hostname                                         |
| `DB_PORT`                     | `3306`                       | MySQL port                                             |
| `DB_USER`                     | `root`                       | MySQL user                                             |
| `DB_PASSWORD`                 | —                            | **Required.** MySQL password                           |
| `DB_NAME`                     | `doc_scanner`                | Database name                                          |
| `JWT_SECRET`                  | —                            | **Required.** Min 32 chars. Use `openssl rand -hex 32` |
| `ACCESS_TOKEN_EXPIRE_MINUTES` | `30`                         | Access token TTL                                       |
| `REFRESH_TOKEN_EXPIRE_DAYS`   | `30`                         | Refresh token TTL                                      |
| `PDF_STORAGE_DIR`             | `./pdfs`                     | Directory where PDFs are saved                         |
| `SCANNER_URL`                 | `http://python-scanner:8001` | Python scanner base URL                                |

---

## How Document Scanning Works

```
Flutter app
    │  multipart/form-data (image file)
    ▼
Go API  POST /scan
    │  forwards image bytes
    ▼
Python scanner  POST /scan
    │
    ├─ 1. Grayscale + Gaussian blur
    ├─ 2. Canny edge detection
    ├─ 3. Find largest 4-sided contour (the document)
    ├─ 4. Perspective warp → flat, rectangular view
    ├─ 5. Denoise + sharpen + adaptive threshold
    └─ 6. Embed into A4 PDF (ReportLab)
    │
    ▼  PDF bytes
Go API
    ├─ Saves PDF to /app/pdfs/user_{id}_{ts}.pdf
    ├─ Inserts scanned_documents row
    └─ Updates user_storage stats
    │
    ▼  201 Created + document JSON
Flutter app
    └─ Shows success · refreshes Documents tab
```

---

## Database Schema

```
users
  id, email (unique), username (unique), password (bcrypt),
  full_name, avatar_url, is_active, is_verified, created_at, updated_at

refresh_tokens
  id, user_id → users, token (unique), expires_at, revoked, created_at

scanned_documents
  id, user_id → users, title, original_filename,
  pdf_path, pdf_size_bytes, page_count, status, created_at, updated_at

user_storage
  user_id → users (PK), total_bytes, document_count, updated_at
```

GORM `AutoMigrate` runs on every Go API startup and creates or alters tables as needed.

---

## Flutter State Flow

```
AuthBloc
  AuthCheckRequested  ──►  AuthAuthenticated  ──►  AppShell
                      ──►  AuthUnauthenticated ──►  LoginPage

  AuthLoginRequested  ──►  AuthLoading ──► AuthAuthenticated
                                      ──► AuthFailure (snackbar)

DocumentsBloc
  DocumentsLoadRequested   ──►  DocumentsLoading ──► DocumentsLoaded
  DocumentDownloadRequested ──► DocumentDownloading ──► DocumentDownloaded
  DocumentDeleteRequested  ──►  updates list in-memory

ScanBloc
  ScanFromCamera  ──►  ScanPickingImage ──► ScanUploading ──► ScanSuccess
  ScanFromGallery ──►                                     ──► ScanFailure
```

The Dio client automatically:

1. Attaches `Authorization: Bearer <token>` to every request
2. On `401`, silently refreshes the token using the stored refresh token
3. Retries the original request with the new token
4. If refresh fails, clears tokens and the `_AuthGate` redirects to login

---

## Production Checklist

- [ ] Set a strong `DB_PASSWORD` and `JWT_SECRET` in `.env`
- [ ] Add TLS to Nginx (Let's Encrypt / Certbot)
- [ ] Set `APP_ENV=production` to enable JSON logging and Gin release mode
- [ ] Mount a persistent volume for `pdfs_data`
- [ ] Set `DB_HOST` to point at your managed database in cloud
- [ ] Remove or restrict the `3306:3306` port mapping in docker-compose
- [ ] Set `flutter_pdfview` permissions in `Info.plist` for iOS

---

## License

MIT

## Prepare By

PHIN SOPHEAKNEATH AKA KIDZZ (for project practicum)
