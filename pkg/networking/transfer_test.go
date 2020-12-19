package networking

import (
	"testing"
)

func TestSendData(t *testing.T) {
	if !SendData("Cheese") {
		t.Errorf("SendData() did not return true")
	}

	if !SendData("bad") {
		t.Errorf("SendData() did not return true")
	}
}
