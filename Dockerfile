FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod .
RUN go mod tidy
RUN go mod download -x
RUN go mod verify

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o bin/url_shortner_service cmd/*.go


FROM scratch
WORKDIR /app

COPY --from=builder /app/bin/url_shortner_service .
COPY queries.sql /app

EXPOSE 3000

ENTRYPOINT [ "/app/url_shortner_service" ]