FROM golang:alpine as builder

RUN apk --no-cache add build-base git bzr mercurial gcc
ADD . /src
RUN cd /src && go build -o cowsay

FROM alpine
WORKDIR /app

COPY --from=builder /src/cowsay /app/
COPY ./conf.yaml /app/

EXPOSE 8080

ENTRYPOINT ./cowsay