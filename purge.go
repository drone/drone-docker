package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/drone/drone-docker/dockerhub"
	"github.com/genuinetools/reg/registry"
)

func purge(ctx context.Context, c *config) error {
	image, err := registry.ParseImage(c.Docker.Image)
	if err != nil {
		return err
	}

	// this is the temporary image that was published with
	// a random identifier. This tag should be deleted
	// once the image is properly re-tagged.
	alias, err := registry.ParseImage(c.Docker.ImageAlias)
	if err != nil {
		return err
	}

	r, err := createClient(
		ctx,
		c.Docker.Auth.Username,
		c.Docker.Auth.Password,
		// TODO(bradrydzewski) this is currently hard-coded
		// to DockerHub and needs to be modified to use the
		// correct registry URL.
		"https://index.docker.io",
	)
	if err != nil {
		return err
	}

	// HACK the default authorization and retry logic does
	// not work for http requests with a body, so we implement
	// a custom transport to work around this issue.
	r.Client.Transport = &retryTransport{
		Transport: r.Client.Transport,
	}

	// STEP 1: get the manifest for the temporary tag.
	manifest, err := r.Manifest(ctx, alias.Path, alias.Reference())
	if err != nil {
		return err
	}

	// STEP 2: upload the manifest for the user-defined tag.
	err = r.PutManifest(ctx, image.Path, image.Reference(), manifest)
	if err != nil {
		return err
	}

	// STEP 4: get the digest for the temporary tag.
	digest, err := r.Digest(ctx, alias)
	if err != nil {
		return err
	}
	if err := alias.WithDigest(digest); err != nil {
		return err
	}

	// STEP 5a: delete the tag from dockerhub using the
	// native dockerhub API. Dockerhub does not appear to
	// support deleting the tag by deleting the manifest,
	// and returns an UNSUPPORTED error.
	if strings.HasPrefix(c.Docker.Auth.Address, "https://index.docker.io") {
		return dockerhub.New(
			c.Docker.Auth.Username,
			c.Docker.Auth.Password,
		).DeleteTag(alias.Path, alias.Tag)
	}
	// STEP 5b: delete the digest from the registry.
	return r.Delete(ctx, alias.Path, digest)
}

// HACK the reg package (github.com/genuinetools/reg) tries to
// authenticate and re-send failed http requests. This does not
// work with http requests that have a body because the reader
// is EOF. So we implement a custom transport and buffer that
// snapshots the body so that it can retry.

type retryTransport struct {
	Transport http.RoundTripper
}

func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body == nil {
		res, err := t.Transport.RoundTrip(req)
		if res != nil && res.StatusCode > 299 {
			dump, _ := httputil.DumpResponse(res, true)
			println(string(dump))
		}
		return res, err
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = &retryReadCloser{
		buf:  bytes.NewBuffer(body),
		body: body,
	}
	return t.Transport.RoundTrip(req)
}

type retryReadCloser struct {
	buf  *bytes.Buffer
	body []byte
}

func (t *retryReadCloser) Read(p []byte) (n int, err error) {
	return t.buf.Read(p)
}

func (t *retryReadCloser) Close() error {
	t.buf.Reset()
	t.buf.Write(t.body)
	return nil
}
