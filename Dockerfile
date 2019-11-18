FROM alpine:3.7
COPY . /app
WORKDIR /app
ENTRYPOINT /app/front