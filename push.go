package main

import (
	"io"
	"os/exec"
)

func push(w io.Writer, c *config) error {
	cmd := exec.Command(
		"docker",
		"push",
		c.Docker.ImageAlias,
	)
	cmd.Stdout = w
	cmd.Stderr = w
	err := cmd.Run()
	if err != nil {
		return err
	}

	// TODO: re-tag the image using c.Docker.Image
	return nil
}
