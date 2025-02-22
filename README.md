Here’s the updated `README.md` file with the new fields (`Age`, `ChronicDisease`, and `DrugAllergy`) added to the `User` struct documentation. I've also ensured that the API endpoint descriptions and examples reflect these changes.

---

# OPH-2025 Backend #Not Update Yet

## Stack
- **Go** (Golang)
- **Fiber** (Go Web Framework)
- **PostgreSQL**

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go 1.22 or later
- Docker
- Makefile
- **Air** (optional for auto-reload)

### Installing

1. **Clone the repository:**

   ```bash
   git clone https://github.com/isd-sgcu/cutu2025-backend.git
   cd cutu2025-backend
   ```

2. **Copy the environment configuration file:**

   ```bash
   cp .env.example .env
   ```

   Fill in the values in the `.env` file for your local environment.

3. **Download dependencies:**

   ```bash
   go mod download
   ```

### Running the Project

#### Database

To start the local database for development, run:

```bash
docker-compose up -d
```

This will launch the PostgreSQL database in a Docker container.

#### Server

Option 1: **Standard Mode**

To run the server normally:

```bash
make server
```

Option 2: **Development Mode (with auto-reload)**

To run the server in development mode with live auto-reload, use **Air** (a Go live reloading tool):

```bash
air -c .air.toml
```

This option will automatically reload the server when you change any Go files.

---

## API Endpoints

### 1. **Get All Users**
**Endpoint:** `/api/users`  
**Method:** `GET`  
**Permission:** BearerAuth (Staff, Admin)

Retrieve a list of all users.

**Response:**
- `200 OK`: Returns a list of users.
- `500 Internal Server Error`: Failed to fetch users.

---

### 2. **Update Account Info**
**Endpoint:** `/api/users`  
**Method:** `PATCH`  
**Permission:** BearerAuth

Update that staff member's personal information.

**Parameters:**
- `user` (body) - User data (JSON).

**Response:**
- `204 No Content`: User successfully updated.
- `400 Bad Request`: Invalid input.
- `401 Unauthorized`: Unauthorized.
- `403 Forbidden`: Forbidden.
- `404 Not Found`: User not found.
- `500 Internal Server Error`: Failed to update user.

---

### 3. **Add Staff by Phone**
**Endpoint:** `/api/users/addstaff/{phone}`  
**Method:** `PATCH`  
**Permission:** BearerAuth (Admin)

Add a staff member by their phone number.

**Parameters:**
- `phone` (path) - The phone number of the user.

**Response:**
- `204 No Content`: Staff added successfully.
- `400 Bad Request`: User is already a staff.
- `500 Internal Server Error`: Failed to add staff.

---

### 4. **Get User by ID**
**Endpoint:** `/api/users/{id}`  
**Method:** `GET`  
**Permission:** BearerAuth

Retrieve a user by its ID.

**Parameters:**
- `id` (path) - The ID of the user.

**Response:**
- `200 OK`: Returns the user details.
- `404 Not Found`: User not found.
- `500 Internal Server Error`: Failed to fetch user.

---

### 5. **Update User by ID**
**Endpoint:** `/api/users/{id}`  
**Method:** `PATCH`  
**Permission:** BearerAuth

Update a user by its ID.

**Parameters:**
- `id` (path) - The ID of the user.
- `user` (body) - User data (JSON).

**Response:**
- `204 No Content`: User successfully updated.
- `400 Bad Request`: Invalid input.
- `401 Unauthorized`: Unauthorized.
- `403 Forbidden`: Forbidden.
- `404 Not Found`: User not found.
- `500 Internal Server Error`: Failed to update user.

---

### 6. **Delete User by ID**
**Endpoint:** `/api/users/{id}`  
**Method:** `DELETE`  
**Permission:** BearerAuth (Admin)

Delete a user by its ID.

**Parameters:**
- `id` (path) - The ID of the user.

**Response:**
- `204 No Content`: User successfully deleted.
- `401 Unauthorized`: Unauthorized.
- `403 Forbidden`: Forbidden.
- `404 Not Found`: User not found.
- `500 Internal Server Error`: Failed to delete user.

---

### 7. **Get QR Code URL for User**
**Endpoint:** `/api/users/qr/{id}`  
**Method:** `GET`  
**Permission:** BearerAuth

Retrieve a QR code URL for a user.

**Parameters:**
- `id` (path) - The ID of the user.

**Response:**
- `200 OK`: Returns the QR code URL.
- `404 Not Found`: User not found.
- `500 Internal Server Error`: Failed to fetch user.

