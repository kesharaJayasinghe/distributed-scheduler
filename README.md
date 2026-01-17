# Distributed Scheduler

A distributed task scheduler built with Go and PostgreSQL.

## Prerequisites

- Go 1.25.6 or higher
- Docker and Docker Compose

## Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd distributed-scheduler
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables**
   
   Copy the example environment file:
   ```bash
   cp .env.example .env
   ```
   
   The `.env` file contains:
   - `DATABASE_URL`: PostgreSQL connection string
   - `PORT`: Application port (default: 8080)
   - `ENV`: Environment mode (development/staging/production)

4. **Start the PostgreSQL database**
   ```bash
   docker-compose up -d
   ```

5. **Run database migrations** (if needed)
   ```bash
   # Apply schema
   docker exec -i scheduler_db psql -U user -d scheduler < schema.sql
   ```

6. **Run the application**
   ```bash
   go run cmd/api/main.go
   ```

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| DATABASE_URL | PostgreSQL connection string | - | Yes |
| PORT | Application port | 8080 | No |
| ENV | Environment mode | development | No |

## Database Configuration

The PostgreSQL database runs in Docker with the following configuration:
- Host: `localhost`
- Port: `5433`
- Database: `scheduler`
- User: `user`
- Password: `password`

## Development

### Start the database
```bash
docker-compose up -d
```

### Stop the database
```bash
docker-compose down
```

### View database logs
```bash
docker-compose logs -f postgres
```

### Connect to database
```bash
docker exec -it scheduler_db psql -U user -d scheduler
```

## Project Structure

```
.
├── cmd/
│   └── api/          # Application entrypoint
├── internal/
│   ├── config/       # Configuration management
│   ├── db/           # Database connection
│   └── task/         # Task-related logic
├── docker-compose.yml
├── schema.sql        # Database schema
└── .env              # Environment variables (not in version control)
```
