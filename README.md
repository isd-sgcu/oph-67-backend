
# OPH-2025 Backend

## Stack
- **Go** (Golang)
- **Fiber** (Go Web Framework)
- **PostgreSQL**

## Getting Started

### Prerequisites
- Go 1.22 or later
- Docker
- Makefile

### Installation
1. **Clone the repository:**
   ```bash
   git clone https://github.com/isd-sgcu/cutu2025-backend.git
   cd cutu2025-backend
   ```

2. **Setup environment:**
   ```bash
   cp .env.example .env
   # Edit .env file with your configurations
   ```

3. **Install dependencies:**
   ```bash
   go mod download
   ```

### Running
**Start Database:**
```bash
docker-compose up -d
```

**Run Server:**
```bash
make server  # Normal mode
```

**Re-create doc**
```bash
swag init -g cmd/main.go  # Normal mode
```

---

## API Documentation

### User Management

#### 1. Register New User *Need Implementation*
**Endpoint:** `POST /api/users/register`  
**Request Format (multipart/form-data):**
```
name: string
phone: string (format: 0812345678)
province: string
school: string
firstInterest: string
secondInterest: string
thirdInterest: string
other...
```

**Success Response (201):**
```json
{
  "accessToken": "string",
  "userId": "string"
}
```

#### 2. Get All Users
**Endpoint:** `GET /api/users`  
**Permissions:** Bearer Token (Staff/Admin)  
**Query Parameters:**
- `name`: Filter by name

**Success Response (200):**
```json
[
  {
    "id": "string",
    "name": "string",
    "phone": "+66812345678",
    "province": "Bangkok",
    "school": "Chulalongkorn University"
  }
]
```

#### 3. Get User by ID
**Endpoint:** `GET /api/users/{id}`  
**Permissions:** Bearer Token

**Success Response (200):**
```json
{
  "id": "string",
  "uid": "string",
  "name": "string",
  "phone": "+66812345678",
  "province": "Bangkok",
  "school": "Chulalongkorn University",
  "firstInterest": "Technology",
  "secondInterest": "Design",
  "thirdInterest": "Business"
}
```

#### 4. Update User
**Endpoint:** `PATCH /api/users/{id}`  
**Permissions:** Bearer Token  
**Request Body:**
```json
{
  "email": "user@example.com",
  "birthDate": "2000-01-01T00:00:00Z",
  "school": "New University"
}
```

**Success Response:** 204 No Content

---

### QR Code Management

#### 5. Scan QR Code
**Endpoint:** `POST /api/users/qr/{id}`  
**Permissions:** Bearer Token (Staff/Admin)

**Success Response (200):**
```json
{
  "id": "string",
  "name": "string",
  "lastEntered": "2024-01-01T12:00:00Z"
}
```

**Error Response (400):**
```json
{
  "error": "User has already entered",
  "message": "2024-01-01 12:00:00 +0000 UTC"
}
```

---

### Admin Operations

#### 6. Add Staff Member
**Endpoint:** `PATCH /api/users/addstaff/{phone}`  
**Permissions:** Bearer Token (Admin)

**Success Response:** 204 No Content

---

## Data Structures

### User Model
```go
type User struct {
    ID              string     `json:"id"`
    UID             string     `json:"uid"`
    Name            string     `json:"name"`
    Email           *string    `json:"email"`
    Phone           string     `json:"phone"`
    BirthDate       *time.Time `json:"birthDate"`
    Role            Role       `json:"role"`
    Province        string     `json:"province"`
    School          string     `json:"school"`
    SelectedSources []string   `json:"selectedSources"`
    OtherSource     *string    `json:"otherSource"`
    FirstInterest   string     `json:"firstInterest"`
    SecondInterest  string     `json:"secondInterest"`
    ThirdInterest   string     `json:"thirdInterest"`
    Objective       string     `json:"objective"`
    RegisteredAt    *time.Time `json:"registerAt"`
    LastEntered     *time.Time `json:"lastEntered"`
}
```

### Enumerations
**User Roles:**
```go
type Role string
const (
    Member Role = "member"
    Staff  Role = "staff"
    Admin  Role = "admin"
)
```

---

## Error Handling

### Common Responses
**400 Bad Request:**
```json
{
  "error": "Invalid phone number format"
}
```

**401 Unauthorized:**
```json
{
  "error": "Missing or invalid token"
}
```

**403 Forbidden:**
```json
{
  "error": "Insufficient permissions"
}
```

**404 Not Found:**
```json
{
  "error": "User not found"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Database connection failed"
}
```
