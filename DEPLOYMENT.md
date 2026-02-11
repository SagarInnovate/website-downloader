# Deployment Guide

## Table of Contents
- [Local Development](#local-development)
- [Production Build](#production-build)
- [Docker Deployment](#docker-deployment)
- [Environment Configuration](#environment-configuration)
- [Troubleshooting](#troubleshooting)

## Local Development

### Backend

1. Navigate to backend directory:
```bash
cd backend
```

2. Install dependencies:
```bash
go mod download
```

3. Run the development server:
```bash
go run main.go
```

The backend will start on `http://localhost:8080`

### Frontend

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start development server:
```bash
npm run dev
```

The frontend will start on `http://localhost:5173`

## Production Build

### Backend

1. Build the binary:
```bash
cd backend
go build -o website-downloader main.go
```

For Windows:
```bash
go build -o website-downloader.exe main.go
```

2. Run the binary:
```bash
./website-downloader
```

### Frontend

1. Build for production:
```bash
cd frontend
npm run build
```

This creates an optimized build in the `dist` folder.

2. Serve with a static file server:

Using Python:
```bash
cd dist
python -m http.server 8000
```

Using Node.js `serve`:
```bash
npm install -g serve
serve -s dist -p 3000
```

Using Nginx (recommended):
```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    root /path/to/frontend/dist;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # Proxy API requests to backend
    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

## Docker Deployment

### Backend Dockerfile

Create `backend/Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o website-downloader .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/website-downloader .
COPY --from=builder /app/config.json .

# Create directories
RUN mkdir -p temp downloads

EXPOSE 8080

CMD ["./website-downloader"]
```

Build and run:
```bash
cd backend
docker build -t website-downloader-backend .
docker run -p 8080:8080 -v $(pwd)/downloads:/root/downloads website-downloader-backend
```

### Frontend Dockerfile

Create `frontend/Dockerfile`:

```dockerfile
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm install

# Copy source code
COPY . .

# Build the app
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy built files
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

Create `frontend/nginx.conf`:

```nginx
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://backend:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

Build and run:
```bash
cd frontend
docker build -t website-downloader-frontend .
docker run -p 80:80 website-downloader-frontend
```

### Docker Compose

Create `docker-compose.yml` in the root directory:

```yaml
version: '3.8'

services:
  backend:
    build: ./backend
    container_name: downloader-backend
    ports:
      - "8080:8080"
    volumes:
      - ./downloads:/root/downloads
      - ./temp:/root/temp
    environment:
      - GIN_MODE=release
    restart: unless-stopped

  frontend:
    build: ./frontend
    container_name: downloader-frontend
    ports:
      - "80:80"
    depends_on:
      - backend
    restart: unless-stopped

volumes:
  downloads:
  temp:
```

Run with Docker Compose:
```bash
docker-compose up -d
```

Stop services:
```bash
docker-compose down
```

View logs:
```bash
docker-compose logs -f
```

## Environment Configuration

### Backend Environment Variables

Create `backend/.env`:

```env
# Server Configuration
GIN_MODE=release
PORT=8080

# Scraper Configuration
MAX_DEPTH=5
MAX_PAGES=1000
WORKER_COUNT=10
REQUEST_TIMEOUT=30
RATE_LIMIT=5

# Directories
TEMP_DIR=./temp
OUTPUT_DIR=./downloads
```

### Frontend Environment Variables

For production, update `frontend/.env.production`:

```env
VITE_API_URL=https://your-domain.com
VITE_WS_URL=wss://your-domain.com
```

## Cloud Deployment

### Deploy to Heroku

Backend:
```bash
cd backend
echo "web: ./website-downloader" > Procfile
git init
heroku create your-app-backend
git add .
git commit -m "Initial commit"
git push heroku main
```

Frontend:
```bash
cd frontend
echo '{ "root": "dist/" }' > static.json
heroku create your-app-frontend
heroku buildpacks:set https://github.com/heroku/heroku-buildpack-static
git add .
git commit -m "Initial commit"
git push heroku main
```

### Deploy to AWS

1. **Backend on EC2**:
   - Launch an EC2 instance (Ubuntu/Amazon Linux)
   - Install Go
   - Clone repository
   - Build and run the backend
   - Use systemd for service management

2. **Frontend on S3 + CloudFront**:
   - Build the frontend (`npm run build`)
   - Upload `dist` folder to S3
   - Configure S3 for static website hosting
   - Create CloudFront distribution
   - Point to S3 bucket

### Deploy to DigitalOcean

1. Create a Droplet
2. Install Docker and Docker Compose
3. Clone repository
4. Run `docker-compose up -d`
5. Configure firewall (allow ports 80, 8080)

## Monitoring & Logs

### Systemd Service (Linux)

Create `/etc/systemd/system/website-downloader.service`:

```ini
[Unit]
Description=Website Downloader Backend
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/website-downloader/backend
ExecStart=/opt/website-downloader/backend/website-downloader
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable website-downloader
sudo systemctl start website-downloader
sudo systemctl status website-downloader
```

View logs:
```bash
sudo journalctl -u website-downloader -f
```

## Troubleshooting

### Backend Issues

**Port already in use:**
```bash
# Find process using port 8080
netstat -ano | findstr :8080  # Windows
lsof -i :8080                 # Linux/Mac

# Kill the process
taskkill /PID <pid> /F        # Windows
kill -9 <pid>                 # Linux/Mac
```

**Go dependencies error:**
```bash
go clean -modcache
go mod download
go mod tidy
```

### Frontend Issues

**Build fails:**
```bash
# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm install
npm run build
```

**WebSocket connection fails:**
- Check that backend is running
- Verify WebSocket URL in `.env`
- Check firewall/proxy settings

### Docker Issues

**Container won't start:**
```bash
docker logs <container-id>
docker inspect <container-id>
```

**Clean up Docker:**
```bash
docker system prune -a
docker volume prune
```

## Security Considerations

1. **Use HTTPS** in production (Let's Encrypt)
2. **Set rate limits** to prevent abuse
3. **Validate input URLs** on backend
4. **Set CORS properly** for your domain
5. **Use environment variables** for secrets
6. **Keep dependencies updated**

## Performance Optimization

1. **Backend**:
   - Adjust worker count based on server resources
   - Use connection pooling
   - Implement caching for repeated requests

2. **Frontend**:
   - Enable gzip compression
   - Use CDN for static assets
   - Optimize bundle size

---

For additional help, see the main [README.md](README.md) or open an issue.
