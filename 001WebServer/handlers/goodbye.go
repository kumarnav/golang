package handlers

import (
	"log"
	"net/http"
)

type Goodbye struct {
	l *log.Logger
}

func NewGoodbye(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}

func (g *Goodbye) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Bye"))
	//g.l.Println("Hello World")
	//data, err := io.ReadAll(r.Body)
	//if err != nil {
	//http.Error(rw, "Oops", http.StatusBadRequest)
	//w.WriteHeader(http.StatusBadRequest)
	//w.Write([]byte("Oops"))
	//	return
	//}
	//rw.Write([]byte(data))
	//fmt.Fprintf(rw, "Hello %s", data)
}
