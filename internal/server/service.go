package server

type Service struct {
	Name       string
	ListenAddr string
	Algorithm  Algorithm
	Endpoints  []string
}
