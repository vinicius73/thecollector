#!/bin/sh
set -o errexit
set -o nounset

OPTIND=1

ENTRYPOINT_ACTION="$1"

BANNER=$(cat <<- "EOF"
     |\_/|
     | @ @   Woof!
     |   <>              _
     |  _/\------____ ((| |))
     |               `--' |
 ____|_       ___|   |___.'
/_/_____/____/_______|
EOF
)

banner () {
  echo "${BANNER}

  ACTION      = ${ENTRYPOINT_ACTION}
  APP_VERSION = ${APP_VERSION}
  GIT_HASH    = ${GIT_HASH}

  DB_HOST     = ${DB_HOST:-}
  DB_PORT     = ${DB_PORT:-}
  DB_USER     = ${DB_USER:-}
  DB_NAME     = ${DB_NAME:-}
  LOG_LEVEL   = ${LOG_LEVEL:-}
  LOG_FORMAT  = ${LOG_FORMAT:-}
  BUCKET_KEY  = ${BUCKET_KEY:-}
  BUCKET_NAME = ${BUCKET_NAME:-}

  TZ          = ${TZ:-}

  SCHEDULE_DUMP_CRON = ${SCHEDULE_DUMP_CRON:-}

  THECOLLECTOR_TARGET_DIR  = ${THECOLLECTOR_TARGET_DIR:-}
  THECOLLECTOR_CONFIG_FILE = ${THECOLLECTOR_CONFIG_FILE:-}

  "
}

{
  update-ca-certificates 2>&1
}

cron_action () {
  thecollector cron
}

case $ENTRYPOINT_ACTION in
  cron)
  banner && cron_action
  ;;
  *)
  exec "$@"
  ;;
esac
