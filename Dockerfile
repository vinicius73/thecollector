# syntax=docker/dockerfile:1

FROM golang:1.21-alpine as builder

RUN apk add --update --no-cache gcc g++ upx

WORKDIR /app

COPY . .

# args
ARG APP_VERSION=unknown
ARG GIT_HASH=unknown
ARG BUILD_DATE=unknown
ARG PKG=github.com/vinicius73/thecollector

ENV CGO_ENABLED=0

RUN go build \
  -a -installsuffix cgo -ldflags '-s -w -extldflags "-static"' \
  -ldflags "-X $PKG/pkg/vars.commit=$GIT_HASH -X $PKG/pkg/vars.version=$APP_VERSION -X $PKG/pkg/vars.buildDate=$BUILD_DATE" \
  -o ./bin/thecollector ./app && \
  upx -8 ./bin/thecollector

FROM alpine:3

RUN apk add --update --no-cache ca-certificates tzdata postgresql15-client

# args
ARG APP_REVISION=unknown
ARG BUILD_DATE=unknown
ARG BUILD_REF=unknown

# Labels.
LABEL name="thecollector" \
  description="The Collector - Database Backup Tool" \
  vcs.url="https://${PKG}" \
  vcs.ref=$GIT_HASH \
  version=$APP_VERSION \
  build.date=$BUILD_DATE \
  build.builder=$BUILDER

LABEL org.opencontainers.image.title="thecollector" \
  org.opencontainers.image.description="The Collector - Database Backup Tool" \
  org.opencontainers.image.url="https://$PKG" \
  org.opencontainers.image.source="https://$PKG" \
  org.opencontainers.image.revision="$BUILD_REF"

# Environment
ENV GIT_HASH=$GIT_HASH \
  APP_VERSION=$APP_VERSION \
  APP_REVISION=$APP_REVISION \
  TZ=America/Sao_Paulo \
  LOG_LEVEL=info \
  SCHEDULE_DUMP_CRON="0 1 * * *"\
  THECOLLECTOR_CONFIG_FILE=/thecollector/config.yml \
  THECOLLECTOR_TARGET_DIR=/thecollector/outputs


COPY docker/docker-entrypoint.sh /docker-entrypoint.sh
COPY thecollector.yml /thecollector/config.yml


ENV UID=1000
ENV GID=1000
ENV UMASK=022

WORKDIR /app

COPY --from=builder /app/bin /sbin

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["cron"]