---

### 8. **Scan QR Code**
**Endpoint:** `/api/users/qr/{id}`  
**Method:** `POST`  
**Permission:** BearerAuth (Staff, Admin)

Scan a QR code and perform associated actions.

**Parameters:**
- `id` (path) - The ID of the user.

**Response:**
- `200 OK`: User scanned successfully with User data including last
- `400 Bad Request`: User has already entered with last enter time.
```
{
    "error": "User has already entered",
    "message": "2025-01-26 18:39:15.10983 +0700 +07"
}
```
- `500 Internal Server Error`: Failed to fetch user.

---

### 9. **Register a New User**
**Endpoint:** `/api/users/register`  
**Method:** `POST`  
**Permission:** No

Register a new user in the system.

**Parameters (form data):**
- `id` (string) - User ID
- `name` (string) - User Name
- `email` (string) - User Email
- `phone` (string) - User Phone
- `university` (string) - User University
- `sizeJersey` (string) - Jersey Size
- `foodLimitation` (string) - Food Limitation
- `invitationCode` (string) - Invitation Code
- `status` (string) - User Status (`chula_student`, `alumni`, `general_public`, `general_student`)
- `image` (file) - User Image
- `age` (string) - User Age
- `chronicDisease` (string) - Chronic Disease
- `drugAllergy` (string) - Drug Allergy
- `graduatedYear` (string) - Graduated Year
- `faculty` (string) - Faculty
- `education` (string) - User Education (`studying`, `graduated`)
- 'isAcrophobia' (bool) - Is User acrophobia (`true`, `false`)

**Response:**
- `201 Created`: User successfully created.
- `400 Bad Request`: Invalid input.
- `401 Unauthorized`: Unauthorized.
- `500 Internal Server Error`: Failed to create user.

---

### 10. **Update User Role by ID**
**Endpoint:** `/api/users/role/{id}`  
**Method:** `PATCH`  
**Permission:** BearerAuth (Admin)

Update a user role by its ID.

**Parameters:**
- `id` (path) - The ID of the user.
- `role` (body) - User role (string).

**Response:**
- `204 No Content`: User role updated successfully.
- `400 Bad Request`: Invalid input.
- `401 Unauthorized`: Unauthorized.
- `403 Forbidden`: Forbidden.
- `404 Not Found`: User not found.
- `500 Internal Server Error`: Failed to update user role.

---

### 11. **SignIn**
**Endpoint:** `/api/users/signin`  
**Method:** `POST`  
**Permission:** No

Authenticate a user and return an access token.

**Parameters:**
- `id` (body) - User ID.

**Response:**
- `200 OK`: Returns an access token.
- `400 Bad Request`: Invalid input.
- `401 Unauthorized`: Unauthorized.
- `500 Internal Server Error`: Failed to sign in.

---

## Error Responses

### Error Response Format
```json
{
  "error": "Error message here"
}
```

### Common Error Codes
- `400 Bad Request`: Invalid input.
- `401 Unauthorized`: Unauthorized access.
- `403 Forbidden`: Forbidden action.
- `404 Not Found`: Resource not found.
- `500 Internal Server Error`: An error occurred on the server.

---

## Definitions

### **Education Enum**
- `studying`: The user is currently studying.
- `graduated`: The user has graduated.

### **Role Enum**
- `member`: A member user.
- `staff`: A staff user.
- `admin`: An admin user.

### **Status Enum**
- `chula_student`: The user is a Chula student.
- `alumni`: The user is an alumni.
- `general_public`: The user is from the general public.
- `general_student`: The user is a general student.

### **TokenResponse**
- `accessToken`: The access token for authentication.
- `userId`: The user ID associated with the token.

### **User**
A user object containing:
- `id`: The user liff ID.
- `uid`: The user UID
- `name`: The user name.
- `email`: The user email.
- `phone`: The user phone number.
- `status`: The user's status.
- `role`: The user's role.
- `education`: The user's education status.
- `imageUrl`: The user's profile image URL.
- `faculty`: The user's faculty.
- `foodLimitation`: The user's food limitations.
- `graduatedYear`: The year the user graduated.
- `invitationCode`: The user's invitation code.
- `lastEntered`: Timestamp for the last QR scan.
- `sizeJersey`: The user's jersey size.
- `university`: The user's university.
- `age`: The user's age.
- `chronicDisease`: The user's chronic disease information.
- `drugAllergy`: The user's drug allergy information.
- `isAcrophobia`: Check if user is acrophobia (bool).
