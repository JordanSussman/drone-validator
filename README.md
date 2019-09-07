# drone-validator

This is a simple PoC of the new [validation extension](https://github.com/drone/boilr-validate) to showcase failing builds if they have steps that utilize `command` or `commands` stanzas with secrets for specific image. This extension is merely meant for example purposes, and should **NOT** be utilized in any sort of production environment.

_Please note this project requires Drone server version 1.4 or higher._

## Installation

Create a shared secret:

```console
$ openssl rand -hex 16
bea26a2221fd8090ea38720fc445eca6
```

Download and run the plugin:

```console
$ docker run -d \
  --publish=3000:3000 \
  --env=DRONE_DEBUG=true \
  --env=DRONE_SECRET=bea26a2221fd8090ea38720fc445eca6 \
  --restart=always \
  --name=drone-validator \
  jsussman/drone-validator
```

Update your Drone server configuration to include the plugin address, shared secret, and list of disallowed docker images.

```text
DRONE_VALIDATE_PLUGIN_ENDPOINT=http://1.2.3.4:3000
DRONE_VALIDATE_PLUGIN_SECRET=bea26a2221fd8090ea38720fc445eca6
DRONE_DISALLOWED_ENTRYPOINT=alpine,plugins/docker
```

## Example

Both examples assume `DRONE_DISALLOWED_ENTRYPOINT=alpine`.

The following `.drone.yml` will fail due to the step utilizing the `alpine` image with `commands` stanza **and** secrets:

```yml
---
kind: pipeline
type: docker
name: default

steps:
- name: failure
  image: alpine
  commands:
  - env
  settings:
    username:
      from_secret: docker_username
```

The following `.drone.yml` will succeed since the image isn't defined as being disallowed from utilizing `commands` stanza with secrets:

```yml
---
kind: pipeline
type: docker
name: default

steps:
- name: success
  image: ubuntu
  commands:
  - env
  settings:
    username:
      from_secret: docker_username
```
