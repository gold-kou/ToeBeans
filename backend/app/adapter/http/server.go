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

const gracefulShutdownTimeoutDefault = 5

var gracefulShutdownTimeout time.Duration
var csrfAuthKey string

func init() {
	t, e := time.ParseDuration(os.Getenv("GRACEFUL_SHUTDOWN_TIMEOUT_SECOND") + "s")
	if e != nil {
		gracefulShutdownTimeout = gracefulShutdownTimeoutDefault
	} else {
		gracefulShutdownTimeout = t
	}

	csrfAuthKey = os.Getenv("CSRF_AUTH_KEY")
	if csrfAuthKey == "" {
		panic(csrfAuthKey)
	}

}

func Serve() {
	r := mux.NewRouter().StrictSlash(true)

	// middleware
	r.Use(middleware.CORSMiddleware)

	csrfMiddleware := csrf.Protect([]byte(csrfAuthKey))
	r.Use(csrfMiddleware)

	r.Use(middleware.CurrentTimeMiddleware)

	l, err := applicationLog.NewLogger()
	if err != nil {
		log.Panic(err)
	}
	r.Use(middleware.NewLoggingMiddleware(l).Middleware)

	r.Use(middleware.AuthMiddleware)

	// routing
	r.HandleFunc("/health/liveness", controller.HealthController)
	r.HandleFunc("/health/readiness", controller.HealthController)
	r.HandleFunc("/csrf-token", controller.CSRFTokenController)
	r.HandleFunc("/login", controller.LoginController)
	r.HandleFunc("/users", controller.UserController)
	r.HandleFunc("/users/{user_name}", controller.UserController)
	r.HandleFunc("/user-activation/{user_name}/{activation_key}", controller.UserController)
	r.HandleFunc("/password", controller.PasswordController)
	r.HandleFunc("/password-reset-email", controller.PasswordController)
	r.HandleFunc("/password-reset", controller.PasswordController)
	r.HandleFunc("/postings", controller.PostingController)
	r.HandleFunc("/postings/{posting_id}", controller.PostingController)
	r.HandleFunc("/likes/{posting_id}", controller.LikeController)
	r.HandleFunc("/comments/{posting_id}", controller.CommentController)
	r.HandleFunc("/comments", controller.CommentController)
	r.HandleFunc("/comments/{comment_id}", controller.CommentController)
	r.HandleFunc("/follows/{followed_user_name}", controller.FollowController)
	r.HandleFunc("/reports/users/{user_name}", controller.ReportController)
	r.HandleFunc("/reports/postings/{posting_id}", controller.ReportController)

	// graceful shutdown
	server := &http.Server{Addr: fmt.Sprintf(":%v", 80), Handler: r}
	idleConnsClosed := make(chan struct{})
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM)
		<-sigCh

		ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
		defer cancel()
		if e := server.Shutdown(ctx); e != nil {
			// Error from closing listeners, or context gracefulShutdownTimeout:
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
