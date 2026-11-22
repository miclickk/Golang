# Password Analyzer — SIS Checkpoint 1

## Domain Description
This application analyzes password strings and evaluates their strength based on length, digit count, uppercase letters, and special characters.

## Validation Rules
- Password must not be empty (`ErrEmptyPassword`).
- Password must not contain spaces (`ErrPasswordHasSpaces`).

## Task Lifecycle and Allowed Transitions
Each task progresses through three states:
1. `Pending` — Initial state assigned when a valid task is created.
2. `Processing` — Intermediate state entered when task processing starts.
3. `Completed` — Final state reached after the password analysis result is stored.

### Allowed State Transitions
- `Pending` -> `Processing` -> `Completed`

Any other state transition (such as attempting to process a task that is already `Completed`) is rejected and returns `ErrInvalidState`.


# Password Analyzer (Checkpoint 2 - Concurrent Processing)

## Overview
An in-memory Go application that validates passwords and evaluates their strength (length, digits, uppercase letters, and special characters) concurrently using a worker pool.

## Architecture
- **Task Creation & Submission:** `CreateTask` validates and stores a task in memory (`Pending`). `SubmitTask` pushes the task ID into a bounded channel queue (`Queued`) and returns immediately.
- **Worker Pool:** `Start(workers, queueSize)` launches a fixed number of worker goroutines (3 workers in `main.go`) that read task IDs from a shared buffered channel (capacity 5 in `main.go`).
- **Full Queue Handling:** Non-blocking rejection. If the queue is full, `SubmitTask` returns `ErrQueueFull` and leaves the task in the `Pending` state.
- **Task Lifecycle:** `Pending -> Queued -> Processing -> Completed`. Duplicate or invalid transitions return `ErrInvalidState`.
- **Synchronization:** `sync.RWMutex` protects concurrent access to the task map and task fields. Password processing runs outside the lock so workers execute concurrently.
- **Shutdown:** `Stop()` marks the manager as stopped, closes the queue channel so workers drain all accepted tasks, and waits for all workers to finish via `sync.WaitGroup`. Submissions after shutdown return `ErrManagerStopped`.

## Running & Testing
```bash
go run .
gofmt -l .
go vet ./...
go test ./...
go test -race ./...
