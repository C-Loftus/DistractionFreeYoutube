FROM node:22.1.0-alpine

# set the working direction
WORKDIR /app

# install app dependencies
COPY package.json ./
COPY package-lock.json ./

# Fix permissions for packages
# RUN npm config set unsafe-perm true

RUN npm install
RUN npm install -g serve

# Bundle app source
COPY . ./

# RUN chown -R node:node /app/node_modules

EXPOSE 5173

# start app
CMD ["npm", "run", "dev", "--", "--host"]