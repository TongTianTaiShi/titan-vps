FROM golang:1.19

WORKDIR /basis

COPY titan-vps/go.mod titan-vps/go.sum ./

RUN go mod download

COPY titan-vps/ .

# when compiling with dynamic link function，don't rely on GLIBC
ENV CGO_ENABLED 0

RUN go build -o basis ./cmd/basis

FROM alpine:3.17.0

WORKDIR /basis
COPY --from=0 /basis/basis ./basis

# host address and port the edge api will listen on
EXPOSE 5577

ENTRYPOINT ["./basis run"]