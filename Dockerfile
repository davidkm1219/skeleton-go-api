ARG REPO_NAME=skeleton-go-api
ARG REPO_PATH=/go/src/github.com/$REPO_NAME

FROM golang:1.22.1-alpine3.18 AS build
ARG REPO_PATH
ARG REPO_NAME

RUN apk add --no-cache bash~=5 make~=4 gcc libc-dev tzdata git

WORKDIR $REPO_PATH

COPY . ./
RUN go mod download
RUN make test
RUN make build

FROM alpine:3.18 AS release
ARG REPO_NAME
ARG REPO_PATH
RUN apk add --no-cache tzdata

COPY --from=build $REPO_PATH/$REPO_NAME /go/bin/$REPO_NAME
ENTRYPOINT ["/go/bin/$REPO_NAME"]
