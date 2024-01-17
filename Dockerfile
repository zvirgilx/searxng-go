FROM node:21-alpine3.19 as NODE_BUILD
WORKDIR /build
ENV NODE_ENV=production
COPY web/package.json .
RUN npm install --production=false
COPY web .
RUN npm run build && npm prune

FROM golang:1.21-alpine3.19 as GO_BUILD
WORKDIR /build
COPY --from=NODE_BUILD /build/.output ./.output
COPY kernel .
RUN CGO_ENABLED=0 go build -v -ldflags "-s -w" -o searxng-go &&\
    mkdir /app &&\
    mv /build/.output /app/.output &&\
    mv /build/searxng-go /app/searxng-go &&\
    find /app -name .git | xargs rm -rf

FROM node:21-alpine3.19
LABEL maintainer="zvirgilx<seacheasy4@gmail.com>"

WORKDIR /app
COPY --from=GO_BUILD /app /app
COPY ./docker-entrypoint.sh .
RUN apk add --no-cache ca-certificates tzdata

ENV RUN_IN_CONTAINER=true

ENTRYPOINT ["/app/docker-entrypoint.sh"]