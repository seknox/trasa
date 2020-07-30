FROM node:alpine

RUN mkdir /build

COPY build build

RUN npm -g install serve

WORKDIR .
EXPOSE 5000

CMD ["serve", "-s", "build"]
