# Reel Farm

A full-stack application with React frontend and Go backend, featuring Supabase authentication and PostgreSQL database.

## 🚀 Quick Start

### Prerequisites

- **Node.js** (v18 or higher)
- **Go** (v1.21 or higher)
- **PostgreSQL** database (or use Supabase)
- **Git**

### Installation

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd reel-farm
   ```

2. **Install dependencies**

   ```bash
   npm run install:all
   ```

3. **Set up environment variables**

   Create a `.env` file in the `server/` directory:

   ```bash
   # Database
   DATABASE_URL=postgresql://username:password@localhost:5432/reel_farm

   # JWT Secret (generate a secure random string)
   JWT_SECRET=your-super-secret-jwt-key-here

   # Server Port (optional, defaults to 3000)
   PORT=3000
   ```

4. **Set up Supabase** (for frontend authentication)

   Create a `.env` file in the `frontend/` directory:

   ```bash
   VITE_SUPABASE_URL=your-supabase-project-url
   VITE_SUPABASE_ANON_KEY=your-supabase-anon-key
   ```

### Development

#### Option 1: Run Everything Together (Recommended)

```bash
npm run dev
```

This will start both frontend and backend with hot reloading:

- **Frontend**: http://localhost:5173 (Vite dev server)
- **Backend**: http://localhost:3000 (Go server with Air hot reload)

#### Option 2: Run Services Separately

**Terminal 1 - Backend:**

```bash
npm run dev:backend
```

**Terminal 2 - Frontend:**

```bash
npm run dev:frontend
```

#### Option 3: Manual Setup

**Backend:**

```bash
cd server
go run cmd/main.go
```

**Frontend:**

```bash
cd frontend
npm run dev
```

## 🔧 Development Features

### Hot Reloading

- **Frontend**: Vite provides instant hot module replacement
- **Backend**: Air automatically rebuilds and restarts the Go server on file changes

### Authentication Modes

The backend supports two authentication modes:

1. **Production Mode** (default): Requires valid JWT tokens

   ```bash
   npm run dev:backend
   ```

2. **Development Mode** (no authentication): Uses hardcoded user ID for testing
   ```bash
   npm run dev:backend:noauth
   ```

### API Endpoints

- **Health Check**: `GET /health` (no auth required)
- **User Account**: `GET /user` (requires authentication)

## 📁 Project Structure

```
reel-farm/
├── frontend/                 # React + Vite frontend
│   ├── src/
│   │   ├── components/      # Reusable UI components
│   │   ├── features/        # Feature-specific components
│   │   ├── contexts/        # React contexts (Auth, etc.)
│   │   ├── hooks/           # Custom React hooks
│   │   ├── lib/             # Utilities and API client
│   │   └── providers/       # Context providers
│   └── package.json
├── server/                  # Go backend
│   ├── cmd/                 # Application entry point
│   ├── internal/            # Private application code
│   │   ├── api/             # Generated OpenAPI code
│   │   ├── handler/         # HTTP handlers
│   │   ├── middleware/      # CORS, Auth middleware
│   │   ├── repository/      # Data access layer
│   │   └── service/         # Business logic
│   ├── db/                  # Generated database code
│   ├── sql/                 # SQL queries and schema
│   └── go.mod
├── openapi/                 # OpenAPI specification
│   └── api.yaml
├── infra/                   # Terraform infrastructure
└── package.json             # Root package.json for scripts
```

## 🛠️ Available Scripts

### Root Level Scripts

```bash
npm run dev                    # Start both frontend and backend
npm run dev:frontend          # Start only frontend
npm run dev:backend           # Start backend with authentication
npm run dev:backend:noauth    # Start backend without authentication
npm run build                  # Build both frontend and backend
npm run build:frontend        # Build only frontend
npm run build:backend         # Build only backend
npm run install:all           # Install all dependencies
npm run clean                 # Clean build artifacts
npm run generate-api          # Generate TypeScript API client
```

### Frontend Scripts

```bash
cd frontend
npm run dev                   # Start Vite dev server
npm run build                 # Build for production
npm run preview               # Preview production build
npm run lint                  # Run ESLint
npm run generate-api          # Generate API client from OpenAPI spec
```

### Backend Scripts

```bash
cd server
go run cmd/main.go           # Run server directly
go run cmd/main.go --noAuth  # Run without authentication
make generate-api            # Generate OpenAPI Go code
make clean                  # Clean generated files
```

## 🔐 Authentication

### Frontend (Supabase)

The frontend uses Supabase for authentication:

1. Users sign in through Supabase Auth
2. Access tokens are stored in localStorage
3. API client automatically includes tokens in requests

### Backend (JWT)

The backend validates JWT tokens:

1. Extracts user ID from JWT token
2. Passes user ID through request context
3. API handlers use user ID for data access

### Development Mode

For local development, you can bypass authentication:

```bash
npm run dev:backend:noauth
```

This uses a hardcoded user ID (`65a950f6-a3b0-4be2-824a-b99051d5a62f`) for testing.

## 🗄️ Database

### Schema

The application uses PostgreSQL with the following main tables:

- `user_accounts`: User profile and subscription information

### Migrations

Database migrations are located in the `migrations/` directory and can be run using:

```bash
./migrate.sh
```

### Code Generation

Database code is generated using SQLC:

```bash
cd server
make generate-db
```

## 🌐 API

### OpenAPI Specification

The API is defined in `openapi/api.yaml` and generates:

- **Go backend**: Types and HTTP handlers
- **TypeScript frontend**: Type-safe API client

### Generating API Code

**Backend (Go):**

```bash
cd server
make generate-api
```

**Frontend (TypeScript):**

```bash
npm run generate-api
```

## 🚀 Deployment

### Infrastructure

The project includes Terraform configuration for AWS deployment:

```bash
cd infra
terraform init
terraform plan
terraform apply
```

### Environment Variables

Set the following environment variables in production:

- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: Secret key for JWT token validation
- `PORT`: Server port (default: 3000)

## 🐛 Troubleshooting

### Common Issues

1. **CORS Errors**: Make sure both frontend and backend are running on the correct ports
2. **Authentication Errors**: Check that JWT_SECRET is set and matches between services
3. **Database Connection**: Verify DATABASE_URL is correct and database is accessible
4. **Hot Reload Not Working**: Ensure Air is installed (`go install github.com/cosmtrek/air@latest`)

### Debug Mode

Enable debug logging:

```bash
# Backend
cd server
go run cmd/main.go --debug

# Frontend
cd frontend
npm run dev -- --debug
```

## 📝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.
