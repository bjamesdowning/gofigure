package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bjamesdowning/web2mongo/controllers"
	"github.com/julienschmidt/httprouter"
)

//basic echo web server. Allows environment variable PORT to dicate listening port.
//user 'export PORT=<port> to set
func main() {
	mux := httprouter.New()
	mux.GET("/", index)
	mux.GET("/api/echo", controllers.Echo)
	mux.POST("/api/books", controllers.BooksHandler)
	mux.POST("/api/books/", controllers.BookHandler)
	http.ListenAndServe(port(), mux)
}

//dynamic listening port
func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = ":8080"
	}
	return ":" + port
}

//responds with http code 200 and message
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Cloud native go course on %s.", os.Getenv("PORT"))
}
