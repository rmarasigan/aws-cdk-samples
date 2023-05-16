package schema

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

type Customer struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type Transaction struct {
	ID              string   `json:"id,omitempty"`
	Customer        Customer `json:"customer"`
	TransactionType string   `json:"transaction_type"`
	Timestamp       string   `json:"timestamp,omitempty"`
}

func (t *Transaction) SetDefaultValues() {
	t.ID = fmt.Sprint(int(time.Now().Unix()) + rand.Intn(12345) + 654321)
	t.Timestamp = time.Now().Format("2006-01-02 15:04:05 Monday")
}

// IsValid validates the transaction details of a customer.
func (t *Transaction) IsValid() error {
	var params []string

	if t.Customer.FirstName == "" {
		params = append(params, "customer.first_name")
	}

	if t.Customer.LastName == "" {
		params = append(params, "customer.last_name")
	}

	if t.Customer.Email == "" {
		params = append(params, "customer.email")
	}

	if t.TransactionType == "" {
		params = append(params, "transaction_type")
	}

	if len(params) > 0 {
		msg := fmt.Sprintf("missing %s field(s.)", strings.Join(params, ", "))
		invalidType := t.validateTransactionType()

		if invalidType != "" {
			msg += fmt.Sprintf(" %s", invalidType)
			return errors.New(msg)
		}

		return errors.New(msg)
	}

	return nil
}

func (t *Transaction) validateTransactionType() string {
	if t.TransactionType != PURCHASE && t.TransactionType != REJECTED {
		return fmt.Sprintf("incorrect 'transaction_type'. received %s", t.TransactionType)
	}

	return ""
}

func (t *Transaction) EmailContent() (string, string) {
	switch t.TransactionType {
	case PURCHASE:
		subject := fmt.Sprintf("COD #%s APPROVED", t.ID)
		body := fmt.Sprintf("Hello %s %s,\nYour Cash on Delivery (COD) payment request for order #%s has been approved. We have notified the seller to start shipping your item(s).\n\n\tORDER DETAILS\n\t\tOrder ID:\t#%s\n\t\tOrder Date:\t%s\n\t\tSeller:\tABC Company Inc.\n\nKindly wait for your shipment and have your cash ready to pay for your products upon delivery. Once you receive and accept the product(s), kindly confirm this with us in ABC App.", t.Customer.FirstName, t.Customer.LastName, t.ID, t.ID, t.Timestamp)

		return subject, body

	default:
		subject := fmt.Sprintf("COD #%s REJECTED", t.ID)
		body := fmt.Sprintf("Hello %s %s, \n\nYour Cash on Delivery (COD) payment request for order #%s has been rejected. We have notified the seller to cancel the shipping of your items.", t.Customer.FirstName, t.Customer.LastName, t.ID)

		return subject, body
	}
}
