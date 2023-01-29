FROM golang:alpine

RUN mkdir "/app"

WORKDIR "/app"

COPY kaoyanyun_auth /app/app_user
COPY config/app.yaml /

ENTRYPOINT ["./app_user"]
