FROM node:18-alpine AS builder

# Add a work directory
WORKDIR /app

COPY package.json .
COPY yarn.lock .

RUN yarn install --network-timeout 1000000

COPY . .

RUN yarn build

COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

RUN rm -rf ./node_modules

FROM node:18-alpine
WORKDIR /app

COPY --from=builder /app .
RUN yarn global add sirv-cli

EXPOSE 8080

ENTRYPOINT [ "/app/entrypoint.sh" ]
