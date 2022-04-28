FROM golang:1.16-alpine

WORKDIR /app

COPY . .

RUN go mod vendor

RUN go build -o aksan-kasir

CMD ./aksan-kasir

## creating image docker
# docker build -t aksan-kasir:v1 .

## create tag or branch to push to docker hub
# docker tag aksan-kasir:v1 aksanart/aksan-kasir:v1

## login to your docker hub via terminal for push
# docker login --username aksanart

## push docker to your repository
# docker push aksanart/aksan-kasir:v1