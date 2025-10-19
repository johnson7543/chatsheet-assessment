# LinkedIn Account Connector via Unipile

A full-stack application that allows users to connect their LinkedIn accounts using Unipile's native authentication API. The app supports both cookie-based authentication and username/password login methods.

## Features

- 🔐 User authentication (register/login)
- 🔗 LinkedIn account connection via Unipile API
- 📊 Support for both cookie auth and username/password methods
- 💾 Persistent storage of connected accounts
- 🎨 Modern, responsive UI

## Tech Stack

### Backend
- **Go 1.22** - Backend API server
- **Gin** - Web framework
- **GORM** - ORM for database operations
- **SQLite** - Database (can be swapped with PostgreSQL/MySQL)
- **JWT** - Authentication tokens

### Frontend
- **React 18** - UI framework
- **Vite** - Build tool
- **Axios** - HTTP client
- **TailwindCSS** - Styling

## Project Structure

```
.
├── backend/                    # Go 1.22 backend (Standard Go Layout)
│   ├── cmd/
│   │   └── api/
│   │       └── main.go        # Application entry point
│   ├── internal/              # Private application code
│   │   ├── config/           # Configuration management
│   │   ├── database/         # Database setup and migrations
│   │   ├── models/           # Data models
│   │   ├── handlers/         # HTTP request handlers
│   │   ├── middleware/       # Authentication middleware
│   │   ├── repository/       # Data access layer (NEW)
│   │   └── service/          # Business logic layer (NEW)
│   ├── configs/              # Configuration files
│   ├── scripts/              # Build and utility scripts
│   └── go.mod                # Go dependencies
├── frontend/                   # React frontend
│   ├── src/
│   │   ├── components/       # React components
│   │   ├── services/         # API service layer
│   │   └── App.jsx          # Main app component
│   ├── package.json          # npm dependencies
│   └── vite.config.js        # Vite configuration
├── README.md                  # This file
└── REFACTORING_SUMMARY.md     # Backend refactoring details
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME,
    updated_at DATETIME
);
```

### Linked Accounts Table
```sql
CREATE TABLE linked_accounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    provider TEXT NOT NULL DEFAULT 'linkedin',
    account_id TEXT NOT NULL,
    account_name TEXT,
    created_at DATETIME,
    updated_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

## Setup Instructions

### Prerequisites
- Go 1.22 or higher
- Node.js 18+ and npm
- Unipile API credentials

### 📝 Understanding .env Files

This project has **3 different `.env.example` files** for different scenarios:

| File Location | Purpose | When to Use |
|--------------|---------|-------------|
| `.env.example` (root) | Docker Compose | Running with `docker-compose up` |
| `frontend/.env.example` | Frontend local dev | Running `npm run dev` |
| `backend/configs/.env.example` | Backend local dev | Running `go run cmd/api/main.go` |

**See [ENV_SETUP_GUIDE.md](ENV_SETUP_GUIDE.md) for detailed explanation.**

### Backend Setup

1. Navigate to the backend directory:
```bash
cd backend
```

2. Install Go dependencies:
```bash
go mod download
```

3. Create a `.env` file with your configuration:
```bash
cp configs/.env.example configs/.env
```

4. Edit `configs/.env` with your settings:
```
PORT=8080
JWT_SECRET=your-secret-key-here
UNIPILE_API_KEY=your-unipile-api-key
UNIPILE_API_URL=https://api.unipile.com/v1
DATABASE_PATH=./linkedin_connector.db
```

5. Run the backend:
```bash
# Quick start
go run cmd/api/main.go

# Or use the run script
./scripts/run.sh

# Or build and run
./scripts/build.sh
./bin/api
```

The backend will start on `http://localhost:8080`

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Create a `.env` file:
```bash
cp .env.example .env
```

4. Edit `.env`:
```
VITE_API_URL=http://localhost:8080
```

5. Run the development server:
```bash
npm run dev
```

The frontend will start on `http://localhost:5173`

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user

### LinkedIn Connection
- `POST /api/linkedin/connect/cookie` - Connect LinkedIn via cookie
- `POST /api/linkedin/connect/credentials` - Connect LinkedIn via username/password
- `GET /api/accounts` - Get all linked accounts for current user

## Usage Flow

1. **Register/Login**: Create an account or login with existing credentials
2. **Connect LinkedIn**: Choose between cookie-based or credentials-based authentication
3. **View Accounts**: See all your connected LinkedIn accounts with their account IDs

### Connecting via Cookie
```json
POST /api/linkedin/connect/cookie
{
  "cookie": "li_at=your-linkedin-cookie-value"
}
```

### Connecting via Credentials
```json
POST /api/linkedin/connect/credentials
{
  "username": "your-email@example.com",
  "password": "your-password"
}
```

## Unipile Integration

This app uses Unipile's native authentication API (not the hosted wizard). The integration handles:

1. Sending authentication credentials to Unipile
2. Receiving and storing the `account_id` returned by Unipile
3. Associating the `account_id` with the logged-in user

## Security Considerations

- Passwords are hashed using bcrypt
- JWT tokens for session management
- Environment variables for sensitive configuration
- CORS enabled for frontend-backend communication
- Input validation on all endpoints

## Development

### Running Tests
```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test
```

### Building for Production

#### Backend
```bash
cd backend
go build -o linkedin-connector
```

#### Frontend
```bash
cd frontend
npm run build
```

## Deployment

### Backend Deployment Options
- **Heroku**: Use the included Procfile
- **Railway/Render**: Auto-detects Go applications
- **Docker**: Build and deploy containerized version
- **VPS**: Run the compiled binary with systemd

### Frontend Deployment Options
- **Vercel**: Connect to GitHub repository
- **Netlify**: Auto-deploys from repository
- **Cloudflare Pages**: Fast global CDN
- **AWS S3 + CloudFront**: Static hosting

### Environment Variables for Production
Make sure to set these in your hosting platform:
- Backend: `PORT`, `JWT_SECRET`, `UNIPILE_API_KEY`, `UNIPILE_API_URL`, `DATABASE_PATH`
- Frontend: `VITE_API_URL`

## Troubleshooting

### Backend won't start
- Check if port 8080 is already in use
- Verify `.env` file exists and is properly configured
- Ensure Go 1.22+ is installed: `go version`

### Frontend can't connect to backend
- Verify backend is running on the specified port
- Check CORS settings in backend
- Ensure `VITE_API_URL` points to correct backend URL

### LinkedIn connection fails
- Verify Unipile API credentials are correct
- Check if cookie/credentials are valid
- Review Unipile API rate limits

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License - feel free to use this project for your own purposes.

## Support

For issues or questions:
- Open an issue on GitHub
- Check Unipile documentation: https://docs.unipile.com

## Acknowledgments

- [Unipile](https://unipile.com) for the LinkedIn API integration
- [Gin](https://gin-gonic.com/) for the excellent Go web framework
- [React](https://react.dev/) and [Vite](https://vitejs.dev/) for the modern frontend stack

