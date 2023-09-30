package internal

import (
    "testing"
    "github.com/arichr/eabi/pkg/eabi"
)

func BenchmarkEncoding(b *testing.B) {
    for i := 0; i < b.N; i++ {
        eabi.Marshal([]any{2, 255, 257, 259, nil})
    }
}
