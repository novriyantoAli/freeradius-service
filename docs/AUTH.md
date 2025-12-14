# Authentication (Auth) Module Documentation

## Overview

The Auth module provides authentication credential management for RADIUS protocol. It enables creating user authentication credentials with atomic transaction support, ensuring data consistency between radcheck and radreply entries.

## Features

- **Atomic Transactions**: All authentication operations are wrapped in database transactions using the TransactionManagerI interface
- **RADIUS Integration**: Creates User-Password radcheck entries and associated RADIUS reply attributes
- **REST API**: HTTP REST endpoint for credential creation with JSON request/response
- **gRPC API**: Protocol Buffer service definition for credential creation over gRPC
- **Comprehensive Logging**: Zap logger integration at all layers (service, handler, gRPC)
- **Input Validation**: Strict validation of username and password requirements
- **Clean Architecture**: Follows DDD pattern with vertical slicing and separation of concerns

## Architecture

### Module Structure

```
internal/application/auth/
├── handler/
│   ├── auth.handler.go          # REST HTTP handlers
│   └── auth.grpc.handler.go     # gRPC service implementation
├── service/
│   ├── auth.service.go          # Business logic
│   └── auth.service_test.go     # Service tests
├── dto/
│   └── auth.dto.go              # Data transfer objects
├── entity/
│   └── auth.entity.go           # Domain entities (if applicable)
└── module.go                    # Dependency injection setup
```

### API Proto

```
api/proto/auth/
├── auth.proto                   # Proto3 service definition
├── auth.pb.go                   # Generated message definitions
└── auth_grpc.pb.go              # Generated gRPC service code
```

## Service Layer

### AuthService Interface

```go
type AuthService interface {
	CreateAuth(ctx context.Context, req *dto.CreateAuthRequest) (*dto.CreateAuthResponse, error)
}
```

### CreateAuth Method

Creates authentication credentials by:
1. Validating username and password (both required)
2. Creating User-Password radcheck entry with `:=` operator
3. Creating additional radcheck attributes (if provided)
4. Creating radreply entries (if provided)
5. All operations executed atomically within a transaction

**Key Features:**
- Transaction atomicity via `txManager.WithinTransaction()`
- Context propagation with `WithTx()` for database operations
- Password masking in response (`***` instead of actual password)
- Default operators: `:=` for radcheck, `+=` for radreply
- Error handling with descriptive messages

**Example:**
```go
req := &dto.CreateAuthRequest{
    Username: "john_doe",
    Password: "secure_password",
    Attributes: []dto.CreateAuthAttribute{
        {Attribute: "Framed-IP-Address", Value: "192.168.1.100", Op: ":="},
    },
    ReplyAttributes: []dto.CreateAuthAttribute{
        {Attribute: "Service-Type", Value: "Framed-User", Op: "+="},
    },
}

response, err := authService.CreateAuth(ctx, req)
```

## REST API Handler

### Endpoint

- **Method**: `POST`
- **Path**: `/api/v1/auth`
- **Status Code**: `201 Created` on success

### Request Body

```json
{
  "username": "john_doe",
  "password": "secure_password",
  "attributes": [
    {
      "attribute": "Framed-IP-Address",
      "value": "192.168.1.100",
      "op": ":="
    }
  ],
  "reply_attributes": [
    {
      "attribute": "Service-Type",
      "value": "Framed-User",
      "op": "+="
    }
  ]
}
```

### Response (201 Created)

```json
{
  "username": "john_doe",
  "password": "***",
  "attributes": [
    {
      "id": 1,
      "attribute": "User-Password",
      "value": "***",
      "op": ":="
    },
    {
      "id": 2,
      "attribute": "Framed-IP-Address",
      "value": "192.168.1.100",
      "op": ":="
    }
  ],
  "reply_attributes": [
    {
      "id": 1,
      "attribute": "Service-Type",
      "value": "Framed-User",
      "op": "+="
    }
  ]
}
```

### Error Responses

**400 Bad Request** - Missing or invalid fields:
```json
{
  "message": "Key: 'CreateAuthRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag"
}
```

**500 Internal Server Error** - Database or service error:
```json
{
  "message": "failed to create radcheck entry"
}
```

## gRPC API Handler

### Service Definition

Located in `api/proto/auth/auth.proto`

