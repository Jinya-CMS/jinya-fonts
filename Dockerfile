FROM library/alpine:latest

COPY /builds/jinya-cms/jinya-fonts/jinya-fonts /app/jinya-fonts

CMD ["/app/jinya-fonts", "serve"]
