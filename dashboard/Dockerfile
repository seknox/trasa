FROM node:alpine

COPY build build

RUN npm -g install serve

EXPOSE 5000

CMD ["serve", "-s", "build"]
