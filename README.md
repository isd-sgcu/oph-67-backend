# OPH-2025 Backend API Documentation

## User Management

### 1. Register New Staff Member
**Endpoint:** `POST /api/staff/register`  
**Request Format (multipart/form-data):**
- `id`: string (required)
- `name`: string (required)
- `phone`: string (required, format: 0812345678)
- `nickname`: string (required)
- `studentId`: string (required)
- `email`: string (required)
- `faculty`: string (optional)
- `year`: int (optional)
- `isCentralStaff`: boolean (optional)

**Success Response (201):**
```json
{
  "accessToken": "string",
  "userId": "string"
}
```

**Error Responses:**
- `400 Bad Request`: Missing required fields or invalid phone format
- `500 Internal Server Error`: Failed to create user

---

### 2. Register New Student
**Endpoint:** `POST /api/student/register`  
**Request Format (multipart/form-data):**
- `id`: string (required)
- `name`: string (required)
- `phone`: string (required, format: 0812345678)
- `email`: string (required)
- `status`: string (optional)
- `otherStatus`: string (optional)
- `birthDate`: string (optional, format 2004-05-02)
- `province`: string (optional)
- `school`: string (optional)
- `selectedSources`: string (comma-separated list, optional)
- `otherSource`: string (optional)
- `firstInterest`: string (optional)
- `secondInterest`: string (optional)
- `thirdInterest`: string (optional)
- `objective`: string (optional)

example
```
name:John Doe
phone:0949823192
email:john.doe@example.com
status:Study
province:Bangkok
school:CU
selectedSources:[Facebook,Website]	
firstInterest:Business
secondInterest:Technology
thirdInterest:Marketing
objective:Learn for skill
id:11345677
```

**Success Response (201):**
```json
{
  "accessToken": "string",
  "userId": "string"
}
```

**Error Responses:**  
Same as Staff Registration

---

### 3. Get All Users
**Endpoint:** `GET /api/users`  
**Permissions:** Bearer Token (Staff/Admin)  
**Query Parameters:**
- `name`: Filter by name (optional)
- `role`: Filter by role (`member`/`staff`/`admin`/`student`)

**Success Response (200):**
```json
[
  {
    "id": "user1",
    "name": "John Doe",
    "phone": "+66812345678",
    "role": "staff",
    "email": "john@example.com",
    "faculty": "Engineering"
  }
]
```

---

### 4. Get User by ID
**Endpoint:** `GET /api/users/{id}`  
**Permissions:** Bearer Token

**Success Response (200):**
```json
{
  "id": "user1",
  "name": "John Doe",
  "phone": "+66812345678",
  "role": "student",
  "school": "Chulalongkorn University",
  "firstInterest": "Technology"
}
```

---

### 5. Update User
**Endpoint:** `PATCH /api/users/{id}`  
**Permissions:** Bearer Token  
**Request Body (JSON):**
```json
{
  "email": "new@example.com",
  "school": "New University"
}
```

**Success Response:** `204 No Content`

---

### 6. Scan QR Code
**Endpoint:** `POST /api/users/qr/{studentId}`  
**Permissions:** Bearer Token (Staff/Admin)

**Success Response (200):**
```json
{
  "id": "user1",
  "name": "John Doe",
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

### 7. Add Staff Member
**Endpoint:** `PATCH /api/admin/addstaff/{phone}`  
**Permissions:** Bearer Token (Admin)

**Success Response:** `204 No Content`

---

### 8. Change Role
**Endpoint:** `PATCH /api/admin/role/{userId}`  
**Permissions:** Bearer Token (Admin)

**Success Response:** `204 No Content`
```json
{
    "role": "admin"
}
```

### 9. Delete
**Endpoint:** `DELETE /api/admin/delete/{userId}`  
**Permissions:** Bearer Token (Admin)

**Success Response:** `204 No Content`

---

Here is the updated **Student Evaluation API Documentation** reflecting your latest route and handler implementation:

---

# Student Evaluation API Documentation

**Base URL:** `/api/student-evaluation`

---

## Authentication

- **JWT Required** for all routes
- **Staff/Admin Role Required** for listing all evaluations

---

## Endpoints

### 1. Create a Student Evaluation  
**POST** `/api/student-evaluation/`  
**Authorization:** Bearer Token (JWT)  
**Role:** Any authenticated user

#### Request Body
```json
{
  "newSources": ["Instagram", "Friend"],
  "overallActivity": 4,
  "interestActivity": 5,
  "receivedFacultyInfoClearly": 4,
  "wouldRecommendCUOpenHouseNextTime": 5,
  "favoriteBooth": "Engineering",
  "activityDiversity": 4,
  "perceivedCrowdDensity": 3,
  "hasFullBoothAccess": 1,
  "facilityConvenienceRating": 4,
  "campusNavigationRating": 4,
  "hesitationLevelAfterDisaster": 2,
  "lineOASignupRating": 5,
  "designBeautyRating": 4,
  "websiteImprovementSuggestions": "Improve mobile responsiveness"
}
```

#### Responses
- `201 Created` – Evaluation successfully created.
- `401 Unauthorized` – Missing or invalid JWT.
- `409 Conflict` – Evaluation already exists for this user.
- `400 Bad Request` – Invalid input.
- `500 Internal Server Error`

---

### 2. Get Student Evaluation by ID  
**GET** `/api/student-evaluation/:id`  
**Authorization:** Bearer Token  
**Role:** Any authenticated user

#### Path Parameters
- `id` – Student ID

#### Responses
- `200 OK` – Returns the student evaluation.
- `400 Bad Request` – Missing student ID.
- `404 Not Found` – Evaluation not found.

---

### 3. Update Student Evaluation  
**PATCH** `/api/student-evaluation/:id`  
**Authorization:** Bearer Token  
**Role:** Any authenticated user

#### Path Parameters
- `id` – Student ID

#### Request Body
Same structure as Create endpoint.


#### Responses
- `200 OK` – Evaluation updated.
- `400 Bad Request` – Missing or invalid input.
- `500 Internal Server Error`

---

### 4. Delete Student Evaluation  
**DELETE** `/api/student-evaluation/:id`  
**Authorization:** Bearer Token  
**Role:** Any authenticated user

#### Path Parameters
- `id` – Student ID

#### Responses
- `204 No Content` – Evaluation deleted.
- `400 Bad Request` – Missing student ID.
- `500 Internal Server Error`

---

### 5. Get All Student Evaluations  
**GET** `/api/student-evaluation/`  
**Authorization:** Bearer Token  
**Role:** Staff or Admin

#### Responses
- `200 OK` – Returns a list of all student evaluations.
- `403 Forbidden` – Insufficient role.
- `500 Internal Server Error`

---

## Data Structures

### User Model
| Field             | Type            | Description                     |
|-------------------|-----------------|---------------------------------|
| `id`              | string          | Unique user identifier          |
| `role`            | Role            | `staff`/`admin`/`student` |
| `selectedSources` | array[string]   | Sources user heard about event  |
| `faculty`         | string          | Staff member's faculty          |
| `isCentralStaff`  | boolean         | Central committee status        |

### Full Role Enumeration
```go
enum Role {
  staff
  admin
  student
}
```

---

## Error Examples
**403 Forbidden (Insufficient Permissions):**
```json
{
  "error": "Only admins can modify user roles"
}
```

**404 Not Found (User Not Found):**
```json
{
  "error": "User not found with ID: user123"
}
```