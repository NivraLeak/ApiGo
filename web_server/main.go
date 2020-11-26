package main

func main() {
	server := NewServer(":3001")
	server.Handle("GET", "/", HandlerRoot)
	server.Handle("POST", "/create", PostRequest)
	server.Handle("POST", "/user", UserPostRequest)
	server.Handle("POST", "/api", server.AddMiddleware(HandleHome, CheckAuth(), Logging()))

	server.Listen()
}
