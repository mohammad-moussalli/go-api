package mascot_test

import (
	"testing"

	"github.com/mohammad-moussalli/go-api.git/mascot"
)

func TestMascot(t *testing.T) {
	if mascot.BestMascot() != "Go Gopher" {
		t.Fatal("Wrong Mascot")
	}
}
