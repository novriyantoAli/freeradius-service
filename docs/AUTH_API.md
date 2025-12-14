# Auth API Reference

## REST API

### Create Authentication Credentials

Creates a user with authentication credentials and optional RADIUS attributes.

#### Request

```http
POST /api/v1/auth HTTP/1.1
Host: api.example.com
Content-Type: application/json

{
  "username": "john_doe",
  "password": "secure_password123",
  "attributes": [
    {
      "attribute": "Framed-IP-Address",
      "value": "192.168.1.100",
      "op": ":="
    },
    {
      "attribute": "Framed-IP-Netmask",
      "value": "255.255.255.0",
      "op": ":="
    }
  ],
  "reply_attributes": [
    {
      "attribute": "Service-Type",
      "value": "Framed-User",
      "op": "+="
    },
    {
      "attribute": "Session-Timeout",
      "value": "3600",
      "op": "+="
    }
  ]
}
```

#### Response (201 Created)

```http
HTTP/1.1 201 Created
Content-Type: application/json

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
    },
    {
      "id": 3,
      "attribute": "Framed-IP-Netmask",
      "value": "255.255.255.0",
      "op": ":="
    }
  ],
  "reply_attributes": [
    {
      "id": 1,
      "attribute": "Service-Type",
      "value": "Framed-User",
      "op": "+="
    },
    {
      "id": 2,
      "attribute": "Session-Timeout",
      "value": "3600",
      "op": "+="
    }
  ]
}
```

#### Error Responses

**400 Bad Request** - Validation error:
```http
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "message": "Key: 'CreateAuthRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag"
}
```

**500 Internal Server Error** - Server error:
```http
HTTP/1.1 500 Internal Server Error
Content-Type: application/json

{
  "message": "failed to create radcheck entry"
}
```

#### Request Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `username` | string | Yes | Username for authentication (max 64 chars) |
| `password` | string | Yes | User password (max 253 chars) |
| `attributes` | array | No | Additional radcheck attributes |
| `reply_attributes` | array | No | RADIUS reply attributes |

#### Attribute Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `attribute` | string | Yes | Attribute name (e.g., "Framed-IP-Address") |
| `value` | string | Yes | Attribute value (max 253 chars) |
| `op` | string | No | Operator (`:=`, `==`, `+=`, etc.) Defaults: `:=` for radcheck, `+=` for radreply |

#### Common Operators

| Operator | Description | Usage |
|----------|-------------|-------|
| `:=` | Assign | radcheck (default), sets exact value |
| `==` | Equal | Comparisons and equality checks |
| `!=` | Not equal | Negative matching |
| `+=` | Append/Add | radreply (default), adds to existing values |
| `-=` | Remove | Removes specific values |
| `=~` | Regular expression | Pattern matching |
| `!~` | Negative regex | Pattern non-matching |
| `>` | Greater than | Numeric comparisons |
| `>=` | Greater or equal | Numeric comparisons |
| `<` | Less than | Numeric comparisons |
| `<=` | Less or equal | Numeric comparisons |

#### Example: Minimal Request

```json
{
  "username": "simple_user",
  "password": "password123"
}
```

**Response:**
```json
{
  "username": "simple_user",
  "password": "***",
  "attributes": [
    {
      "id": 1,
      "attribute": "User-Password",
      "value": "***",
      "op": ":="
    }
  ],
  "reply_attributes": []
}
```

#### Example: Full Request with Multiple Attributes

```json
{
  "username": "advanced_user",
  "password": "secure_pass_2024",
  "attributes": [
    {
      "attribute": "Framed-IP-Address",
      "value": "10.0.0.100",
      "op": ":="
    },
    {
      "attribute": "Framed-IP-Netmask",
      "value": "255.255.255.0",
      "op": ":="
    },
    {
      "attribute": "NAS-Port-Type",
      "value": "Ethernet",
      "op": ":="
    }
  ],
  "reply_attributes": [
    {
      "attribute": "Service-Type",
      "value": "Framed-User",
      "op": "+="
    },
    {
      "attribute": "Session-Timeout",
      "value": "7200",
      "op": "+="
    },
    {
      "attribute": "Idle-Timeout",
      "value": "1800",
      "op": "+="
    }
  ]
}
```

---

## gRPC API

### Service Definition

```protobuf
package auth;

service AuthService {
  rpc CreateAuth(CreateAuthRequest) returns (CreateAuthResponse);
}
```

### Messages

#### CreateAuthRequest

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
  string op = 3;         // Operator (optional)
}
```

#### CreateAuthResponse

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

| Code | Description | Cause |
|------|-------------|-------|
| `3 (INVALID_ARGUMENT)` | Invalid argument | Missing required fields (username/password) |
| `13 (INTERNAL)` | Internal error | Server error, database failure, or service error |

### Go Client Example

```go
package main

