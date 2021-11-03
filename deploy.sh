#!/bin/bash

REGISTRY=""
REGISTRY_PATH="ghotos/ghotos"

VERSION=$(head -n 1 VERSION)
VERSION=$(./version.sh $VERSION bug)


IMAGE=${REGISTRY}${REGISTRY_PATH}

docker login  \
&& docker build -f docker/app/prod.Dockerfile  -t ${IMAGE}:${VERSION} .  \
&& docker tag ${IMAGE}:${VERSION} ${IMAGE}:latest  \
&& docker push ${IMAGE}:${VERSION}  \
&& docker push ${IMAGE}:latest  \
&& echo $VERSION > VERSION \
&& git add . \
&& git commit -m "Deployed: ${VERSION}" \
&& git push 