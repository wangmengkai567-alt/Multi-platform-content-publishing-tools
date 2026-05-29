package domain

import "time"

type ContentInput struct {
	Title      string   `json:"title"`
	Summary    string   `json:"summary"`
	Body       string   `json:"body"`
	Tags       []string `json:"tags"`
	CoverImage string   `json:"coverImage"`
	Tone       string   `json:"tone"`
}

type PlatformDescriptor struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	StyleHints  []string `json:"styleHints"`
	Supports    []string `json:"supports"`
}

type PlatformPreview struct {
	PlatformID     string   `json:"platformId"`
	PlatformName   string   `json:"platformName"`
	Title          string   `json:"title"`
	FormattedBody  string   `json:"formattedBody"`
	RecommendedTags []string `json:"recommendedTags"`
	Notes          []string `json:"notes"`
	Warnings       []string `json:"warnings"`
}

type PreviewRequest struct {
	Content   ContentInput `json:"content"`
	Platforms []string     `json:"platforms"`
}

type PreviewResponse struct {
	GeneratedAt time.Time         `json:"generatedAt"`
	Previews    []PlatformPreview `json:"previews"`
}

type PublishRequest struct {
	Content         ContentInput `json:"content"`
	Platforms       []string     `json:"platforms"`
	Simulate        bool         `json:"simulate"`
	ScheduledAt     string       `json:"scheduledAt"`
	EnableAnalytics bool         `json:"enableAnalytics"`
}

type PublishResult struct {
	PlatformID    string    `json:"platformId"`
	PlatformName  string    `json:"platformName"`
	Status        string    `json:"status"`
	Message       string    `json:"message"`
	PublishedAt   time.Time `json:"publishedAt"`
	ExternalRef   string    `json:"externalRef"`
	Simulation    bool      `json:"simulation"`
}

type PublishResponse struct {
	RequestID string          `json:"requestId"`
	Results   []PublishResult `json:"results"`
}

type PublishPayload struct {
	Content         ContentInput
	Preview         PlatformPreview
	Simulate        bool
	ScheduledAt     string
	EnableAnalytics bool
}
