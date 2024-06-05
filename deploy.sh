#!/bin/bash
IMAGE_NAME="go-openapi"
IMAGE_TAG="latest"
docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .
if docker ps -a | grep -q ${IMAGE_NAME}; then
  docker stop ${IMAGE_NAME}
  docker rm ${IMAGE_NAME}
fi
docker run --name ${IMAGE_NAME} --network=scripts_lou2-net -d -e VIRTUAL_HOST=api.lou2.kr -e VIRTUAL_PORT=80 -e LETSENCRYPT_HOST=api.lou2.kr -e LETSENCRYPT_EMAIL=lee@lou2.kr ${IMAGE_NAME}:${IMAGE_TAG}
