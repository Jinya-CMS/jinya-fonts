FROM quay.imanuel.dev/dockerhub/library---golang:1.17-alpine
WORKDIR /app
COPY . .

ENV DATA_DIR=/jinya-fonts-data

RUN go build -o jinya-fonts
RUN mkdir /jinya-fonts

CMD ["/app/jinya-fonts"]