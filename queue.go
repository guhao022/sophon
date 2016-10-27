package sophon

import "sync"

type Queue interface {

}



type Queue struct {

	// signal channels
	closed, cancelled, done chan struct{}

	wg sync.WaitGroup
}
