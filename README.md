# motion-bot

A telegram bot to control and receive captures from [Motion](https://motion-project.github.io/)

## Features

* Receive pictures and videos captured by Motion.
* Sends commands to get snapshot, status, start and pause motion detection.
* Manage authorized users and chats/groups to control and see your data.
* Prometheus metrics instrumented

## Install
Installations steps for Raspberry Pi with a Pi camera running motion-bot, Motion and Prometheus with docker-compose file in `/deployments`.
### Requirements:

* Telegram API token. Found more info [here](https://core.telegram.org/bots#3-how-do-i-create-a-bot)
* Your user ID. You can use the [userinfo](https://telegram.me/userinfobot)
* Docker Compose

Fill the `<required>` parameters of the `deployments/etc/*.sample` configuration files and run:
```bash
cd deployments
docker-compose up -d
```
## Build and Test

If you have Go installed in your environment, you can compile and test this project locally with:
```bash
make
```

Or with docker (linux):

```bash
docker build . -t motion-bot
```
