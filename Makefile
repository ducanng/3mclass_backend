generate:
	swag init -o ./api --ot go,yaml -d ./
	swag fmt ./api