package submod

import (
	"testing"
)

func TestQuack(t *testing.T) {
	str := Quack()
	if str != "Quack!" {
		t.Errorf("Quack(): %q, want %q", str, "Quack!")
	}
}
