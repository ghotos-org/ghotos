# Build client
# -----------------
FROM node:14.17-alpine3.13 as client-env
#RUN apk update && apk add --no-cache 
# add `/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH

# set working directory
WORKDIR /app_client

COPY ./client /app_client

RUN yarn global add @vue/cli@4.5.13 
RUN yarn install
RUN yarn build --mode prod

# Build app environment
# -----------------
FROM golang:1.16-alpine as build-env
WORKDIR /app_ghotos
RUN apk update && apk add --no-cache gcc musl-dev vips-dev 

#COPY go.mod go.sum ./
#RUN go mod download

COPY . .
COPY --from=client-env /app_client/dist /app_ghotos/server/app/client
#COPY --from=client-env /app_client/dist /app_ghotos/server/router/client
RUN go build  -mod vendor  -ldflags '-w -s' -a -o ./bin/app 
#RUN go build  -mod vendor  -a -o ./bin/app 
# Prod environment
# ----------------------
FROM alpine
RUN apk update && apk add --no-cache bash su-exec vips-dev

COPY --from=build-env /app_ghotos/bin/app /app_ghotos/

COPY ./docker/app/bin/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

EXPOSE 8080

CMD /app_ghotos/app
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

