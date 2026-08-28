# Distributed Task Queue

A small Go task-queue example backed by Redis and MongoDB. The HTTP API stores each submitted task in MongoDB and places its ID on a Redis list. A separate worker process reads task IDs, simulates processing, and records status updates in MongoDB.

## Features

- Submit tasks over HTTP and receive a task ID.
- Store task data and status in MongoDB.
- Queue task IDs in Redis.
- Run API and worker processes separately.
- Inspect one task or list all stored tasks.
- Simulate task duration and occasional failure/retry in the worker.

## Requirements

- Go 1.21.1 or later
- Redis
- MongoDB

The repository does not include a Docker Compose file; start Redis and MongoDB separately before running the services.

## Configuration

The API and worker both read these variables from the process environment. If a `.env` file exists in the current working directory, it is loaded as well.

| Variable | Required | Description | Example |
| --- | --- | --- | --- |
| `REDIS_ADDR` | Yes | Redis host and port | `localhost:6379` |
| `MONGO_URI` | Yes | MongoDB connection URI | `mongodb://localhost:27017` |
| `MONGO_DB_NAME` | No | MongoDB database name; defaults to `tasks` | `tasks` |

Create a `.env` file in the project root, or export the variables in your shell:

```dotenv
REDIS_ADDR=localhost:6379
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=tasks
```

## Run locally

From the repository root, install dependencies:

```bash
go mod download
```

Start the API in one terminal:

```bash
go run ./cmd/api
```

Start the worker in another terminal:

```bash
go run ./cmd/worker
```

The API listens on `http://localhost:8080`. Redis and MongoDB must be reachable using the configured addresses before either process starts.

## API

### Submit a task

`POST /submit-task` accepts a JSON task. The server assigns the ID, creation time, and initial `queued` status; clients provide the task type, payload, and processing duration in seconds.

```bash
curl -X POST http://localhost:8080/submit-task \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "payload": {"to": "person@example.com", "subject": "Hello"},
    "timeToComplete": 2
  }'
```

Example response (`202 Accepted`):

```json
{
  "message": "Task submitted",
  "task_id": "a-task-uuid"
}
```

### Get a task's status

`GET /task/status/:task_id` returns the ID and current status.

```bash
curl http://localhost:8080/task/status/a-task-uuid
```

```json
{
  "task_id": "a-task-uuid",
  "status": "in-progress"
}
```

### List all tasks

`GET /tasks/status` returns all task documents currently stored in MongoDB.

```bash
curl http://localhost:8080/tasks/status
```

## Task lifecycle

Tasks are stored in the MongoDB `tasks` collection. The initial status is `queued`; workers update it to `in-progress` and then `success` or `failed`. Successful completion also records elapsed time in `completion_time` (seconds). The task payload is an arbitrary JSON object.

The worker simulates processing by sleeping for `timeToComplete` seconds. Its failure behavior is also simulated: when the current Unix second is divisible by 10, it marks the task failed and inserts a new task record for retry. This is example behavior, not a production retry policy: retries have no limit, and the retried record receives a new task ID.

## Project layout

```text
api/                 HTTP routes and handlers
cmd/api/             API service entry point
cmd/worker/          Worker service entry point
config/              Environment and Redis/MongoDB initialization
internal/models/     Task model
internal/queue/      Redis queue operations
internal/repository/ MongoDB task operations
internal/worker/     Task polling and simulated processing
pkg/logger/          Logger utility
```

## Notes

- The API port is currently fixed at `8080`.
- The Redis queue is a list named `tasks_queue`.
- The worker entry point starts worker loops in groups of one, two, and three (six loops total).
- Task listing currently returns every task without pagination.
