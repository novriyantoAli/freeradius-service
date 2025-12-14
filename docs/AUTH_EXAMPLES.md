# Auth Module - Code Examples

## Table of Contents

1. [Service Layer Examples](#service-layer-examples)
2. [Handler Usage Examples](#handler-usage-examples)
3. [gRPC Client Examples](#grpc-client-examples)
4. [Testing Examples](#testing-examples)
5. [Integration Examples](#integration-examples)

---

## Service Layer Examples

### Using AuthService Directly

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/novriyantoAli/freeradius-service/internal/application/auth/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/service"
)

func CreateAuthCredentials(authService service.AuthService) error {
	ctx := context.Background()

	// Prepare request
	req := &dto.CreateAuthRequest{
		Username: "john_doe",
		Password: "secure_password_123",
		Attributes: []dto.CreateAuthAttribute{
			{
				Attribute: "Framed-IP-Address",
				Value:     "192.168.1.100",
				Op:        ":=",
			},
			{
				Attribute: "Framed-IP-Netmask",
				Value:     "255.255.255.0",
				Op:        ":=",
			},
		},
		ReplyAttributes: []dto.CreateAuthAttribute{
			{
				Attribute: "Service-Type",
				Value:     "Framed-User",
				Op:        "+=",
			},
			{
				Attribute: "Session-Timeout",
				Value:     "3600",
				Op:        "+=",
			},
		},
	}

	// Create authentication
	response, err := authService.CreateAuth(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create auth: %w", err)
	}

	// Log results
	log.Printf("Created authentication for: %s", response.Username)
	log.Printf("Attributes created: %d", len(response.Attributes))
	log.Printf("Reply attributes created: %d", len(response.ReplyAttributes))

	return nil
}
```

### Error Handling in Service

```go
func CreateAuthWithErrorHandling(authService service.AuthService) {
	ctx := context.Background()

	req := &dto.CreateAuthRequest{
		Username: "",  // Empty username - will fail validation
		Password: "password123",
	}

	response, err := authService.CreateAuth(ctx, req)
	if err != nil {
		switch err.Error() {
		case "username is required":
			fmt.Println("Username validation failed")
		case "password is required":
			fmt.Println("Password validation failed")
		default:
			fmt.Printf("Service error: %v\n", err)
		}
		return
	}

	fmt.Printf("Success: %v\n", response)
}
```

### Batch Creating Multiple Users

```go
func CreateMultipleUsers(authService service.AuthService) error {
	users := []struct {
		Username string
		Password string
		IP       string
	}{
		{"user1", "pass1", "10.0.0.1"},
		{"user2", "pass2", "10.0.0.2"},
		{"user3", "pass3", "10.0.0.3"},
	}

	for _, user := range users {
		req := &dto.CreateAuthRequest{
			Username: user.Username,
			Password: user.Password,
			Attributes: []dto.CreateAuthAttribute{
				{
					Attribute: "Framed-IP-Address",
					Value:     user.IP,
					Op:        ":=",
				},
			},
		}

		response, err := authService.CreateAuth(context.Background(), req)
		if err != nil {
			log.Printf("Failed to create user %s: %v", user.Username, err)
			continue
		}

		log.Printf("Created user %s with %d attributes", response.Username, len(response.Attributes))
	}

	return nil
}
```

---

## Handler Usage Examples

### REST API with cURL

#### Create Basic User

```bash
curl -X POST http://localhost:8080/api/v1/auth \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "alice_secure_pass"
  }'
```

#### Create User with IP Address

```bash
curl -X POST http://localhost:8080/api/v1/auth \
  -H "Content-Type: application/json" \
  -d '{
    "username": "bob",
    "password": "bob_password",
    "attributes": [
      {
        "attribute": "Framed-IP-Address",
        "value": "192.168.100.50"
      }
    ]
  }'
```

#### Create VPN User with Multiple Attributes

```bash
curl -X POST http://localhost:8080/api/v1/auth \
  -H "Content-Type: application/json" \
  -d '{
    "username": "vpn_user",
    "password": "vpn_secure_password",
    "attributes": [
      {
        "attribute": "Service-Type",
        "value": "Framed-User",
        "op": ":="
      },
      {
        "attribute": "Framed-Protocol",
        "value": "PPP",
        "op": ":="
      },
      {
        "attribute": "Framed-IP-Address",
        "value": "10.8.0.100",
        "op": ":="
      }
    ],
    "reply_attributes": [
      {
        "attribute": "Session-Timeout",
        "value": "86400",
        "op": "+="
      }
    ]
  }'
```

### Go HTTP Client

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/novriyantoAli/freeradius-service/internal/application/auth/dto"
)

func CreateAuthViaHTTP(baseURL string) error {
	req := dto.CreateAuthRequest{
		Username: "http_user",
		Password: "http_password",
		Attributes: []dto.CreateAuthAttribute{
			{
				Attribute: "Framed-IP-Address",
				Value:     "10.0.0.50",
			},
		},
	}

	body, _ := json.Marshal(req)

	httpReq, err := http.NewRequest("POST", baseURL+"/api/v1/auth", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var response dto.CreateAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("Successfully created auth for: %s\n", response.Username)
	fmt.Printf("Attributes: %v\n", response.Attributes)

	return nil
}
```

---

## gRPC Client Examples

### Go gRPC Client

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/novriyantoAli/freeradius-service/api/proto/auth"
	"google.golang.org/grpc"
)

