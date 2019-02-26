package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/events"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
)

const (
	disbursementMsg = "Your wallet has been built successfully!. Your balance has been set to %d TCRP. Respond with 'help' to see what you can do with these tokens."
)

type incomingDM struct {
	Type             string `json:"type"`
	ID               string `json:"id"`
	CreatedTimestamp string `json:"created_timestamp"`
	MessageCreated   struct {
		SenderID    string `json:"sender_id"`
		MessageData struct {
			Text string `json:"text"`
		} `json:"message_data"`
	} `json:"message_create"`
}

type user struct {
	ID         int64  `json:"id"`
	ScreenName string `json:"screen_name"`
}

type incomingTweet struct {
	IDStr string `json:"id_str"`
	Text  string `json:"text"`
	User  user   `json:"user"`
}

type incomingFollow struct {
	Source struct {
		ScreenName string `json:"screen_name"`
		ID         string `json:"id"`
	} `json:"source"`
}

type incomingWebhook struct {
	ForUserID           string           `json:"for_user_id"`
	DirectMessageEvents []incomingDM     `json:"direct_message_events"`
	TweetCreateEvents   []incomingTweet  `json:"tweet_create_events"`
	FollowEvents        []incomingFollow `json:"follow_events"`
}

// Server holds relevant data for running an API server
type Server struct {
	errChan    chan<- error
	eventsChan chan<- *events.TwitterEvent
}

func getAccountFromRequestID(body string) (*models.Account, error) {
	args := strings.Split(body, ",")
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return nil, err
	}

	account, err := models.FindAccountByID(id)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func getIntArgFromRequest(body string) (int64, error) {
	args := strings.Split(body, ",")
	fmt.Println(args)
	if len(args) != 2 {
		return 0, fmt.Errorf("not enough arguments")
	}

	val, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (server *Server) processDMs(dms []incomingDM) {
	botToken, err := models.FindOAuthTokenByHandle(os.Getenv("VIP_BOT_HANDLE"))
	if err != nil || botToken == nil {
		return
	}

	for _, dm := range dms {
		fromID, err := strconv.ParseInt(dm.MessageCreated.SenderID, 10, 64)
		if err != nil {
			server.errChan <- err
			continue
		}

		// Ignore outgoing DM events
		if fromID == botToken.TwitterID {
			continue
		}

		account, err := models.FindAccountByTwitterID(fromID)
		sourceHandle := ""
		if err != nil {
			server.errChan <- err
			continue
		} else if account != nil {
			sourceHandle = account.TwitterHandle
		}

		server.eventsChan <- &events.TwitterEvent{
			EventType:    events.TwitterEventTypeDM,
			Time:         time.Now().UTC(),
			SourceHandle: sourceHandle,
			SourceID:     fromID,
			Message:      dm.MessageCreated.MessageData.Text,
		}
	}
}

func (server *Server) processMentions(tweets []incomingTweet) {
	for _, tweet := range tweets {
		server.eventsChan <- &events.TwitterEvent{
			EventType:    events.TwitterEventTypeMention,
			ObjectID:     tweet.IDStr,
			Time:         time.Now().UTC(),
			SourceHandle: tweet.User.ScreenName,
			SourceID:     tweet.User.ID,
			Message:      tweet.Text,
		}
	}
}

func (server *Server) processFollows(follows []incomingFollow) {
	botToken, err := models.FindOAuthTokenByHandle(os.Getenv("VIP_BOT_HANDLE"))
	if err != nil || botToken == nil {
		return
	}

	for _, follow := range follows {
		id, err := strconv.ParseInt(follow.Source.ID, 10, 64)
		if err != nil {
			server.errChan <- err
			continue
		}

		// Ignore outgoing DM events
		if id == botToken.TwitterID {
			continue
		}

		server.eventsChan <- &events.TwitterEvent{
			EventType:    events.TwitterEventTypeFollow,
			Time:         time.Now().UTC(),
			SourceHandle: follow.Source.ScreenName,
			SourceID:     id,
		}
	}
}

func (server *Server) handleTwitterWebhook(w http.ResponseWriter, r *http.Request) {
	// A GET request signals that Twitter is attempting a CRC request
	if r.Method == "GET" {
		keys, ok := r.URL.Query()["crc_token"]
		if !ok || len(keys) < 1 {
			w.WriteHeader(400)
			w.Write([]byte("Bad request"))
			return
		}

		mac := hmac.New(sha256.New, []byte(os.Getenv("TWITTER_CONSUMER_SECRET")))
		mac.Write([]byte(keys[0]))

		token := "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("{\"response_token\": \"" + token + "\"}"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	data := &incomingWebhook{}
	err := decoder.Decode(data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		w.Write([]byte("Bad request"))
		return
	}

	if len(data.DirectMessageEvents) > 0 {
		go server.processDMs(data.DirectMessageEvents)
	}

	if len(data.TweetCreateEvents) > 0 {
		go server.processMentions(data.TweetCreateEvents)
	}

	if len(data.FollowEvents) > 0 {
		go server.processFollows(data.FollowEvents)
	}

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (server *Server) createWebhook(w http.ResponseWriter, r *http.Request) {
	webhookID, err := models.GetKey("webhookID")
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	// If we don't already have a webhook ID we should create it
	if webhookID == "" {
		id, err := twitter.CreateWebhook()
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Error: " + err.Error()))
			return
		}

		log.Printf("Webhook %s created successfully", id)
		models.SetKey("webhookID", id)
	}

	// And subscribe to TCRPartyVIP's DMs
	if err := twitter.CreateSubscription(); err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	log.Printf("Subscription created successfully")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func authenticateUser(user string, w http.ResponseWriter, r *http.Request) {
	reqToken, err := models.GetKey(models.TwitterRequestTokenKey)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	if reqToken == "" {
		request := &twitter.OAuthRequest{
			Handle: user,
		}

		url, err := request.GetOAuthURL()
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Error: " + err.Error()))
			return
		}

		err = models.SetKey(models.TwitterRequestTokenKey, request.RequestToken)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Error: " + err.Error()))
			return
		}

		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("Go to %s to get pin", url)))
	} else {
		// Convert body to string (this will be our PIN)
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		pin := buf.String()

		// Authenticate with Twitter
		request := &twitter.OAuthRequest{
			Handle:       user,
			RequestToken: reqToken,
			PIN:          pin,
		}

		if err = request.ReceivePIN(); err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Error: " + err.Error()))
			return
		}

		if err = models.ClearKey(models.TwitterRequestTokenKey); err != nil {
			w.WriteHeader(400)
			w.Write([]byte("DB Error: " + err.Error()))
			return
		}

		w.WriteHeader(200)
		w.Write([]byte("Success!"))
	}
}

