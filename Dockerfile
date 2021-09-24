FROM golang:1.17.1-alpine AS build_image
WORKDIR /tmp/go
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 go test -v ./server
RUN go build -o ./out/go-app ./server/main.go

# separate run container => smaller image
FROM alpine:latest
COPY --from=build_image /tmp/go/out/go-app /app/go-app

#EXPOSE 8080
CMD ["/app/go-app"]
