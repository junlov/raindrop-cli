package enrich

import (
	"strings"
	"time"

	"github.com/dedene/raindrop-cli/internal/api"
)

// Summary stores structured AI summary fields for an enriched link.
type Summary struct {
	What      string   `json:"what"`
	Why       string   `json:"why"`
	Takeaways []string `json:"takeaways"`
	Audience  string   `json:"audience"`
}

// Record is the v1 enrichment schema for a single raindrop.
type Record struct {
	ID               int      `json:"id"`
	URL              string   `json:"url"`
	Title            string   `json:"title"`
	Excerpt          string   `json:"excerpt"`
	Note             string   `json:"note"`
	CreatedAt        string   `json:"created_at"`
	CanonicalURL     string   `json:"canonical_url"`
	SourceDomain     string   `json:"source_domain"`
	ContentType      string   `json:"content_type"`
	Topic            string   `json:"topic"`
	Subtopic         string   `json:"subtopic"`
	Tags             []string `json:"tags"`
	Summary          Summary  `json:"summary"`
	RevisitAfterDays int      `json:"revisit_after_days"`
	QualityScore     int      `json:"quality_score"`
	TrustScore       int      `json:"trust_score"`
	EnrichedAt       string   `json:"enriched_at"`
}

// FromRaindrop creates a scaffolded enrich record from a raindrop.
func FromRaindrop(r *api.Raindrop, now time.Time) Record {
	normalized := NormalizeURL(r.Link)

	canonicalURL := normalized.Canonical
	domain := normalized.Domain
	if canonicalURL == "" {
		canonicalURL = r.Link
	}
	if domain == "" {
		domain = strings.ToLower(strings.TrimSpace(r.Domain))
	}

	return Record{
		ID:           r.ID,
		URL:          r.Link,
		Title:        r.Title,
		Excerpt:      r.Excerpt,
		Note:         r.Note,
		CreatedAt:    r.Created.UTC().Format(time.RFC3339),
		CanonicalURL: canonicalURL,
		SourceDomain: domain,
		ContentType:  fallbackValue(r.Type, "unknown"),
		Topic:        "unclassified",
		Subtopic:     "unclassified",
		Tags:         nil,
		Summary: Summary{
			What:      "",
			Why:       "",
			Takeaways: nil,
			Audience:  "",
		},
		RevisitAfterDays: 30,
		QualityScore:     0,
		TrustScore:       0,
		EnrichedAt:       now.UTC().Format(time.RFC3339),
	}
}

func fallbackValue(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}

	return value
}
