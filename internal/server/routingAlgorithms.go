package server

type Algorithm int

const (
	ROUND_ROBIN Algorithm = iota
	LEAST_CONNECTION
	IP_HASH
)

