FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./go-sample ./main.go


FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/go-sample .
EXPOSE 8080
ENTRYPOINT ["./go-sample"]
