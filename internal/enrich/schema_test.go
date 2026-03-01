package enrich

import (
	"testing"
	"time"

	"github.com/dedene/raindrop-cli/internal/api"
)

func TestFromRaindropScaffoldDefaults(t *testing.T) {
	r := &api.Raindrop{
		ID:      123,
		Link:    "https://example.com/article",
		Title:   "Example",
		Type:    "",
		Created: time.Date(2026, 3, 1, 1, 2, 3, 0, time.UTC),
	}
	now := time.Date(2026, 3, 1, 4, 5, 6, 0, time.UTC)

	rec := FromRaindrop(r, now)

	if rec.ID != 123 {
		t.Fatalf("unexpected ID: %d", rec.ID)
	}

	if rec.ContentType != "unknown" {
		t.Fatalf("expected default content type, got %q", rec.ContentType)
	}

	if rec.Topic != "unclassified" || rec.Subtopic != "unclassified" {
		t.Fatalf("unexpected topic defaults: topic=%q subtopic=%q", rec.Topic, rec.Subtopic)
	}

	if rec.EnrichedAt != "2026-03-01T04:05:06Z" {
		t.Fatalf("unexpected enriched timestamp: %q", rec.EnrichedAt)
	}
}