```protobuf
service AuthService {
  rpc CreateAuth(CreateAuthRequest) returns (CreateAuthResponse);
}
```

### Request Message

```protobuf
message CreateAuthRequest {
  string username = 1;                              // Required
  string password = 2;                              // Required
  repeated CreateAuthAttribute attributes = 3;     // Optional
  repeated CreateAuthAttribute reply_attributes = 4; // Optional
}

message CreateAuthAttribute {
  string attribute = 1;  // Attribute name
  string value = 2;      // Attribute value
  string op = 3;         // Operator (optional, defaults in service)
}
```

### Response Message

```protobuf
message CreateAuthResponse {
  string username = 1;
  string password = 2;
  repeated AuthCreateAttrResponse attributes = 3;
  repeated AuthCreateAttrResponse reply_attributes = 4;
}

message AuthCreateAttrResponse {
  uint32 id = 1;         // Database ID
  string attribute = 2;  // Attribute name
  string value = 3;      // Attribute value
  string op = 4;         // Operator used
}
```

### gRPC Error Codes

- **`INVALID_ARGUMENT`** (3): Username or password is missing/empty
- **`INTERNAL`** (13): Database operation failed or service error

### Example gRPC Client Usage

```go
client := auth.NewAuthServiceClient(conn)
req := &auth.CreateAuthRequest{
    Username: "john_doe",
    Password: "secure_password",
    Attributes: []*auth.CreateAuthAttribute{
        {Attribute: "Framed-IP-Address", Value: "192.168.1.100", Op: ":="},
    },
}

resp, err := client.CreateAuth(ctx, req)
if err != nil {
    // Handle error
}
```

## Data Transfer Objects (DTOs)

### CreateAuthRequest

```go
type CreateAuthRequest struct {
    Username        string                  `json:"username" binding:"required"`
    Password        string                  `json:"password" binding:"required"`
    Attributes      []CreateAuthAttribute   `json:"attributes"`
    ReplyAttributes []CreateAuthAttribute   `json:"reply_attributes"`
}
```

### CreateAuthAttribute

```go
type CreateAuthAttribute struct {
    Attribute string `json:"attribute" binding:"required"`
    Value     string `json:"value" binding:"required"`
    Op        string `json:"op"`  // Optional, defaults in service
}
```

### CreateAuthResponse

```go
type CreateAuthResponse struct {
    Username        string                    `json:"username"`
    Password        string                    `json:"password"`
    Attributes      []AuthCreateAttrResponse  `json:"attributes"`
    ReplyAttributes []AuthCreateAttrResponse  `json:"reply_attributes"`
}
```

### AuthCreateAttrResponse

```go
type AuthCreateAttrResponse struct {
    ID        uint   `json:"id"`
    Attribute string `json:"attribute"`
    Value     string `json:"value"`
    Op        string `json:"op"`
}
```

## Dependency Injection Setup

### Module Provider

The auth module is configured in `internal/application/auth/module.go`:

```go
var Module = fx.Module(
    "auth",
    fx.Provide(
        provideAuthService,
        provideAuthHandler,
        provideAuthGrpcHandler,
    ),
)
```

### Service Provider

```go
func provideAuthService(
    radcheckRepo radcheckrepo.RadcheckRepository,
    radreplyRepo radreplyrepo.RadreplyRepository,
    txManager database.TransactionManagerI,
) service.AuthService {
    return service.NewAuthService(radcheckRepo, radreplyRepo, txManager)
}
```

### Handler Providers

```go
func provideAuthHandler(authService service.AuthService) *handler.AuthHandler {
    return handler.NewAuthHandler(authService)
}

func provideAuthGrpcHandler(authService service.AuthService, logger *zap.Logger) *handler.AuthGrpcHandler {
    return handler.NewAuthGrpcHandler(authService, logger)
}
```

## Transaction Management

### TransactionManagerI Interface

Enables testable transaction management:

```go
type TransactionManagerI interface {
    WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
```

### Transaction Flow

1. **Begin Transaction**: `txManager.WithinTransaction()` starts a database transaction
2. **Inject Context**: Context is wrapped with `WithTx(ctx, tx)` 
3. **Execute Operations**: All repository calls use the injected transaction context
4. **Commit/Rollback**: Automatic on success or error