func CreateAuthViaGRPC(address string) error {
	// Connect to gRPC server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	defer conn.Close()

	client := auth.NewAuthServiceClient(conn)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Prepare request
	req := &auth.CreateAuthRequest{
		Username: "grpc_user",
		Password: "grpc_secure_password",
		Attributes: []*auth.CreateAuthAttribute{
			{
				Attribute: "Framed-IP-Address",
				Value:     "192.168.1.150",
				Op:        ":=",
			},
		},
		ReplyAttributes: []*auth.CreateAuthAttribute{
			{
				Attribute: "Service-Type",
				Value:     "Framed-User",
				Op:        "+=",
			},
		},
	}

	// Call gRPC method
	resp, err := client.CreateAuth(ctx, req)
	if err != nil {
		return fmt.Errorf("gRPC call failed: %w", err)
	}

	log.Printf("Created user: %s\n", resp.Username)
	log.Printf("Attributes: %d\n", len(resp.Attributes))
	log.Printf("Reply Attributes: %d\n", len(resp.ReplyAttributes))

	return nil
}
```

### Python gRPC Client

```python
import grpc
import sys
import time
sys.path.insert(0, 'api/proto')

from auth import auth_pb2, auth_pb2_grpc

def create_auth_via_grpc(address='localhost:50051'):
    """Create authentication credentials via gRPC"""
    try:
        # Connect to server
        channel = grpc.secure_channel(address, grpc.ssl_channel_credentials())
        stub = auth_pb2_grpc.AuthServiceStub(channel)

        # Prepare request
        request = auth_pb2.CreateAuthRequest(
            username='python_user',
            password='python_secure_pass',
            attributes=[
                auth_pb2.CreateAuthAttribute(
                    attribute='Framed-IP-Address',
                    value='192.168.1.200',
                    op=':='
                )
            ],
            reply_attributes=[
                auth_pb2.CreateAuthAttribute(
                    attribute='Service-Type',
                    value='Framed-User',
                    op='+='
                )
            ]
        )

        # Call gRPC method
        response = stub.CreateAuth(request, timeout=5)

        print(f"Created user: {response.username}")
        print(f"Attributes: {len(response.attributes)}")
        print(f"Reply Attributes: {len(response.reply_attributes)}")

        return response

    except grpc.RpcError as e:
        print(f"gRPC error: {e.code()} - {e.details()}")
        return None
    finally:
        channel.close()

