FROM golang as build

COPY . /src

WORKDIR /src

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o captcha_server



FROM alpine as app

LABEL maintainer="ivan.kuchin@gmail.com"
LABEL org.opencontainers.image.source="https://github.com/IvanKuchin/${BUILD_REPO}"

WORKDIR /workdir

COPY --from=build /src/captcha_server .
COPY --from=build /src/config ./config

EXPOSE 80

CMD ["./captcha_server"]
