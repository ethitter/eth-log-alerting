#!/bin/sh
### BEGIN INIT INFO
# Provides:          eth-log-alerting
# Required-Start:    $network $local_fs
# Required-Stop:     $network $local_fs
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Log alerting via Mattermost
# Description:       Pipe myriad logs to Mattermost webhooks
### END INIT INFO

. /lib/lsb/init-functions

NAME=eth-log-alerting
USER=root
DAEMON="/usr/local/bin/${NAME}/${NAME}"
DAEMON_ARGS="-config /usr/local/bin/${NAME}/config.json"
PIDFILE="/var/run/${NAME}.pid"

# Overrides
[ -f "/etc/default/$NAME" ] && . /etc/default/$NAME

# If the daemon is not there, then exit.
test -x $DAEMON || exit 5

case $1 in
 start)
  # Checked the PID file exists and check the actual status of process
  if [ -e $PIDFILE ]; then
   status_of_proc -p $PIDFILE $DAEMON "$NAME process" && status="0" || status="$?"
   # If the status is SUCCESS then don't need to start again.
   if [ $status = "0" ]; then
    exit # Exit
   fi
  fi
  # Start the daemon.
  log_daemon_msg "Starting the process" "$NAME"
  # Start the daemon with the help of start-stop-daemon
  # Log the message appropriately
  if start-stop-daemon --start --chuid $USER --background --oknodo --pidfile $PIDFILE --make-pidfile --exec $DAEMON -- $DAEMON_ARGS; then
   log_end_msg 0
  else
   log_end_msg 1
  fi
  ;;
 stop)
  # Stop the daemon.
  if [ -e $PIDFILE ]; then
   status_of_proc -p $PIDFILE $DAEMON "Stoppping the $NAME process" && status="0" || status="$?"
   if [ "$status" = 0 ]; then
    start-stop-daemon --stop --quiet --oknodo --pidfile $PIDFILE
    /bin/rm -rf $PIDFILE
   fi
  else
   log_daemon_msg "$NAME process is not running"
   log_end_msg 0
  fi
  ;;
 restart)
  # Restart the daemon.
  $0 stop && sleep 2 && $0 start
  ;;
 status)
  # Check the status of the process.
  if [ -e $PIDFILE ]; then
   status_of_proc -p $PIDFILE $DAEMON "$NAME process" && exit 0 || exit $?
  else
   log_daemon_msg "$NAME Process is not running"
   log_end_msg 0
  fi
  ;;
 reload)
  # Reload the process. Basically sending some signal to a daemon to reload
  # it configurations.
  if [ -e $PIDFILE ]; then
   start-stop-daemon --stop --signal USR1 --quiet --pidfile $PIDFILE --name $NAME
   log_success_msg "$NAME process reloaded successfully"
  else
   log_failure_msg "$PIDFILE does not exists"
  fi
  ;;
 *)
  # For invalid arguments, print the usage message.
  echo "Usage: $0 {start|stop|restart|reload|status}"
  exit 2
  ;;
esac
