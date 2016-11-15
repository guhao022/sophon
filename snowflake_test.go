package sophon

import (
	"testing"
	"fmt"
)

func TestGenerate(t *testing.T) {
	node, _ := NewNode(1023)

	for n := 0; n < 12; n++ {
		id := node.Generate()

		fmt.Println(id.Base2())
	}
}

func BenchmarkGenerate(b *testing.B) {

	node, _ := NewNode(1)

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = node.Generate()
	}
}

