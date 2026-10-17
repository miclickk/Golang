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