#!/bin/sh

tail -Fq log.log | while read x ; do php im.php "$x" "#errors"; done
