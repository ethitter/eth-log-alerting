#!/bin/sh

DAEMON="/usr/local/eth-log-alerting/pipe.sh"
NAME="eth_log_alerting"

# Check if DAEMON binary exist
[ -f $DAEMON ] || exit 0

# Check if config exists
if [ -f "/etc/default/$NAME" ]
then
  . /etc/default/$NAME
else
  exit 0
fi

# On with it!
$DAEMON $DAEMON_OPTIONS
