
# defines a secret resource that points to external
# docker credentials.
local SecretDocker = {
  kind: "secret",
  type: "external",
  external_data: {
    "username": {
        "path": "drone/docker",
        "name": "username",
    },
    "password": {
        "path": "drone/docker",
        "name": "password",
    },
  },
};

# returns a Pipeline resource that builds, tests and
# publishes Docker images for the Linux architecture.
local PipelineLinux(os='linux', arch='amd64', variant='') = {
  kind: "pipeline",
  name: os + '-' + arch,
  platform: {
    os: os,
    arch: arch,
  },
  steps: [
    # this step is responsible for building
    # and testing the drone docker binary.
    {
      name: "build",
      image: "golang:1.11",
      commands: [
        "go test -v",
        "./scripts/build_"+os+"_"+arch+".sh",
      ],
    },
    # this step is responsible for building
    # and testing the drone docker image.
    {
      name: "publish",
      image: "plugins/docker",
      settings: {
        auto_tag: true,
        auto_tag_suffix: os + "-" + arch,
        dockerfile: "docker/Dockerfile." + os + "." + arch,
        repo: "drone/docker",
        username: "drone",
        password: { "from_secret": "password" },
      },
      when: {
        event: [ "push", "tag" ],
      },
    },
  ],
};

# returns a Pipeline resource that pushes a Docker
# manifest for the server and aganet.
local PipelineManifest = {
  kind: "pipeline",
  name: "manifest",
  depends_on: [
    "build-linux-arm",
    "build-linux-arm64",
    "build-linux-amd64",
  ],
  steps: [
    {
      name: "publish",
      image: "plugins/manifest:1",
      settings: {
        username: { "from_secret": "username" },
        password: { "from_secret": "password" },
        spec: "docker/manifest.tmpl",
        ignore_missing: true,
      },
      when: {
        event: [ "push" ]
      }
    },
  ],
};

[
    PipelineLinux('linux', 'amd64'),
    PipelineLinux('linux', 'arm64'),
    PipelineLinux('linux', 'arm'),
    PipelineManifest,
    SecretDocker,
]