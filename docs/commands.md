# Alexandria Commands Reference

## Alternative Build Methods

### Running with Docker

**Build the Docker image:**
```bash
docker build -t alexandria .
```

**Run with Docker:**
```bash
# Create a volume for persistent database storage
docker volume create alexandria-data

# Create a ticket
docker run -v alexandria-data:/root/work/DB/Alexandria alexandria create \
  --title "Fix login bug" --project "Alexandria" --type bug --priority high

# List tickets
docker run -v alexandria-data:/root/work/DB/Alexandria alexandria list

# View a ticket
docker run -v alexandria-data:/root/work/DB/Alexandria alexandria view \
  --id "1762471479992286465"

# Update a ticket
docker run -v alexandria-data:/root/work/DB/Alexandria alexandria update \
  --id "1762471479992286465" --status "in-progress" --priority high

# Delete a ticket
docker run -v alexandria-data:/root/work/DB/Alexandria alexandria delete \
  --id "1762471479992286465"
```

**Using docker-compose:**

The repository includes a `docker-compose.yml` file. To use it:
```bash
# Build and run
docker-compose run alexandria create --title "New ticket" --project "Alexandria"
docker-compose run alexandria list
docker-compose run alexandria view --id "123"
docker-compose run alexandria update --id "123" --status "in-progress"
```

**Shell Alias (optional convenience):**

Add this to your `.bashrc` or `.zshrc` for easier usage:
```bash
# For local installation
alias Alexandria='./Alexandria'

# For Docker installation
alias Alexandria='docker run -v alexandria-data:/root/work/DB/Alexandria alexandria'

# For docker-compose installation
alias Alexandria='docker-compose run alexandria'
```

Then you can use:
```bash
Alexandria create --title "My task" --project "Alexandria"
Alexandria list
Alexandria view --id "123"
Alexandria update --id "123" --status "in-progress"
Alexandria delete --id "123"
```

### Running Locally with Go Build

**Build:**
```bash
go build -o Alexandria
```

**Run:**
```bash
./Alexandria [command] [flags]
```

**Note:** This method requires you to run Alexandria from the project directory. For system-wide installation, use the Makefile method described in the main README.

---

## Command Reference

**Note:** If you installed using `make install`, use `alexandria` (lowercase). If running locally with `go build`, use `./Alexandria`.

### Create a Ticket

```bash
alexandria create --title "Ticket title" --project "ProjectName" [options]
```

**Options:**
- `--title, -t` - Ticket title (required)
- `--project` - Project name (required)
- `--description, -d` - Ticket description
- `--type` - Ticket type: bug, feature, task (default: task)
- `--priority, -p` - Priority: undefined, low, medium, high (default: undefined)
- `--criticalpath, -c` - Mark as critical path (default: false)
- `--assigned-to, -a` - Assign to user
- `--created-by` - Ticket creator
- `--tags` - Comma-separated list of tags

**Example:**
```bash
alexandria create --title "Fix login bug" --project "Alexandria" --description "Users unable to login" --type bug --priority high --criticalpath --tags "security,urgent"
```

### List Tickets

```bash
alexandria list [options]
```

**Options:**
- `--status` - Filter by status: open, in-progress, closed
- `--type` - Filter by type: bug, feature, task
- `--priority` - Filter by priority: undefined, low, medium, high
- `--assigned-to` - Filter by assigned user
- `--tags` - Filter by tags (comma-separated)
- `--output, -o` - Output format: json, table, summary (default: table)

**Examples:**
```bash
# List all tickets
alexandria list

# List open bugs
alexandria list --status open --type bug

# List high priority tickets
alexandria list --priority high

# List with JSON output
alexandria list --output json

# List tickets with specific tags
alexandria list --tags "security,urgent"
```

### View a Ticket

```bash
alexandria view [--id ID | --title "Ticket Title"] [--project "ProjectName"]
```

**Options:**
- `--id, -i` - Ticket ID to view (required if title not provided)
- `--title, -t` - Ticket title to view (required if ID not provided)
- `--project, -p` - Project name (optional filter)

**Note:** Either `--id` or `--title` must be provided.

**Examples:**
```bash
# View ticket by ID
alexandria view --id "1699564789123456789"

# View ticket by title
alexandria view --title "Fix login bug"

# View ticket by ID with project filter (more precise if IDs might overlap)
alexandria view --id "1699564789123456789" --project "Alexandria"

# View ticket by title with project filter (recommended when titles might not be unique)
alexandria view --title "Fix login bug" --project "Alexandria"

# Using short flags
alexandria view -i "1699564789123456789" -p "Alexandria"
```

