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

# View a ticket
docker run -v alexandria-data:/root/work/DB/Alexandria alexandria view \
  --project "Alexandria" --id "1762471479992286465"

# Update a ticket
docker run -v alexandria-data:/root/work/DB/Alexandria alexandria update \
  --project "Alexandria" --id "1762471479992286465" --status "in-progress" --priority high

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
docker-compose run  alexandria view --project "Alexandria" --id "123"
docker-compose run  alexandria update --project "Alexandria" --id "123" --status "in-progress"
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
Alexandria view --project "Alexandria" --id "123"
Alexandria update --project "Alexandria" --id "123" --status "in-progress"
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

### View a Ticket

```bash
./Alexandria view --project "ProjectName" [--id ID | --title "Ticket Title"]
```

Options:
- `--project, -p` - Project name (required)
- `--id, -i` - Ticket ID to view
- `--title, -t` - Ticket title to view

**Note:** Either `--id` or `--title` must be provided (not both).

Examples:
```bash
# View ticket by ID
./Alexandria view --project "Alexandria" --id "1699564789123456789"

# View ticket by title
./Alexandria view --project "Alexandria" --title "Fix login bug"

# Using short flags
./Alexandria view -p "Alexandria" -i "1699564789123456789"
```

The command outputs the full ticket details in JSON format, including all fields, tags, files, and comments.

### Update a Ticket

```bash
./Alexandria update --project "ProjectName" [--id ID | --title "Ticket Title"] [options]
```

Options:
- `--project` - Project name (required)
- `--id, -i` - Ticket ID to update
- `--title, -t` - Find ticket by title to update
- `--new-title` - New title for the ticket
- `--description, -d` - New description for the ticket
- `--type` - New type: bug, feature, task
- `--status` - New status: open, in-progress, closed
- `--priority, -p` - New priority: undefined, low, medium, high
- `--criticalpath, -c` - Mark ticket as critical path (boolean flag)
- `--assigned-to, -a` - Assign ticket to user
- `--created-by` - Update ticket creator
- `--tags` - Comma-separated list of tags (replaces existing)
- `--files` - Comma-separated list of file paths (replaces existing)
- `--comments` - Comma-separated list of comments to add

**Note:** `--project`, and either `--id` or `--title` must be provided to identify the ticket. At least one field to update must be specified.

Examples:
```bash
# Update ticket status by ID
./Alexandria update --project "Alexandria" --id "1699564789123456789" --status "in-progress"

# Update multiple fields by title
./Alexandria update --project "Alexandria" --title "Fix login bug" --status "closed" --priority high

# Change ticket assignment and add tags
./Alexandria update --project "Alexandria" --id "1699564789123456789" --assigned-to "john@example.com" --tags "security,urgent,reviewed"

# Update description and mark as critical
./Alexandria update --project "Alexandria" --id "1699564789123456789" --description "Updated requirements" --criticalpath

# Add comments to a ticket
./Alexandria update --project "Alexandria" --id "1699564789123456789" --comments "Fixed in PR #123,Ready for review"

# Change the ticket title
./Alexandria update --project "Alexandria" --title "Fix login bug" --new-title "Fix authentication issue"

# Using short flags
./Alexandria update --project "Alexandria" -i "1699564789123456789" -a "jane@example.com" -p high
```

**Behavior:**
- Only specified fields are updated; unspecified fields remain unchanged
- Tags and files are replaced entirely when specified (not appended)
- Comments are added to existing comments (not replaced)
- The `updated_at` timestamp is automatically set to the current time

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

### Switch Database Source

```bash
./Alexandria source [sqlite|turso]
```

Options:
- `--status` - Show current database configuration

**Note:** When switching to Turso, ensure `TURSO_URL` and `TURSO_AUTH_TOKEN` environment variables are set.

Examples:
```bash
# Switch to local SQLite database
./Alexandria source sqlite

# Switch to Turso cloud database
./Alexandria source turso

# Show current database configuration
./Alexandria source --status
```

```

