FROM nekkkkitch/docker AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux go build -a -o emailSender ./services/emailSender/cmd/app/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/emailSender .
COPY ./services/emailSender/cfg.yml .
CMD ["./emailSender"]