import (
	"context"
	"log"

	"github.com/novriyantoAli/freeradius-service/api/proto/auth"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := auth.NewAuthServiceClient(conn)

	req := &auth.CreateAuthRequest{
		Username: "john_doe",
		Password: "secure_password",
		Attributes: []*auth.CreateAuthAttribute{
			{Attribute: "Framed-IP-Address", Value: "192.168.1.100", Op: ":="},
		},
		ReplyAttributes: []*auth.CreateAuthAttribute{
			{Attribute: "Service-Type", Value: "Framed-User", Op: "+="},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := client.CreateAuth(ctx, req)
	if err != nil {
		log.Fatalf("gRPC error: %v", err)
	}

	log.Printf("Created user: %s", resp.Username)
	log.Printf("Attributes: %v", resp.Attributes)
}
```

### Python Client Example

```python
import grpc
import sys
sys.path.insert(0, 'api/proto')

from auth import auth_pb2, auth_pb2_grpc

def main():
    channel = grpc.insecure_channel('localhost:50051')
    stub = auth_pb2_grpc.AuthServiceStub(channel)

    request = auth_pb2.CreateAuthRequest(
        username='john_doe',
        password='secure_password',
        attributes=[
            auth_pb2.CreateAuthAttribute(
                attribute='Framed-IP-Address',
                value='192.168.1.100',
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

    try:
        response = stub.CreateAuth(request)
        print(f"Created user: {response.username}")
        print(f"Attributes: {response.attributes}")
    except grpc.RpcError as e:
        print(f"gRPC error: {e.code()} - {e.details()}")

if __name__ == '__main__':
    main()
```

### cURL Examples

#### Minimal Request

```bash
curl -X POST http://localhost:8080/api/v1/auth \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass123"
  }'
```

#### Full Request

```bash
curl -X POST http://localhost:8080/api/v1/auth \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "secure_password",
    "attributes": [
      {
        "attribute": "Framed-IP-Address",
        "value": "192.168.1.100",
        "op": ":="
      },
      {
        "attribute": "Framed-IP-Netmask",
        "value": "255.255.255.0",
        "op": ":="
      }
    ],
    "reply_attributes": [
      {
        "attribute": "Service-Type",
        "value": "Framed-User",
        "op": "+="
      },
      {
        "attribute": "Session-Timeout",
        "value": "3600",
        "op": "+="
      }
    ]
  }'
```

---

## HTTP Status Codes

| Status Code | Description | Reason |
|------------|-------------|--------|
| `201 Created` | Success | Authentication credentials created successfully |
| `400 Bad Request` | Client error | Invalid request body, missing required fields, validation failure |
| `500 Internal Server Error` | Server error | Database error, transaction failure, unexpected server error |

---

## Request/Response Examples

### Scenario 1: Create Basic User

**Request:**
```bash
POST /api/v1/auth
Content-Type: application/json

{
  "username": "alice",
  "password": "pass@2024"
}
```

**Response (201):**
```json
{
  "username": "alice",
  "password": "***",
  "attributes": [
    {
      "id": 101,
      "attribute": "User-Password",
      "value": "***",
      "op": ":="
    }
  ],
  "reply_attributes": []
}
```

### Scenario 2: Create User with IP Assignment

**Request:**
```bash
POST /api/v1/auth
Content-Type: application/json

{
  "username": "bob",
  "password": "bobpass123",
  "attributes": [
    {
      "attribute": "Framed-IP-Address",
      "value": "10.20.30.40"
    }
  ],
  "reply_attributes": [
    {
      "attribute": "Service-Type",
      "value": "Framed-User"
    }
  ]
}
```

**Response (201):**
```json
{
  "username": "bob",
  "password": "***",
  "attributes": [
    {
      "id": 102,
      "attribute": "User-Password",
      "value": "***",
      "op": ":="
    },
    {
      "id": 103,
      "attribute": "Framed-IP-Address",
      "value": "10.20.30.40",
      "op": ":="
    }
  ],
  "reply_attributes": [
    {
      "id": 51,
      "attribute": "Service-Type",
      "value": "Framed-User",
      "op": "+="
    }
  ]
}
```

### Scenario 3: Validation Error

**Request (missing password):**
```bash
POST /api/v1/auth
Content-Type: application/json

{
  "username": "charlie"
}
```

**Response (400):**
```json
{
  "message": "Key: 'CreateAuthRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"
}
```

---

## Best Practices

1. **Always provide both username and password** - Both fields are required
2. **Use HTTPS in production** - Ensure credentials are transmitted securely
3. **Set realistic timeout values** - Use appropriate session and idle timeouts
4. **Validate before submission** - Check username/password requirements client-side
5. **Handle errors gracefully** - Implement proper error handling in clients
6. **Use gRPC for high-performance scenarios** - gRPC is more efficient than REST
7. **Batch operations efficiently** - Consider creating multiple users in appropriate batches
8. **Monitor logs** - Enable logging to track authentication operations
9. **Test with both APIs** - Verify functionality with both REST and gRPC
10. **Document custom attributes** - If using non-standard RADIUS attributes, document them
