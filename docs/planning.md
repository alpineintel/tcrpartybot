# TCR Party Bot
Note that these are pretty raw notes taken during planning conversations!

## Network Design
* Stake tokens to initiate a poll, (Y/N)
* People must register in order for their vote to count. Twitter poll results
  would be bad.
* List is capped to X (100, probs 50 for beta) accounts
* DM bot for token balance, DM bot to stake tokens to add yourself.
* Bot posts Y/N
* Anyone that owns a token is a governor on the list
* Run on testnet
* you can only nominate yourself.
* bot manages custody
* have two accounts: one for bot to interact with (polls, nomination, etc.),
  one for list.
* @TCRParty, @TCRPartyVIP

### Initial distribution:
* Anyone that wants to be a token holder needs to DM the bot.
* Follow @TCRBot, tweet "Let's party, @TCRBot", TCR Bot DMs them to do dance.
* Look at past 50 tweets to see their engagement to distribute tokens:
    * # of followers, likes, and replies, retweets. Interactions / total
      followers.
* Signup period, collects all of them. After it, does a snapshot and
  distributes proportionally.
* Send a DM with balance at the end.

### Nominate to list:
* Limit how many people can be nominated in a day.
* User DMs the bot to nominate themselves.
* Any nomination is 200.

### Remove from list:
* User A stakes X tokens (opens poll to remove user b).
* 1a. If user A wins the poll, user a keeps stake, user b is off the list.
* 1b. If user A loses the poll, user A loses stake user B gets the stake.
* 2. The voting losers pay the voting winners proportionally.

### Faucet:
* Everyone gets 1,000 tokens. 200 tokens to nominate (add or remove). Voting is
  100.
* If you lose a few times and lose everything, there is a faucet system (you
  can party again with the TCR bot) get 25 tokens/day. It's inflationary.
* Faucet is universally available, as 25 tokens' marginal value decreases over
  time.
* Badge of honor is having a lot of tokens.

## Application Design
* SqlX
* Expose a simple HTTP API with header authorization.
* One high-level process that spins off goroutines.
    * HTTP server routine 
    * One constantly polls twitter and pushes to an events channel 
    * Task-specific goroutines listen to events channel, picking off their
      relevant events and running logic based on them.
* Geth lookups are important 
