# Deployment Guide

This guide covers various deployment options for the LinkedIn Connector application.

## Table of Contents
- [Local Development](#local-development)
- [Docker Deployment](#docker-deployment)
- [Cloud Deployment](#cloud-deployment)
  - [Backend Deployment](#backend-deployment)
  - [Frontend Deployment](#frontend-deployment)

---

## Local Development

### Prerequisites
- Go 1.22+
- Node.js 18+
- Unipile API credentials

### Backend Setup

1. Navigate to backend directory:
```bash
cd backend
```

2. Create `.env` file:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Install dependencies:
```bash
go mod download
```

4. Run the server:
```bash
go run main.go
```

Backend will be available at `http://localhost:8080`

### Frontend Setup

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Create `.env` file:
```bash
cp .env.example .env
# Edit .env with your backend URL
```

3. Install dependencies:
```bash
npm install
```

4. Run development server:
```bash
npm run dev
```

Frontend will be available at `http://localhost:5173`

---

## Docker Deployment

### Using Docker Compose (Recommended for local/staging)

1. Create `.env` file in project root:
```bash
cp .env.example .env
```

2. Edit `.env` with your configuration:
```env
JWT_SECRET=your-super-secret-jwt-key
UNIPILE_API_KEY=your-unipile-api-key
UNIPILE_API_URL=https://api.unipile.com/v1
FRONTEND_URL=http://localhost:3000
VITE_API_URL=http://localhost:8080
```

3. Build and run:
```bash
docker-compose up -d
```

4. View logs:
```bash
docker-compose logs -f
```

5. Stop containers:
```bash
docker-compose down
```

Application will be available:
- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`

### Building Individual Docker Images

#### Backend
```bash
docker build -f Dockerfile.backend -t linkedin-connector-backend .
docker run -p 8080:8080 \
  -e JWT_SECRET=your-secret \
  -e UNIPILE_API_KEY=your-key \
  linkedin-connector-backend
```

#### Frontend
```bash
docker build -f Dockerfile.frontend -t linkedin-connector-frontend .
docker run -p 3000:80 linkedin-connector-frontend
```

---

## Cloud Deployment

### Backend Deployment

#### Option 1: Railway

1. Create account at [railway.app](https://railway.app)
2. Install Railway CLI:
```bash
npm i -g @railway/cli
```

3. Login and initialize:
```bash
railway login
railway init
```

4. Add environment variables in Railway dashboard:
   - `PORT=8080`
   - `JWT_SECRET=<your-secret>`
   - `UNIPILE_API_KEY=<your-key>`
   - `UNIPILE_API_URL=https://api.unipile.com/v1`
   - `FRONTEND_URL=<your-frontend-url>`

5. Deploy:
```bash
railway up
```

#### Option 2: Render

1. Create account at [render.com](https://render.com)
2. Create new Web Service
3. Connect your GitHub repository
4. Configure:
   - **Build Command**: `cd backend && go build -o linkedin-connector`
   - **Start Command**: `./backend/linkedin-connector`
5. Add environment variables in Render dashboard

#### Option 3: Heroku

1. Create `Procfile` in backend directory:
```
web: ./linkedin-connector
```

2. Deploy:
```bash
cd backend
heroku create your-app-name
heroku buildpacks:set heroku/go
heroku config:set JWT_SECRET=<your-secret>
heroku config:set UNIPILE_API_KEY=<your-key>
git push heroku main
```

#### Option 4: AWS EC2 / VPS

1. SSH into your server
2. Install Go 1.22:
```bash
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

3. Clone and build:
```bash
git clone <your-repo>
cd linkedin-connector/backend
go build -o linkedin-connector
```

4. Create systemd service (`/etc/systemd/system/linkedin-connector.service`):
```ini
[Unit]
Description=LinkedIn Connector API
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/linkedin-connector/backend
ExecStart=/opt/linkedin-connector/backend/linkedin-connector
Restart=on-failure
Environment="PORT=8080"
Environment="JWT_SECRET=your-secret"
Environment="UNIPILE_API_KEY=your-key"

[Install]
WantedBy=multi-user.target
```

5. Start service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable linkedin-connector
sudo systemctl start linkedin-connector
```

6. Setup nginx as reverse proxy:
```nginx
server {
    listen 80;
    server_name api.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Frontend Deployment

#### Option 1: Vercel (Recommended)

1. Create account at [vercel.com](https://vercel.com)
2. Install Vercel CLI:
```bash
npm i -g vercel
```

3. Deploy:
```bash
cd frontend
vercel
```

4. Set environment variable:
   - `VITE_API_URL=<your-backend-url>`

#### Option 2: Netlify

1. Create account at [netlify.com](https://netlify.com)
2. Install Netlify CLI:
```bash
npm i -g netlify-cli
```

3. Build and deploy:
```bash
cd frontend
npm run build
netlify deploy --prod --dir=dist
```

4. Set environment variable in Netlify dashboard:
   - `VITE_API_URL=<your-backend-url>`

#### Option 3: Cloudflare Pages

1. Create account at [pages.cloudflare.com](https://pages.cloudflare.com)
2. Connect GitHub repository
3. Configure build settings:
   - **Build command**: `cd frontend && npm run build`
   - **Build output directory**: `frontend/dist`
4. Add environment variable:
   - `VITE_API_URL=<your-backend-url>`

#### Option 4: AWS S3 + CloudFront

1. Build the frontend:
```bash
cd frontend
npm run build
```

2. Create S3 bucket and enable static website hosting

3. Upload files:
```bash
aws s3 sync dist/ s3://your-bucket-name/
```

4. Create CloudFront distribution pointing to S3 bucket

5. Update DNS to point to CloudFront

---

## Production Checklist

### Security
- [ ] Change `JWT_SECRET` to a strong random value
- [ ] Use HTTPS for both frontend and backend
- [ ] Set up CORS properly
- [ ] Enable rate limiting
- [ ] Set up firewall rules
- [ ] Use environment variables for all secrets
- [ ] Enable database backups

### Performance
- [ ] Enable gzip compression
- [ ] Set up CDN for frontend
- [ ] Configure caching headers
- [ ] Optimize database queries
- [ ] Set up connection pooling

### Monitoring
- [ ] Set up error logging (e.g., Sentry)
- [ ] Configure uptime monitoring
- [ ] Set up performance monitoring
- [ ] Configure alerts for errors

### Database
- [ ] For production, consider PostgreSQL instead of SQLite
- [ ] Set up automated backups
- [ ] Configure connection limits
- [ ] Enable SSL for database connections

---

## Database Migration from SQLite to PostgreSQL

If you want to use PostgreSQL instead of SQLite:

1. Update `backend/go.mod`:
```go
require (
    // ... existing imports
    gorm.io/driver/postgres v1.5.4
)
```

2. Update `backend/database/database.go`:
```go
import (
    "gorm.io/driver/postgres"
    // ...
)

func InitDatabase(dbPath string) error {
    dsn := "host=localhost user=postgres password=yourpass dbname=linkedin_connector port=5432 sslmode=disable"
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    // ...
}
```

3. Add PostgreSQL configuration to `.env`:
```env
DATABASE_URL=postgres://user:password@host:port/dbname
```

---

## Troubleshooting

### Backend Issues

**Port already in use:**
```bash
# Kill process on port 8080
lsof -ti:8080 | xargs kill -9
```

**Database locked:**
```bash
# Check if another process is using the database
lsof backend/linkedin_connector.db
```

### Frontend Issues

**Build fails:**
```bash
# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm install
```

**API connection fails:**
- Check `VITE_API_URL` in `.env`
- Verify CORS settings in backend
- Check if backend is running

### Docker Issues

**Permission denied:**
```bash
sudo usermod -aG docker $USER
# Logout and login again
```

**Port conflicts:**
```bash
# Change ports in docker-compose.yml
# Or stop conflicting containers
docker ps
docker stop <container-id>
```

---

## Support

For issues or questions:
- Check the main README.md
- Review Unipile documentation: https://docs.unipile.com
- Open an issue on GitHub

