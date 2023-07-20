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

# Current blocking error

The runner picks up the job successfully, scales up the statefulset, polls the created agent for ready status, but as soon as it gets into ready state, there is a misteryous nil pointer exception happening:

```shell
2023-07-20T14:05:08.657Z [INFO]  plugin.fleeting-plugin-k8s: updating instance state: state=creating statefulset=docker-runner-autoscaler id=gitlab/docker-runner-autoscaler-0 namespace=gitlab timestamp=2023-07-20T14:05:08.657Z
2023-07-20T14:05:09.659Z [INFO]  plugin.fleeting-plugin-k8s: updating instance state: id=gitlab/docker-runner-autoscaler-0 namespace=gitlab state=running statefulset=docker-runner-autoscaler timestamp=2023-07-20T14:05:09.659Z
2023-07-20T14:05:09.659Z [INFO]  instance update: group=gitlab/docker-runner-autoscaler id=gitlab/docker-runner-autoscaler-0 state=running
2023-07-20T14:05:09.667Z [INFO]  plugin.fleeting-plugin-k8s: connect info: namespace=gitlab statefulset=docker-runner-autoscaler id=gitlab/docker-runner-autoscaler-0 info="map[ExternalAddr:10.240.0.45 ID:gitlab/docker-runner-autoscaler-0 InternalAddr:10.240.0.45 arch:amd64 expires:2023-07-20T14:07:09.666832002Z keepalive:6e+10 key:<nil> os:linux password:blub protocol:ssh timeout:3e+11 use_static_credentials:true username:root]" timestamp=2023-07-20T14:05:09.666Z
2023-07-20T14:05:09.667Z [INFO]  ready: instance=gitlab/docker-runner-autoscaler-0 took=7.635436ms
WARNING: Job failed (system failure): panic: runtime error: invalid memory address or nil pointer dereference
  duration_s=62.018934829 job=39606 project=96 runner=wAAhgyYd
Appending trace to coordinator...ok                 code=202 job=39606 job-log=0-548 job-status=running runner=wAAhgyYd sent-log=383-547 status=202 Accepted update-interval=1m0s
Updating job...                                     bytesize=548 checksum=crc32:e9e3f2cb job=39606 runner=wAAhgyYd
Submitting job to coordinator...accepted, but not yet completed  bytesize=548 checksum=crc32:e9e3f2cb code=202 job=39606 job-status= runner=wAAhgyYd update-interval=1s
Feeding runners to channel                          builds=1
Feeding runner to channel                           builds=1 runner=wAAhgyYd
Processing runner                                   builds=1 runner=wAAhgyYd
```