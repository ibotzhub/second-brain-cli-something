# Changelog

All notable changes to this project will be documented in this file.

## [v0.1.0] - 2026-01-29

### Fixed
- **Critical: Fixed note duplication bug** - Notes were being duplicated in memory on every restart. Now JSON is the canonical source of truth, and the vector store is properly cleared and rebuilt on startup.
- **Critical: Fixed local embedder normalization** - Was using 1/sum instead of 1/sqrt(sum) for L2 normalization, which caused incorrect similarity scores in offline mode.
- **Fixed context search implementation** - Context command now actually uses git keywords, project name, and directory context instead of just searching for the literal word "context". Also adds 20% similarity boost for notes from the same project.
- **Fixed documentation** - Removed reference to non-existent `vectors.db` file in README.

### Added
- **`brain list` command** - List all notes sorted by most recent, with optional tag filtering and limit.
- Context type moved to brain package for proper encapsulation.

### Technical Details

#### Duplication Bug Fix
Previously, the flow was:
1. Load notes from JSON → Add to vector store
2. User adds note → Add to vector store → Save all notes from vector store to JSON
3. Restart → Load notes from JSON → Add to vector store **again**

Now:
1. Load notes from JSON → **Clear** vector store → Add to vector store
2. User adds note → Add to vector store → Save all notes from vector store to JSON
3. Restart → **Clear** vector store → Load notes from JSON → Add to vector store (no duplicates)

#### Normalization Bug Fix
Old (incorrect):
```go
norm := float32(1.0) / float32(sum)  // sum of squares, not magnitude
```

New (correct):
```go
norm := float32(1.0) / float32(math.Sqrt(float64(sum)))  // proper L2 norm
```

#### Context Search Fix
Old:
```go
return b.Search("context", 5, nil)  // literally searches for "context"
```

New:
```go
// Builds query from: project name + directory + git keywords
query := strings.Join([project, description, keywords...], " ")
// Then boosts notes from same project by 20%
```

### Known Limitations
- Vector store is in-memory only (embeddings regenerated on startup)
- Linear search O(n) - fine for thousands of notes, but could use HNSW for scale
- Local embedder is basic - works offline but not as good as OpenAI

### Next Version Priorities
- Persistent vector store (SQLite + embeddings stored)
- Better local embeddings (sentence-transformers integration)
- LLM integration for `brain ask` to synthesize answers
- Export to Markdown/Obsidian format
