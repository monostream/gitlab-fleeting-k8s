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
hecking for jobs...nothing                         runner=i5hxLjcN
2023-06-21T08:12:40.672Z [INFO]  plugin.fleeting-plugin-k8s: updating instance state: id=gitlab/docker-runner-autoscaler-0 namespace=gitlab state=creating statefulset=docker-runner-autoscaler timestamp=2023-06-21T08:12:40.672Z
2023-06-21T08:12:41.672Z [INFO]  plugin.fleeting-plugin-k8s: updating instance state: statefulset=docker-runner-autoscaler id=gitlab/docker-runner-autoscaler-0 namespace=gitlab state=running timestamp=2023-06-21T08:12:41.672Z
2023-06-21T08:12:41.672Z [INFO]  instance update: group=gitlab/docker-runner-autoscaler id=gitlab/docker-runner-autoscaler-0 state=running
2023-06-21T08:12:41.680Z [INFO]  plugin.fleeting-plugin-k8s: connect info: id=gitlab/docker-runner-autoscaler-0 info="map[ExternalAddr: ID: InternalAddr:10.240.0.21 arch:amd64 expires:2023-06-21T08:14:41.680063997Z keepalive:3e+10 key:<nil> os:linux password: protocol:ssh timeout:6e+11 use_static_credentials:true username:root]" namespace=gitlab statefulset=docker-runner-autoscaler timestamp=2023-06-21T08:12:41.680Z
2023-06-21T08:12:41.680Z [INFO]  ready: instance=gitlab/docker-runner-autoscaler-0 took=7.820779ms
WARNING: Job failed (system failure): panic: runtime error: invalid memory address or nil pointer dereference
  duration_s=62.01549636 job=34991 project=25 runner=i5hxLjcN
Appending trace to coordinator...ok                 code=202 job=34991 job-log=0-548 job-status=running runner=i5hxLjcN sent-log=383-547 status=202 Accepted update-interval=1m0s
Updating job...                                     bytesize=548 checksum=crc32:76dac24d job=34991 runner=i5hxLjcN
Feeding runners to channel                          builds=1
Feeding runner to channel                           builds=1 runner=i5hxLjcN
Processing runner                                   builds=1 runner=i5hxLjcN
Acquiring executor from provider                    builds=1 runner=i5hxLjcN
Acquiring job slot                                  builds=1 runner=i5hxLjcN
Acquiring request slot                              builds=1 runner=i5hxLjcN
Dialing: tcp gitlab.tribbles.cloud:443 ...         
Checking for jobs...nothing                         runner=i5hxLjcN
Submitting job to coordinator...accepted, but not yet completed  bytesize=548 checksum=crc32:76dac24d code=202 job=34991 job-status= runner=i5hxLjcN update-interval=1s
Updating job...                                     bytesize=548 checksum=crc32:76dac24d job=34991 runner=i5hxLjcN
Submitting job to coordinator...ok                  bytesize=548 checksum=crc32:76dac24d code=200 job=34991 job-status= runner=i5hxLjcN update-interval=0s
WARNING: Failed to process runner                   builds=0 error=panic: runtime error: invalid memory address or nil pointer dereference executor=docker-autoscaler runner=i5hxLjcN
```