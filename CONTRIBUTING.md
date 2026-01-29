# Contributing to Brain CLI

Thanks for your interest in contributing! This project is a learning experiment and all contributions are welcome.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/brain-cli`
3. Create a branch: `git checkout -b feature/your-feature`
4. Make your changes
5. Run tests: `make test`
6. Format code: `make fmt`
7. Commit: `git commit -am 'Add some feature'`
8. Push: `git push origin feature/your-feature`
9. Create a Pull Request

## Development Setup

```bash
# Install dependencies
make deps

# Build the project
make build

# Run tests
make test

# Install locally for testing
make install-local
```

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing (or use `make fmt`)
- Write tests for new functionality
- Keep functions small and focused
- Add comments for exported functions and types

## Project Structure

```
brain-cli/
├── cmd/                 # CLI commands (cobra)
│   ├── root.go
│   ├── add.go
│   ├── search.go
│   └── context.go
├── internal/
│   └── brain/          # Core brain logic
│       ├── brain.go    # Main Brain type
│       ├── embedder.go # Embedding generation
│       └── vectorstore.go # Vector storage & search
├── main.go             # Entry point
└── README.md
```

## Adding a New Command

1. Create a new file in `cmd/` (e.g., `cmd/list.go`)
2. Define the command using cobra:

```go
package cmd

import (
    "github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List all notes",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Your implementation
        return nil
    },
}

func init() {
    rootCmd.AddCommand(listCmd)
}
```

3. Add tests in `cmd/list_test.go`
4. Update README.md

## Adding a New Feature

1. Open an issue first to discuss the feature
2. Write tests for the feature
3. Implement the feature
4. Update documentation
5. Submit a PR

## Testing

Write tests for new functionality. We use Go's built-in testing:

```go
func TestMyFeature(t *testing.T) {
    // Setup
    // Exercise
    // Verify
    // Teardown
}
```

Run tests with:
```bash
make test
# or
go test ./...
```

## Ideas for Contributions

- **Easy**:
  - Add a `brain list` command to show all notes
  - Add a `brain delete <id>` command
  - Add color output with a library like fatih/color
  - Add shell completion (cobra supports this)

- **Medium**:
  - Implement better local embeddings using sentence-transformers
  - Add export to Markdown/Obsidian format
  - Add `brain ask` for conversational queries
  - Improve context detection (file type analysis, package.json, etc.)

- **Hard**:
  - Add sync between machines (p2p or cloud)
  - Build a web UI
  - Create VS Code extension
  - Add graph visualization of note relationships

## Questions?

Open an issue! This is a learning project, so all questions are welcome.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
