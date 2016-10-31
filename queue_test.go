package sophon

import (
	"testing"
	"fmt"
)

func TestNewQueue(t *testing.T) {
	q := NewQueue(func() {
		fmt.Println("Hello Queue!")
	})

	err := q.Close()

	if err != nil {
		fmt.Println(err.Error())
	}
}
