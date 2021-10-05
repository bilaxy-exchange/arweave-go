package transactor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/bilaxy-exchange/arweave-go"
	"github.com/bilaxy-exchange/arweave-go/api"
	"github.com/bilaxy-exchange/arweave-go/tx"
)

// defaultPort of the arweave client
const defaultPort = "1984"

// defaultURL is the local host url
const defaultURL = "http://127.0.0.1" + ":" + defaultPort

// ClientCaller is the base interface needed to create a Transactor
type ClientCaller interface {
	TxAnchor(ctx context.Context) (string, error)
	LastTransaction(ctx context.Context, address string) (string, error)
	GetReward(ctx context.Context, data []byte, target string) (string, error)
	Commit(ctx context.Context, data []byte) (string, error)
	GetTransaction(ctx context.Context, txID string) (*tx.Transaction, error)
}

// Transactor type, allows one to create transactions
type Transactor struct {
	Client ClientCaller
}

// NewTransactor creates a new arweave transactor. You need to pass in a context and a url
// If sending an empty string, the default url is localhosts
func NewTransactor(fullURL string) (*Transactor, error) {
	if fullURL == "" {
		c, err := api.Dial(defaultURL)
		if err != nil {
			return nil, err
		}
		return &Transactor{
			Client: c,
		}, nil
	}
	u, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}
	formattedURL := fullURL
	if u.Scheme == "" {
		formattedURL = fmt.Sprintf("http://%s:%s", formattedURL, defaultPort)
	}
	c, err := api.Dial(formattedURL)
	if err != nil {
		return nil, err
	}

	return &Transactor{
		Client: c,
	}, nil
}

// CreateTransaction creates a brand new transaction
func (tr *Transactor) CreateTransaction(ctx context.Context, w arweave.WalletSigner, amount string, data []byte, target string, minFee int64, maxFee int64, includeFee bool) (*tx.Transaction, error) {
	lastTx, err := tr.Client.TxAnchor(ctx)
	if err != nil {
		return nil, err
	}

	price, err := tr.Client.GetReward(ctx, []byte(data), target)
	if err != nil {
		return nil, err
	}

	fee, err := strconv.ParseInt(price, 10, 64)
	if err != nil {
		return nil, err
	}
	if fee < minFee {
		fee = minFee
	} else if fee > maxFee {
		return nil, fmt.Errorf("transfer fee exceeds limit -> %d ", maxFee)
	}

	if includeFee {
		amountInt, err := strconv.ParseInt(amount, 10, 64)
		if err != nil {
			return nil, err
		}
		amountInt = amountInt - fee
		amount = fmt.Sprintf("%d", amountInt)
	}

	// Non encoded transaction fields
	return tx.NewTransaction(
		lastTx,
		w.PubKeyModulus(),
		amount,
		target,
		data,
		fmt.Sprintf("%d", fee),
	), nil
}

func (tr *Transactor) CreateTransactionWithFee(ctx context.Context, w arweave.WalletSigner, amount string, data []byte, target string, fee int64, includeFee bool) (*tx.Transaction, error) {
	lastTx, err := tr.Client.TxAnchor(ctx)
	if err != nil {
		return nil, err
	}

	if includeFee {
		amountInt, err := strconv.ParseInt(amount, 10, 64)
		if err != nil {
			return nil, err
		}
		amountInt = amountInt - fee
		amount = fmt.Sprintf("%d", amountInt)
	}

	// Non encoded transaction fields
	return tx.NewTransaction(
		lastTx,
		w.PubKeyModulus(),
		amount,
		target,
		data,
		fmt.Sprintf("%d", fee),
	), nil
}

// SendTransaction formats the transactions (base64url encodes the necessary fields)
// marshalls the Json and sends it to the arweave network
func (tr *Transactor) SendTransaction(ctx context.Context, tx *tx.Transaction) (string, error) {
	if len(tx.Signature()) == 0 {
		return "", errors.New("transaction missing signature")
	}
	serialized, err := json.Marshal(tx)
	if err != nil {
		return "", err
	}
	return tr.Client.Commit(ctx, serialized)
}

// WaitMined waits for the transaction to be mined
func (tr *Transactor) WaitMined(ctx context.Context, tx *tx.Transaction) (*tx.Transaction, error) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		receipt, err := tr.Client.GetTransaction(ctx, tx.Hash())
		if receipt != nil {
			return receipt, nil
		}
		if err != nil {
			fmt.Printf("Error retrieving transaction %s \n", err.Error())
		} else {
			fmt.Printf("Transaction not yet mined \n")
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}