if __name__ == '__main__':
    create_auth_via_grpc()
```

### JavaScript/TypeScript gRPC Client

```typescript
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';

const PROTO_PATH = './api/proto/auth/auth.proto';

async function createAuthViaGRPC(address: string = 'localhost:50051') {
  const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
  });

  const authProto = grpc.loadPackageDefinition(packageDefinition).auth;
  const client = new authProto.AuthService(address, grpc.credentials.createInsecure());

  const request = {
    username: 'js_user',
    password: 'js_secure_password',
    attributes: [
      {
        attribute: 'Framed-IP-Address',
        value: '192.168.1.250',
        op: ':='
      }
    ],
    reply_attributes: [
      {
        attribute: 'Service-Type',
        value: 'Framed-User',
        op: '+='
      }
    ]
  };

  return new Promise((resolve, reject) => {
    client.createAuth(request, (error: any, response: any) => {
      if (error) {
        console.error('gRPC error:', error);
        reject(error);
      } else {
        console.log('Created user:', response.username);
        console.log('Attributes:', response.attributes.length);
        resolve(response);
      }
    });
  });
}
```

---

## Testing Examples

### Unit Test with Mocked TransactionManager

```go
package service_test

import (
	"context"
	"testing"

	"github.com/novriyantoAli/freeradius-service/internal/application/auth/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/service"
	radcheckrepo "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/repository"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/database"
	"github.com/stretchr/testify/assert"
)

func TestCreateAuthSuccess(t *testing.T) {
	// Mock transaction manager
	mockTxManager := &database.MockTransactionManager{
		WithinTransactionFunc: func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx) // Execute function directly
		},
	}

	// Mock repositories
	mockRadcheckRepo := &radcheckrepo.MockRadcheckRepository{
		CreateFunc: func(ctx context.Context, rc *radcheckentity.Radcheck) error {
			rc.ID = 1 // Simulate ID assignment
			return nil
		},
	}

	// Create service
	authService := service.NewAuthService(
		mockRadcheckRepo,
		nil, // radreplyRepo
		mockTxManager,
	)

	// Test
	req := &dto.CreateAuthRequest{
		Username: "testuser",
		Password: "testpass",
	}

	resp, err := authService.CreateAuth(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "testuser", resp.Username)
	assert.Equal(t, "***", resp.Password)
	assert.Greater(t, len(resp.Attributes), 0)
}

func TestCreateAuthMissingUsername(t *testing.T) {
	authService := service.NewAuthService(nil, nil, nil)

	req := &dto.CreateAuthRequest{
		Username: "",
		Password: "testpass",
	}

	resp, err := authService.CreateAuth(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "username is required", err.Error())
}
```

### Handler Test with HTTP Client

```go
package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/handler"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/service"
	"github.com/stretchr/testify/assert"
)

func TestCreateAuthHandler(t *testing.T) {
	// Mock service
	mockService := &MockAuthService{
		CreateAuthFunc: func(ctx context.Context, req *dto.CreateAuthRequest) (*dto.CreateAuthResponse, error) {
			return &dto.CreateAuthResponse{
				Username: req.Username,
				Password: "***",
			}, nil
		},
	}

	// Create handler
	authHandler := handler.NewAuthHandler(mockService)

	// Setup router
	router := gin.New()
	authHandler.RegisterRoutes(router.Group("/api/v1"))

	// Create request
	body := dto.CreateAuthRequest{
		Username: "testuser",
		Password: "testpass",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(
		"POST",
		"/api/v1/auth",
		bytes.NewBuffer(jsonBody),
	)
	req.Header.Set("Content-Type", "application/json")

	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dto.CreateAuthResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "testuser", resp.Username)
}
```

---

## Integration Examples

### Full Application Integration

```go
package main

