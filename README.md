# Alexandria

Alexandria is a project to create a better approach to managing work, it's engineer lead and use a CLI interface.

## Build

```bash
go build -o Alexandria
```

## Commands

### Create a Ticket

```bash
./Alexandria create --title "Ticket title" --project "ProjectName" [options]
```

Options:
- `--title, -t` - Ticket title (required)
- `--project` - Project name (required)
- `--description, -d` - Ticket description
- `--type` - Ticket type: bug, feature, task (default: task)
- `--priority, -p` - Priority: undefined, low, medium, high (default: undefined)
- `--criticalpath, -c` - Mark as critical path (default: false)
- `--assigned-to, -a` - Assign to user
- `--created-by` - Ticket creator
- `--tags` - Comma-separated list of tags

Example:
```bash
./Alexandria create --title "Fix login bug" --project "Alexandria" --description "Users unable to login" --type bug --priority high --criticalpath --tags "security,urgent"
```

### List Tickets

```bash
./Alexandria list [options]
```

Options:
- `--status` - Filter by status: open, in-progress, closed
- `--type` - Filter by type: bug, feature, task
- `--priority` - Filter by priority: undefined, low, medium, high
- `--assigned-to` - Filter by assigned user
- `--tags` - Filter by tags (comma-separated)
- `--output, -o` - Output format: json, table, summary (default: table)

Examples:
```bash
# List all tickets
./Alexandria list

# List open bugs
./Alexandria list --status open --type bug

# List high priority tickets
./Alexandria list --priority high

# List with JSON output
./Alexandria list --output json

# List tickets with specific tags
./Alexandria list --tags "security,urgent"
```

## Database

Tickets are stored in a SQLite database at:
```
~/work/DB/Alexandria/tickets.db
```

