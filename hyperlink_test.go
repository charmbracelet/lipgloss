package lipgloss

import "testing"

func TestHyperlinkGetter(t *testing.T) {
	for i, tt := range []struct {
		style          Style
		expectedURL    string
		expectedParams map[string]string
	}{
		{
			style:          NewStyle().Hyperlink("https://charm.sh"),
			expectedURL:    "https://charm.sh",
			expectedParams: nil,
		},
		{
			style:          NewStyle().Hyperlink("https://charm.sh", "id", "IDK"),
			expectedURL:    "https://charm.sh",
			expectedParams: map[string]string{"id": "IDK"},
		},
	} {
		url, params := tt.style.GetHyperlink()
		if url != tt.expectedURL {
			t.Errorf("Test %d: expected URL %q, got %q", i, tt.expectedURL, url)
		}
		if len(params) != len(tt.expectedParams) {
			t.Errorf("Test %d: expected %d params, got %d", i, len(tt.expectedParams), len(params))
		}
		for k, v := range tt.expectedParams {
			if params[k] != v {
				t.Errorf("Test %d: expected param %q to be %q, got %q", i, k, v, params[k])
			}
		}
	}
}
