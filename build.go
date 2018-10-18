package main

import (
	"io"
	"os/exec"
)

func build(w io.Writer, c *config) error {
	var flags []string
	flags = append(flags, "build")

	for k, v := range c.Docker.Args {
		arg := k + "=" + v
		flags = append(flags, "--build-arg", arg)
	}

	for _, v := range c.Docker.CacheFrom {
		flags = append(flags, "--cache-from", v)
	}

	if v := c.Docker.File; v != "" {
		flags = append(flags, "-f", v)
	}

	for k, v := range c.Docker.Labels {
		label := k + "=" + v
		flags = append(flags, "--label", label)
	}

	for k, v := range createLabels(c) {
		label := k + "=" + v
		flags = append(flags, "--label", label)
	}

	for k, v := range createDroneLabels(c) {
		label := k + "=" + v
		flags = append(flags, "--label", label)
	}

	flags = append(flags, "-t", c.Docker.ImageAlias)
	flags = append(flags, c.Docker.Context)

	cmd := exec.Command("docker", flags...)
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}
