FROM alpine:latest
COPY . /app
WORKDIR /app
CMD ["echo", "Hello Docker! v1.0.4"]
