FROM golang:1.21.4 AS builder

LABEL maintainer="Nguyen Duc An <ducan172002@gmail.com>"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/main .

FROM alpine:3.18
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/main .
EXPOSE 8080
CMD ["./main"]