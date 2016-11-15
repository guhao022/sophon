package sophon

import (
	"testing"
	"sophon/storage"
)

func TestStar(t *testing.T) {
	s := NewStar()

	s.Discovery("地球", StarFixedStar)

	sto, err := storage.New("./tmp/", s.StarName()+"star", false)

	if err != nil {
		t.Error(err)
	}

	sto.Store(s)

}

func BenchmarkNewStar(b *testing.B) {

	b.ReportAllocs()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		s := NewStar()

		s.Discovery("地球", StarFixedStar)
	}
}

