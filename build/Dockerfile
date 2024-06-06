FROM golang:1.22.1-alpine3.19 as builder

WORKDIR /workspace
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY app-config.yaml ./
COPY cmd/   cmd/
COPY internal/ internal/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o http_serv ./cmd/httpserv/main.go

FROM scratch
WORKDIR /
COPY --from=builder /workspace/http_serv .
COPY --from=builder /workspace/app-config.yaml .

ENTRYPOINT ["/http_serv"]
