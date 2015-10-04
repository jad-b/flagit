package flagit

import (
	"net/url"
	"testing"
)

func TestURLFlag(t *testing.T) {
	testURLs := []string{
		"localhost:18000",
		"https://localhost:18000",
		"/a/path",
	}
	urlFlag := &URLFlag{}
	earl := (*url.URL)(urlFlag)
	for _, tURL := range testURLs {
		err := urlFlag.Set(tURL)
		if err != nil {
			t.Error(err)
		}
		if earl.Opaque != "" {
			t.Errorf("Opaque is set; you probably didn't want that: %#v", urlFlag)
		}
	}
}