func (server *Server) authenticateVIP(w http.ResponseWriter, r *http.Request) {
	authenticateUser(os.Getenv("VIP_BOT_HANDLE"), w, r)
}

func (server *Server) authenticateParty(w http.ResponseWriter, r *http.Request) {
	authenticateUser(os.Getenv("PARTY_BOT_HANDLE"), w, r)
}

func (server *Server) redeployWallet(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	account, err := getAccountFromRequestID(buf.String())

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	tx, identifier, err := contracts.DeployWallet()
	if err != nil {
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	err = account.SetMultisigFactoryIdentifier(identifier)
	if err != nil {
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(tx.Hash().Hex()))
}

func (server *Server) voteBalance(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	account, err := getAccountFromRequestID(buf.String())
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	if !account.MultisigAddress.Valid {
		w.WriteHeader(400)
		w.Write([]byte("Account does not have a multisig address"))
		return
	}

	balance, err := contracts.PLCRFetchBalance(account.MultisigAddress.String)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	humanBalance := contracts.GetHumanTokenAmount(balance).Int64()
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Voting balance: %d", humanBalance)))
}

func (server *Server) walletBalance(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	account, err := getAccountFromRequestID(buf.String())
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	if !account.MultisigAddress.Valid {
		w.WriteHeader(400)
		w.Write([]byte("Account does not have a multisig address"))
		return
	}

	balance, err := contracts.GetTokenBalance(account.MultisigAddress.String)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	humanBalance := contracts.GetHumanTokenAmount(balance).Int64()
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Wallet balance: %d", humanBalance)))
}

func (server *Server) voteWithdraw(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	account, err := getAccountFromRequestID(buf.String())
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	if !account.MultisigAddress.Valid {
		w.WriteHeader(400)
		w.Write([]byte("Account does not have a multisig address"))
		return
	}

	amount, err := getIntArgFromRequest(buf.String())
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	toWithdraw := contracts.GetAtomicTokenAmount(amount)
	balance, err := contracts.PLCRFetchBalance(account.MultisigAddress.String)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	if balance.Cmp(toWithdraw) == -1 {
		w.WriteHeader(400)
		w.Write([]byte("Insufficient funds in PLCR contract"))
		return
	}

	tx, err := contracts.PLCRWithdraw(account.MultisigAddress.String, toWithdraw)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(tx.Hash().Hex()))
}

func (server *Server) voteDeposit(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	account, err := getAccountFromRequestID(buf.String())
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	if !account.MultisigAddress.Valid {
		w.WriteHeader(400)
		w.Write([]byte("Account does not have a multisig address"))
		return
	}

	amount, err := getIntArgFromRequest(buf.String())
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	toDeposit := contracts.GetAtomicTokenAmount(amount)
	balance, err := contracts.GetTokenBalance(account.MultisigAddress.String)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	if balance.Cmp(toDeposit) == -1 {
		w.WriteHeader(400)
		w.Write([]byte("Not enough tokens in wallet"))
		return
	}

	tx, err := contracts.PLCRDeposit(account.MultisigAddress.String, toDeposit)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(tx.Hash().Hex()))
}

