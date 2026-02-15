# MetaManager

MetaManager is a CLI tool for managing your files metadata. It supports both local filesystem and Google Drive, allowing you to track, tag, search, and organize your files across different contexts.

## Features

- **Context Management**: Create and switch between multiple named environments (local, Google Drive, work, etc.)
- **File Tracking**: Track local directories or Google Drive folders
- **Tagging System**: Add and manage tags for files and directories
- **ID Management**: Assign unique IDs to nodes for quick access
- **Search Capabilities**: Find files and directories using regex patterns
- **Google Drive Integration**: Browse and manage Google Drive files directly from the CLI
- **Metadata Storage**: Persistent storage of file metadata and directory structures

## Prerequisites

- **Go 1.23.0** or later
- **Google Cloud Credentials** (for Google Drive features) - `credentials.json` file

## Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd MetaManager
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Configure Google Drive (Optional)

If you want to use Google Drive features:

1. **Create OAuth 2.0 Client ID** in Google Cloud Console:
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Navigate to "APIs & Services" > "Credentials"
   - Click "Create Credentials" > "OAuth client ID"
   - Choose "Desktop app" as the application type
   - Download the `credentials.json` file

2. Place the `credentials.json` file in the project root directory

3. The credentials will be embedded into the binary during build

**Note:** The `credentials.json` file contains a public OAuth client ID (not a secret). This is the standard way to configure OAuth for desktop applications and is safe to include in your repository. It's not a credential leak - Google designed these credentials to be public for desktop apps. The actual authentication happens through user login via `MetaManager login`.

### 4. Build the Application

```bash
go build -o MetaManager
```

Or install it:

```bash
go install
```

### 5. Initialize Your First Context

```bash
# Create a local context
./MetaManager context create local --type local --root /path/to/your/directory

# Or create a Google Drive context (requires login first)
./MetaManager login
./MetaManager context create gdrive --type gdrive --root "Drive Folder Name"
```

## Quick Start

### Basic Usage

```bash
# Set the current context
./MetaManager context set local

# Track a directory
./MetaManager track /path/to/directory

# List files in current directory
./MetaManager ls

# Search for files
./MetaManager search searchNode "pattern"

# Add tags to a file
./MetaManager tag add /path/to/file tag1 tag2

# Assign an ID to a file
./MetaManager id set /path/to/file my-unique-id

# Jump to a file by ID
./MetaManager id jump my-unique-id
```

### Google Drive Commands

```bash
# Login to Google Drive
./MetaManager login

# List Google Drive files
./MetaManager gdrive list

# Navigate Google Drive
./MetaManager gdrive cd "Folder Name"
./MetaManager gdrive ls
./MetaManager gdrive pwd
```

## Documentation

Complete command documentation is available in the `docs/` directory. All documentation is auto-generated from the CLI commands using Cobra.

### Main Documentation Files

- **[MetaManager.md](docs/MetaManager.md)** - Main command overview and all available commands
- **[context.md](docs/MetaManager_context.md)** - Context management commands
- **[track.md](docs/MetaManager_track.md)** - File tracking commands
- **[tag.md](docs/MetaManager_tag.md)** - Tagging commands
- **[id.md](docs/MetaManager_id.md)** - ID management commands
- **[gdrive.md](docs/MetaManager_gdrive.md)** - Google Drive commands
- **[search.md](docs/MetaManager_search.md)** - Search commands

### Viewing Documentation

You can view any command's documentation using:

```bash
# View help for any command
./MetaManager <command> --help

# Or read the markdown files directly
cat docs/MetaManager_<command>.md
```

### Regenerating Documentation

To regenerate the documentation files:

```bash
cd docs
go run gen_doc.go
```

This will update all markdown files in the `docs/` directory based on the current command structure.

## Project Structure

```
MetaManager/
├── cmd/              # CLI command implementations
├── docs/             # Auto-generated command documentation
├── internal/         # Internal packages
│   ├── data/        # Data management
│   ├── ds/          # Data structures
│   ├── file/        # File operations
│   ├── repository/  # Storage repositories
│   └── services/    # Service layer
├── main.go          # Application entry point
├── go.mod           # Go module definition
└── README.md        # This file
```

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
# Build for current platform
go build -o MetaManager

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o MetaManager-linux
```

### Debug Mode

Enable debug logging for troubleshooting:

```bash
./MetaManager --debug <command>
```

## Common Commands

| Command | Description |
|---------|-------------|
| `context create` | Create a new context |
| `context set` | Set the current context |
| `context list` | List all contexts |
| `track <path>` | Start tracking a directory |
| `untrack <path>` | Stop tracking a directory |
| `tag add <path> <tags...>` | Add tags to a file/directory |
| `tag list` | List all tags |
| `id set <path> <id>` | Assign an ID to a node |
| `id get <path>` | Get the ID of a node |
| `id jump <id>` | Print path for a given ID |
| `search searchNode <pattern>` | Search for files/directories |
| `gdrive list` | List Google Drive files |
| `login` | Authenticate with Google |

## Troubleshooting

### Google Drive Authentication Issues

If you encounter authentication issues:

1. Ensure `credentials.json` is in the project root
2. Run `./MetaManager login` to re-authenticate
3. Check that the token file has proper permissions

### Context Not Found

If you get "context not found" errors:

1. List contexts: `./MetaManager context list`
2. Create a context if needed: `./MetaManager context create <name> --type <local|gdrive> --root <path>`
3. Set the context: `./MetaManager context set <name>`

### Debug Mode

For detailed error messages, use debug mode:

```bash
./MetaManager --debug <command>
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

Copyright © 2023

## Support

For issues and questions, please open an issue in the repository.
