package main

func main() {
	server := NewServer(":3001")
	server.Listen()
}
