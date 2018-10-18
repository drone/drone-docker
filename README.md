Drone utility for build and publishing images. The yaml build and push directives are translated into an action that runs this image using the yaml and runtime configuration parameters. If you are using Drone, this is not an image you would ever use directly.


Build the plugin binaries:

```
./scripts/build.sh
```

Build the Docker images:

```
docker build -t drone/docker -f docker/Dockerfile.linux.amd64 .
```

Build an image using a temporary tag:

```
docker run --rm \
  -e DRONE_BUILD_NUMBER=1 \
  -e DRONE_BUILD_EVENT=push \
  -e DRONE_COMMIT_SHA=d8dbe4d94f15fe89232e0402c6e8a0ddf21af3ab \
  -e DRONE_COMMIT_REF=refs/heads/master \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_COMMIT_AUTHOR=octocat \
  -e DOCKER_BUILD_IMAGE=octocat/hello-world:latest \
  -e DOCKER_BUILD_IMAGE_ALIAS=octocat/hello-world:dd45e66c1934 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  drone/docker --build
```

Publish and re-tag the image:

```
docker run --rm \
  -e DOCKER_ADDRESS=docker.io \
  -e DOCKER_USERNAME=octocat \
  -e DOCKER_PASSWORD=correct-horse-battery-staple \
  -e DRONE_BUILD_NUMBER=1 \
  -e DRONE_BUILD_EVENT=push \
  -e DRONE_COMMIT_SHA=d8dbe4d94f15fe89232e0402c6e8a0ddf21af3ab \
  -e DRONE_COMMIT_REF=refs/heads/master \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_COMMIT_AUTHOR=octocat \
  -e DOCKER_BUILD_IMAGE=octocat/hello-world:latest \
  -e DOCKER_BUILD_IMAGE_ALIAS=octocat/hello-world:dd45e66c1934 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  drone/docker --login --push
```