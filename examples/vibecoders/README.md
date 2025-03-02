# VibeCoder

A developer-focused social platform built with Go Echo framework, SQLite, Tailwind CSS, and esbuild.

## Features

- User authentication (registration and login)
- User profiles with bio, LinkedIn URL, GitHub URL, and profile photo
- Dark-themed, responsive UI built with Tailwind CSS
- REST API endpoints

## Technology Stack

- **Backend**: Go with Echo framework
- **Database**: SQLite
- **Frontend**: JavaScript with esbuild bundler
- **Styling**: Tailwind CSS

## Getting Started

### Prerequisites

- Go 1.18 or higher
- Node.js 16 or higher
- npm or yarn

### Installation

1. Clone this repository
```bash
git clone <repository-url>
cd vibecoders
```

2. Install backend dependencies
```bash
go mod tidy
```

3. Install frontend dependencies
```bash
npm install
```

4. Build the frontend assets
```bash
npm run build
npm run build:css
```

### Running the Application

1. Start the server
```bash
go run cmd/server/main.go
```

2. For development with hot-reloading:
```bash
npm run dev
```

3. Visit `http://localhost:3000` in your browser

## API Endpoints

- `POST /api/register` - Register a new user
- `POST /api/login` - Log in a user
- `GET /api/users/:id` - Get user profile
- `PUT /api/users/:id` - Update user profile

## Project Structure

```
vibecoders/
├── api/
│   └── handlers/       # HTTP handlers
├── cmd/
│   └── server/         # Application entry point
├── models/             # Database models
├── static/             # Static assets
│   ├── src/            # Frontend source files
│   └── dist/           # Bundled frontend files
└── templates/          # HTML templates
```

## License

MIT