package helper

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func CORS(router *mux.Router) http.Handler {
	var options []handlers.CORSOption
	if os.Getenv("APP_ENV") == "development" {
		options = []handlers.CORSOption{
			handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE"}),
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedOrigins([]string{"*"}),
		}
	}
	return handlers.CORS(options...)(router)
}
