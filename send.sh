#!/bin/sh

MSG=${1:-Oops}
CHANNEL=${2:-#errors}
COLOR=${3:-default}

# TODO: source these from a common config
USERNAME="Log bot"
CONTEXT="New log entry"
ICON_URL=""
WEBHOOK_URL=""

/usr/bin/curl \
    -X POST \
    -s \
    --data-urlencode "payload={ \
        \"channel\": \"$CHANNEL\", \
        \"username\": \"$USERNAME\", \
        \"attachments\": [ { \
            \"fallback\": \"New log entry\", \
            \"pretext\": \"$CONTEXT\", \
            \"text\": \"$MSG\", \
            \"color\": \"$COLOR\" \
        } ],\
        \"icon_url\": \"$ICON_URL\" \
    }" \
    $WEBHOOK_URL
