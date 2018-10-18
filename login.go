package main

import (
	"io"
	"os/exec"
)

func login(w io.Writer, c *config) error {
	cmd := exec.Command(
		"docker", "login",
		"-u", c.Docker.Auth.Username,
		"-p", c.Docker.Auth.Password,
		c.Docker.Auth.Address,
	)
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}
