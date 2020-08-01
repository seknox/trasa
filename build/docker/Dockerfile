FROM golang:1.14 as gobuilder

WORKDIR /go/src/seknox/trasa

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN go build


FROM node:13.12.0-alpine as dashbuilder

WORKDIR /trasa
ENV PATH /trasa/node_modules/.bin:$PATH

# install app dependencies
COPY dashboard/package.json ./
COPY dashboard/package-lock.json ./
RUN yarn install --silent

COPY dashboard ./

RUN ls -hal

RUN yarn run build

FROM ubuntu:latest

WORKDIR /trasa

COPY --from=gobuilder /go/src/seknox/trasa/trasa .
COPY --from=dashbuilder /trasa/build /etc/trasa/build
COPY build/etc/trasa /etc/trasa

RUN cd /etc/trasa && ls -hal
ENTRYPOINT ["/trasa/trasa"]
