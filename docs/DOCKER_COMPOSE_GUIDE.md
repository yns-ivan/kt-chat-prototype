# Docker Compose Command Guide

This guide explains the differences between `docker-compose` and `docker compose` commands and how to use them in the KT Chat project.

## 🔄 Two Different Commands

### `docker compose` (New - Recommended)
- **Type**: Docker CLI plugin
- **Availability**: Since Docker 20.10.13
- **Installation**: Comes with Docker Desktop and Docker CLI
- **Status**: **Latest and recommended**
- **Performance**: Better performance and integration
- **Features**: More features and better Docker ecosystem integration

### `docker-compose` (Legacy)
- **Type**: Standalone binary
- **Availability**: Traditional installation
- **Installation**: Needs to be installed separately
- **Status**: Still supported but being phased out
- **Performance**: Older, less optimized
- **Features**: Basic functionality

## 🧪 How to Check Which One You Have

### Check for `docker compose` (new)
```bash
docker compose --version
```

### Check for `docker-compose` (legacy)
```bash
docker-compose --version
```

### Check both
```bash
# Check new version
if command -v docker compose &> /dev/null; then
    echo "✅ docker compose (new) is available"
fi

# Check legacy version
if command -v docker-compose &> /dev/null; then
    echo "✅ docker-compose (legacy) is available"
fi
```

## 🚀 Using Docker Compose in KT Chat

### Automatic Detection
The project includes scripts that automatically detect and use the correct command:

```bash
# Use the wrapper script (recommended)
./scripts/docker-compose-wrapper.sh up -d

# Or use the test script to see which command is detected
./scripts/test-setup.sh
```

### Manual Usage

#### If you have `docker compose` (new):
```bash
# Start services
docker compose up -d

# Stop services
docker compose down

# View logs
docker compose logs -f

# Build and start
docker compose up -d --build
```

#### If you have `docker-compose` (legacy):
```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f

# Build and start
docker-compose up -d --build
```

## 🛠️ Using Makefile Commands

The Makefile automatically uses the correct command:

```bash
cd backend

# Start services
make docker-run

# Stop services
make docker-stop

# View logs
make logs
```

## 📋 Command Comparison

| Feature | `docker compose` | `docker-compose` |
|---------|------------------|------------------|
| **Installation** | Comes with Docker CLI | Separate installation |
| **Performance** | Better | Slower |
| **Integration** | Native Docker CLI | External tool |
| **Updates** | With Docker updates | Manual updates |
| **Future** | Actively developed | Maintenance mode |
| **Compatibility** | Backward compatible | Legacy |

## 🔧 Migration Guide

### From `docker-compose` to `docker compose`

1. **Update Docker**: Ensure you have Docker 20.10.13 or later
2. **Test the new command**: `docker compose --version`
3. **Update scripts**: Replace `docker-compose` with `docker compose`
4. **Update documentation**: Update any references

### Example Migration

**Before:**
```bash
docker-compose up -d
docker-compose logs -f
docker-compose down
```

**After:**
```bash
docker compose up -d
docker compose logs -f
docker compose down
```

## 🎯 Best Practices

### 1. Use the Wrapper Script
```bash
# Instead of guessing which command to use
./scripts/docker-compose-wrapper.sh up -d
```

### 2. Use Makefile Commands
```bash
# The Makefile handles the command detection
make docker-run
make docker-stop
make logs
```

### 3. Check Your Environment
```bash
# Run the test script to verify your setup
./scripts/test-setup.sh
```

### 4. Update Your Scripts
If you have custom scripts, update them to use the wrapper:

```bash
#!/bin/bash
# Instead of hardcoding docker-compose
# Use the wrapper script
./scripts/docker-compose-wrapper.sh "$@"
```

## 🐛 Troubleshooting

### Command Not Found
```bash
# Check if Docker is installed
docker --version

# Check if docker compose is available
docker compose --version

# Check if docker-compose is available
docker-compose --version
```

### Permission Issues
```bash
# Make sure the wrapper script is executable
chmod +x scripts/docker-compose-wrapper.sh
```

### Version Conflicts
If you have both commands installed:
1. The wrapper script will prefer `docker compose` (newer)
2. You can manually specify which one to use
3. Consider uninstalling the legacy version to avoid confusion

## 📚 Additional Resources

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Docker CLI Plugin Documentation](https://docs.docker.com/engine/reference/commandline/compose/)
- [Migration Guide](https://docs.docker.com/compose/migrate/)

## 🎉 Summary

- **Use `docker compose`** (new) when possible - it's the future
- **Use the wrapper script** for automatic detection
- **Use Makefile commands** for convenience
- **Both commands work** - the project supports both
- **Test your setup** with `./scripts/test-setup.sh`

The KT Chat project is designed to work with both commands, so you can use whichever one is available on your system! 