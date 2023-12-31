package app

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/rogpeppe/go-internal/cache"
	"github.com/rs/cors"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/ducanng/3mclass_backend/config"
	"github.com/ducanng/3mclass_backend/handler"
	"github.com/ducanng/3mclass_backend/handler/authhandler"
	"github.com/ducanng/3mclass_backend/handler/userhandler"
	"github.com/ducanng/3mclass_backend/helper"
	"github.com/ducanng/3mclass_backend/internal/model/assignmentdm"
	"github.com/ducanng/3mclass_backend/internal/model/classdm"
	"github.com/ducanng/3mclass_backend/internal/model/gradedm"
	"github.com/ducanng/3mclass_backend/internal/model/userdm"
	"github.com/ducanng/3mclass_backend/internal/repository"
	"github.com/ducanng/3mclass_backend/internal/service/userservice"
	"github.com/ducanng/3mclass_backend/pkg/logutil"
)

type App struct {
	DB     *gorm.DB
	Cache  cache.Cache
	Router *chi.Mux
	cfg    *config.Config

	authHandler handler.Handler
	userHandler handler.Handler
}

func NewApp(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) InitializeApp() {
	logger := logutil.GetLogger()
	logger.Infof("Initializing DB conection...")
	a.InitializeDBConn()
	a.Router = chi.NewRouter()
	a.Router.Use(middleware.RequestID)
	a.Router.Use(middleware.RealIP)
	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.Recoverer)
	a.InitializeRoutes()
}

func (a *App) InitializeDBConn() {
	logger := logutil.GetLogger()
	db, err := gorm.Open(postgres.Open(a.cfg.DSN), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Error),
	})
	if err != nil {
		logger.Fatal("Not able to start Application.")
	}
	if a.cfg.GormDebug {
		db = db.Debug()
	}
	a.DB = db
	if a.cfg.DBAutoMigrate {
		a.MigrateDB()
	}
}

func (a *App) MigrateDB() {
	a.DB.AutoMigrate(
		userdm.User{},
		userdm.Student{},
		classdm.Class{},
		classdm.ClassMember{},
		classdm.ClassInvitation{},
		gradedm.Grade{},
		gradedm.GradeStructure{},
		assignmentdm.Assignment{},
	)
}

func (a *App) Run(ctx context.Context) (err error) {
	logger := logutil.GetLogger()
	logger.Info("Starting server request...")
	srv := a.InitializeServer()
	<-ctx.Done()
	logger.Info("Shutting down server")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	go func() {
		if err := srv.Shutdown(ctxShutDown); err != nil {
			logger.Errorf("server shutdown failed: %v", err.Error())
		}
	}()
	var wgShutDown sync.WaitGroup
	wgShutDown.Add(1)
	go func() {
		defer wgShutDown.Done()
		if err := srv.Shutdown(ctxShutDown); err != nil {
			logger.Errorf("server shutdown failed: %v", err.Error())
		}
	}()
	wgShutDown.Wait()
	logger.Info("server exited properly")

	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}
	return
}

func (a *App) InitializeServer() *http.Server {
	var (
		logger              = logutil.GetLogger()
		addr                = fmt.Sprintf(":%s", a.cfg.HTTPPort)
		router http.Handler = a.Router
	)

	if a.cfg.CORS.Enabled {
		opts := cors.Options{
			AllowedOrigins:   a.cfg.CORS.Origins,
			AllowedMethods:   a.cfg.CORS.AllowedMethods,
			Debug:            a.cfg.CORS.Debug,
			AllowCredentials: a.cfg.CORS.AllowCredentials,
		}

		if len(a.cfg.CORS.ExposedHeaders) != 0 {
			opts.ExposedHeaders = a.cfg.CORS.ExposedHeaders
		}

		if len(a.cfg.CORS.AllowedHeaders) != 0 {
			opts.AllowedHeaders = a.cfg.CORS.AllowedHeaders
		}
		c := cors.New(opts)
		router = c.Handler(a.Router)
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	go func() {
		logger.Info("start listening on port: ", a.cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("listen:%s", err)
		}
	}()

	return srv
}

func (a *App) InitializeRoutes() {
	publicKey, privateKey := getJWTKey(a.cfg.JWT.PublicKey, a.cfg.JWT.PrivateKey)
	jwt := jwtauth.New(string(jwa.RS256), privateKey, publicKey)
	jwtHelper := helper.NewJWTHelper(*jwt, time.Duration(a.cfg.JWT.ExpiryTime))
	baseHost := a.cfg.BaseHost
	a.Router.Get("/", handler.IndexHandler)
	a.Router.Get("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("./docs/swagger.yaml")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/yaml")
		w.Write(content)
	})

	a.Router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(baseHost+"/swagger.yaml"),
	))
	userRepo := repository.NewUserRepository(a.DB)
	userService := userservice.NewUserService(a.cfg, userRepo)
	a.authHandler = authhandler.NewAuthHandler(a.cfg, userService, jwtHelper)
	a.userHandler = userhandler.NewUserHandler(a.cfg, userService, jwtHelper)

	a.Router.Route("/v1/public/auth", func(r chi.Router) {
		a.authHandler.Register(r)
	})
	a.Router.Group(func(r chi.Router) {
		r.Use(jwtHelper.Verifier())
		r.Use(jwtauth.Authenticator(jwt))
		r.Route("/v1/public/u/user", func(r chi.Router) {
			a.userHandler.Register(r)
		})

	})
}

func getJWTKey(publicKeyStr, privateKeyStr string) (publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) {
	logger := logutil.GetLogger()
	derPub, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		logger.Fatalf("Could not load public key string, err: %s", err.Error())
	}

	derPrivate, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		logger.Fatalf("Could not load private key string, err: %s", err.Error())
	}

	privateKeyInf, err := x509.ParsePKCS8PrivateKey(derPrivate)
	if err != nil {
		logger.Fatalf("Could not parse private key, err: %s", err.Error())
	}

	privateKey = privateKeyInf.(*rsa.PrivateKey)
	genericPublicKeyInf, err := x509.ParsePKIXPublicKey(derPub)
	if err != nil {
		logger.Fatalf("Could not parse public key, err: %s", err.Error())
	}
	publicKey = genericPublicKeyInf.(*rsa.PublicKey)
	return publicKey, privateKey
}
