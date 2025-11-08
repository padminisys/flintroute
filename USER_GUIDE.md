# FlintRoute User Guide

## Welcome to FlintRoute

FlintRoute is a modern BGP route management system with a web-based interface for monitoring and managing BGP routing configurations.

## Getting Started

### Starting the Application

1. **Start the FlintRoute server:**
   ```bash
   ./bin/flintroute --config configs/config.yaml
   ```

2. **Verify the server is running:**
   - The server will start on port 8080 by default
   - You should see log messages indicating successful startup
   - The process will run in the foreground (use `&` to run in background)

### Accessing the Web UI

1. **Open your web browser** and navigate to:
   ```
   http://localhost:8080
   ```

2. **You should see the FlintRoute login page**

### Logging In

**Default Credentials:**
- **Username:** `admin`
- **Password:** `admin`

> ⚠️ **Security Note:** Please change the default password after your first login in a production environment.

**Login Steps:**
1. Enter username: `admin`
2. Enter password: `admin`
3. Click the "Login" button
4. You will be redirected to the dashboard upon successful authentication

## Features

### Dashboard
After logging in, you'll have access to the main dashboard which provides:
- **BGP Peer Status:** View all configured BGP peers and their connection status
- **Route Information:** Monitor BGP routes and routing tables
- **Real-time Updates:** WebSocket-based live updates for BGP state changes
- **System Health:** Overview of the routing system health and statistics

### BGP Management
- **View BGP Peers:** See all configured BGP neighbors
- **Monitor Routes:** Track BGP route advertisements and withdrawals
- **Configuration:** Manage BGP settings and peer configurations
- **Alerts:** Receive notifications for BGP state changes

### API Access
FlintRoute provides a RESTful API for programmatic access:

**Authentication:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'
```

**Using the API with JWT Token:**
```bash
# Save the access_token from login response
TOKEN="your_access_token_here"

# Make authenticated API calls
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/bgp/peers
```

**Available API Endpoints:**
- `GET /api/v1/bgp/peers` - List all BGP peers
- `GET /api/v1/bgp/routes` - List BGP routes
- `POST /api/v1/config/reload` - Reload configuration
- `GET /api/v1/health` - System health check

## Configuration

### Configuration File
The main configuration file is located at `configs/config.yaml`. Key settings include:

- **Server Port:** Default is 8080
- **Database:** SQLite database location
- **FRR Integration:** FRRouting daemon connection settings
- **Authentication:** JWT token settings and user management
- **Logging:** Log level and output configuration

### Customizing Configuration
1. Edit `configs/config.yaml` with your preferred settings
2. Restart the FlintRoute server for changes to take effect

## Stopping the Application

### If running in foreground:
Press `Ctrl+C` to stop the server gracefully

### If running in background:
```bash
# Find the process ID
ps aux | grep flintroute

# Stop the process
kill <PID>

# Or use pkill
pkill flintroute
```

## Troubleshooting

### Server Won't Start
- **Check if port 8080 is already in use:**
  ```bash
  lsof -i :8080
  ```
- **Verify configuration file exists:**
  ```bash
  ls -la configs/config.yaml
  ```
- **Check file permissions:**
  ```bash
  ls -la bin/flintroute
  ```

### Cannot Access Web UI
- **Verify server is running:**
  ```bash
  curl http://localhost:8080/
  ```
- **Check firewall settings** if accessing from another machine
- **Ensure correct URL:** `http://localhost:8080` (not https)

### Login Issues
- **Verify credentials:** Default is admin/admin
- **Check browser console** for JavaScript errors (F12 in most browsers)
- **Clear browser cache** and try again

### API Authentication Errors
- **Verify token is valid:** Tokens expire after 15 minutes
- **Check Authorization header format:** Must be `Bearer <token>`
- **Ensure token is not expired:** Request a new token if needed

### FRR Integration Issues
- **Verify FRR is installed and running:**
  ```bash
  systemctl status frr
  ```
- **Check FRR configuration** in `configs/frr/`
- **Review logs** for connection errors

## Advanced Usage

### Running with Custom Configuration
```bash
./bin/flintroute --config /path/to/custom/config.yaml
```

### Enabling Debug Logging
Edit `configs/config.yaml` and set:
```yaml
log:
  level: debug
```

### Using Docker
```bash
docker-compose up -d
```

### Monitoring Logs
```bash
# If running in background, check system logs
journalctl -u flintroute -f

# Or redirect output to a file when starting
./bin/flintroute --config configs/config.yaml > flintroute.log 2>&1 &
```

## Security Best Practices

1. **Change Default Password:** Immediately change the admin password
2. **Use HTTPS:** Configure a reverse proxy (nginx/Apache) with SSL/TLS
3. **Firewall Rules:** Restrict access to port 8080 to trusted networks
4. **Regular Updates:** Keep FlintRoute and dependencies up to date
5. **Backup Configuration:** Regularly backup your configuration files
6. **Monitor Access:** Review authentication logs regularly

## Getting Help

- **Documentation:** Check the `docs/` directory for detailed documentation
- **Logs:** Review application logs for error messages
- **GitHub Issues:** Report bugs or request features on the project repository
- **Configuration Examples:** See `configs/config.example.yaml` for reference

## Next Steps

1. **Explore the Dashboard:** Familiarize yourself with the web interface
2. **Configure BGP Peers:** Add your BGP neighbors in the configuration
3. **Set Up Monitoring:** Configure alerts for important BGP events
4. **Integrate with FRR:** Connect to your FRRouting daemon
5. **Customize Settings:** Adjust configuration to match your network requirements

---

**Version:** 1.0.0  
**Last Updated:** November 2025  
**License:** See LICENSE file for details