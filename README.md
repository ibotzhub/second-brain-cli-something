# 🧠 Brain CLI

A second brain for developers. Save your ideas, learnings, and insights, then surface them exactly when you need them using semantic search and context awareness.

## Why Brain?

Ever had that "I know I learned this before..." moment? Brain remembers for you. It's not just a note-taking tool—it understands what you're working on and surfaces relevant knowledge automatically.

### Features

- 🔍 **Semantic Search**: Find notes by meaning, not just keywords
- 📍 **Context Awareness**: Get relevant notes based on your current project/directory
- 🏷️ **Smart Tagging**: Organize with tags and projects
- ⚡ **Lightning Fast**: Built in Go, instant responses
- 🔒 **Local First**: Your data stays on your machine
- 🤖 **AI-Powered**: Uses embeddings for intelligent search (OpenAI or local fallback)

## Installation

```bash
go install github.com/yourusername/brain-cli@latest
```

Or clone and build:

```bash
git clone https://github.com/yourusername/brain-cli
cd brain-cli
go build -o brain
sudo mv brain /usr/local/bin/
```

## Quick Start

```bash
# Add a note
brain add "Redis caching reduced our API latency by 60%"

# Add with tags
brain add "Always use context.WithTimeout for API calls" --tags go,best-practices

# Add to a project
brain add "Team decided on tabs over spaces" --project myapp

# Search semantically
brain search "making APIs faster"

# Get notes relevant to your current directory/project
brain context
```

## Commands

### `brain add`

Save a note to your brain.

```bash
brain add "Your insight here"
brain add "Note content" --tags tag1,tag2
brain add "Note content" --project myproject
brain add "Note content" --tags go,performance --project api-service
```

### `brain search`

Search your notes semantically.

```bash
brain search "performance optimization"
brain search "database" --limit 10
brain search "API design" --tags go
```

The search understands meaning—searching for "making things faster" will find notes about "performance optimization" even if they don't contain those exact words.

### `brain context`

Show notes relevant to what you're currently working on.

```bash
cd ~/projects/my-api
brain context
# Shows notes related to this project, recent commits, etc.
```

Brain analyzes:
- Current directory name
- Git repository (if present)
- Recent commit messages
- Project metadata

## Configuration

### OpenAI Embeddings (Recommended)

For best results, set your OpenAI API key:

```bash
export OPENAI_API_KEY="sk-..."
```

Brain uses `text-embedding-3-small` which costs ~$0.02 per 1M tokens. For a typical note, that's less than $0.0001.

### Local Embeddings (Fallback)

If no API key is set, Brain falls back to a simple local embedder. It works but won't be as accurate for semantic search.

## How It Works

1. **You save a note**: Brain generates an embedding (semantic vector) of your note
2. **You search**: Brain generates an embedding of your query
3. **Magic happens**: Brain finds notes with similar embeddings using cosine similarity
4. **You get results**: Ranked by relevance, not just keyword matching

## Data Storage

All data is stored locally in `~/.brain/`:
- `notes.json`: Your notes and metadata
- `vectors.db`: Embeddings for fast search

## Examples

**Learning from a bug fix:**
```bash
brain add "Goroutine leak fixed by using context cancellation in HTTP client" --tags go,debugging
```

**Later, working on HTTP code:**
```bash
brain search "http client issues"
# Finds your note about goroutine leaks
```

**Team decision:**
```bash
brain add "Decided to use PostgreSQL over MySQL for better JSON support" --project myapp
```

**Months later in that project:**
```bash
cd ~/projects/myapp
brain context
# Reminds you about the PostgreSQL decision
```

## Roadmap

- [ ] Interactive `brain ask` command (chat with your notes)
- [ ] Auto-tagging with LLMs
- [ ] Export to Markdown/Obsidian
- [ ] Sync between machines
- [ ] Browser extension for saving web insights
- [ ] Integration with IDE (VS Code extension)
- [ ] Link detection between related notes
- [ ] Spaced repetition reminders

## Contributing

Contributions welcome! This is a learning project and an experiment in building useful dev tools.

## License

MIT

## Inspiration

Inspired by the concept of a "second brain" and tools like Obsidian, but built specifically for developers who live in the terminal.

---

**Built with Go because josh liked it and he has good taste.**
