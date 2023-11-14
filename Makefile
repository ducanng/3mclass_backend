dep:
	go mod tidy
	go mod download
build:
	go build -v -o ./bin/main .
build_deploy:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/main .
run:
	./bin/main
generate:
	swag init --ot go,yaml -d ./
	swag fmt
clean:
	rm -rf ./bin
all:
	make clean
	make generate
	make dep
	make build
	make run
deploy:
	make clean
	make dep
	go get github.com/go-chi/chi
	go get github.com/go-chi/chi/middleware
	go get github.com/go-chi/jwtauth/v5
	go get github.com/go-playground/validator/v10
	go get github.com/lestrrat-go/jwx/v2/jwa
	go get github.com/natefinch/lumberjack
	go get github.com/rogpeppe/go-internal/cache
	go get github.com/rs/cors
	go get github.com/spf13/viper
	go get github.com/swaggo/http-swagger/example/go-chi/docs
	go get github.com/swaggo/http-swagger/v2
	go get go.uber.org/zap
	go get go.uber.org/zap/zapcore
	go get golang.org/x/crypto/bcrypt
	go get gorm.io/driver/mysql
	go get gorm.io/gorm
	go get gorm.io/gorm/logger
	make build_deploy
	make run