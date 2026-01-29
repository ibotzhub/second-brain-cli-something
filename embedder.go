package brain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
)

// Embedder interface for generating embeddings
type Embedder interface {
	Embed(text string) ([]float32, error)
}

// OpenAIEmbedder uses OpenAI's API for embeddings
type OpenAIEmbedder struct {
	apiKey string
	model  string
}

func NewOpenAIEmbedder() (*OpenAIEmbedder, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}

	return &OpenAIEmbedder{
		apiKey: apiKey,
		model:  "text-embedding-3-small", // Cheaper and faster
	}, nil
}

func (e *OpenAIEmbedder) Embed(text string) ([]float32, error) {
	reqBody := map[string]interface{}{
		"input": text,
		"model": e.model,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return result.Data[0].Embedding, nil
}

// LocalEmbedder is a simple fallback that uses basic text features
// In a real implementation, you'd use sentence-transformers or similar
type LocalEmbedder struct{}

func NewLocalEmbedder() *LocalEmbedder {
	return &LocalEmbedder{}
}

// Simple bag-of-words style embedding (just for demonstration)
// In production, you'd want to use a proper model
func (e *LocalEmbedder) Embed(text string) ([]float32, error) {
	// This is a very naive implementation
	// You'd want to use a proper embedding model here
	// For now, just create a simple hash-based vector
	
	embedding := make([]float32, 384) // Standard size
	
	// Simple character-based hashing (not good, but works as fallback)
	for i, char := range text {
		idx := int(char) % len(embedding)
		embedding[idx] += float32(i+1) / float32(len(text))
	}
	
	// Normalize (L2 normalization for cosine similarity)
	var sum float32
	for _, val := range embedding {
		sum += val * val
	}
	if sum > 0 {
		norm := float32(1.0) / float32(math.Sqrt(float64(sum)))
		for i := range embedding {
			embedding[i] *= norm
		}
	}
	
	return embedding, nil
}
