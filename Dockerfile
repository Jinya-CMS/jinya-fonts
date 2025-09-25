ARG CI_DEPENDENCY_PROXY_GROUP_IMAGE_PREFIX

FROM docker.io/library/alpine:latest

COPY jinya-fonts /app/jinya-fonts

CMD ["/app/jinya-fonts", "serve"]
