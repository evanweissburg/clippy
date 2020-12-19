package client

import (
	"fmt"
	"github.com/evanweissburg/clippy/pkg/networking"
)

func Execute() {
	fmt.Println("Clippy client sending file!")
	networking.SendData("test.txt")
}
