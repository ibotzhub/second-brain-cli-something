package brain

import (
	"testing"
	"time"
)

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float64
		delta    float64
	}{
		{
			name:     "identical vectors",
			a:        []float32{1.0, 0.0, 0.0},
			b:        []float32{1.0, 0.0, 0.0},
			expected: 1.0,
			delta:    0.001,
		},
		{
			name:     "orthogonal vectors",
			a:        []float32{1.0, 0.0},
			b:        []float32{0.0, 1.0},
			expected: 0.0,
			delta:    0.001,
		},
		{
			name:     "opposite vectors",
			a:        []float32{1.0, 0.0},
			b:        []float32{-1.0, 0.0},
			expected: -1.0,
			delta:    0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cosineSimilarity(tt.a, tt.b)
			if diff := result - tt.expected; diff < -tt.delta || diff > tt.delta {
				t.Errorf("cosineSimilarity() = %v, want %v (±%v)", result, tt.expected, tt.delta)
			}
		})
	}
}

func TestHasAnyTag(t *testing.T) {
	tests := []struct {
		name       string
		noteTags   []string
		filterTags []string
		expected   bool
	}{
		{
			name:       "matching tag",
			noteTags:   []string{"go", "performance"},
			filterTags: []string{"go"},
			expected:   true,
		},
		{
			name:       "no matching tags",
			noteTags:   []string{"go", "performance"},
			filterTags: []string{"rust", "python"},
			expected:   false,
		},
		{
			name:       "multiple matching tags",
			noteTags:   []string{"go", "performance", "api"},
			filterTags: []string{"api", "database"},
			expected:   true,
		},
		{
			name:       "empty filter",
			noteTags:   []string{"go"},
			filterTags: []string{},
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasAnyTag(tt.noteTags, tt.filterTags)
			if result != tt.expected {
				t.Errorf("hasAnyTag() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSimpleVectorStore(t *testing.T) {
	store := &SimpleVectorStore{
		notes: make([]*Note, 0),
	}

	// Test adding a note
	note := &Note{
		ID:        "test-1",
		Content:   "Test note",
		Tags:      []string{"test"},
		Timestamp: time.Now(),
		Embedding: []float32{1.0, 0.0, 0.0},
	}

	err := store.Add(note)
	if err != nil {
		t.Fatalf("Failed to add note: %v", err)
	}

	// Test getting all notes
	notes := store.GetAllNotes()
	if len(notes) != 1 {
		t.Errorf("Expected 1 note, got %d", len(notes))
	}

	// Test search
	queryEmbedding := []float32{0.9, 0.1, 0.0} // Similar to the note
	results, err := store.Search(queryEmbedding, 10, nil)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	if results[0].Note.ID != "test-1" {
		t.Errorf("Expected note test-1, got %s", results[0].Note.ID)
	}

	// Test tag filtering
	note2 := &Note{
		ID:        "test-2",
		Content:   "Another test note",
		Tags:      []string{"other"},
		Timestamp: time.Now(),
		Embedding: []float32{0.0, 1.0, 0.0},
	}
	store.Add(note2)

	results, err = store.Search(queryEmbedding, 10, []string{"test"})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 filtered result, got %d", len(results))
	}
}

func TestLocalEmbedder(t *testing.T) {
	embedder := NewLocalEmbedder()
	
	text := "This is a test"
	embedding, err := embedder.Embed(text)
	
	if err != nil {
		t.Fatalf("Embedding failed: %v", err)
	}
	
	if len(embedding) != 384 {
		t.Errorf("Expected embedding length 384, got %d", len(embedding))
	}
	
	// Test that same text produces same embedding
	embedding2, _ := embedder.Embed(text)
	
	for i := range embedding {
		if embedding[i] != embedding2[i] {
			t.Error("Same text should produce same embedding")
			break
		}
	}
}
