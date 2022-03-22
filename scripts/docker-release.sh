#!/bin/bash
if [[ -z $1 ]];
then
    echo "${0##*/}:  No username passed"
    exit 1
else
    username=$1
fi
if [[ -z $2 ]];
then
    echo "${0##*/}:  No password passed"
    exit 2
else
    password=$2
fi



docker login -u "$username" -p "$password"

BINARY_NAME=minitwit-go-dev
DOCKER_REGISTRY=groupddevops/
VERSION=$(git rev-parse --short HEAD)

echo "tagging:" $BINARY_NAME $DOCKER_REGISTRY$BINARY_NAME:$VERSION
echo "pushing " $DOCKER_REGISTRY$BINARY_NAME:$VERSION "to dockerhub"


docker tag $BINARY_NAME $DOCKER_REGISTRY$BINARY_NAME:$VERSION
# Push the docker images
docker push $DOCKER_REGISTRY$BINARY_NAME:$VERSION


if [[ -z $3 ]];
then
    echo "Done"
else
    echo "production release choosen, also tagging with latest"
    docker tag $BINARY_NAME $DOCKER_REGISTRY$BINARY_NAME:latest
    docker push $DOCKER_REGISTRY$BINARY_NAME:latest
    echo "Done"
fi