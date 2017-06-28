Logs to alerts
==============

Pipe logs to Mattermost or Slack webhooks

# Requirements
* Mattermost or Slack instance
* `curl`, `jq`, `tail`

# Installation
1. `git clone https://git.ethitter.com/debian/eth-log-alerting.git /usr/local/eth-log-alerting`
2. `cp /usr/local/eth-log-alerting/init.sh /etc/init.d/eth_log_alerting`
3. `cp /usr/local/eth-log-alerting/defaults /etc/default/eth_log_alerting`
4. Edit `/etc/default/eth_log_alerting`
5. `sudo update-rc.d eth_log_alerting defaults`
6. `/etc/init.d/eth_log_alerting start`
