# Quick Start Guide

Get started with Brain CLI in 5 minutes.

## Installation

### Option 1: Pre-built Binary (Easiest)

Download the latest release for your platform from the [releases page](https://github.com/yourusername/brain-cli/releases).

```bash
# macOS/Linux
chmod +x brain-*
sudo mv brain-* /usr/local/bin/brain

# Test it
brain --help
```

### Option 2: Install with Go

```bash
go install github.com/yourusername/brain-cli@latest
```

### Option 3: Build from Source

```bash
git clone https://github.com/yourusername/brain-cli
cd brain-cli
make build
sudo make install
```

## Setup

### For Best Results (Recommended)

Set your OpenAI API key for high-quality semantic search:

```bash
export OPENAI_API_KEY="sk-..."
```

Add this to your `~/.bashrc` or `~/.zshrc` to make it permanent.

### Without OpenAI (Works Offline)

Brain will automatically fall back to a local embedder. It's less accurate but still functional.

## First Steps

### 1. Add Your First Note

```bash
brain add "Redis caching reduced our API response time by 60%"
```

### 2. Add More Notes

```bash
brain add "Use context.WithTimeout for all HTTP client calls" --tags go,best-practices

brain add "Team decided on PostgreSQL over MySQL for better JSON support" --project myapp

brain add "Learned that goroutine leaks happen when channels aren't closed properly" --tags go,debugging
```

### 3. Search Your Brain

```bash
# Semantic search - finds conceptually similar notes
brain search "making APIs faster"

# Search with tags
brain search "concurrency" --tags go

# Limit results
brain search "database" --limit 3
```

### 4. Ask Questions

```bash
brain ask "What have I learned about Go best practices?"
brain ask "How did we improve performance?"
```

### 5. Get Contextual Notes

```bash
cd ~/projects/my-project
brain context
```

This shows notes relevant to your current directory/git repo.

## Daily Workflow

### As You Learn

Whenever you learn something useful, immediately save it:

```bash
# After fixing a bug
brain add "Bug was caused by goroutine reading from closed channel" --tags go,debugging

# After a code review
brain add "Team prefers early returns over nested ifs for readability" --tags style

# After reading docs
brain add "PostgreSQL EXPLAIN ANALYZE shows actual query execution time" --tags postgresql,performance
```

### When You Need It

Instead of searching Google or Stack Overflow for something you know you've learned:

```bash
brain search "postgres query optimization"
brain ask "How do I debug slow queries?"
```

### Project-Specific Notes

```bash
# When starting work on a project
cd ~/projects/api-service
brain context

# Add project-specific learnings
brain add "Auth service URL is https://auth.internal.example.com" --project api-service
```

## Tips & Tricks

### Use Descriptive Notes

❌ Bad: "Fixed it"
✅ Good: "Fixed N+1 query by using eager loading with .includes(:author)"

### Tag Consistently

Create a personal tagging system:
- Language: `go`, `python`, `rust`
- Category: `performance`, `security`, `debugging`
- Team: `team-decision`, `architecture`

### Review Regularly

```bash
# See all your notes (once we add list command)
brain list --recent

# Or search broadly
brain search "performance"
```

### Context is King

Always run `brain context` when starting work - you'll be surprised what's relevant.

## Common Use Cases

### Personal Knowledge Base

```bash
brain add "TIL: TypeScript 'satisfies' keyword for type checking without widening"
brain add "React useCallback only needed when passing to memoized children"
```

### Team Knowledge

```bash
brain add "Deployment process: merge to main → auto-deploy to staging → manual prod" --project myapp
brain add "Production database backups run at 2 AM UTC daily"
```

### Code Snippets

```bash
brain add "Retry logic: for i := 0; i < 3; i++ { if err = fn(); err == nil { break }}" --tags go,pattern
```

### Debugging Lessons

```bash
brain add "Memory leak was caused by keeping references in closure" --tags javascript,debugging
```

## Next Steps

- Read the [full README](README.md) for all features
- Check [ARCHITECTURE.md](ARCHITECTURE.md) to understand how it works
- See [CONTRIBUTING.md](CONTRIBUTING.md) to add features you want

## Need Help?

Open an issue on GitHub or search your brain for help:

```bash
brain search "how do I"
```

Just kidding on that last one - but once you've used it enough, it'll actually work!
