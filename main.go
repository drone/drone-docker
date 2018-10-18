package main

import (
	"flag"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
)

// TODO consider rewriting the Dockerfile so that we can
// detect references to images created by drone and rewrite
// to their aliased names.

// TODO decide how we should enable automatic tagging,
// including multiple tags using semver and platform-specific
// prefixes.

type config struct {
	Build struct {
		Number int    `envconfig:"DRONE_BUILD_NUMBER"`
		Event  string `envconfig:"DRONE_BUILD_EVENT"`
	}

	Commit struct {
		SHA    string `envconfig:"DRONE_COMMIT_SHA"`
		REF    string `envconfig:"DRONE_COMMIT_REF"`
		Branch string `envconfig:"DRONE_COMMIT_BRANCH"`
		Author string `envconfig:"DRONE_COMMIT_AUTHOR"`
	}

	Project struct {
		Name      string `envconfig:"DRONE_REPO_NAME"`
		Namespace string `envconfig:"DRONE_REPO_NAMESPACE"`
		Link      string `envconfig:"DRONE_REPO_LINK"`
		Source    string `envconfig:"DRONE_GIT_HTTP_URL"`
	}

	Docker struct {
		Args      map[string]string `envconfig:"DOCKER_BUILD_ARGS"`
		CacheFrom []string          `envconfig:"DOCKER_BUILD_CACHE_FROM"`
		Context   string            `envconfig:"DOCKER_BUILD_CONTEXT" default:"."`
		File      string            `envconfig:"DOCKER_BUILD_DOCKERFILE" default:"Dockerfile"`
		Image     string            `envconfig:"DOCKER_BUILD_IMAGE"`
		Labels    map[string]string `envconfig:"DOCKER_BUILD_LABELS"`

		// Temporary image name. It uses a unique identifier
		// as the image tag to prevent overwriting existing
		// images on the host machine
		ImageAlias string `envconfig:"DOCKER_BUILD_IMAGE_ALIAS"`

		Auth struct {
			Address  string `envconfig:"DOCKER_ADDRESS"`
			Username string `envconfig:"DOCKER_USERNAME"`
			Password string `envconfig:"DOCKER_PASSWORD"`
		}
	}
}

var (
	dockerBuild bool
	dockerLogin bool
	dockerPush  bool
)

func main() {
	flag.BoolVar(&dockerBuild, "build", false, "docker build")
	flag.BoolVar(&dockerLogin, "login", false, "docker login")
	flag.BoolVar(&dockerPush, "push", false, "docker push")
	flag.Parse()

	c := new(config)
	err := envconfig.Process("", c)
	if err != nil {
		log.Fatal(err.Error())
	}

	if dockerBuild {
		err = build(os.Stdout, c)
		if err != nil {
			os.Exit(1)
		}
		return
	}

	if dockerLogin {
		err := login(os.Stdout, c)
		if err != nil {
			os.Exit(1)
		}
	}

	if dockerPush {
		err := push(os.Stdout, c)
		if err != nil {
			os.Exit(1)
		}
	}
}
