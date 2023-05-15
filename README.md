# celestia-lightnode-monitor

A simple monitor for the Celestia Light node. 
# Features
* It will send alert in discord channel when your celestia light node is down, not synced, or lack of balance. 
* It will check the status every 5 minutes.


# Pre-requisites

1. You have go 1.20+ installed.
2. If you do not run monitor on the same machine as your light node, you need to open the gateway using flag `--gateway --gateway.addr 0.0.0.0 --gateway.port 26659`, also the port should be open.
3. You have set a discord webhook. You can follow this guide: https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks

# How to use

## Input your webhook url in .env file
`cp .env.example .env`
then edit .env file

## Edit config.toml
* `standardRPC`: Celestia Network Public RPC address
* `gatewayAPI`: the Light Node Gateway API address,default is http://localhost:26659 , you can input as many as you need to monitor
* `minimumBalance`: Minimum balance of your light node, if your balance is lower than this value, you will receive an alert. Denom is utia.
* `alertuserid`: Discord user id. If you want be @ when receive alert in discord, set this.  https://support.discord.com/hc/en-us/articles/206346498-Where-can-I-find-my-User-Server-Message-ID-
* `alertroleid`: Discord role id. If you want the people who have special discord role be @ when receive alert in discord, set this. https://www.itgeared.com/how-to-get-role-id-on-discord/

## Build 
```shell
 go build -o celestia-lightnode-monitor  ./cmd
```


## Run
```shell
./celestia-lightnode-monitor
```

You can check the log in `celestia-lightnode-monitor.log`

The alert is like this:
![](https://i.imgur.com/pKokq3j.jpg)