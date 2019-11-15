package zaplog

import "testing"

func TestParseKeyFromURL(t *testing.T) {
	tests := []struct {
		url string
		key string
	}{
		{
			url: "file://a.log?maxsize=5g",
			key: "file://a.log",
		},
		{
			url: "a.log?maxsize=5g",
			key: "file://a.log",
		},
		{
			url: "file://localhost/a.log?maxsize=5g",
			key: "file://localhost/a.log",
		},
	}
	for i, tt := range tests {
		key, err := parseKeyFromURL(tt.url)
		if err != nil {
			t.Errorf("%d: parse key from %q url: %v", i, tt.url, err)
			continue
		}
		if got, want := key, tt.key; got != want {
			t.Errorf("%d: key: got %q, want %q", i, got, want)
			continue
		}
		t.Logf("%d: key: got %q", i, key)
	}
}
