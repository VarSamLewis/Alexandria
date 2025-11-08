# Alexandria

Alexandria is a project to create a better approach to managing work, it's engineer lead and use a CLI interface.

## Installation & Running

### Running Locally

**Build:**
```bash
go build -o Alexandria
```

**Run:**
```bash
./Alexandria [command] [flags]
```

### Running with Docker

**Build the Docker image:**
```bash
docker build -t alexandria .
```

**Run with Docker:**
```bash
# Create a volume for persistent database storage
docker volume create alexandria-data

**Docker Examples:**
```bash
# Create a ticket
docker run  -v alexandria-data:/root/work/DB/Alexandria alexandria create \
  --title "Fix login bug" --project "Alexandria" --type bug --priority high

# List tickets
docker run -v alexandria-data:/root/work/DB/Alexandria alexandria list

# Delete a ticket
docker run -v alexandria-data:/root/work/DB/Alexandria alexandria delete \
  --project "Alexandria" --id "1762471479992286465"
```

**Using docker-compose (recommended):**

The repository includes a `docker-compose.yml` file. To use it:
```bash
# Build and run
docker-compose run  alexandria create --title "New ticket" --project "Alexandria"
docker-compose run  alexandria list
```

**Shell Alias (optional convenience):**

Add this to your `.bashrc` or `.zshrc` for easier usage:
```bash
# For local installation
alias Alexandria='./Alexandria'

# For Docker installation
alias Alexandria='docker run  -v alexandria-data:/root/work/DB/Alexandria alexandria'

# For docker-compose installation
alias Alexandria='docker-compose run alexandria'
```

Then you can use:
```bash
Alexandria create --title "My task" --project "Alexandria"
Alexandria list
Alexandria delete --project "Alexandria" --id "123"
```

### Local Commands

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

### Delete a Ticket

```bash
./Alexandria delete --project "ProjectName" [--id ID | --title "Ticket Title"]
```

Options:
- `--project, -p` - Project name (required)
- `--id, -i` - Ticket ID to delete
- `--title, -t` - Ticket title to delete

**Note:** Either `--id` or `--title` must be provided (not both).

Examples:
```bash
# Delete by ID
./Alexandria delete --project "Alexandria" --id "1699564789123456789"

# Delete by title
./Alexandria delete --project "Alexandria" --title "Fix login bug"

# Using short flags
./Alexandria delete -p "Alexandria" -i "1699564789123456789"
```

**Warning:** This command will permanently delete the ticket and all related data including tags, files, and comments.

```

