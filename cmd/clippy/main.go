package main

import (
    "fmt"
    "github.com/evanweissburg/clippy/pkg/networking"
)

func main() {
    fmt.Println("clippy client")
    networking.SendData("test data")
}
