# Password Analyzer — SIS Checkpoint 2 (Concurrent Processing)

## Overview
An in-memory Go application that validates passwords and evaluates their strength (length, digits, uppercase letters, and special characters) concurrently using a worker pool.

## Validation Rules
- Password must not be empty (`ErrEmptyPassword`).
- Password must not contain spaces (`ErrPasswordHasSpaces`).

## Concurrency Architecture
- **Task Creation & Submission:** `CreateTask` validates input, assigns a unique ID, and stores the task in memory (`Pending`). `SubmitTask` places the task ID into a bounded channel queue (`Queued`) and returns immediately without waiting for processing.
- **Worker Pool:** `Start(workers, queueSize)` launches a fixed pool of worker goroutines (3 workers in `main.go`, minimum 2 enforced) that continuously consume task IDs from a shared buffered channel (capacity 5 in `main.go`).
- **Full Queue Handling:** Non-blocking rejection. If the queue is full, `SubmitTask` returns `ErrQueueFull` and leaves the task in the `Pending` state.
- **Synchronization:** `sync.RWMutex` protects concurrent access to the in-memory `map` and mutable task fields. The actual password analysis (`processor.Process`) runs outside the lock so workers process tasks in parallel.
- **Shutdown:** `Stop()` marks the manager as stopped, closes the queue channel so workers finish all accepted tasks, and waits for all worker goroutines to exit via `sync.WaitGroup`. Calling `Stop()` multiple times is safe, and submissions after shutdown return `ErrManagerStopped`.

## Task Lifecycle and Allowed Transitions
1. `Pending` — Task created, waiting to be submitted or processed.
2. `Queued` — Task accepted into the bounded queue.
3. `Processing` — Worker is currently analyzing the password.
4. `Completed` — Processing finished and result is stored.

Allowed transitions:
- Concurrent path: `Pending -> Queued -> Processing -> Completed`
- Synchronous path (`ProcessTask`): `Pending -> Processing -> Completed`

Duplicate or invalid transitions return `ErrInvalidState`.

## Running & Testing
```bash
go run .
gofmt -l .
go vet ./...
go test ./...
go test -race ./...
```