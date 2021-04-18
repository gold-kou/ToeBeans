package http

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/csrf"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/controller"

	"github.com/gorilla/mux"
)

func Serve() {
	r := mux.NewRouter().StrictSlash(true)

	r.Use(helper.CORSMiddleware)
	csrfMiddleware := csrf.Protect([]byte(os.Getenv("CSRF_AUTH_KEY")))
	r.Use(csrfMiddleware)

	r.HandleFunc("/health", controller.HealthController)
	r.HandleFunc("/csrf-token", controller.CSRFTokenController)
	r.HandleFunc("/user", controller.UserController)
	r.HandleFunc("/user-activation/{user_name}/{activation_key}", controller.UserActivationController)
	r.HandleFunc("/password", controller.PasswordController)
	r.HandleFunc("/password-reset-email", controller.PasswordResetEmailController)
	r.HandleFunc("/password-reset", controller.PasswordResetController)
	r.HandleFunc("/login", controller.LoginController)
	r.HandleFunc("/posting", controller.PostingController)
	r.HandleFunc("/postings", controller.PostingsController)
	r.HandleFunc("/posting/{posting_id}", controller.PostingPostingIDController)
	r.HandleFunc("/like", controller.LikeController)
	r.HandleFunc("/like/{posting_id}", controller.LikePostingIDController)
	r.HandleFunc("/comment", controller.CommentController)
	r.HandleFunc("/comments", controller.CommentsController)
	r.HandleFunc("/comment/{comment_id}", controller.CommentCommentIDController)
	r.HandleFunc("/follow", controller.FollowController)
	r.HandleFunc("/follow/{followed_user_name}", controller.FollowUserNameController)

	log.Println("Server started!")
	log.Fatal(http.ListenAndServe(":8080", r))
}
