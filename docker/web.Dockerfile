FROM node:latest

WORKDIR /web


# install app dependencies
COPY ./web/package.json ./
COPY ./web/package-lock.json ./
RUN npm install --silent

# Copy the code into the container
COPY ./web/ ./

ENTRYPOINT ["npm", "start"]