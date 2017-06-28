#!/bin/bash

# Defaults and aliases
LOG_FILE=$1
WEBHOOK_URL=$2
USERNAME=${3:-logbot}
CHANNEL=${4:-#logs}
COLOR=${5:-default}
ICON_URL=$6
GREP=$7

tail -n0 -F "$LOG_FILE" | while read LINE; do
    (echo "$LINE" | grep -e "$GREP") && jq -n --arg line_encoded "    $LINE"  \ "{ \
        channel: \"$CHANNEL\", \
        username: \"$USERNAME\", \
        attachments: [ { \
            fallback: \"New entry in $LOG_FILE\", \
            pretext: \"\`$LOG_FILE\`\", \
            text: \$line_encoded, \
            color: \"$COLOR\" \
        } ],\
        icon_url: \"$ICON_URL\" \
    }" | cat <(echo "payload=") <(cat -) | curl \
    -X POST \
    -s \
    -d@- \
    $WEBHOOK_URL > /dev/null;
done
