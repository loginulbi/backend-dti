FROM golang:latest as build

WORKDIR /usr/src/app

COPY . .
RUN go build -o main main.go

FROM debian:stable-slim
WORKDIR /usr/src/app

COPY --from=build /usr/src/app/main .
COPY --from=build /usr/src/app/.env .env

RUN apt-get update && apt-get install -y dumb-init

EXPOSE 8080

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD [ "/usr/src/app/main" ]