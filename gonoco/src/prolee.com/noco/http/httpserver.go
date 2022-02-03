package http

import (
	"net/http"
)

func sayHello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello"))
	//fmt.Fprintf(w, "hello world, \n URL %s!", r.URL.Path[1:])
	//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func HttpServer() {
	http.HandleFunc("/hello", sayHello)
	http.ListenAndServe(":8080", nil)
}
