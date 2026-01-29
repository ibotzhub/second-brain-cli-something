# Testing Guide

Manual testing checklist to verify Brain CLI works correctly.

## Prerequisites

```bash
# Build the project
make build

# Or run directly
go run . --help
```

## Test 1: Basic Add & List

```bash
# Add a note
./brain add "Test note 1"

# Verify it was added
./brain list

# Expected: Shows 1 note with ID and timestamp
```

## Test 2: Tags & Projects

```bash
# Add with tags
./brain add "Go note" --tags go,programming

# Add with project
./brain add "Project note" --project testproject

# Add with both
./brain add "Complex note" --tags go,test --project testproject

# List all
./brain list

# Expected: Shows 4 notes total with tags/projects displayed
```

## Test 3: Tag Filtering

```bash
# Filter by tag
./brain list --tags go

# Expected: Shows only the 2 notes tagged with "go"

# Filter with limit
./brain list --tags go --limit 1

# Expected: Shows only 1 note
```

## Test 4: Semantic Search

```bash
# Add some related notes
./brain add "Redis caching improved performance by 60%" --tags performance
./brain add "Database query optimization using indexes" --tags database,performance
./brain add "API response time reduced with compression" --tags api,performance

# Search semantically (should find all performance-related notes)
./brain search "making things faster"

# Expected: Shows all 3 performance notes, ranked by relevance

# Search with specific term
./brain search "database"

# Expected: Database note should rank highest
```

## Test 5: Context Detection

```bash
# Initialize a git repo for testing
mkdir /tmp/test-brain-context
cd /tmp/test-brain-context
git init
git config user.email "test@example.com"
git config user.name "Test User"
echo "test" > README.md
git add .
git commit -m "Add database caching layer"

# Add a note related to this context
~/path/to/brain add "Database caching improves query performance" --project test-brain-context

# Test context
~/path/to/brain context

# Expected: Should show the database note as relevant
# Should display git repo name and recent commit
```

## Test 6: Ask Command

```bash
# Ask a question
./brain ask "What have I learned about performance?"

# Expected: Shows relevant performance-related notes
```

## Test 7: Persistence

```bash
# Add a note
./brain add "Persistence test note"

# Exit and restart (simulate)
# Just run list again - this loads from disk
./brain list

# Expected: Still shows all notes including "Persistence test note"
```

## Test 8: Duplication Bug (Fixed)

```bash
# Count notes
./brain list | grep -c "^\d\."

# Remember the count (let's say it's 8)

# "Restart" by running list again (loads from disk)
./brain list | grep -c "^\d\."

# Expected: SAME count as before (no duplicates)
# If count doubled, duplication bug is back!
```

## Test 9: Local Embedder (Offline Mode)

```bash
# Unset API key
unset OPENAI_API_KEY

# Add and search
./brain add "Local embedder test"
./brain search "embedder"

# Expected: Should work, though search quality may be lower
# Should not crash or error
```

## Test 10: Large Dataset

```bash
# Add many notes quickly
for i in {1..100}; do
  ./brain add "Test note $i about performance and optimization" --tags test
done

# Search should still be fast
time ./brain search "optimization"

# Expected: Returns results in < 1 second

# List with limit
./brain list --limit 10 --tags test

# Expected: Shows exactly 10 notes
```

## Unit Tests

```bash
# Run Go tests
go test ./...

# Expected: All tests pass
```

## Known Limitations to Verify

1. **Embeddings regenerate on startup**: First command after build/restart may be slow
2. **Linear search**: With 10,000+ notes, search might slow down (still under 1s though)
3. **Local embedder quality**: Without OpenAI, semantic search is less accurate

## Cleanup

```bash
# Remove test data
rm -rf ~/.brain

# Remove test git repo
rm -rf /tmp/test-brain-context
```

## Critical Bug Checks

### Duplication Bug (Fixed)
Run list multiple times - count should stay same.

### Normalization Bug (Fixed)
With local embedder, similarity scores should be 0.0-1.0, not weird values.

### Context Stub (Fixed)
Context command should show actual git/directory info, not just search for "context".

## Success Criteria

- All commands run without errors
- Notes persist across "restarts"
- Search returns relevant results
- No duplicate notes in storage
- Context uses actual git/directory data
- Tests pass
