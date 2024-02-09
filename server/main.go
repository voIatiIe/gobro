package main

import gb "gobro/gobro"

func main() {
	server := gb.NewServer(
		gb.WithAddr(":8010"),
		gb.WithUrl("/ws"),
	)
	server.Start()
}
