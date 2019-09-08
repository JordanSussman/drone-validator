# drone-validator

This is a simple PoC of the new [validation extension](https://github.com/drone/boilr-validate) to showcase failings builds if the specified external secret isn't allowed to be used with non whitelisted images, and if the step utilize `command` or `commands` stanza. This extension is merely meant for example purposes, and should **NOT** be utilized in any sort of production environment.

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
  --env=DRONE_CONFIG_FILE=/config.yml \
  --restart=always \
  --volume=${PWD}/config.yml:/config.yml \
  --name=drone-validator \
  jsussman/drone-validator
```

Update your Drone server configuration to include the plugin address, shared secret, and path to configuration file.

```text
DRONE_VALIDATE_PLUGIN_ENDPOINT=http://1.2.3.4:3000
DRONE_VALIDATE_PLUGIN_SECRET=bea26a2221fd8090ea38720fc445eca6
```

The configuration file should be in the following format:

```yml
secrets:
  <full path to external secret>:
    - <name of first allowed image>
  <full path to external secret>:
    - <name of first allowed image>
    - <name of second allowed image>
```

## Examples

All examples assume the configuration file looks like:

```yml
secrets:
  secret/data/docker/username:
    - alpine
```

The following `.drone.yml` will fail due to the step utilizing the `docker_username` with the `ubuntu` image

```yml
kind: pipeline
type: docker
name: default

steps:
- name: build
  image: ubuntu
  environment:
    USERNAME:
      from_secret: docker_username

---
kind: secret
name: docker_username
get:
  path: secret/data/docker
  name: username
```

The following `.drone.yml` will fail since the step utilizes the `docker_username` secret with the `commands` stanza

```yml
kind: pipeline
type: docker
name: default

steps:
- name: build
  image: alpine
  commands:
    - env
  environment:
    USERNAME:
      from_secret: docker_username

---
kind: secret
name: docker_username
get:
  path: secret/data/docker
  name: username
```

The following `.drone.yml` will succeed since the `alpine` image is allowed to utilize the `docker_username` secret and `commands` stanza isn't being utilized

```yml
kind: pipeline
type: docker
name: default

steps:
- name: build
  image: alpine
  environment:
    USERNAME:
      from_secret: docker_username

---
kind: secret
name: docker_username
get:
  path: secret/data/docker
  name: username
```
