FROM library/node:latest AS build-frontend
WORKDIR /app
COPY angular/frontend .

RUN npm install
RUN npm run build

FROM library/node:latest AS build-admin-web
WORKDIR /app
COPY admin/web .

RUN npm install
RUN npm run build

FROM library/golang:1.22-alpine AS build
WORKDIR /app
COPY . .

COPY --from=build-frontend /app/dist /app/angular/frontend/dist
COPY --from=build-admin-web /app/dist /app/admin/web/dist

RUN go build -o /app/jinya-fonts

FROM library/alpine:latest

COPY --from=build /app/jinya-fonts /app/jinya-fonts

CMD ["/app/jinya-fonts", "serve"]
