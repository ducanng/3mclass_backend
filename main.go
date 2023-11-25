package main

import (
	"context"

	"github.com/ducanng/3mclass_backend/app"
	"github.com/ducanng/3mclass_backend/config"
	"github.com/ducanng/3mclass_backend/pkg/logutil"

	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

//	@title		NoName API
//	@version	1.0
//	@description
//	@termsOfService	http://swagger.io/terms/
//	@contact.name	API Support
//	@contact.email
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//	@host			localhost:8080
//	@schemes		http https

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

// @BasePath	/

func main() {
	logutil.InitLog()
	logger := logutil.GetLogger()
	ctx, cancel := context.WithCancel(context.Background())
	cfg := config.Load()
	defer cancel()
	a := app.NewApp(cfg)
	a.InitializeApp()
	if err := a.Run(ctx); err != nil {
		logger.Fatalf("failed to serve:+%v", err)
	}
}
