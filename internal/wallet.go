package internal

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tjvr/go-monzo"
	"golang.org/x/oauth2"
	"net/http"
	"os/exec"
	"runtime"
)

type Account struct {
	ID            string
	Created       string
	AccountNumber string
	Balance       float64
	Transactions  []*Transaction
}

type Transaction struct {
	ID       string
	Amount   float64
	Created  string
	Merchant string
}

type Accounts []*Account

type Wallet struct {
	cl              *monzo.Client
	accounts        Accounts
	SelectedAccount *Account
}

func (w *Wallet) LoadedWallet() bool {
	return w.cl != nil && w.accounts != nil
}

func (w *Wallet) Shutdown() {
	w.cl = nil
	w.accounts = nil
}

func (w *Wallet) Connect(conf *oauth2.Config) (*oauth2.Token, error) {

	respCh := make(chan string)

	srv := http.Server{Addr: ":8080"}
	srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			respCh <- code
		}
		w.WriteHeader(200)
	})

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			logrus.Info(err)
		}
	}()
	defer func(srv *http.Server, ctx context.Context) {
		err := srv.Shutdown(ctx)
		if err != nil {
			logrus.Info(err)
		}
	}(&srv, context.Background())

	url := conf.AuthCodeURL(generateRandomState())
	if err := openbrowser(url); err != nil {
		return nil, err
	}

	return conf.Exchange(context.Background(), <-respCh)

}

func (w *Wallet) FetchAccounts(token string) error {
	w.cl = &monzo.Client{
		BaseURL:     "https://api.monzo.com",
		AccessToken: token,
	}
	accounts, err := w.cl.Accounts("uk_retail")
	if err != nil {
		return err
	}

	w.accounts = make([]*Account, 0, len(accounts))

	for _, account := range accounts {

		balance, err := w.cl.Balance(account.ID)
		if err != nil {
			return err
		}

		mTransactions, err := w.cl.Transactions(account.ID, false)
		if err != nil {
			return err
		}

		transactions := make([]*Transaction, 0, len(mTransactions))
		for _, transaction := range mTransactions {
			transactions = append(transactions, &Transaction{
				ID:       transaction.ID,
				Amount:   float64(transaction.Amount),
				Created:  transaction.Created,
				Merchant: transaction.Merchant.Name,
			})
		}

		w.accounts = append(w.accounts, &Account{
			ID:            account.ID,
			Created:       account.Created,
			AccountNumber: account.AccountNumber,
			Balance:       float64(balance.Balance),
			Transactions:  transactions,
		})
	}

	return nil
}

func (w *Wallet) AccountsList() Accounts {
	return w.accounts
}

func openbrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}

func generateRandomState() string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)

	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(randomBytes)
}
