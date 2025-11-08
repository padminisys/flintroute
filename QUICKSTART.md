# FlintRoute Quick Start Guide

## Prerequisites Check

Before starting, ensure you have the following installed:

- **Go 1.21 or later**: Check with `go version`
- **Node.js 18 or later**: Check with `node --version`
- **npm**: Check with `npm --version`
- **FRR (Free Range Routing)**: Optional for BGP functionality
  - Install on Debian: `sudo apt-get install frr`

## Quick Setup (5 minutes)

### 1. Clone and Navigate

```bash
git clone <repository-url>
cd flintroute
```

### 2. Backend Setup

```bash
# Install Go dependencies
go mod tidy

# Create configuration file
cp configs/config.example.yaml configs/config.yaml

# Build the backend
go build -o bin/flintroute ./cmd/flintroute

# Or use the Makefile
make build
```

### 3. Frontend Setup

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Create environment file
cp .env.example .env

# Build for production
npm run build

# Return to root directory
cd ..
```

## Running the Application

### Option 1: Production Mode (Recommended for Testing)

```bash
# Start the backend (from project root)
./bin/flintroute --config configs/config.yaml
```

The server will start on `http://localhost:8080`

### Option 2: Development Mode

**Terminal 1 - Backend:**
```bash
# From project root
go run ./cmd/flintroute --config configs/config.yaml
```

**Terminal 2 - Frontend:**
```bash
# From frontend directory
cd frontend
npm run dev
```

Frontend dev server will start on `http://localhost:5173`

### Option 3: Using Docker Compose (with FRR)

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## Accessing the Application

### Web Interface

1. Open your browser and navigate to:
   - **Production**: `http://localhost:8080`
   - **Development**: `http://localhost:5173`

2. **Login with default credentials:**
   - Username: `admin`
   - Password: `admin`

   ⚠️ **IMPORTANT**: Change the default password immediately after first login!

### API Endpoints

The backend exposes the following endpoints:

- **Health Check**: `GET http://localhost:8080/api/health`
- **Login**: `POST http://localhost:8080/api/auth/login`
- **BGP Sessions**: `GET http://localhost:8080/api/bgp/sessions`
- **WebSocket**: `ws://localhost:8080/ws`

Test the health endpoint:
```bash
curl http://localhost:8080/api/health
```

## Default Configuration

The application uses the following defaults:

- **Server Address**: `0.0.0.0:8080`
- **Database**: SQLite at `./data/flintroute.db`
- **JWT Secret**: `your-secret-key-change-this-in-production`
- **FRR VTY**: `localhost:2605`
- **Admin User**: `admin` / `admin`

## File Locations

- **Configuration**: `configs/config.yaml`
- **Database**: `./data/flintroute.db`
- **Binary**: `./bin/flintroute`
- **Frontend Build**: `./frontend/dist/`
- **Logs**: Console output (configure file logging in config.yaml)

## Common Troubleshooting

### Backend won't start

**Issue**: Port 8080 already in use
```bash
# Check what's using port 8080
sudo lsof -i :8080

# Kill the process or change port in configs/config.yaml
```

**Issue**: Database permission error
```bash
# Ensure data directory exists and is writable
mkdir -p data
chmod 755 data
```

### Frontend build fails

**Issue**: Node modules not installed
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

**Issue**: TypeScript errors
```bash
cd frontend
npm run build -- --force
```

### Cannot connect to FRR

**Issue**: FRR not running
```bash
# Check FRR status
sudo systemctl status frr

# Start FRR
sudo systemctl start frr

# Enable FRR on boot
sudo systemctl enable frr
```

**Issue**: VTY connection refused
```bash
# Check FRR VTY configuration in /etc/frr/vtysh.conf
# Ensure zebra daemon is enabled in /etc/frr/daemons
```

### Login fails

**Issue**: Invalid credentials
- Use default credentials: `admin` / `admin`
- Check backend logs for authentication errors

**Issue**: JWT token errors
- Ensure JWT secret is set in `configs/config.yaml`
- Clear browser cookies and try again

## Next Steps

1. **Change Default Password**: 
   - Login and navigate to Settings → Change Password

2. **Configure FRR Integration**:
   - Edit `configs/config.yaml` to match your FRR setup
   - Restart the backend

3. **Add BGP Sessions**:
   - Use the web interface to configure BGP neighbors
   - Monitor session status in real-time

4. **Explore Documentation**:
   - Read `docs/architecture/overview.md` for system architecture
   - Check `docs/development/setup.md` for development details
   - Review `docs/phase1/roadmap.md` for feature roadmap

## Quick Commands Reference

```bash
# Build everything
make build

# Run tests
make test

# Clean build artifacts
make clean

# Start backend
./bin/flintroute --config configs/config.yaml

# Start frontend dev server
cd frontend && npm run dev

# Build frontend for production
cd frontend && npm run build

# View backend logs (if running in background)
tail -f logs/flintroute.log

# Check database
sqlite3 data/flintroute.db ".tables"
```

## Getting Help

- **Documentation**: Check the `docs/` directory
- **Issues**: Report bugs on the project issue tracker
- **Logs**: Check console output or log files for errors
- **Configuration**: Review `configs/config.example.yaml` for all options

## Security Notes

⚠️ **Before deploying to production:**

1. Change the default admin password
2. Update JWT secret in configuration
3. Use HTTPS/TLS for production deployments
4. Configure proper firewall rules
5. Review and update CORS settings
6. Enable authentication for FRR VTY access
7. Regular security updates for dependencies

---

**Ready to start?** Run `./bin/flintroute --config configs/config.yaml` and open `http://localhost:8080` in your browser!