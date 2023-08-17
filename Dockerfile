FROM golang:1.19

WORKDIR /mall

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# when compiling with dynamic link function，don't rely on GLIBC
ENV CGO_ENABLED 0

RUN go build -o mall ./cmd/mall

FROM alpine:3.17.0

WORKDIR /mall
COPY --from=0 /mall/mall ./

# host address and port the edge api will listen on
EXPOSE 5577

ENTRYPOINT ["./mall","run"]