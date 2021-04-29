package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// graceful shutdown
	server := &http.Server{Addr: fmt.Sprintf(":%v", 8080), Handler: r}
	idleConnsClosed := make(chan struct{})
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM)
		<-sigCh

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5))
		defer cancel()
		if e := server.Shutdown(ctx); e != nil {
			// Error from closing listeners, or context timeout:
			log.Panic("Failed to gracefully shutdown ", e)
		}
		close(idleConnsClosed)
	}()

	log.Println("Server started!")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Panic(err)
	}
	<-idleConnsClosed
}
