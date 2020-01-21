package main

import (
    "encoding/binary"
    "fmt"
    "bytes"
)

func main() {

    var a = 1
    b := make([]byte, 1)
    binary.PutVarint(b, int64(a))
    reader := bytes.NewReader(b)
    i, err := binary.ReadVarint(reader)
    if err != nil {
        panic(err)
    }
    fmt.Println(i)
}
