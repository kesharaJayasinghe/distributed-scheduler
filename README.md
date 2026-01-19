# Distributed Task Scheduler (Go + Postgres)

A robust, fault-tolerant distributed task scheduler built in Go. This system is designed to execute tasks asynchronously with "at-least-once" delivery guarantees, handling concurrency across multiple worker nodes and recovering from node failures automatically.

## System Architecture

The system consists of two distinct microservices sharing a single PostgreSQL database:

- **API Service (Producer)**: Accepts HTTP requests to schedule tasks (e.g., "send email at 3:00 PM").
- **Scheduler Service (Consumer/Worker)**: A background process that polls for due tasks, executes them, and manages their lifecycle.

## Key Features

- **Concurrency Control**: Implements Distributed Locking using PostgreSQL's `FOR UPDATE SKIP LOCKED`. This allows multiple scheduler instances to run simultaneously without race conditions (no double-execution of tasks).
- **Fault Tolerance**: Includes a "Zombie Task" rescue mechanism. If a worker node crashes mid-task, a background reaper detects the stale state and resets the task to `PENDING` for another worker to pick up.
- **Idempotency Ready**: Designed to support idempotent task execution to ensure safety under "at-least-once" delivery semantics.
- **Horizontal Scalability**: The stateless worker design allows for adding more worker nodes simply by spinning up new containers.

## Tech Stack

- **Language**: Go (Golang) 1.21+
- **Database**: PostgreSQL 16 (Persistence & Concurrency Control)
- **Driver**: pgx/v5 (High-performance connection pooling)
- **Infrastructure**: Docker & Docker Compose

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose

### 1. Start the Infrastructure

Spin up the PostgreSQL container:

```bash
docker compose up -d
```

### 2. Run the API (Producer)

In a new terminal:

```bash
# Downloads dependencies and starts the HTTP server on :8080
go mod tidy
go run cmd/api/main.go
```

### 3. Run the Scheduler (Consumer)

In a separate terminal:

```bash
# Starts the polling worker
go run cmd/scheduler/main.go
```

### 4. Schedule a Task

Send a task scheduled for execution 10 seconds in the future:

```bash
# Note: Ensure the timestamp is in the future relative to UTC
curl -X POST http://localhost:8080/tasks \
     -H "Content-Type: application/json" \
     -d '{
           "due_at": "2026-01-20T12:00:00Z",
           "payload": {"action": "email_user", "user_id": 42}
         }'
```

## Design Decisions & Trade-offs

### Why Postgres for a Queue?

**Decision**: Use PostgreSQL over a dedicated message broker (like RabbitMQ or Kafka) or Redis.

**Reasoning**: For high-value tasks (like payments or emails) where data loss is unacceptable, the ACID compliance of Postgres is superior to Redis.

**The "Skip Locked"**: Historically, using SQL DBs as queues caused massive lock contention. However, PostgreSQL's `SELECT ... FOR UPDATE SKIP LOCKED` feature allows us to fetch available rows without blocking other transactions, making it performant enough for thousands of tasks per second without the operational overhead of managing a Kafka cluster.

### Handling Node Failures

**Challenge**: If a worker crashes after marking a task `RUNNING` but before marking it `COMPLETED`, the task remains stuck forever.

**Solution**: A "Visibility Timeout" pattern. A background goroutine checks for tasks that have been `RUNNING` > 2 minutes and resets them to `PENDING`.

**Trade-off**: This implies "At-Least-Once" delivery. If a worker is just slow (not dead) and the rescue kicks in, the task might execute twice. This necessitates that the task logic itself be idempotent.

## Future Improvements

If this were going to production with millions of users, the following can be added:

- **Exponential Backoff**: If a task fails (e.g., 3rd party API is down), retry it with increasing delays (1m, 5m, 15m) instead of failing immediately.
- **Dead Letter Queue (DLQ)**: After 5 failed retries, move the task to a separate `failed_tasks` table for manual inspection so it doesn't clog the queue.
- **gRPC**: Replace REST with gRPC for internal communication if we split the architecture further.