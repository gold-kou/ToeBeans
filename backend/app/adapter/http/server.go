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

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/controller"
	applicationLog "github.com/gold-kou/ToeBeans/backend/app/adapter/http/log"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/middleware"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func Serve() {
	r := mux.NewRouter().StrictSlash(true)

	// middleware
	r.Use(middleware.CORSMiddleware)

	csrfMiddleware := csrf.Protect([]byte(os.Getenv("CSRF_AUTH_KEY")))
	r.Use(csrfMiddleware)

	r.Use(middleware.CurrentTimeMiddleware)

	l, err := applicationLog.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	r.Use(middleware.NewLoggingMiddleware(l).Middleware)

	// routing
	r.HandleFunc("/health/liveness", controller.HealthController)
	r.HandleFunc("/health/readiness", controller.HealthController)
	r.HandleFunc("/csrf-token", controller.CSRFTokenController)
	r.HandleFunc("/login", controller.LoginController)
	r.HandleFunc("/user", controller.UserController)
	r.HandleFunc("/user-activation/{user_name}/{activation_key}", controller.UserController)
	r.HandleFunc("/password", controller.PasswordController)
	r.HandleFunc("/password-reset-email", controller.PasswordController)
	r.HandleFunc("/password-reset", controller.PasswordController)
	r.HandleFunc("/posting", controller.PostingController)
	r.HandleFunc("/postings", controller.PostingController)
	r.HandleFunc("/posting/{posting_id}", controller.PostingController)
	r.HandleFunc("/like", controller.LikeController)
	r.HandleFunc("/like/{posting_id}", controller.LikeController)
	r.HandleFunc("/comment", controller.CommentController)
	r.HandleFunc("/comments", controller.CommentController)
	r.HandleFunc("/comment/{comment_id}", controller.CommentController)
	r.HandleFunc("/follow", controller.FollowController)
	r.HandleFunc("/follow/{followed_user_name}", controller.FollowController)

	// graceful shutdown
	server := &http.Server{Addr: fmt.Sprintf(":%v", 80), Handler: r}
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

	// launch
	log.Println("Server started!")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Panic(err)
	}
	<-idleConnsClosed
}
