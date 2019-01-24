# TCR Party Bot!

This repository implements the brains behind the TCR Party Bot. This bot
facilitates all interactions between Twitter and the underlying smart contracts
on the blockchain in addition to retweeting all members of the live TCR.

## Contributing
Since you're reading this you've probably realized that TCR Party is an open
source project! If you're feeling particularly frisky and have a feature that
you'd like to see added or a bug that needs to be fixed we'd love for you to
contribute.

The best place to get an inside look into our development process is by
checking out our issue [board](https://gitlab.com/alpinefresh/tcr-party/tcrpartybot/boards) and [list](https://gitlab.com/alpinefresh/tcr-party/tcrpartybot/issues), which
keep track of tasks we're looking to complete. If you notice your request
hasn't been already mentioned, feel free to open a ticket to begin discussion.
We'd ask that you **please open a ticket to discuss feature requests/issues
before opening a PR**, as we would like to make sure that there is no
wasted/duplicated efforts. Nobody likes spending hours building a PR just to
see it go stale in the tracker!

Additionally, we require all discussions around TCR Party follow Recurse
Center's wonderful [code of conduct](https://www.recurse.com/code-of-conduct)
in order to ensure a welcoming and fun environment for everybody.

## Dev environment setup
### Prepare dependencies
In order to get started you'll need to first deploy our contracts to a local
blockchain. You can find more info on how to do this in the
[contracts repository](https://gitlab.com/alpinefresh/tcr-party/contracts).

Additionally, you'll need to have a working installation of
[Golang](https://golang.org/) and [Dep](https://github.com/golang/dep) and a
local Postgres database to work with.

### Set up your environment file
We'll need to set up some variables dependant on your environment to get
started. Create a new file in your cloned repository `tcrpartybot/.env` with
the following contents:

```bash
API_TOKEN_HASH='some-bcrypt-hash-here'
VIP_BOT_HANDLE=tcrpartyvip
PARTY_BOT_HANDLE=tcrpartybot
DATABASE_URL=postgres://localhost:5432/tcrparty?sslmode=disable
TWITTER_CONSUMER_KEY=
TWITTER_CONSUMER_SECRET=
TWITTER_ENV=
SEND_TWITTER_INTERACTIONS=false
PREREGISTRATION=false
INITIAL_DISTRIBUTION_AMOUNT=3000

BASE_URL=https://bbd828f3.ngrok.io
SERVER_HOST=0.0.0.0:8080

ETH_NODE_URI=http://localhost:8545
MASTER_PRIVATE_KEY= # Note that this should omit the opening 0x of the private key

TOKEN_ADDRESS=
WALLET_FACTORY_ADDRESS=
TCR_ADDRESS=
START_BLOCK=0
```

A few notes on the fields:
* You can generate the API_TOKEN_HASH by using [this tool](https://bcrypt-generator.com/).
* You'll need to populate `TOKEN_ADDRESS`, `WALLET_FACTORY_ADDRESS` and `TCR_ADDRESS` using the contract addresses that are spit out when you run `truffle migrate` in the [contracts repository](https://gitlab.com/alpinefresh/tcr-party/contracts).
* You can leave `START_BLOCK` at 0 unless you're seeing performance issues when running the initial sync (ie if you're deploying to mainnet or testnet).
* `MASTER_PRIVATE_KEY` should be a key with some ETH in its balance, otherwise
  the bot will run into issues when creating transactions. Ganache provides a
  list of addresses with 100ETH when you start your development chain.
* If you want to test with Twitter you'll need to get a consumer key, secret,
  and env from [their developer portal](https://developer.twitter.com). This is
  an annoying process and isn't necessary to develop, so don't worry about it
  if you just want to make some smaller changes or play around.
* On that note, you should probably keep `SEND_TWITTER_INTERACTIONS` set to
  false, as this will tell the bot to avoid attempting to send data to
  Twitter's API and instead just echo interactions to the console for
  debugging purposes.

### Database setup
Connect to your postgres instance and create a new database entitled `tcrparty`:

```SQL
CREATE DATABASE tcrparty;
```

Ensure whatever role you're using has permission to connect to and modify this database:

```SQL
GRANT ALL PRIVILEGES ON DATABASE tcrparty TO steve;
```

Now let's set up our schema! We use [Migrate](https://github.com/golang-migrate/migrate/tree/master/cli) for
database migrations, which you can find installation instructions for in their
README. You can find a conveniently preconfigured script for running migration
commands in `bin/migrate`.

Once you have migrate set up, you can set up your database by running:

```bash
$ bin/migrate up
```

### Starting the server
Once you have everything set up (reminder: did you clone the contracts
repository and deploy them onto your local Ganache instance?) you should be able
to run the party bot by `cd`ing into the `tcrpartybot` directory and running:

```bash
$ go run *.go -repl
```

Notice the `-repl` flag, which tells the bot to start an interactive command
line for you to simulate various interactions such as DMs, mentions, etc.
without needing to use Twitter's API.
