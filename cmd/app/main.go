package main

import (
	"github.com/ubahwin/vdovin-auth/internal/api"
	authapi "github.com/ubahwin/vdovin-auth/internal/api/auth"
	"github.com/ubahwin/vdovin-auth/internal/api/user"
	"github.com/ubahwin/vdovin-auth/internal/core/auth"
	"github.com/ubahwin/vdovin-auth/internal/core/user"
	qrstorage "github.com/ubahwin/vdovin-auth/internal/storage/qr"
	sessionstorage "github.com/ubahwin/vdovin-auth/internal/storage/session"
	userstorage "github.com/ubahwin/vdovin-auth/internal/storage/user"
	"github.com/ubahwin/vdovin-auth/pkg/phasher"
	"github.com/ubahwin/vdovin-auth/pkg/router"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	accessTokenLength  = 128
	refreshTokenLength = 128
	accessTokenTTL     = time.Hour
	qrCodeLength       = 128
	qrCodeTTL          = 10 * time.Minute
)

func main() {
	logger := log.New(os.Stdout, "[LOG]", log.Lshortfile)

	userStorage := userstorage.NewMem()
	sessionStorage := sessionstorage.NewMem(accessTokenLength, refreshTokenLength, accessTokenTTL)
	userManager := user.NewManager(userStorage, sessionStorage, phasher.NewBcrypt())
	userAPIGroup := userapi.New(logger, userManager)

	qrStorage := qrstorage.NewMem(qrCodeLength, qrCodeTTL)
	authorizer := auth.NewAuthorizer(sessionStorage, qrStorage)
	authAPIGroup := authapi.New(logger, authorizer)

	r := router.New(logger)
	r.Add(
		router.NewGroup("/user",
			router.POST("/signUp", userAPIGroup.SignUp),
			router.POST("/signIn", userAPIGroup.SignIn),
			router.POST("/info", userAPIGroup.UserInfo),
		).SetPreHandler(api.CORS).SetErrHandler(api.ErrHandler),
		router.POST("/auth", authAPIGroup.Auth).SetPreHandler(api.CORS).SetErrHandler(api.ErrHandler),
		router.POST("/code", authAPIGroup.Code).SetPreHandler(api.CORS).SetErrHandler(api.ErrHandler),
	)

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}

	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
