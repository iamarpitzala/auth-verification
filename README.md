# Go Authentication Backend

A production-ready Go backend for email-based authentication with verification using Google Apps Script for email delivery.

## 🏗️ Architecture

```
auth-backend/
├── cmd/
│   └── server/           # Application entrypoint
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and migrations
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models and request/response structs
│   └── services/        # Business logic
├── pkg/
│   └── jwt/             # JWT utilities (reusable package)
├── scripts/             # Utility scripts
├── deployments/         # Deployment configurations
└── docs/                # Documentation
```

## ✨ Features

- **Clean Architecture**: Organized with proper separation of concerns
- **Email Verification Flow**: 6-digit code verification system
- **JWT Authentication**: Secure token-based authentication
- **PostgreSQL Integration**: Full database support with migrations
- **Docker Support**: Complete containerization setup
- **Google Apps Script Integration**: Email delivery via webhook
- **IPv6 Compatible**: Works with Supabase and modern networks
- **Production Ready**: Proper error handling, logging, and security

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL (or Supabase account)
- Google Apps Script setup for email delivery

### 1. Clone and Setup
```bash
git clone <repository-url>
cd auth-backend
cp .env.example .env
# Edit .env with your actual values
```

### 2. Install Dependencies
```bash
make deps
```

### 3. Run the Application
```bash
# Development
make dev

# Or build and run
make build
make run
```

### 4. Test the API
```bash
chmod +x scripts/test_api.sh
./scripts/test_api.sh
```

## 🔧 Configuration

### Environment Variables
```env
DATABASE_URL=postgresql://postgres:[PASSWORD]@db.gwfotffujwjyskeciund.supabase.co:5432/postgres
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
PORT=8080
EMAIL_WEBHOOK_URL=https://script.google.com/macros/s/[YOUR-SCRIPT-ID]/exec
```

### Google Apps Script Setup
1. Go to [Google Apps Script](https://script.google.com/)
2. Create a new project
3. Copy code from `deployments/google-apps-script.js`
4. Deploy as web app with "Anyone" permissions
5. Update `EMAIL_WEBHOOK_URL` in your `.env`

## API Endpoints

### Public Routes

1. **Request Verification Code**
   ```
   POST /auth/request-verification
   Content-Type: application/json
   
   {
     "email": "user@example.com"
   }
   ```

2. **Verify Code**
   ```
   POST /auth/verify-code
   Content-Type: application/json
   
   {
     "email": "user@example.com",
     "code": "123456"
   }
   ```

3. **Set Password**
   ```
   POST /auth/set-password
   Content-Type: application/json
   
   {
     "email": "user@example.com",
     "password": "your-password"
   }
   ```

4. **Login**
   ```
   POST /auth/login
   Content-Type: application/json
   
   {
     "email": "user@example.com",
     "password": "your-password"
   }
   ```

### Protected Routes

1. **Get Profile**
   ```
   GET /api/profile
   Authorization: Bearer <jwt-token>
   ```

## 🐳 Docker Support

### Using Docker Compose
```bash
# Start the application with PostgreSQL
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

### Manual Docker Build
```bash
# Build image
docker build -t auth-backend .

# Run container
docker run -p 8080:8080 \
  -e DATABASE_URL="your-db-url" \
  -e JWT_SECRET="your-secret" \
  auth-backend
```

## 🛠️ Development

### Available Make Commands
```bash
make build      # Build the application
make run        # Build and run
make dev        # Run with hot reload (requires air)
make test       # Run tests
make clean      # Clean build artifacts
make deps       # Install dependencies
make fmt        # Format code
make lint       # Lint code (requires golangci-lint)
```

### Database Schema
The application automatically creates required tables:
- `users`: User information and verification status
- `verification_codes`: Temporary verification codes with expiration

## 📝 License
This project is licensed under the MIT License.