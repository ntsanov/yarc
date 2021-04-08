package main

type Config struct {

	// Max Concurent connections at once
	MaxConnections int

	// Timeout in seconds
	Timeout int

	// NodeURL to query
	NodeURL string
}
