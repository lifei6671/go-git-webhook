#!/bin/sh

export CI_BUILD_SHA="x${CIRCLE_WORKFLOW_ID:0:8}"
export IMAGE_TMP_NAME=$(echo "$CIRCLE_PROJECT_REPONAME" | tr '[:upper:]' '[:lower:]')

# auto login
docker() {
	if [ ! -e ~/.docker/config.json ]; then
		command docker login -u $DOCKERCLOUD_USER -p $DOCKERCLOUD_APIKEY
	fi

	command docker "$@"
}

# push to docker cloud
docker_push() {
	echo Tag image...
	docker tag $IMAGE_TMP_NAME:$CI_BUILD_SHA $1
	echo Push to docker cloud...
	docker push $1
}

# build image
docker_build() {
	echo Build image...
	docker build -t $IMAGE_TMP_NAME:$CI_BUILD_SHA .
}
