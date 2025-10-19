# API Documentation

Complete API reference for the LinkedIn Connector backend.

## Base URL

```
http://localhost:8080/api
```

## Authentication

Most endpoints require authentication using JWT tokens. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

---

## Endpoints

### Health Check

#### GET /health

Check if the API is running.

**Request:**
```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok",
  "message": "LinkedIn Connector API is running"
}
```

---

## Authentication Endpoints

### Register User

#### POST /api/auth/register

Create a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com"
  }
}
```

**Error Response (409 Conflict):**
```json
{
  "error": "User already exists"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

---

### Login User

#### POST /api/auth/login

Authenticate a user and receive a JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com"
  }
}
```

**Error Response (401 Unauthorized):**
```json
{
  "error": "Invalid email or password"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

---

## LinkedIn Connection Endpoints

### Connect LinkedIn with Cookie

#### POST /api/linkedin/connect/cookie

Connect a LinkedIn account using cookie authentication.

**Headers:**
```
Authorization: Bearer <your-jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "cookie": "AQEDATEa...your-li_at-cookie-value..."
}
```

**Response (200 OK):**
```json
{
  "message": "LinkedIn account connected successfully",
  "account_id": "linkedin:12345678",
  "account": {
    "id": 1,
    "user_id": 1,
    "provider": "linkedin",
    "account_id": "linkedin:12345678",
    "account_name": "John Doe",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Unipile API error: Invalid cookie"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/linkedin/connect/cookie \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "cookie": "AQEDATEa...your-cookie..."
  }'
```

**How to get LinkedIn cookie:**
1. Log in to LinkedIn in your browser
2. Open Developer Tools (F12)
3. Go to Application/Storage tab
4. Find Cookies → linkedin.com
5. Copy the value of `li_at` cookie

---

### Connect LinkedIn with Credentials

#### POST /api/linkedin/connect/credentials

Connect a LinkedIn account using username and password.

**Headers:**
```
Authorization: Bearer <your-jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "username": "your-email@example.com",
  "password": "your-linkedin-password"
}
```

**Response (200 OK):**
```json
{
  "message": "LinkedIn account connected successfully",
  "account_id": "linkedin:87654321",
  "account": {
    "id": 2,
    "user_id": 1,
    "provider": "linkedin",
    "account_id": "linkedin:87654321",
    "account_name": "Jane Smith",
    "created_at": "2024-01-15T11:00:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Unipile API error: Invalid credentials"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/linkedin/connect/credentials \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your-email@example.com",
    "password": "your-password"
  }'
```

---

## Account Management Endpoints

### Get All Linked Accounts

#### GET /api/accounts

Retrieve all linked accounts for the authenticated user.

**Headers:**
```
Authorization: Bearer <your-jwt-token>
```

**Response (200 OK):**
```json
{
  "accounts": [
    {
      "id": 1,
      "user_id": 1,
      "provider": "linkedin",
      "account_id": "linkedin:12345678",
      "account_name": "John Doe",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": 2,
      "user_id": 1,
      "provider": "linkedin",
      "account_id": "linkedin:87654321",
      "account_name": "Jane Smith",
      "created_at": "2024-01-15T11:00:00Z",
      "updated_at": "2024-01-15T11:00:00Z"
    }
  ],
  "count": 2
}
```

**Example:**
```bash
curl -X GET http://localhost:8080/api/accounts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### Delete Linked Account

#### DELETE /api/accounts/:id

Remove a linked account.

**Headers:**
```
Authorization: Bearer <your-jwt-token>
```

**URL Parameters:**
- `id` (required): The ID of the account to delete

**Response (200 OK):**
```json
{
  "message": "Account deleted successfully"
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "Account not found"
}
```

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/accounts/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
Invalid request body or parameters.
```json
{
  "error": "Invalid request format"
}
```

### 401 Unauthorized
Missing or invalid authentication token.
```json
{
  "error": "Authorization header required"
}
```

### 404 Not Found
Resource not found.
```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error
Server-side error.
```json
{
  "error": "Internal server error"
}
```

---

## Unipile Integration

This API integrates with Unipile's native authentication API. Here's how it works:

### Cookie Authentication Flow
1. Client sends LinkedIn cookie to `/api/linkedin/connect/cookie`
2. Backend forwards cookie to Unipile API: `POST https://api.unipile.com/v1/accounts`
3. Unipile validates cookie and returns `account_id`
4. Backend stores `account_id` in database
5. Backend returns success response with account details

### Credentials Authentication Flow
1. Client sends LinkedIn username/password to `/api/linkedin/connect/credentials`
2. Backend forwards credentials to Unipile API: `POST https://api.unipile.com/v1/accounts`
3. Unipile validates credentials and returns `account_id`
4. Backend stores `account_id` in database
5. Backend returns success response with account details

### Unipile API Request Format

The backend sends requests to Unipile in this format:

```json
{
  "provider": "linkedin",
  "cookie": "cookie-value",  // For cookie auth
  "username": "email",       // For credentials auth
  "password": "password",    // For credentials auth
  "type": "cookie"           // or "credentials"
}
```

### Unipile API Response Format

Unipile returns responses in this format:

```json
{
  "account_id": "linkedin:12345678",
  "provider": "linkedin",
  "name": "John Doe",
  "username": "john.doe",
  "status": "active"
}
```

---

## Rate Limiting

Currently, there is no rate limiting implemented. For production use, consider implementing rate limiting to prevent abuse.

Example using middleware:
```go
// Limit to 100 requests per minute per IP
rateLimiter := limiter.Rate{
    Period: 1 * time.Minute,
    Limit:  100,
}
```

---

## CORS Configuration

The API is configured to accept requests from:
- `http://localhost:5173` (development)
- The URL specified in `FRONTEND_URL` environment variable

Allowed methods: GET, POST, PUT, DELETE, OPTIONS

---

## Testing the API

### Using cURL

```bash
# Register
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | jq -r '.token')

# Connect LinkedIn with cookie
curl -X POST http://localhost:8080/api/linkedin/connect/cookie \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"cookie":"your-cookie-here"}'

# Get accounts
curl http://localhost:8080/api/accounts \
  -H "Authorization: Bearer $TOKEN"
```

### Using Postman

1. Import the API collection
2. Set environment variable `BASE_URL` to `http://localhost:8080`
3. Register/login to get a token
4. Add token to Authorization header for protected endpoints

---

## Webhooks (Future Enhancement)

Consider implementing webhooks to notify your application when:
- Account connection status changes
- LinkedIn account is disconnected
- API rate limits are reached

---

## Best Practices

1. **Always use HTTPS in production**
2. **Store JWT tokens securely** (httpOnly cookies or secure storage)
3. **Implement token refresh mechanism** for better UX
4. **Validate all inputs** on both client and server
5. **Handle errors gracefully** and log them properly
6. **Set up monitoring** for API health and performance
7. **Implement rate limiting** to prevent abuse
8. **Use environment variables** for sensitive configuration
9. **Keep dependencies updated** for security patches
10. **Document API changes** when making updates

---

## Support

For questions about:
- **This API**: Check the main README.md or open an issue
- **Unipile API**: Visit https://docs.unipile.com
- **Authentication issues**: Review the security section in README.md

