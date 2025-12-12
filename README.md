# Alexandria

Alexandria is a project to create a better approach to managing work, it's engineer lead and uses a CLI interface. Database supplied is either SQLite locally or Turso,  by default all tickets are written to 1 table and grouped by project. 

For more detailed commands see the docs dir.

## Installation

The easiest way to build and install Alexandria is using the provided Makefile.

**Build and install to ~/.local/bin:**
```bash
make install
```

This will build the binary and install it to `~/.local/bin`, allowing you to run `alexandria` from any directory.

**Important:** Make sure `~/.local/bin` is in your PATH. Add this to your `~/.bashrc` or `~/.zshrc`:
```bash
export PATH="$HOME/.local/bin:$PATH"
```

Then reload your shell:
```bash
source ~/.bashrc  # or source ~/.zshrc
```

## Quick Start

After installation, you can use Alexandria from any directory:

```bash
# Create a ticket
alexandria create --title "Fix login bug" --project "MyProject" --type bug --priority high

# List all tickets
alexandria list

# View a ticket
alexandria view --project "MyProject" --id "1699564789123456789"

# Update a ticket
alexandria update --project "MyProject" --id "1699564789123456789" --status "in-progress"

# Delete a ticket
alexandria delete --project "MyProject" --id "1699564789123456789"
```

## Documentation

- **[Commands Reference](docs/commands.md)** - Detailed documentation for all commands and options
- **[Alternative Build Methods](docs/commands.md#alternative-build-methods)** - Docker and Go build instructions

## Features

- **Multiple ticket types**: bug, feature, task
- **Priority management**: undefined, low, medium, high
- **Critical path tracking**: Mark important tickets
- **Tags and assignments**: Organize and assign work
- **Flexible output formats**: table, JSON, summary
- **Database options**: Local SQLite or cloud Turso
- **Comments and file attachments**: Full ticket context

## Configuration

Alexandria stores its configuration and data in `~/Alexandria/`.

### Database Options

Switch between local SQLite and cloud Turso databases:

```bash
# Use local SQLite (default)
alexandria source sqlite

# Use Turso cloud database
alexandria source turso

# Check current database configuration
alexandria source --status
```

**Note:** When using Turso, set these environment variables:
- `TURSO_URL` - Your Turso database URL
- `TURSO_AUTH_TOKEN` - Your Turso authentication token

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

See [LICENSE](LICENSE) for details.
