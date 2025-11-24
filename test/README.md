# Alexandria E2E Tests

End-to-end tests that simulate real user behavior by executing the actual binary.

## Running Tests

```bash
cd test
go test -v
```

## What's Tested

- **Binary Build**: Verifies the application compiles successfully
- **Help Command**: Tests help output
- **Database Switching**:
  - `source --status` - shows current database
  - `source sqlite` - switches to SQLite
  - `source turso` - switches to Turso (requires env vars)
- **Ticket Management**:
  - `create` - creates tickets with various options
  - `list` - lists tickets (table and JSON formats)
  - `list --filter` - filters by status, type, etc.
  - `view` - views ticket details
  - `update` - updates ticket fields
  - `delete` - deletes tickets
- **Complete Workflow**: Creates tickets, switches databases, verifies behavior

## Test Environment

- Tests run in isolated temporary directories
- Builds a separate test binary (`alexandria-test`)
- Cleans up automatically after tests complete
- No mocks - uses real database connections and commands
