# Build environment
# -----------------
FROM golang:1.16-alpine as build-env
WORKDIR /app_ghotos
COPY . .

RUN apk update && apk add --no-cache gcc musl-dev git bash su-exec vips-dev 

RUN go get github.com/codegangsta/gin



CMD bash -c "go mod vendor && gin -i --appPort 4090 run main.go"
