# DevBoard 1.0

A project management and productivity API built with **Go**, **Gin**, **GORM**, and **SQLite**.

DevBoard allows authenticated users to create and manage projects, organize tasks, track time spent on their work, and view productivity reports.

## 🚀 Features

* User registration and authentication
* JWT-based authentication
* Protected API routes
* Project management
* Task management
* Task status and priority management
* Time tracking
* Productivity reports
* User-specific data authorization
* SQLite database
* GORM database management
* RESTful API architecture

## 🛠️ Tech Stack

| Technology | Purpose                      |
| ---------- | ---------------------------- |
| Go         | Backend programming language |
| Gin        | HTTP web framework           |
| GORM       | ORM and database management  |
| SQLite     | Database                     |
| JWT        | Authentication               |
| bcrypt     | Password hashing             |

## 📁 Project Structure

```text
devboard-1.0/
│
├── cmd/
│   └── main.go
│
├── internal/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── routes/
│   └── services/
│
├── go.mod
├── go.sum
└── README.md
```

## 🔐 Authentication

DevBoard uses JWT authentication to protect private API endpoints.

After logging in, the client receives a JWT token.

Protected requests use:

```http
Authorization: Bearer <your-token>
```

Users can only access resources belonging to their account.

## 📌 API Endpoints

### Authentication

```text
POST /api/register
POST /api/login
GET  /api/profile
```

### Projects

```text
POST   /api/projects
GET    /api/projects
GET    /api/projects/:id
PUT    /api/projects/:id
DELETE /api/projects/:id
```

### Tasks

```text
POST   /api/projects/:project_id/tasks
GET    /api/projects/:project_id/tasks
GET    /api/tasks/:id
PUT    /api/tasks/:id
DELETE /api/tasks/:id
```

### Time Tracking

```text
POST /api/time/start
POST /api/time/stop/:id
GET  /api/time
GET  /api/projects/:project_id/time
GET  /api/tasks/:task_id/time
```

### Reports

```text
GET /api/reports/projects/:project_id/time
GET /api/reports/tasks/productivity
GET /api/reports/activity
GET /api/reports/dashboard
```

## ⚙️ Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/Demiladeolorunsola/devboard-1.0.git
```

### 2. Enter the project

```bash
cd devboard-1.0
```

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Configure environment variables

Create a `.env` file in the project root.

Example:

```env
PORT=8080
JWT_SECRET=your-secret-key
DATABASE_URL=devboard.db
```

> Do not commit your `.env` file or expose your JWT secret publicly.

### 5. Run the server

Depending on the project's current entry point:

```bash
go run ./cmd
```

The API should then be available at:

```text
http://localhost:8080
```

## 🗄️ Database

DevBoard currently uses **SQLite** with **GORM**.

The database stores:

* Users
* Projects
* Tasks
* Time entries

Database migrations are handled by the application during startup.

## 🔒 Security

The project includes several security measures:

* JWT authentication
* Password hashing
* Protected API routes
* User ownership checks
* Validation of project/task status values
* Validation of task priorities
* Protection against accessing another user's resources
* Prevention of multiple active timers for the same user

## 📊 Example Workflow

A typical DevBoard workflow looks like:

```text
Register
   ↓
Login
   ↓
Receive JWT
   ↓
Create Project
   ↓
Create Tasks
   ↓
Start Timer
   ↓
Complete Tasks
   ↓
Stop Timer
   ↓
View Productivity Reports
```

## 🎯 Project Goals

DevBoard was created to demonstrate practical backend development with Go, including:

* REST API development
* Authentication and authorization
* Database design
* Service-layer architecture
* HTTP request handling
* CRUD operations
* Data validation
* Time tracking
* Reporting
* Backend security

## 🔮 Future Improvements

Planned improvements include:

* React frontend
* User profile management
* Team collaboration
* Project members
* Notifications
* Advanced analytics
* API documentation with Swagger/OpenAPI
* Automated unit and integration tests
* Docker support
* Production deployment

## 👨🏾‍💻 Author

**Oluwatidemilade Olorunsola**

Software Engineer

GitHub: `Demiladeolorunsola`

## 📄 License

This project is currently intended as a portfolio and learning project.
