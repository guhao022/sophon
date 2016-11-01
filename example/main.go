package main

import "sophon"

func main() {
	s := sophon.New()

	w := s.Start()

	w.SendStringGet("http://golune.com")

	w.Close()
}
