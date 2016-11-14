package sophon

import (
	"testing"
	"fmt"
)

func TestSnowFlake(t *testing.T) {
	fmt.Println("start generate")
	iw, _ := NewIdWorker(2)
	var prevId int64 = 0
	for i := 0; i < 1000; i++ {
		id, err := iw.NextId()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(id)
		}
		if prevId >= id {
			panic("prevId >= id")
		} else {
			prevId = id
		}
	}
	fmt.Println("end generate")
}
