# TCR Party Bot!

This repository implements the brains behind the TCR Party Bot. This bot
facilitates all interactions between Twitter and the underlying smart contracts
on the blockchain in addition to retweeting all members of the live TCR.

_More docs to come..._

## Environment variables
```
VIP_BOT_HANDLE=tcrpartyvip
PARTY_BOT_HANDLE=tcrpartybot
DATABASE_URI=data.db
TWITTER_CONSUMER_KEY=
TWITTER_CONSUMER_SECRET=

BASE_URL=https://bbd828f3.ngrok.io
SERVER_HOST=0.0.0.0:8080
```

## Setup
Warning: hastily written documentation ahead. This will be improved before
release:

1. Spin up the binary.
2. Run `auth-vip` and `auth-party` to set up Twitter OAuth credentials.
3. Run `create-webhook` to create the webhook that allows for receiving DMs

## Migrations
We use [Migrate](https://github.com/golang-migrate/migrate/tree/master/cli) for
database migrations.
You should read the repository's documentation for more details. As a shortcut,
you may use the script `bin/migrate` to access a preconfigured CLI for handling
migrations.
