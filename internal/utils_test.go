package cromp

import "testing"

func TestParseHeader(t *testing.T) {
	const orgHeader = `
* Thing with many heads
`
	h, err := ParseHeader(orgHeader)
	if err != nil {
		t.Error("parsing", err)
	}

	if h.Title != "Thing with many heads" {
		t.Errorf("expected 'Thing with many heads'; got: %q\n", h.Title)
	}
}
