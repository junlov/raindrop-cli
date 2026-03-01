package enrich

import (
	"net"
	"net/url"
	"strings"
)

// NormalizedURL stores canonicalized URL information.
type NormalizedURL struct {
	Original  string `json:"original"`
	Canonical string `json:"canonical"`
	Domain    string `json:"domain"`
}

var trackingQueryKeys = map[string]struct{}{
	"fbclid": {},
	"gclid":  {},
	"mc_cid": {},
	"mc_eid": {},
	"ref":    {},
}

// NormalizeURL normalizes a URL into a canonical form for dedupe/indexing.
func NormalizeURL(raw string) NormalizedURL {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return NormalizedURL{
			Original:  raw,
			Canonical: "",
			Domain:    "",
		}
	}

	parsed, ok := parseURLWithFallback(trimmed)
	if !ok {
		return NormalizedURL{
			Original:  raw,
			Canonical: trimmed,
			Domain:    "",
		}
	}

	parsed.Fragment = ""
	parsed.Scheme = strings.ToLower(parsed.Scheme)
	parsed.Host = normalizedHost(parsed)
	filterTrackingQuery(parsed)

	return NormalizedURL{
		Original:  raw,
		Canonical: parsed.String(),
		Domain:    parsed.Hostname(),
	}
}

func parseURLWithFallback(raw string) (*url.URL, bool) {
	parsed, err := url.Parse(raw)
	if err == nil && parsed != nil && parsed.Host != "" {
		if parsed.Scheme == "" {
			parsed.Scheme = "https"
		}

		return parsed, true
	}

	if strings.Contains(raw, "://") {
		return nil, false
	}

	parsed, err = url.Parse("https://" + raw)
	if err != nil || parsed == nil || parsed.Host == "" {
		return nil, false
	}

	return parsed, true
}

func normalizedHost(parsed *url.URL) string {
	hostname := strings.ToLower(parsed.Hostname())
	port := parsed.Port()

	if port == "" {
		return hostname
	}

	if (parsed.Scheme == "http" && port == "80") || (parsed.Scheme == "https" && port == "443") {
		return hostname
	}

	return net.JoinHostPort(hostname, port)
}

func filterTrackingQuery(parsed *url.URL) {
	if parsed == nil {
		return
	}

	query := parsed.Query()
	for key := range query {
		lowerKey := strings.ToLower(key)
		if strings.HasPrefix(lowerKey, "utm_") {
			query.Del(key)
			continue
		}

		if _, ok := trackingQueryKeys[lowerKey]; ok {
			query.Del(key)
		}
	}

	parsed.RawQuery = query.Encode()
}
