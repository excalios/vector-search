package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-app/domain"
	"go-app/internal/logging"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/pgvector/pgvector-go"
)

type EmbeddingHTTPRepository struct {
	c     *http.Client
	aiURL string
}

func NewEmbeddingHTTPRepository() *EmbeddingHTTPRepository {
	aiURL := os.Getenv("AI_API_URL")

	client := &http.Client{
		Timeout: 180 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 90 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
		},
	}

	return &EmbeddingHTTPRepository{
		c:     client,
		aiURL: aiURL,
	}
}

func (r *EmbeddingHTTPRepository) GetGeneralEmbedding(
	ctx context.Context,
	sentence string,
	vType domain.VectorType,
) (*pgvector.Vector, error) {
	url := fmt.Sprintf("%s/embedding/general", r.aiURL)

	data := domain.EmbeddingInput{
		Sentence: sentence,
		Type:     vType,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		logging.LogError(ctx, err, "Error marshaling json")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		logging.LogError(ctx, err, "Failed to create request")
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := r.c.Do(req)
	if err != nil {
		logging.LogError(ctx, err, "Failed to send request")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			return nil, fmt.Errorf("Error reading http response: %v", err)
		}

		errResp := string(bodyBytes)
		log.Printf("🪵CS8 errResp: %v CS8\n", errResp)

		return nil, fmt.Errorf("failed to upload image, status code: %d", resp.StatusCode)
	}

	var response domain.EmbeddingOutput
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		logging.LogError(ctx, err, "Failed to decode response body")
		return nil, err
	}

	return &response.Data, nil
}
