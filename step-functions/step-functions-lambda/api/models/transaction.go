package models

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var (
	REJECTED = "REJECTED"
	PURCHASE = "PURCHASE"
)

type Transaction struct {
	ID              string `json:"id"`
	CustomerId      string `json:"customerId"`
	CardNumber      string `json:"card_number"`
	TransactionType string `json:"transaction_type"`
	Timestamp       string `json:"timestamp"`
}

func (t *Transaction) SetDefaultTransactionValues() {
	t.ID = fmt.Sprint(int(time.Now().Unix()) + rand.Intn(12345) + 654321)
	t.Timestamp = time.Now().Format("2006-01-02 15:04:05 Monday")
}

func (t *Transaction) Validate() error {
	var (
		err_msg []string
	)

	if t.CustomerId == "" {
		err_msg = append(err_msg, "'customerId'")
	}

	if t.CardNumber == "" {
		err_msg = append(err_msg, "'card_number'")
	}

	if t.TransactionType == "" {
		err_msg = append(err_msg, "'transaction_type'")
	}

	if len(err_msg) > 0 {
		msg := fmt.Sprintf("missing %s field(s).", strings.Join(err_msg, ", "))
		invalidType := t.validateTransactionType()

		if len(invalidType) > 0 {
			msg += fmt.Sprintf(" %s", invalidType)
			return errors.New(msg)
		}

		return errors.New(msg)
	}

	if len(t.validateTransactionType()) > 0 {
		return errors.New(t.validateTransactionType())
	}

	return nil
}

func (t *Transaction) validateTransactionType() string {
	if t.TransactionType != PURCHASE && t.TransactionType != REJECTED {
		return fmt.Sprintf("incorrect 'transaction_type'. received %s", t.TransactionType)
	}

	return ""
}
