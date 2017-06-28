#!/bin/sh

### BEGIN INIT INFO
# Provides:          eth_log_alerting
# Required-Start:    $network
# Required-Stop:     $network
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Log entry alerting
### END INIT INFO

set -e

. /lib/lsb/init-functions

DAEMON="/usr/local/eth-log-alerting/pipe.sh"
NAME="eth_log_alerting"
DESC="log alerting"
DAEMON_OPTIONS=
PID="/run/$NAME.pid"

# Check if DAEMON binary exist
[ -f $DAEMON ] || exit 0

[ -f "/etc/default/$NAME" ] && . /etc/default/$NAME

daemon_not_configured () {
  if [ "$1" != "stop" ]
  then
    log_daemon_msg "Configuration required! Update /etc/default/$NAME"
    exit 0
  fi
}

config_checks () {
  # Check that log is configured
  if [ -z "$DAEMON_OPTIONS" ]
  then
    daemon_not_configured "$1"
  fi
}

case "$1" in
  start)
    log_daemon_msg "Starting $DESC" "$NAME"
    config_checks "$1"
    if start-stop-daemon --start --quiet --oknodo --pidfile $PID --exec $DAEMON -- $DAEMON_OPTIONS 1>/dev/null
    then
      log_end_msg 0
    else
      log_end_msg 1
    fi
    ;;
  stop)
    log_daemon_msg "Stopping $DESC" "$NAME"
    if start-stop-daemon --retry TERM/5/KILL/5 --oknodo --stop --quiet --pidfile $PID 1>/dev/null
    then
      log_end_msg 0
    else
      log_end_msg 1
    fi
    ;;
  restart)
    log_daemon_msg "Restarting $DESC" "$NAME"
    start-stop-daemon --retry TERM/5/KILL/5 --oknodo --stop --quiet --pidfile $PID 1>/dev/null
    if start-stop-daemon --start --quiet --oknodo --pidfile $PID --exec $DAEMON -- $DAEMON_OPTIONS 1>/dev/null
    then
      log_end_msg 0
    else
      log_end_msg 1
    fi
    ;;
  status)
    status_of_proc -p $PID $DAEMON $NAME
    ;;
  *)
    log_action_msg "Usage: /etc/init.d/$NAME {start|stop|restart|status}"
    ;;
esac

exit 0
