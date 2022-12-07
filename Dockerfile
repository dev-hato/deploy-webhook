FROM golang:1.19.4-bullseye AS base

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
ARG TARGETPLATFORM

RUN apt-get update \
    && apt-get install -y --no-install-recommends lsb-release \
    && mkdir -p /etc/apt/keyrings \
    && curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg \
    && echo \
         "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
         $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null \
    && apt-get update \
    && apt-get install -y --no-install-recommends docker-ce docker-ce-cli containerd.io docker-compose-plugin \
    && apt-get remove -y lsb-release \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists \
    && find / -type f -perm /u+s -ignore_readdir_race -exec chmod u-s {} \; \
    && find / -type f -perm /g+s -ignore_readdir_race -exec chmod g-s {} \; \
    && useradd -l -m -s /bin/bash -N -u "1000" "nonroot"

WORKDIR /go/app

RUN chown -R nonroot .

COPY go.mod go.sum main.go restart.sh healthcheck.sh ./

RUN mapfile -t PLATFORM < <(echo "${TARGETPLATFORM}" | tr '/' ' ') \
    && CGO_ENABLED=1 GOOS=linux GOARCH=${PLATFORM[2]} go build -o deploy-webhook .

HEALTHCHECK --interval=5s --retries=20 CMD ["./healthcheck.sh"]

FROM base AS main

USER nonroot

ENTRYPOINT ["./deploy-webhook"]

FROM base AS develop

COPY .air.toml .

RUN go install github.com/cosmtrek/air

USER nonroot

CMD ["air", "-c", ".air.toml"]
