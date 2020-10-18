FROM golang:1.15.3
WORKDIR /go/src/telebot
COPY . .
RUN CGO_ENABLED=0 GOOS=linux make compile

FROM alpine:latest  
WORKDIR /root/
COPY --from=0 /go/src/telebot/bin/bot .
CMD ["./bot"]  