package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"content-publisher/backend/internal/adapter"
	"content-publisher/backend/internal/domain"
)

type ContentService struct {
	registry *adapter.Registry
}

func NewContentService(registry *adapter.Registry) *ContentService {
	return &ContentService{registry: registry}
}

func (s *ContentService) ListPlatforms() []domain.PlatformDescriptor {
	return s.registry.List()
}

func (s *ContentService) GeneratePreviews(request domain.PreviewRequest) (domain.PreviewResponse, error) {
	if err := validateContent(request.Content); err != nil {
		return domain.PreviewResponse{}, err
	}

	adapters, err := s.resolveAdapters(request.Platforms)
	if err != nil {
		return domain.PreviewResponse{}, err
	}

	previews := make([]domain.PlatformPreview, 0, len(adapters))
	for _, item := range adapters {
		previews = append(previews, item.BuildPreview(request.Content))
	}

	return domain.PreviewResponse{
		GeneratedAt: time.Now(),
		Previews:    previews,
	}, nil
}

func (s *ContentService) Publish(request domain.PublishRequest) (domain.PublishResponse, error) {
	if err := validateContent(request.Content); err != nil {
		return domain.PublishResponse{}, err
	}

	adapters, err := s.resolveAdapters(request.Platforms)
	if err != nil {
		return domain.PublishResponse{}, err
	}

	results := make([]domain.PublishResult, 0, len(adapters))
	for _, item := range adapters {
		preview := item.BuildPreview(request.Content)
		results = append(results, item.Publish(domain.PublishPayload{
			Content:         request.Content,
			Preview:         preview,
			Simulate:        request.Simulate,
			ScheduledAt:     request.ScheduledAt,
			EnableAnalytics: request.EnableAnalytics,
		}))
	}

	return domain.PublishResponse{
		RequestID: newRequestID(),
		Results:   results,
	}, nil
}

func (s *ContentService) resolveAdapters(platforms []string) ([]adapter.PlatformAdapter, error) {
	targets := platforms
	if len(targets) == 0 {
		for _, descriptor := range s.registry.List() {
			targets = append(targets, descriptor.ID)
		}
	}

	adapters := make([]adapter.PlatformAdapter, 0, len(targets))
	for _, id := range targets {
		item, ok := s.registry.Get(id)
		if !ok {
			return nil, errors.New("unsupported platform: " + id)
		}
		adapters = append(adapters, item)
	}

	return adapters, nil
}

func validateContent(content domain.ContentInput) error {
	if strings.TrimSpace(content.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(content.Body) == "" {
		return errors.New("body is required")
	}
	return nil
}

func newRequestID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "req-fallback"
	}
	return "req-" + hex.EncodeToString(buf)
}
