#!/bin/sh

LISTENER_BIN=/usr/local/eth-log-alerting/pipe.sh
test -x $LISTENER_BIN || exit 5
PIDFILE=/var/run/eth_log_alerting.pid

case "$1" in
      start)
            echo -n "Starting log alerting.... "
            startproc -f -p $PIDFILE $LISTENER_BIN
            echo "running"
            ;;
          *)
            echo "Usage: $0 start"
            exit 1
            ;;
esac
