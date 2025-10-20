# OpenAPI Setup

This project uses OpenAPI 3.0.3 for API specification and code generation.

## Structure

- `openapi/api.yaml` - OpenAPI specification
- `openapi/config.yaml` - Code generation configuration
- `server/internal/api/generated.go` - Generated Go code (do not edit)

## Commands

### Generate API Code

```bash
cd server
make generate-api
```

### Clean Generated Files

```bash
cd server
make clean
```

## How It Works

1. **OpenAPI Spec**: Define your API endpoints, request/response schemas in `openapi/api.yaml`
2. **Code Generation**: Run `make generate-api` to generate Go types and server interface
3. **Implementation**: Implement the `ServerInterface` in your handlers
4. **Integration**: Use `api.HandlerFromMux()` to create HTTP handlers

## Generated Types

The code generator creates:

- `api.HealthResponse` - Health check response
- `api.UserAccount` - User account model
- `api.ErrorResponse` - Error response model
- `api.ServerInterface` - Interface to implement

## Type Conversion

Note: Database models may have different pointer levels than API models. For example:

- Database: `**time.Time` (double pointer)
- API: `*time.Time` (single pointer)

Handle conversions in your handlers as needed.