import (
	"context"
	"fmt"

	"github.com/novriyantoAli/freeradius-service/internal/application/auth/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/module"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	// Create DI container with auth module
	app := fx.New(
		// ... other modules ...
		module.Module, // Auth module
		fx.Invoke(runApplication),
	)

	app.Run()
}

func runApplication(authService service.AuthService, logger *zap.Logger) error {
	ctx := context.Background()

	// Create authentication
	req := &dto.CreateAuthRequest{
		Username: "app_user",
		Password: "app_secure_pass",
		Attributes: []dto.CreateAuthAttribute{
			{Attribute: "Framed-IP-Address", Value: "10.0.0.1"},
		},
	}

	resp, err := authService.CreateAuth(ctx, req)
	if err != nil {
		logger.Error("Failed to create auth", zap.Error(err))
		return err
	}

	logger.Info("Auth created successfully",
		zap.String("username", resp.Username),
		zap.Int("attributes", len(resp.Attributes)),
	)

	return nil
}
```

### Database Transaction with Context

```go
func ExampleTransactionFlow(authService service.AuthService) {
	ctx := context.Background()

	// The service automatically wraps operations in a transaction
	req := &dto.CreateAuthRequest{
		Username: "tx_user",
		Password: "tx_password",
		Attributes: []dto.CreateAuthAttribute{
			{Attribute: "Framed-IP-Address", Value: "192.168.0.50"},
		},
		ReplyAttributes: []dto.CreateAuthAttribute{
			{Attribute: "Session-Timeout", Value: "3600"},
		},
	}

	// All operations (radcheck create, radreply create) are atomic
	resp, err := authService.CreateAuth(ctx, req)
	if err != nil {
		// On error, transaction is automatically rolled back
		fmt.Printf("Transaction failed, all operations rolled back: %v\n", err)
		return
	}

	// On success, transaction is committed
	fmt.Printf("Transaction committed: %d attributes created\n", len(resp.Attributes))
}
```

---

## Monitoring and Logging Examples

### Logging Output

```
{"level":"info","ts":1702513200.123,"caller":"handler/auth.handler.go:40","msg":"CreateAuth request","username":"john_doe"}
{"level":"info","ts":1702513200.124,"caller":"service/auth.service.go:50","msg":"Creating authentication credentials","username":"john_doe"}
{"level":"info","ts":1702513200.125,"caller":"repository/radcheck.repo.go:60","msg":"Created radcheck entry","id":1,"attribute":"User-Password"}
{"level":"info","ts":1702513200.126,"caller":"handler/auth.handler.go:50","msg":"CreateAuth response sent","status":201}
```

### Error Logging

```
{"level":"warn","ts":1702513201.100,"caller":"handler/auth.handler.go:35","msg":"CreateAuth request validation failed","error":"username is required"}
{"level":"error","ts":1702513202.200,"caller":"service/auth.service.go:80","msg":"Database transaction failed","error":"connection refused"}
```

---

## Real-World Scenarios

### Scenario: RADIUS AAA Server Integration

```go
func IntegrationWithRADIUSServer(authService service.AuthService) error {
	// When a new user registers, create auth credentials
	users := []struct {
		username string
		password string
		groupID  string
	}{
		{"user1", "pass1", "premium"},
		{"user2", "pass2", "basic"},
	}

	for _, user := range users {
		req := &dto.CreateAuthRequest{
			Username: user.username,
			Password: user.password,
			Attributes: []dto.CreateAuthAttribute{
				{
					Attribute: "Group-Name",
					Value:     user.groupID,
					Op:        ":=",
				},
			},
			ReplyAttributes: []dto.CreateAuthAttribute{
				{
					Attribute: "Acct-Interim-Interval",
					Value:     "60",
					Op:        "+=",
				},
			},
		}

		_, err := authService.CreateAuth(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to create auth for %s: %w", user.username, err)
		}
	}

	return nil
}
```

This comprehensive examples guide demonstrates the Auth module's capabilities across different layers and use cases.