func (server *Server) activate(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	id, err := strconv.ParseInt(buf.String(), 10, 64)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	account, err := models.FindAccountByID(id)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	// Get their wallet address
	if account.MultisigAddress == nil || (account.MultisigAddress != nil && !account.MultisigAddress.Valid) {
		w.WriteHeader(400)
		w.Write([]byte("User does not have multisig wallet"))
		return
	}

	balance, err := contracts.GetTokenBalance(account.MultisigAddress.String)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Balance error: " + err.Error()))
		return
	}

	// Should we give them more tokens?
	initialAmountInt, err := strconv.ParseInt(os.Getenv("INITIAL_DISTRIBUTION_AMOUNT"), 10, 64)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}
	diff := contracts.GetAtomicTokenAmount(initialAmountInt)
	diff = diff.Sub(balance, diff)

	if diff.Cmp(big.NewInt(0)) == -1 {
		diff = diff.Mul(diff, big.NewInt(-1))
		log.Printf("%s has less than the initial distribution amount. Minting %d tokens!", account.TwitterHandle, diff.Int64())

		tx, err := contracts.MintTokens(account.MultisigAddress.String, diff)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Error minting tokens: " + err.Error()))
			return
		}

		if _, err = contracts.AwaitTransactionConfirmation(tx.Hash()); err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Error awaiting mint tx: " + err.Error()))
			return
		}

		// Refresh their balance
		balance, err = contracts.GetTokenBalance(account.MultisigAddress.String)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Balance error: " + err.Error()))
			return
		}
	}

	// Make sure they have enough tokens to vote
	plcrBalance, err := contracts.PLCRFetchBalance(account.MultisigAddress.String)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	diff = contracts.GetAtomicTokenAmount(50)
	diff = diff.Sub(plcrBalance, diff)
	if diff.Cmp(big.NewInt(0)) == -1 {
		diff = diff.Mul(diff, big.NewInt(-1))
		log.Printf("%s has less than 50 tokens locked up for voting, depositing %d tokens!", account.TwitterHandle, diff.Int64())
		tx, err := contracts.PLCRDeposit(account.MultisigAddress.String, diff)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Error: " + err.Error()))
			return
		}

		if _, err = contracts.AwaitTransactionConfirmation(tx.Hash()); err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Error PLCR deposit tx: " + err.Error()))
			return
		}
	}

	// Activate their account
	err = account.MarkActivated()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Activation error: " + err.Error()))
		return
	}

	// Send 'em a message
	humanBalance := contracts.GetHumanTokenAmount(balance)
	err = twitter.SendDM(account.TwitterID, fmt.Sprintf(disbursementMsg, humanBalance))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Twitter error: " + err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("Ok"))
}

func renderServerError(err error, w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
}

func (server *Server) syncList(w http.ResponseWriter, r *http.Request) {
	if err := twitter.SyncList(); err != nil {
		renderServerError(err, w)
		return
	}

	w.Write([]byte("Ok"))
}

func (server *Server) showList(w http.ResponseWriter, r *http.Request) {
	listings, err := models.FindWhitelistedRegistryListings()
	if err != nil {
		renderServerError(err, w)
		return
	}

	jsonText, err := json.Marshal(listings)
	if err != nil {
		renderServerError(err, w)
		return
	}

	w.WriteHeader(200)
	w.Write(jsonText)
}

// StartServer spins up a webserver for the API
func StartServer(eventsChan chan<- *events.TwitterEvent, errChan chan<- error) *Server {
	server := &Server{
		eventsChan: eventsChan,
		errChan:    errChan,
	}

	http.HandleFunc("/webhooks/twitter", server.handleTwitterWebhook)
	http.HandleFunc("/api/list", server.showList)
	http.HandleFunc("/admin/sync-list", requireAuth(server.syncList))
	http.HandleFunc("/admin/create-webhook", requireAuth(server.createWebhook))
	http.HandleFunc("/admin/authenticate-vip", requireAuth(server.authenticateVIP))
	http.HandleFunc("/admin/authenticate-party", requireAuth(server.authenticateParty))
	http.HandleFunc("/admin/redeploy-wallet", requireAuth(server.redeployWallet))
	http.HandleFunc("/admin/activate", requireAuth(server.activate))
	http.HandleFunc("/admin/wallet-balance", requireAuth(server.walletBalance))
	http.HandleFunc("/admin/vote-balance", requireAuth(server.voteBalance))
	http.HandleFunc("/admin/vote-deposit", requireAuth(server.voteDeposit))
	http.HandleFunc("/admin/vote-withdraw", requireAuth(server.voteWithdraw))

	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		errChan <- err
	}

	return server
}
