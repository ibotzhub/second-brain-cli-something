package brain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	Project   string    `json:"project,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Embedding []float32 `json:"-"` // Don't serialize, computed on demand
}

type SearchResult struct {
	Note       *Note
	Similarity float64
}

type Context struct {
	Directory   string
	Project     string
	Description string
	Keywords    []string
}

type Brain struct {
	dataDir    string
	notesPath  string
	embedder   Embedder
	vectorStore VectorStore
}

// New creates a new Brain instance
func New() (*Brain, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dataDir := filepath.Join(homeDir, ".brain")
	notesPath := filepath.Join(dataDir, "notes.json")

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	// Initialize embedder (using OpenAI by default, can be configured)
	embedder, err := NewOpenAIEmbedder()
	if err != nil {
		// Fall back to local embedder if OpenAI fails
		embedder = NewLocalEmbedder()
	}

	// Initialize vector store
	vectorStore, err := NewSimpleVectorStore(dataDir)
	if err != nil {
		return nil, err
	}

	b := &Brain{
		dataDir:     dataDir,
		notesPath:   notesPath,
		embedder:    embedder,
		vectorStore: vectorStore,
	}

	// Load existing notes into vector store
	if err := b.loadNotes(); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Brain) AddNote(note *Note) error {
	// Generate ID if not set
	if note.ID == "" {
		note.ID = uuid.New().String()
	}

	// Generate embedding
	embedding, err := b.embedder.Embed(note.Content)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}
	note.Embedding = embedding

	// Add to vector store
	if err := b.vectorStore.Add(note); err != nil {
		return err
	}

	// Save to disk
	return b.saveNotes()
}

func (b *Brain) Search(query string, limit int, tags []string) ([]SearchResult, error) {
	// Generate embedding for query
	embedding, err := b.embedder.Embed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Search vector store
	return b.vectorStore.Search(embedding, limit, tags)
}

func (b *Brain) GetContextualNotes(ctx interface{}) ([]SearchResult, error) {
	// Type assert to get our Context struct
	context, ok := ctx.(Context)
	if !ok {
		// Fallback if wrong type
		return b.Search("context", 5, nil)
	}

	// Build a query from context information
	queryParts := []string{context.Project, context.Description}
	queryParts = append(queryParts, context.Keywords...)
	
	// Filter out empty strings
	var validParts []string
	for _, part := range queryParts {
		if part != "" {
			validParts = append(validParts, part)
		}
	}
	
	if len(validParts) == 0 {
		return []SearchResult{}, nil
	}
	
	// Join all parts into a single query string
	query := strings.Join(validParts, " ")
	
	// Search with the constructed query
	results, err := b.Search(query, 5, nil)
	if err != nil {
		return nil, err
	}
	
	// If we have a project name, boost results that match it
	if context.Project != "" {
		for i := range results {
			if results[i].Note.Project == context.Project {
				results[i].Similarity *= 1.2 // 20% boost for same project
				if results[i].Similarity > 1.0 {
					results[i].Similarity = 1.0 // Cap at 1.0
				}
			}
		}
		
		// Re-sort after boosting
		sort.Slice(results, func(i, j int) bool {
			return results[i].Similarity > results[j].Similarity
		})
	}
	
	return results, nil
}

func (b *Brain) ListNotes(tags []string) ([]*Note, error) {
	allNotes := b.vectorStore.GetAllNotes()
	
	if len(tags) == 0 {
		return allNotes, nil
	}
	
	// Filter by tags (use the hasAnyTag function from vectorstore)
	var filtered []*Note
	for _, note := range allNotes {
		for _, filterTag := range tags {
			hasTag := false
			for _, noteTag := range note.Tags {
				if noteTag == filterTag {
					hasTag = true
					break
				}
			}
			if hasTag {
				filtered = append(filtered, note)
				break
			}
		}
	}
	
	return filtered, nil
}

func (b *Brain) loadNotes() error {
	// Check if notes file exists
	if _, err := os.Stat(b.notesPath); os.IsNotExist(err) {
		return nil // No notes yet, that's fine
	}

	data, err := os.ReadFile(b.notesPath)
	if err != nil {
		return err
	}

	var notes []*Note
	if err := json.Unmarshal(data, &notes); err != nil {
		return err
	}

	// Clear the vector store first to avoid duplicates
	b.vectorStore = &SimpleVectorStore{notes: make([]*Note, 0)}

	// Load each note into vector store
	for _, note := range notes {
		// Generate embedding if not present
		if len(note.Embedding) == 0 {
			embedding, err := b.embedder.Embed(note.Content)
			if err != nil {
				continue // Skip notes we can't embed
			}
			note.Embedding = embedding
		}
		
		b.vectorStore.Add(note)
	}

	return nil
}

func (b *Brain) saveNotes() error {
	notes := b.vectorStore.GetAllNotes()
	
	data, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(b.notesPath, data, 0644)
}
