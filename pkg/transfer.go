package moov

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TransferStatus int

const (
	TransferStatusCreated TransferStatus = iota
	TransferStatusPending
	TransferStatusCompleted
	TransferStatusFailed
	TransferStatusReversed
	TransferStatusQueued
	TransferStatusCanceled
)

var TransferStatusStrings = map[TransferStatus]string{
	TransferStatusCreated:   "created",
	TransferStatusPending:   "pending",
	TransferStatusCompleted: "completed",
	TransferStatusFailed:    "failed",
	TransferStatusReversed:  "reversed",
	TransferStatusQueued:    "queued",
	TransferStatusCanceled:  "canceled",
}

type SynchronousTransfer struct {
	TransferID     string            `json:"transferID,omitempty"`
	CreatedOn      time.Time         `json:"createdOn,omitempty"`
	CompletedOn    time.Time         `json:"completedOn,omitempty"`
	Status         string            `json:"status,omitempty"`
	FailureReason  string            `json:"failureReason,omitempty"`
	Amount         Amount            `json:"amount,omitempty"`
	Description    string            `json:"description,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	FacilitatorFee FacilitatorFee    `json:"facilitatorFee,omitempty"`
	MoovFee        int               `json:"moovFee,omitempty"`
	MoovFeeDecimal string            `json:"moovFeeDecimal,omitempty"`
	MoovFeeDetails MoovFeeDetails    `json:"moovFeeDetails,omitempty"`
	GroupID        string            `json:"groupID,omitempty"`
	RefundedAmount Amount            `json:"refundedAmount,omitempty"`
	Refunds        []Refund          `json:"refunds,omitempty"`
	DisputedAmount Amount            `json:"disputedAmount,omitempty"`
	Disputes       []Dispute         `json:"disputes,omitempty"`
	Source         Source            `json:"source,omitempty"`
	Destination    Destination       `json:"destination,omitempty"`
}

type AsynchronousTransfer struct {
	TransferID string    `json:"transferID,omitempty"`
	CreatedOn  time.Time `json:"createdOn,omitempty"`
}

type Amount struct {
	Currency string `json:"currency,omitempty"`
	Value    int    `json:"value,omitempty"`
}

type FacilitatorFee struct {
	Total         int    `json:"total,omitempty"`
	TotalDecimal  string `json:"totalDecimal,omitempty"`
	Markup        int    `json:"markup,omitempty"`
	MarkupDecimal string `json:"markupDecimal,omitempty"`
}

type MoovFeeDetails struct {
	CardScheme     string `json:"cardScheme,omitempty"`
	Interchange    string `json:"interchange,omitempty"`
	MoovProcessing string `json:"moovProcessing,omitempty"`
}

type Refund struct {
	RefundID    string      `json:"refundID,omitempty"`
	CreatedOn   time.Time   `json:"createdOn,omitempty"`
	UpdatedOn   time.Time   `json:"updatedOn,omitempty"`
	Status      string      `json:"status,omitempty"`
	FailureCode string      `json:"failureCode,omitempty"`
	Amount      Amount      `json:"amount,omitempty"`
	CardDetails CardDetails `json:"cardDetails,omitempty"`
}

type Source struct {
	PaymentMethodID   string          `json:"paymentMethodID,omitempty"`
	PaymentMethodType string          `json:"paymentMethodType,omitempty"`
	Account           TransferAccount `json:"account,omitempty"`
	BankAccount       BankAccount     `json:"bankAccount,omitempty"`
	Wallet            Wallet          `json:"wallet,omitempty"`
	Card              Card            `json:"card,omitempty"`
	ApplePay          ApplePay        `json:"applePay,omitempty"`
	AchDetails        AchDetails      `json:"achDetails,omitempty"`
	CardDetails       CardDetails     `json:"cardDetails,omitempty"`
	TransferID        string          `json:"transferID,omitempty"`
}

type TransferAccount struct {
	AccountID   string `json:"accountID,omitempty"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type Destination struct {
	PaymentMethodID   string          `json:"paymentMethodID,omitempty"`
	PaymentMethodType string          `json:"paymentMethodType,omitempty"`
	Account           TransferAccount `json:"account,omitempty"`
	BankAccount       BankAccount     `json:"bankAccount,omitempty"`
	Wallet            Wallet          `json:"wallet,omitempty"`
	Card              Card            `json:"card,omitempty"`
	ApplePay          ApplePay        `json:"applePay,omitempty"`
	AchDetails        AchDetails      `json:"achDetails,omitempty"`
	CardDetails       CardDetails     `json:"cardDetails,omitempty"`
}

type SearchQueryPayload struct {
	AccountIDs    []string  `json:"accountIDs,omitempty"`
	Status        string    `json:"status,omitempty"`
	StartDateTime time.Time `json:"startDateTime,omitempty"`
	EndDateTime   time.Time `json:"endDateTime,omitempty"`
	GroupID       string    `json:"groupID,omitempty"`
	Count         int       `json:"count,omitempty"`
	Skip          int       `json:"skip,omitempty"`
	Refunded      bool      `json:"refunded,omitempty"`
	Disputed      bool      `json:"disputed,omitempty"`
}

type MetaDataPayload struct {
	Metadata map[string]string `json:"metadata,omitempty"`
}

type RefundPayload struct {
	Amount int `json:"amount,omitempty"`
}

type TransferOptionsSourcePayload struct {
	PaymentMethodID string `json:"paymentMethodID,omitempty"`
	AccountID       string `json:"accountID,omitempty"`
}

type TransferOptionsDestinationPayload struct {
	PaymentMethodID string `json:"paymentMethodID,omitempty"`
	AccountID       string `json:"accountID,omitempty"`
}

type TransferOptionsPayload struct {
	Source      TransferOptionsSourcePayload      `json:"source,omitempty"`
	Destination TransferOptionsDestinationPayload `json:"destination,omitempty"`
	Amount      Amount                            `json:"amount,omitempty"`
}

type CreatedTransferOptions struct {
	SourceOptions      []Source `json:"sourceOptions,omitempty"`
	DestinationOptions []Source `json:"destinationOptions,omitempty"`
}

type RefundStatus struct {
	Status    string    `json:"status,omitempty"`
	CreatedOn time.Time `json:"createdOn,omitempty"`
}

type CanceledTransfer struct {
	Cancellation RefundStatus `json:"cancellation,omitempty"`
	Refund       Refund       `json:"refund,omitempty"`
}

type CreateTransfer struct {
	Source         Source            `json:"source,omitempty"`
	Destination    Destination       `json:"destination,omitempty"`
	Amount         Amount            `json:"amount,omitempty"`
	FacilitatorFee FacilitatorFee    `json:"facilitatorFee,omitempty"`
	Description    string            `json:"description,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// CreateTransfer creates a new transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/createTransfer
func (c Client) CreateTransfer(ctx context.Context, transfer CreateTransfer, isSync bool) (*SynchronousTransfer, *AsynchronousTransfer, error) {
	args := []callArg{AcceptJson(), JsonBody(transfer)}
	if isSync {
		args = append(args, WaitFor("rail-response"))
	}

	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathTransfers), args...)
	if err != nil {
		return nil, nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		st, err := UnmarshalObjectResponse[SynchronousTransfer](resp)
		return st, nil, err
	case StatusStarted:
		st, err := UnmarshalObjectResponse[AsynchronousTransfer](resp)
		return nil, st, err
	case StatusStateConflict:
		return nil, nil, ErrXIdempotencyKey
	default:
		return nil, nil, resp.Error()
	}
}

// ListTransfers lists all transfers
// https://docs.moov.io/api/index.html#tag/Transfers/operation/listTransfers
func (c Client) ListTransfers(payload SearchQueryPayload) ([]SynchronousTransfer, error) {
	var respTransfers []SynchronousTransfer

	values := url.Values{}
	// Convert time values to ISO8601 format
	startDateTimeStr := payload.StartDateTime.Format(time.RFC3339)
	endDateTimeStr := payload.EndDateTime.Format(time.RFC3339)

	// Add non-empty fields to the query string
	if len(payload.AccountIDs) > 0 {
		values.Add("accountIDs", strings.Join(payload.AccountIDs, ","))
	}
	if payload.Status != "" {
		values.Add("status", payload.Status)
	}
	if !payload.StartDateTime.IsZero() {
		values.Add("startDateTime", startDateTimeStr)
	}
	if !payload.EndDateTime.IsZero() {
		values.Add("endDateTime", endDateTimeStr)
	}
	if payload.GroupID != "" {
		values.Add("groupID", payload.GroupID)
	}
	if payload.Count > 0 {
		values.Add("count", fmt.Sprint(payload.Count))
	}
	if payload.Skip > 0 {
		values.Add("skip", fmt.Sprint(payload.Skip))
	}
	if payload.Refunded {
		values.Add("refunded", "true")
	}
	if payload.Disputed {
		values.Add("disputed", "true")
	}

	urlStr := fmt.Sprintf("%s/%s?%s", baseURL, pathTransfers, values.Encode())

	body, statusCode, err := c.GetHTTPResponse(http.MethodGet, urlStr, nil, nil)
	if err != nil {
		return respTransfers, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &respTransfers)
		if err != nil {
			return respTransfers, err
		}
		return respTransfers, nil
	case http.StatusTooManyRequests:
		return respTransfers, ErrRateLimit
	}
	return respTransfers, ErrDefault(statusCode)
}

// GetTransfer retrieves a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/getTransfer
func (c Client) GetTransfer(transferID string, accountID string) (SynchronousTransfer, error) {
	var respTransfer SynchronousTransfer

	values := url.Values{}
	if accountID != "" {
		values.Add("accountID", accountID)
	}

	urlStr := fmt.Sprintf("%s/%s/%s?%s", baseURL, pathTransfers, transferID, values.Encode())

	body, statusCode, err := c.GetHTTPResponse(http.MethodGet, urlStr, nil, nil)
	if err != nil {
		return respTransfer, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &respTransfer)
		if err != nil {
			return respTransfer, err
		}
		return respTransfer, nil
	case http.StatusTooManyRequests:
		return respTransfer, ErrRateLimit
	}
	return respTransfer, ErrDefault(statusCode)
}

// UpdateTransferMetaData updates the metadata for a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/patchTransfer
func (c Client) UpdateTransferMetaData(transferID string, accountID string, metadata map[string]string) (SynchronousTransfer, error) {
	var respTransfer SynchronousTransfer

	values := url.Values{}
	if accountID != "" {
		values.Add("accountID", accountID)
	}
	urlStr := fmt.Sprintf("%s/%s/%s?%s", baseURL, pathTransfers, transferID, values.Encode())
	metaDataPayload := MetaDataPayload{
		Metadata: metadata,
	}

	body, statusCode, err := c.GetHTTPResponse(http.MethodPatch, urlStr, metaDataPayload, nil)
	if err != nil {
		return respTransfer, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &respTransfer)
		if err != nil {
			return respTransfer, err
		}
		return respTransfer, nil
	case http.StatusTooManyRequests:
		return respTransfer, ErrRateLimit
	}
	return respTransfer, ErrDefault(statusCode)
}

// TransferOptions lists all transfer options between a source and destination
// https://docs.moov.io/api/#tag/Transfers/operation/createTransferOptions
func (c Client) TransferOptions(payload TransferOptionsPayload) (CreatedTransferOptions, error) {
	var respOptions CreatedTransferOptions
	urlStr := fmt.Sprintf("%s/%s", baseURL, pathTransferOptions)

	body, statusCode, err := c.GetHTTPResponse(http.MethodPost, urlStr, payload, nil)
	if err != nil {
		return respOptions, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &respOptions)
		if err != nil {
			return respOptions, err
		}
		return respOptions, nil
	case http.StatusTooManyRequests:
		return respOptions, ErrRateLimit
	}
	return respOptions, ErrDefault(statusCode)
}

// RefundTransfer refunds a transfer
// https://docs.moov.io/api/#tag/Transfers/operation/refundTransfer
func (c Client) RefundTransfer(transferID string, isSync bool, amount int) (Refund, error) {
	var respRefund Refund

	urlStr := fmt.Sprintf("%s/%s/%s/refunds", baseURL, pathTransfers, transferID)
	uuidV4 := uuid.NewString()
	header := map[string]string{"X-Idempotency-Key": uuidV4}
	if isSync {
		header["X-Wait-For"] = "rail-response"
	}

	refundPayload := RefundPayload{
		Amount: amount,
	}
	body, statusCode, err := c.GetHTTPResponse(http.MethodPost, urlStr, refundPayload, header)

	if err != nil {
		return respRefund, err
	}

	switch statusCode {
	case http.StatusOK, http.StatusAccepted:
		err = json.Unmarshal(body, &respRefund)
		if err != nil {
			return respRefund, err
		}
		return respRefund, nil
	case http.StatusBadRequest:
		var err error
		_ = json.Unmarshal(body, &err)
		return respRefund, err
	case http.StatusConflict:
		return respRefund, ErrXIdempotencyKey
	case http.StatusUnprocessableEntity:
		return respRefund, ErrRequestBody
	case http.StatusTooManyRequests:
		return respRefund, ErrRateLimit
	}
	return respRefund, ErrDefault(statusCode)
}

// ListRefunds lists all refunds for a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/getRefunds
func (c Client) ListRefunds(transferID string) ([]Refund, error) {
	var respRefunds []Refund

	urlStr := fmt.Sprintf("%s/%s/%s/refunds", baseURL, pathTransfers, transferID)

	body, statusCode, err := c.GetHTTPResponse(http.MethodGet, urlStr, nil, nil)
	if err != nil {
		return respRefunds, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &respRefunds)
		if err != nil {
			return respRefunds, err
		}
		return respRefunds, nil
	case http.StatusTooManyRequests:
		return respRefunds, ErrRateLimit
	}
	return respRefunds, ErrDefault(statusCode)
}

// GetRefund retrieves a refund for a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/getRefund
func (c Client) GetRefund(transferID string, refundID string) (Refund, error) {
	var respRefund Refund

	urlStr := fmt.Sprintf("%s/%s/%s/refunds/%s", baseURL, pathTransfers, transferID, refundID)

	body, statusCode, err := c.GetHTTPResponse(http.MethodGet, urlStr, nil, nil)
	if err != nil {
		return respRefund, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &respRefund)
		if err != nil {
			return respRefund, err
		}
		return respRefund, nil
	case http.StatusTooManyRequests:
		return respRefund, ErrRateLimit
	}
	return respRefund, ErrDefault(statusCode)
}

// ReverseTransfer reverses a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/reverseTransfer
func (c Client) ReverseTransfer(transferID string, amount int) (CanceledTransfer, error) {
	respTransfer := CanceledTransfer{}

	urlStr := fmt.Sprintf("%s/%s/%s/reversals", baseURL, pathTransfers, transferID)
	uuidV4 := uuid.NewString()
	header := map[string]string{"X-Idempotency-Key": uuidV4}

	refundPayload := RefundPayload{
		Amount: amount,
	}

	body, statusCode, err := c.GetHTTPResponse(http.MethodPost, urlStr, refundPayload, header)

	if err != nil {
		return respTransfer, err
	}

	switch statusCode {
	case http.StatusOK, http.StatusAccepted:
		err = json.Unmarshal(body, &respTransfer)
		if err != nil {
			return respTransfer, err
		}
		return respTransfer, nil
	case http.StatusBadRequest:
		var err error
		_ = json.Unmarshal(body, &err)
		return respTransfer, err
	case http.StatusConflict:
		return respTransfer, ErrXIdempotencyKey
	case http.StatusUnprocessableEntity:
		return respTransfer, ErrRequestBody
	case http.StatusTooManyRequests:
		return respTransfer, ErrRateLimit
	}
	return respTransfer, ErrDefault(statusCode)
}
