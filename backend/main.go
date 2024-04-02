package main

import "go-jwt/cmd"

func main() {
	s := cmd.NewServer()
	s.Start()
}
