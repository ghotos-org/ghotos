FROM node:14.17-alpine3.13
RUN apk update && apk add --no-cache bash
# set working directory
WORKDIR /app

# add `/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH

RUN yarn global add @vue/cli@4.5.13 
CMD bash -c "yarn install && yarn serve"
