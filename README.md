# Fleeting Plugin Kubernetes

This is a go plugin for fleeting on Kubernetes. It is intended to be run by
[fleeting](https://gitlab.com/gitlab-org/fleeting/fleeting), and cannot be run directly.

## Known issues

- Everytime the runner pod restarts a it registers itself as a new GitLab runner, because the runner token isn't persisted yet. So make sure to clean up all the registered runners from time to time when expereimenting with this

## Building

To build the runner image:

```shell
docker build . --platform=linux/amd64 -f Dockerfile -t monostream/gitlab-runner-docker-autoscaler:0.0.1 && docker push monostream/gitlab-runner-docker-autoscaler:0.0.1
```

To build the agent image:

```shell
docker build . --platform=linux/amd64 -f Dockerfile-agent -t monostream/gitlab-runner-docker-autoscaler-agent:0.0.1 && docker push monostream/gitlab-runner-docker-autoscaler-agent:0.0.1
```

# Installation

install with helm:

```shell
helm upgrade --install docker-runner-autoscaler ./chart -n gitlab -f docker-runner-autoscaler.yaml
```