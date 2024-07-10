package ego

import "testing"

func TestRawClient(t *testing.T) {
	var a = 2 + 5
	if a != 7 {
		t.Errorf("算数%d", 123)
	}
}
