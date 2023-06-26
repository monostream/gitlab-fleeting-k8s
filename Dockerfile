# builder
FROM golang:1.20-alpine3.18 as builder

WORKDIR /src

COPY . .

WORKDIR /src/cmd/fleeting-plugin-k8s

ENV CGO_ENABLED=0
RUN go mod download
RUN go build

FROM gitlab/gitlab-runner:v16.1.0

ENV XDG_RUNTIME_DIR /run/user/999

RUN mkdir -p $XDG_RUNTIME_DIR \
   && chown 999:999 $XDG_RUNTIME_DIR

# Docker CLI
ENV DOCKER_CLIENT_VERSION="23.0.1"
RUN curl -fsSL -o - https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKER_CLIENT_VERSION}.tgz | tar -zxf - --strip=1 -C /usr/local/bin/ docker/docker

# Docker Buildx Plugin
ENV DOCKER_BUILDX_VERSION="v0.11.0"
RUN mkdir -p /usr/libexec/docker/cli-plugins && curl -fsSL "https://github.com/docker/buildx/releases/download/${DOCKER_BUILDX_VERSION}/buildx-${DOCKER_BUILDX_VERSION}.linux-amd64" -o /usr/libexec/docker/cli-plugins/docker-buildx && chmod +x /usr/libexec/docker/cli-plugins/docker-buildx

# Docker-in-Docker daemon socket
RUN touch ${XDG_RUNTIME_DIR}/docker.sock && ln -s ${XDG_RUNTIME_DIR}/docker.sock /run/docker.sock

WORKDIR /app/

ENV PATH "/app:${PATH}"

COPY --from=builder --chown=nobody:nobody /src/cmd/fleeting-plugin-k8s/fleeting-plugin-k8s .