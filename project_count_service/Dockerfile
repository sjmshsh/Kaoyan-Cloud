FROM golang:alpine

RUN mkdir "/app"

WORKDIR "/app"

COPY kaoyanyun_auth /app/app_count
COPY config/app.yaml /

ENTRYPOINT ["./app_count"]
