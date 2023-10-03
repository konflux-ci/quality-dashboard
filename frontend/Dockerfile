FROM node:18-alpine AS builder

# Add a work directory
WORKDIR /app

COPY package.json .
COPY yarn.lock .
COPY entrypoint.sh .

COPY . .

RUN chmod +x entrypoint.sh

RUN chgrp -R 0 /app && \
    chmod -R g=u /app
EXPOSE 8080

ENTRYPOINT [ "/app/entrypoint.sh" ]
CMD [ "sh" ]
