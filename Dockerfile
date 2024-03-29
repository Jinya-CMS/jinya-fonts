FROM library/node:latest AS build-frontend
WORKDIR /app
COPY . .

WORKDIR /app/angular/frontend

RUN npm install
RUN npm run build

FROM library/golang:1.22-alpine AS build
WORKDIR /app
COPY . .

COPY --from=build-frontend /app/angular/frontend/dist /app/angular/frontend/dist

RUN go build -o /app/jinya-fonts

FROM library/alpine:latest

COPY --from=build /app/jinya-fonts /app/jinya-fonts

CMD ["/app/jinya-fonts", "serve"]
