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

FROM ubuntu:xenial-20200706

WORKDIR /trasa
ENV GUACENC_INSTALLED=true
RUN apt-get update
RUN apt-get install -y --no-install-recommends ca-certificates
RUN update-ca-certificates
COPY --from=gobuilder /go/src/seknox/trasa/server/server .
COPY --from=dashbuilder /trasa/build /var/trasa/dashboard
COPY --from=docker.pkg.github.com/seknox/guacamole-server/guacd:0.0.1 /usr/local/guacamole/bin/guacenc /usr/local/guacamole/bin/guacenc
COPY build/etc/trasa /etc/trasa
COPY build/docker/wait-for-it.sh .

CMD ["/trasa/wait-for-it.sh","db:5432", "--","/trasa/server"]
