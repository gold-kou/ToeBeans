package middleware

import (
	"net/http"

	"github.com/gold-kou/ToeBeans/backend/app"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, X-CSRF-Token")
		// 異なるオリジンへのリクエストでもCookieを許可する
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if app.IsProduction() {
			w.Header().Set("Access-Control-Allow-Origin", "https://toebeans.ml")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		}
		// アカウントのアクティベーションはユーザのメーラからリクエストされるので全許可にする
		// TODO いらないかも
		// if strings.HasPrefix(r.URL.Path, "/user-activation/") && r.Method == http.MethodGet {
		// 	w.Header().Set("Access-Control-Allow-Origin", "*")
		// }
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
		return
	})
}
