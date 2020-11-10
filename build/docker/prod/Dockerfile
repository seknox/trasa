FROM golang:1.14 as gobuilder

WORKDIR /go/src/seknox/trasa

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY server server
WORKDIR /go/src/seknox/trasa/server

RUN go build


FROM node:13.12.0-alpine as dashbuilder

WORKDIR /trasa
ENV PATH /trasa/node_modules/.bin:$PATH

# install app dependencies
COPY dashboard/package.json ./
RUN yarn install --silent

COPY dashboard ./


RUN yarn run build

FROM debian:stable-slim
USER root
WORKDIR /trasa
ENV GUACENC_INSTALLED=true
COPY --from=seknox/guacd:v0.0.1 /usr/local/guacamole/bin/guacenc /usr/local/guacamole/bin/guacenc
COPY --from=seknox/guacd:v0.0.1 /usr/local/guacamole/lib/ /usr/local/guacamole/lib/
COPY --from=seknox/guacd:v0.0.1 /usr/local/guacamole/DEPENDENCIES /usr/local/guacamole/

RUN apt-get update
RUN apt-get install -y --no-install-recommends ca-certificates
RUN apt-get install -y  libavcodec-dev libavformat-dev libavutil-dev libswscale-dev ffmpeg
RUN apt-get install -y  $(cat /usr/local/guacamole/DEPENDENCIES)

RUN update-ca-certificates
COPY --from=gobuilder /go/src/seknox/trasa/server/server .
COPY --from=dashbuilder /trasa/build /var/trasa/dashboard
COPY build/etc/trasa /etc/trasa
CMD ["/trasa/server"]
