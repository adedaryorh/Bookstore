package routes

import (
	"net/http"

	"github.com/adedaryorh/bookstore-app/pkg/controllers"
	"github.com/julienschmidt/httprouter"
)

func RegisterRoutes(r *httprouter.Router) {
	r.GET("/book", controllers.GetBooks)
	r.POST("/book", controllers.CreateBook)
	r.GET("/book/:bookId", controllers.GetBookByID)
	r.GET("/books", controllers.GetAllBooks)
	r.PUT("/book/:bookId", controllers.UpdateBook)
	r.DELETE("/book/:bookId", controllers.DeleteBook)
	r.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}
