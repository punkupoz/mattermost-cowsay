package main

func (s * server) routes() {
	s.router.HandleFunc("/cowsay", s.handleCowsay())
}