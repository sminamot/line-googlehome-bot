FROM golang:1.13 as builder
ARG GOARCH=amd64

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o line-googlehome-bot

FROM alpine:3.11
WORKDIR /app
COPY --from=builder /app/line-googlehome-bot .

ENTRYPOINT ["/app/line-googlehome-bot"]
