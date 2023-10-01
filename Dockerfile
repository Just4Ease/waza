FROM golang:1.21 AS builder
RUN mkdir -p /go/src/waza
WORKDIR /go/src/waza
COPY . .

RUN GIT_TERMINAL_PROMPT=1 \
    GOARCH=amd64 \
    GOOS=linux \
    CGO_ENABLED=0 \
    go build -v --installsuffix cgo --ldflags="-s" -o waza
FROM alpine:3.13

# convert build-arg to env variables
RUN apk add --no-cache tzdata
ENV TZ Africa/Lagos
RUN mkdir -p /svc/
COPY --from=builder /go/src/waza/waza /svc/

WORKDIR /svc/

CMD ["./waza"]
