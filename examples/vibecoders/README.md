# VibeCoders

A platform to find software engineers that are very good at vibecoding.

> "Vibe coding is the art of leveraging AI tools to their fullest potential in software development, creating a seamless fusion between human creativity and machine intelligence. It represents a paradigm shift where developers orchestrate AI systems rather than writing every line manually."

## Technology Stack

### Frontend
- React
- TailwindCSS
- Vite (for building and development)

### Backend
- Golang
- Echo framework
- SQLite

## Features

- Homepage showcasing top "vibecoders"
- User authentication (register, login, logout)
- User profiles with customizable information
- Responsive design for mobile and desktop

## Development

### Prerequisites

- Go 1.19+
- Node.js 16+
- npm

### Setup

1. Clone the repository
2. Install Go dependencies:
   ```
   go mod tidy
   ```
3. Install Node.js dependencies:
   ```
   npm install
   ```

### Running the application

1. Start the frontend development server:
   ```
   npm run dev
   ```

2. Start the backend server:
   ```
   go run main.go
   ```

3. Access the application at `http://localhost:3000`

## API Endpoints

- `POST /api/login` - Login user
- `DELETE /api/logout` - Logout user
- `POST /api/register` - Register new user
- `PATCH /api/user` - Update user profile
- `GET /api/homepage-users` - Get users for homepage
- `GET /api/user` - Get current user information

## Database

SQLite database with the following tables:

### users
- id (primary key)
- username (unique)
- bio
- linked_in_url
- github_url
- photo_url
- password

### sessions
- id (primary key)
- user_id (foreign key)
- token (random uuid)

### Database Migrations

This project uses Flyway for database migrations:

1. Install Flyway CLI: https://flywaydb.org/documentation/usage/commandline/

2. Initialize the database (creates database file and runs migrations):
   ```
   npm run db:init
   ```

3. Run migrations only:
   ```
   npm run db:migrate
   ```

4. Check migration status:
   ```
   npm run db:info
   ```

5. Clean database (delete all tables):
   ```
   npm run db:clean
   ```

Migration files are located in `db/migration` and follow Flyway naming conventions:
- `V1__Create_initial_schema.sql` - Creates the database schema
- `V2__Insert_seed_data.sql` - Inserts initial seed data