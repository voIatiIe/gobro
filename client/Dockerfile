FROM node:21-alpine3.18

RUN apk add --no-cache make gcc g++ python3

COPY . /client
WORKDIR /client

COPY ./package.json ./package.json

RUN yarn install

ENTRYPOINT [ "yarn", "start", "--watch" ]
