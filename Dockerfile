FROM library/golang:1.22-alpine AS build
WORKDIR /app
COPY . .

RUN go build -o /app/jinya-fonts

FROM library/alpine:latest

COPY --from=build /app/jinya-fonts /app/jinya-fonts

CMD ["/app/jinya-fonts", "serve"]
