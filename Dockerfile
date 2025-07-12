FROM docker.io/library/alpine:latest

COPY jinya-fonts /app/jinya-fonts

CMD ["/app/jinya-fonts", "serve"]
