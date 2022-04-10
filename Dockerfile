# --- Section 1 Start
# FROM golang:1.16-alpine as build-dev
FROM golang:1.16 as builder

RUN mkdir /app
COPY . /app
WORKDIR /app

# build
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o main


EXPOSE 8080
# 587 for smtp email
CMD ["./main"]