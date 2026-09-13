# Chirpy API Documentation

Base URL: `http://localhost:8080`

## Authentication

Protected endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

Tokens are issued via the `/api/login` endpoint and validated using HS256 with the JWT secret.

---

## Public Endpoints

### Health Check
**GET** `/api/healthz`

Returns server readiness status.

**Response** `200 OK`
```
OK
```

---

### Create User
**POST** `/api/users`

Registers a new user.

**Request Body**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response** `201 Created`
```json
{
  "id": "uuid",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "email": "user@example.com"
}
```

**Errors**
- `400 Bad Request` - Invalid JSON or missing fields
- `500 Internal Server Error` - Database or hashing failure

---

### Login User
**POST** `/api/login`

Authenticates user and returns JWT token.

**Request Body**
```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "expires_in_seconds": 3600  // optional, max 3600
}
```

**Response** `200 OK`
```json
{
  "id": "uuid",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "email": "user@example.com",
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Errors**
- `400 Bad Request` - Invalid JSON
- `401 Unauthorized` - Invalid credentials
- `500 Internal Server Error` - Token generation failure

---

## Protected Endpoints (Require Authentication)

### Create Chirp
**POST** `/api/chirps`

Creates a new chirp (max 140 characters, profanity filtered).

**Headers**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body**
```json
{
  "body": "Hello world!"
}
```

**Response** `201 Created`
```json
{
  "id": "uuid",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "body": "Hello world!",
  "user_id": "uuid"
}
```

**Errors**
- `400 Bad Request` - Chirp too long (>140 chars) or empty
- `401 Unauthorized` - Missing or invalid token
- `500 Internal Server Error` - Database failure

---

### Get All Chirps
**GET** `/api/chirps`

Returns all chirps ordered by creation date.

**Response** `200 OK`
```json
[
  {
    "id": "uuid",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "body": "Hello world!",
    "user_id": "uuid"
  }
]
```

**Errors**
- `500 Internal Server Error` - Database failure

---

### Get Chirp by ID
**GET** `/api/chirps/{chirpID}`

Returns a single chirp by UUID.

**Path Parameters**
- `chirpID` (UUID) - Chirp identifier

**Response** `200 OK`
```json
{
  "id": "uuid",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "body": "Hello world!",
  "user_id": "uuid"
}
```

**Errors**
- `400 Bad Request` - Invalid UUID format
- `404 Not Found` - Chirp not found
- `500 Internal Server Error` - Database failure

---

## Admin Endpoints

### Get Metrics
**GET** `/admin/metrics`

Returns server hit count for the file server.

**Response** `200 OK`
```html
<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited 42 times!</p></body></html>
```

---

### Reset Database
**POST** `/admin/reset`

Resets the database (development only).

**Response** `200 OK`
```
OK
```

---

## Error Response Format

All error responses follow this format:
```json
{
  "error": "Error message description"
}
```

## Data Models

### Chirp
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| created_at | ISO8601 datetime | Creation timestamp |
| updated_at | ISO8601 datetime | Last update timestamp |
| body | string | Chirp content (max 140 chars) |
| user_id | UUID | Author's user ID |

### User
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| created_at | ISO8601 datetime | Creation timestamp |
| updated_at | ISO8601 datetime | Last update timestamp |
| email | string | User's email |

### UserLoginResponse
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| created_at | ISO8601 datetime | Creation timestamp |
| updated_at | ISO8601 datetime | Last update timestamp |
| email | string | User's email |
| token | string | JWT access token |

---

## Rate Limits & Constraints

- Chirp body: Maximum 140 characters
- JWT token expiry: Maximum 3600 seconds (1 hour)
- Password hashing: Argon2id (default params)
- JWT signing: HS256