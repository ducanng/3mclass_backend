package main

import (
	"context"

	"github.com/ducanng/no-name/app"
	"github.com/ducanng/no-name/config"
	"github.com/ducanng/no-name/pkg/logutil"
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
