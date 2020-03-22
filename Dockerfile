FROM scratch
COPY . /app
WORKDIR /app
CMD ["/app/front"]