The command outputs the full ticket details in JSON format, including all fields, tags, files, and comments.

### Update a Ticket

```bash
alexandria update [--id ID | --title "Ticket Title"] [--project "ProjectName"] [options]
```

**Options:**
- `--id, -i` - Ticket ID to update (required if title not provided)
- `--title, -t` - Find ticket by title to update (required if ID not provided)
- `--project` - Project name (optional filter)
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

**Note:** Either `--id` or `--title` must be provided to identify the ticket. At least one field to update must be specified.

**Examples:**
```bash
# Update ticket status by ID
alexandria update --id "1699564789123456789" --status "in-progress"

# Update multiple fields by title
alexandria update --title "Fix login bug" --status "closed" --priority high

# Update with project filter for precision
alexandria update --id "1699564789123456789" --project "Alexandria" --status "in-progress"

# Change ticket assignment and add tags
alexandria update --id "1699564789123456789" --assigned-to "john@example.com" --tags "security,urgent,reviewed"

# Update description and mark as critical
alexandria update --id "1699564789123456789" --description "Updated requirements" --criticalpath

# Add comments to a ticket
alexandria update --id "1699564789123456789" --comments "Fixed in PR #123,Ready for review"

# Change the ticket title with project filter (recommended for title lookups)
alexandria update --title "Fix login bug" --project "Alexandria" --new-title "Fix authentication issue"

# Using short flags
alexandria update -i "1699564789123456789" -a "jane@example.com" -p high
```

**Behavior:**
- Only specified fields are updated; unspecified fields remain unchanged
- Tags and files are replaced entirely when specified (not appended)
- Comments are added to existing comments (not replaced)
- The `updated_at` timestamp is automatically set to the current time

### Delete a Ticket

```bash
alexandria delete [--id ID | --title "Ticket Title"] [--project "ProjectName"]
```

**Options:**
- `--id, -i` - Ticket ID to delete (required if title not provided)
- `--title, -t` - Ticket title to delete (required if ID not provided)
- `--project, -p` - Project name (optional filter)

**Note:** Either `--id` or `--title` must be provided.

**Examples:**
```bash
# Delete by ID
alexandria delete --id "1699564789123456789"

# Delete by title
alexandria delete --title "Fix login bug"

# Delete with project filter (recommended for safer deletion)
alexandria delete --id "1699564789123456789" --project "Alexandria"

# Delete by title with project filter (highly recommended to avoid deleting wrong ticket)
alexandria delete --title "Fix login bug" --project "Alexandria"

# Using short flags
alexandria delete -i "1699564789123456789" -p "Alexandria"
```

**Warning:** This command will permanently delete the ticket and all related data including tags, files, and comments. When using title-based deletion without a project filter, the first matching ticket will be deleted. It's strongly recommended to use the `--project` filter for safer deletions.

### Switch Database Source

```bash
alexandria source [sqlite|turso]
```

**Options:**
- `--status` - Show current database configuration

**Examples:**
```bash
# Switch to local SQLite database (default)
alexandria source sqlite

# Switch to Turso cloud database
alexandria source turso

# Show current database configuration
alexandria source --status
```

#### Setting up Turso Database

To use Turso cloud database, you need to configure environment variables.

**1. Create a `.env` file from the example:**
```bash
cp .env.example .env
```

**2. Install Turso CLI (if not already installed):**
```bash
# macOS/Linux
curl -sSfL https://get.tur.so/install.sh | bash

# Or using Homebrew
brew install tursodatabase/tap/turso
```

**3. Authenticate with Turso:**
```bash
turso auth login
```

**4. Create a database (or use existing):**
```bash
# Create a new database
turso db create alexandria

# List existing databases
turso db list
```

**5. Get your database credentials:**
```bash
# Get database URL
turso db show alexandria

# Create an authentication token
turso db tokens create alexandria
```

**6. Update your `.env` file:**
```bash
TURSO_URL=libsql://your-db.turso.io
TURSO_AUTH_TOKEN=your-token-here
```

**7. Switch Alexandria to use Turso:**
```bash
# Make sure your .env is loaded (or export the variables)
source .env  # or: export TURSO_URL=... && export TURSO_AUTH_TOKEN=...

# Switch to Turso
alexandria source turso
```

**Alternative - Set environment variables permanently:**

Add to your `~/.bashrc` or `~/.zshrc`:
```bash
export TURSO_URL="libsql://your-db.turso.io"
export TURSO_AUTH_TOKEN="your-token-here"
```

Then reload your shell:
```bash
source ~/.bashrc  # or source ~/.zshrc
```
