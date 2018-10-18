package main

import (
	"fmt"
	"time"
)

// TODO support for org.opencontainers.image.version
// TODO support for org.opencontainers.image.licenses
// TODO support for org.opencontainers.image.ref.name
// TODO support for org.opencontainers.image.title
// TODO support for org.opencontainers.image.description

func createLabels(config *config) map[string]string {
	return map[string]string{
		"org.opencontainers.image.created":       time.Now().Format(time.RFC3339),
		"org.opencontainers.image.authors":       config.Commit.Author,
		"org.opencontainers.image.url":           config.Project.Link,
		"org.opencontainers.image.documentation": config.Project.Link,
		"org.opencontainers.image.source":        config.Project.Source,
		"org.opencontainers.image.revision":      config.Commit.SHA,
		"org.opencontainers.image.vendor":        config.Project.Namespace,
	}
}

func createDroneLabels(config *config) map[string]string {
	return map[string]string{
		"io.drone.image.created":        time.Now().Format(time.RFC3339),
		"io.drone.image.build.number":   fmt.Sprint(config.Build.Number),
		"io.drone.image.build.event":    config.Build.Event,
		"io.drone.image.repo.namespace": config.Project.Namespace,
		"io.drone.image.repo.name":      config.Project.Name,
	}
}
