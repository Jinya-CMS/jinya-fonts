FROM quay.imanuel.dev/dockerhub/library---golang:1.20-alpine
WORKDIR /app
COPY . .

RUN go build -o jinya-fonts

CMD ["/app/jinya-fonts", "serve"]