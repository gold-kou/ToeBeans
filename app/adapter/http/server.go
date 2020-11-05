package http

import (
	"log"
	"net/http"

	"github.com/gold-kou/ToeBeans/app/adapter/http/controller"

	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

func Serve() {
	r := mux.NewRouter(mux.WithServiceName("ToeBeans")).StrictSlash(true)
	r.HandleFunc("/health", controller.HealthController)
	r.HandleFunc("/user", controller.UserController)
	r.HandleFunc("/user-activation/{user_name}/{activation_key}", controller.UserActivationController)
	r.HandleFunc("/password-reset-email", controller.PasswordResetEmailController)
	r.HandleFunc("/password-reset", controller.PasswordReset)
	r.HandleFunc("/login", controller.LoginController)
	r.HandleFunc("/posting", controller.PostingController)
	r.HandleFunc("/postings", controller.PostingsController)
	r.HandleFunc("/posting/{posting_id}", controller.PostingPostingIDController)
	r.HandleFunc("/like", controller.LikeController)
	r.HandleFunc("/like/{like_id}", controller.LikeLikeIDController)
	log.Println("Server started!")
	log.Fatal(http.ListenAndServe(":8080", r))
}