**Example Implementation:**
```go
err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
    // All repo calls use txCtx which contains the transaction
    passwordRadcheck := &radcheckentity.Radcheck{...}
    return s.radcheckRepo.Create(txCtx, passwordRadcheck)
})
```

## Testing

### Unit Tests

Located in:
- `internal/application/auth/service/auth.service_test.go`
- `internal/application/auth/handler/auth.handler_test.go`

#### Test Coverage

**Service Tests:**
1. `TestAuthService_CreateAuth_Success` - Validates successful creation with mocked dependencies
2. `TestAuthService_CreateAuth_MissingUsername` - Validates username validation
3. `TestAuthService_CreateAuth_MissingPassword` - Validates password validation

**Handler Tests:**
1. `TestAuthHandler_CreateAuth_Success` - Validates HTTP request/response handling
2. `TestAuthHandler_CreateAuth_MissingUsername` - Validates request binding and validation

#### Mock Setup

```go
mockTxManager := &database.MockTransactionManager{
    WithinTransactionFunc: func(ctx context.Context, fn func(context.Context) error) error {
        return fn(ctx) // Execute function directly in tests
    },
}
```

### Running Tests

```bash
# Run all auth tests
go test ./internal/application/auth/... -v

# Run with coverage
go test ./internal/application/auth/... -v -cover

# Run specific test
go test ./internal/application/auth/service -v -run TestAuthService_CreateAuth_Success
```

## Integration Points

### Radcheck Module

Creates `User-Password` entries and additional radcheck attributes:
- **Repository**: `internal/application/radcheck/repository/radcheck.repo.go`
- **Entity**: `internal/application/radcheck/entity/radcheck.entity.go`

### Radreply Module

Creates RADIUS reply attributes:
- **Repository**: `internal/application/radreply/repository/radreply.repo.go`
- **Entity**: `internal/application/radreply/entity/radreply.entity.go`

### Database Module

Provides transaction management:
- **TransactionManagerI**: `internal/pkg/database/database.go`
- **WithinTransaction**: Wraps operations in database transactions

### Logger Module

Provides structured logging:
- **Zap Logger**: `internal/pkg/logger/logger.go`
- **Usage**: Logging at handler and gRPC layers for debugging

## Security Considerations

1. **Password Masking**: Passwords always returned as `***` in API responses
2. **Validation**: Username and password are required and validated before processing
3. **Transaction Safety**: Atomic operations prevent partial state creation
4. **Error Messages**: Generic error messages in responses to prevent information leakage

## Error Handling

### Common Errors

| Error | Cause | Resolution |
|-------|-------|-----------|
| `username is required` | Username not provided or empty | Include valid username in request |
| `password is required` | Password not provided or empty | Include valid password in request |
| `failed to create radcheck entry` | Database error | Check database connectivity and logs |
| `failed to create radreply entry` | Database error | Check database connectivity and logs |
| `INVALID_ARGUMENT` (gRPC) | Missing required fields | Validate username/password in request |
| `INTERNAL` (gRPC) | Service/database error | Check server logs and database |

## Debugging

### Enable Detailed Logging

All operations log to zap logger at multiple levels:
- **Info**: Service method execution, handler requests
- **Warn**: Validation failures, missing fields
- **Error**: Database errors, transaction failures

### Example Log Output

```
{"level":"info","ts":1702513200.123,"caller":"handler/auth.grpc.handler.go:35","msg":"CreateAuth gRPC request","username":"john_doe"}
{"level":"info","ts":1702513200.124,"caller":"service/auth.service.go:50","msg":"Creating authentication credentials","username":"john_doe"}
{"level":"info","ts":1702513200.125,"caller":"handler/auth.grpc.handler.go:65","msg":"CreateAuth gRPC response sent successfully"}
```

## Future Enhancements

- [ ] Add UpdateAuth endpoint to modify existing credentials
- [ ] Add DeleteAuth endpoint to remove credentials
- [ ] Add ListAuth endpoint with pagination
- [ ] Add GetAuth endpoint to retrieve specific credentials
- [ ] Implement gRPC middleware for authentication/authorization
- [ ] Add rate limiting for credential creation
- [ ] Add audit logging for security events
- [ ] Add bulk creation endpoint with transaction support
- [ ] Add password hashing strategies (bcrypt, Argon2)
- [ ] Add session management support
