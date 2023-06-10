package aiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"svalka-service/internal/config"
	logicentities "svalka-service/internal/logic/entities"
	logicinterfaces "svalka-service/internal/logic/interfaces"
	"time"
)

type aiClient struct {
	config config.AiClientConfig
}

func NewAiClient(_ context.Context) (logicinterfaces.AiClient, error) {
	config, err := config.NewAiClientConfig()
	if err != nil {
		return nil, err
	}
	return aiClient{
		config: config,
	}, nil
}

// AnalyzeLandfill ...
func (c aiClient) AnalyzeLandfill(cdnResult logicentities.CdnResult) (*logicentities.AiResult, error) {
	requestBody := struct {
		Paths []string `json:"paths"`
	}{
		Paths: cdnResult.RawImagesPaths,
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при преобразовании в JSON:%s", err.Error())
	}

	req, err := http.NewRequest("POST", c.config.AiAddr(), bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("Ошибка при создании запроса:%s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 20 * time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при выполнении запроса:%s", err.Error())
	}
	defer resp.Body.Close()
	var response struct {
		ImagePaths     []string `json:"image_paths"`
		FullImagePaths []string `json:"full_image_paths"`
		Violations     []uint64 `json:"violations"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении ответа:%s", err.Error())
	}

	violations := []logicentities.Violation{}
	for _, v := range response.Violations {
		violations = append(violations, logicentities.Violation{
			ID:     v,
			Status: true,
		})
	}
	return &logicentities.AiResult{
		AiImagesPaths: response.ImagePaths,
		Violations:    violations,
	}, nil
}
