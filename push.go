package main

import (
	"context"
	"io"
	"os/exec"
	"time"

	"github.com/genuinetools/reg/registry"
	"github.com/genuinetools/reg/repoutils"
)

func push(ctx context.Context, w io.Writer, c *config) error {
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

	image, err := registry.ParseImage(c.Docker.Image)
	if err != nil {
		return err
	}

	alias, err := registry.ParseImage(c.Docker.ImageAlias)
	if err != nil {
		return err
	}

	r, err := createClient(ctx,
		c.Docker.Auth.Username,
		c.Docker.Auth.Password,
		alias.Domain,
	)
	if err != nil {
		return err
	}

	// STEP 1: get the manifest for the temporary (aliased)
	// image tag.
	manifest, err := r.Manifest(ctx, alias.Path, alias.Reference())
	if err != nil {
		return err
	}
	// STEP 2: upload the manifest with the user-defined
	// image tag.
	err = r.PutManifest(ctx, image.Path, image.Reference(), manifest)
	if err != nil {
		return err
	}

	// STEP 4: get the digest for the temporary (aliased)
	// image tag.
	digest, err := r.Digest(ctx, alias)
	if err != nil {
		return err
	}
	if err := alias.WithDigest(digest); err != nil {
		return err
	}
	// STEP 5: delete the digest for the temporary (aliased)
	// image tag.
	return r.Delete(ctx, alias.Path, digest)
}

func createClient(ctx context.Context, username, password, domain string) (*registry.Registry, error) {
	auth, err := repoutils.GetAuthConfig(username, password, domain)
	if err != nil {
		return nil, err
	}

	// Create the registry client.
	return registry.New(ctx, auth, registry.Opt{
		Insecure: false,
		Debug:    false,
		SkipPing: false,
		Timeout:  time.Minute,
	})
}
