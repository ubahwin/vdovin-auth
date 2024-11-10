package main

import (
	"github.com/shuryak/api-wrappers/pkg/router"
	"github.com/ubahwin/vdovin-auth/internal/api"
	authapi "github.com/ubahwin/vdovin-auth/internal/api/auth"
	"github.com/ubahwin/vdovin-auth/internal/api/user"
	"github.com/ubahwin/vdovin-auth/internal/core/auth"
	"github.com/ubahwin/vdovin-auth/internal/core/user"
	qrstorage "github.com/ubahwin/vdovin-auth/internal/storage/qr"
	sessionstorage "github.com/ubahwin/vdovin-auth/internal/storage/session"
	userstorage "github.com/ubahwin/vdovin-auth/internal/storage/user"
	"github.com/ubahwin/vdovin-auth/pkg/phasher"
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
		).SetErrHandler(api.ErrHandler),
		router.GET("/auth", authAPIGroup.Auth).SetErrHandler(api.ErrHandler),
		router.GET("/code", authAPIGroup.Code).SetErrHandler(api.ErrHandler),
	)

	srv := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: r,
	}

	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
