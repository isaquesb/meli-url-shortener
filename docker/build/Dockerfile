FROM golang:1.23 as BUILD_IMAGE

WORKDIR /app
ENV TZ America/Sao_Paulo

COPY go.mod go.sum ./

RUN go mod download

COPY dev /app

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o shortener -ldflags '-w -s' main.go

FROM gcr.io/distroless/base

COPY --from=BUILD_IMAGE /app/shortener /

ENTRYPOINT [ "./shortener", "api" ]
