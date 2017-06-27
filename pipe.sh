#!/bin/sh
# TODO: path and config for send.sh
/usr/bin/tail -Fq log.log | while read x ; do send.sh "$x" "#errors" "danger"; done
