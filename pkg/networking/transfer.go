package networking

import (
    "fmt"
)

func SendData(data string) bool {
    if data == "bad" {
        return false
    }

    fmt.Println(data)
    return true
}
