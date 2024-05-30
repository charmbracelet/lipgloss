package lipgloss

import "testing"

func TestHyperlinkGetter(t *testing.T) {
	for i, tt := range []struct {
		style Style
		url   string
	}{
		{
			style: NewStyle(),
			url:   "https://charm.sh",
		},
		{
			style: NewStyle().Bold(true),
			url:   "https://charm.sh/blog/",
		},
	} {
		// Check that URL is set correctly.
		s := tt.style.Hyperlink(tt.url)
		url := s.GetHyperlink()
		if url != tt.url {
			t.Errorf("Test %d: expected URL %q, got %q", i, tt.url, url)
		}

		if len(s.hyperlink) < 3 {
			t.Errorf("Test %d: hyperlink parameters missing (we're looking for an ID)", i)
		}

		// slice to map
		params := make(map[string]string)
		for n := 1; n < len(s.hyperlink); n += 2 {
			params[s.hyperlink[n]] = s.hyperlink[n+1]
		}

		// Check that ID is set.
		id, ok := params["id"]
		if !ok {
			t.Errorf("Test %d: ID key missing in hyperlink data", i)
		}
		if id == "" {
			t.Errorf("Test %d: value missing in hyperlink data", i)
		}
	}
}
