FROM quay.imanuel.dev/dockerhub/library---golang:1.20-alpine AS build
WORKDIR /app
COPY . .

RUN go build -o /app/jinya-fonts

FROM quay.imanuel.dev/dockerhub/library---alpine:latest

COPY --from=build /app/jinya-fonts /app/jinya-fonts
COPY --from=build /app/admin/templates /app/admin/templates
COPY --from=build /app/admin/static /app/admin/static
COPY --from=build /app/webapp /app/webapp

CMD ["/app/jinya-fonts", "serve"]
