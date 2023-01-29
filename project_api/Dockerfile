FROM golang:alpine

RUN mkdir "/app"

WORKDIR "/app"

COPY kaoyanyun_api /app/app
COPY config/app.yaml /

ENTRYPOINT ["./app"]
