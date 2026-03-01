package api

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCreateRaindropRequestPleaseParseMarshalsAsObject(t *testing.T) {
	req := CreateRaindropRequest{
		Link: "https://example.com",
	}
	req.PleaseParse = &struct{}{}

	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	if !strings.Contains(string(raw), `"pleaseParse":{}`) {
		t.Fatalf("expected pleaseParse object, got %s", string(raw))
	}
}

func TestCreateRaindropRequestPleaseParseOmittedWhenNil(t *testing.T) {
	req := CreateRaindropRequest{
		Link: "https://example.com",
	}

	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	if strings.Contains(string(raw), `"pleaseParse"`) {
		t.Fatalf("did not expect pleaseParse when nil, got %s", string(raw))
	}
}
