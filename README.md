Log Alerting
============

Pipe logs to Mattermost (or Slack) webhooks

## Requirements

* Go 1.8.3

## Installation

1. `git clone https://git.ethitter.com/debian/eth-log-alerting.git /usr/local/bin/eth-log-alerting`
1. `cd /usr/local/bin/eth-log-alerting`
1. `go get github.com/42wim/matterbridge/matterhook`
1. `go get github.com/asaskevich/govalidator`
1. `go get github.com/hpcloud/tail`
1. `go build eth-log-alerting.go`
1. `cp /usr/local/bin/eth-log-alerting/init.sh /etc/init.d/eth-log-alerting`
1. `chmod +x /etc/init.d/eth-log-alerting`
1. `cp /usr/local/bin/eth-log-alerting/config-sample.json /usr/local/bin/eth-log-alerting/config.json`
1. Edit `/usr/local/bin/eth-log-alerting/config.json`
1. `cp /usr/local/bin/eth-log-alerting/defaults /etc/default/eth-log-alerting`
1. Edit `/etc/default/eth_log_alerting`
1. `update-rc.d eth-log-alerting defaults`
1. `/etc/init.d/eth-log-alerting start`
