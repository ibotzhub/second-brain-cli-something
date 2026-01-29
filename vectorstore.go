package brain

import (
	"math"
	"sort"
	"sync"
)

// VectorStore interface for storing and searching embeddings
type VectorStore interface {
	Add(note *Note) error
	Search(embedding []float32, limit int, tags []string) ([]SearchResult, error)
	GetAllNotes() []*Note
}

// SimpleVectorStore is an in-memory vector store
type SimpleVectorStore struct {
	mu    sync.RWMutex
	notes []*Note
}

func NewSimpleVectorStore(dataDir string) (*SimpleVectorStore, error) {
	return &SimpleVectorStore{
		notes: make([]*Note, 0),
	}, nil
}

func (s *SimpleVectorStore) Add(note *Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.notes = append(s.notes, note)
	return nil
}

func (s *SimpleVectorStore) Search(embedding []float32, limit int, tags []string) ([]SearchResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]SearchResult, 0)

	for _, note := range s.notes {
		// Filter by tags if specified
		if len(tags) > 0 && !hasAnyTag(note.Tags, tags) {
			continue
		}

		// Calculate cosine similarity
		similarity := cosineSimilarity(embedding, note.Embedding)
		
		results = append(results, SearchResult{
			Note:       note,
			Similarity: similarity,
		})
	}

	// Sort by similarity (highest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	// Limit results
	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

func (s *SimpleVectorStore) GetAllNotes() []*Note {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Return a copy to prevent external modifications
	notesCopy := make([]*Note, len(s.notes))
	copy(notesCopy, s.notes)
	return notesCopy
}

// Helper functions

func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	
	for i := range a {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

func hasAnyTag(noteTags, filterTags []string) bool {
	tagSet := make(map[string]bool)
	for _, tag := range noteTags {
		tagSet[tag] = true
	}

	for _, tag := range filterTags {
		if tagSet[tag] {
			return true
		}
	}

	return false
}
