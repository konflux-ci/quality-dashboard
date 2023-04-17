FROM node:16-alpine AS builder

# Add a work directory
WORKDIR /app

COPY package.json .
COPY yarn.lock .
COPY entrypoint.sh .

RUN yarn install --network-timeout 1000000

COPY . .

RUN chmod +x entrypoint.sh
RUN yarn build

RUN chgrp -R 0 /app && \
    chmod -R g=u /app

ENTRYPOINT [ "/app/entrypoint.sh" ]
CMD [ "sh" ]
