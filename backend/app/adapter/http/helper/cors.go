package helper

import (
	"net/http"
)

//func CORS(router *mux.Router) http.Handler {
//	var options []handlers.CORSOption
//	if os.Getenv("APP_ENV") == "development" {
//		options = []handlers.CORSOption{
//			handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE"}),
//			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
//			handlers.AllowCredentials(),
//			//handlers.AllowedOrigins([]string{"http://localhost:3000"}),
//		}
//	}
//	return handlers.CORS(options...)(router)
//}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
		return
	})
}
