# Brain CLI - Technical Overview

## Architecture

Brain CLI is structured as a simple but extensible CLI application with semantic search capabilities.

### Core Components

```
┌─────────────────────────────────────────────────────────┐
│                     CLI Layer (Cobra)                    │
│  add │ search │ ask │ context │ ...                     │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    Brain (Core Logic)                    │
│  • Note management                                       │
│  • Orchestration between components                      │
└─────────────────────────────────────────────────────────┘
                           │
              ┌────────────┴────────────┐
              ▼                         ▼
┌──────────────────────┐    ┌──────────────────────┐
│   Embedder           │    │   VectorStore        │
│  • OpenAI API        │    │  • In-memory storage │
│  • Local fallback    │    │  • Cosine similarity │
└──────────────────────┘    └──────────────────────┘
```

### Data Flow

1. **Adding a Note**:
   ```
   User input → Parse → Generate embedding → Store in VectorStore → Save to disk
   ```

2. **Searching**:
   ```
   Query → Generate embedding → Search VectorStore → Rank by similarity → Display results
   ```

3. **Context Detection**:
   ```
   Detect git repo/directory → Extract keywords → Generate embeddings → Search → Display relevant notes
   ```

## Implementation Details

### Embeddings

Embeddings are vector representations of text that capture semantic meaning. Similar texts have similar embeddings (measured by cosine similarity).

**OpenAI Embeddings (Primary)**:
- Model: `text-embedding-3-small` (1536 dimensions)
- Cost: ~$0.02 per 1M tokens
- Quality: High semantic understanding

**Local Embeddings (Fallback)**:
- Simple character-based hashing (384 dimensions)
- Zero cost, works offline
- Lower quality but functional

### Vector Store

The `SimpleVectorStore` is an in-memory implementation:
- Stores notes with their embeddings
- Performs linear search (fine for thousands of notes)
- Uses cosine similarity for ranking

**Future improvements**:
- HNSW index for faster search (when note count grows)
- Persistent vector database (qdrant, weaviate)
- Quantization for smaller memory footprint

### Persistence

Notes are stored in `~/.brain/notes.json`:
```json
[
  {
    "id": "uuid-here",
    "content": "Note content",
    "tags": ["tag1", "tag2"],
    "project": "project-name",
    "timestamp": "2026-01-28T12:00:00Z"
  }
]
```

Embeddings are regenerated on startup (cached in VectorStore).

## Performance Characteristics

- **Add**: O(1) for storage, O(1) for embedding generation (API call)
- **Search**: O(n) where n is total notes (linear scan)
- **Memory**: ~6KB per note (with embedding)
- **Disk**: ~500 bytes per note (without embedding)

For 1000 notes:
- Memory: ~6MB
- Search time: <10ms
- Disk space: ~500KB

## Extension Points

### Adding New Commands

Commands are defined in `cmd/` using Cobra. Each command:
1. Parses arguments and flags
2. Calls Brain methods
3. Formats and displays output

Example structure:
```go
var myCmd = &cobra.Command{
    Use:   "my-command",
    Short: "Description",
    RunE: func(cmd *cobra.Command, args []string) error {
        b, _ := brain.New()
        // Do something with b
        return nil
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

### Adding New Embedders

Implement the `Embedder` interface:
```go
type MyEmbedder struct{}

func (e *MyEmbedder) Embed(text string) ([]float32, error) {
    // Generate embedding
    return embedding, nil
}
```

Then use it in `brain.New()`.

### Adding New Vector Stores

Implement the `VectorStore` interface:
```go
type MyVectorStore struct{}

func (s *MyVectorStore) Add(note *Note) error { ... }
func (s *MyVectorStore) Search(embedding []float32, limit int, tags []string) ([]SearchResult, error) { ... }
func (s *MyVectorStore) GetAllNotes() []*Note { ... }
```

## Future Architecture Ideas

### LLM Integration

Add a `brain ask` command that uses an LLM:
```
Question → Retrieve relevant notes → Send to LLM with context → Generate answer
```

This would make Brain conversational rather than just search-based.

### Auto-tagging

When adding a note, use an LLM to suggest tags:
```go
func (b *Brain) AutoTag(note *Note) ([]string, error) {
    // Call LLM with note content
    // Return suggested tags
}
```

### Note Linking

Detect relationships between notes:
- Similar embeddings (already have this)
- Explicit mentions (parse "@note-id")
- Topic clustering

### Sync

Options for syncing between machines:
1. **Git-based**: Store notes in a git repo
2. **P2P**: Use libp2p for direct sync
3. **Cloud**: Optional cloud backup (encrypted)

## Testing Strategy

- **Unit tests**: Core logic (embeddings, similarity, filtering)
- **Integration tests**: Full commands with mock Brain
- **E2E tests**: Actual CLI invocation (optional)

Run tests:
```bash
go test ./...
```

## Build & Release

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Release
git tag v0.1.0
git push origin v0.1.0
goreleaser release
```

## Security Considerations

- API keys stored in environment variables (not in code)
- Local data is unencrypted (user's responsibility to encrypt disk)
- No telemetry or external calls except to OpenAI (opt-in)

## License

MIT - See LICENSE file
