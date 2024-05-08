package domain

import "context"

type SmartlinkRetrieveTraceEntity struct {
	Field0           string `db:"FIELD_0"`
	Field2           string `db:"FIELD_2"`
	Field3           string `db:"FIELD_3"`
	Field4           string `db:"FIELD_4"`
	Field7           string `db:"FIELD_7"`
	Field11          string `db:"FIELD_11"`
	Field12          string `db:"FIELD_12"`
	Field13          string `db:"FIELD_13"`
	Field15          string `db:"FIELD_15"`
	Field18          string `db:"FIELD_18"`
	Field19          string `db:"FIELD_19"`
	Field22          string `db:"FIELD_22"`
	Field25          string `db:"FIELD_25"`
	Field32          string `db:"FIELD_32"`
	Field37          string `db:"FIELD_37"`
	Field38          string `db:"FIELD_38"`
	Field39          string `db:"FIELD_39"`
	Field41          string `db:"FIELD_41"`
	Field42          string `db:"FIELD_42"`
	Field43          string `db:"FIELD_43"`
	Field48          string `db:"FIELD_48"`
	Field49          string `db:"FIELD_49"`
	Field60          string `db:"FIELD_60"`
	Field62          string `db:"FIELD_62"`
	Field100         string `db:"FIELD_100"`
	Field102         string `db:"FIELD_102"`
	Field103         string `db:"FIELD_103"`
	Field104         string `db:"FIELD_104"`
	Field128OCB      string `db:"FIELD_128_OCB"`
	Channel          string `db:"CHANNEL"`
	TraceType        string `db:"TRACE_TYPE"`
	TransID          string `db:"TRANS_ID"`
	ClientID         string `db:"CLIENT_ID"`
	CardMasking      string `db:"CARDMASKING"`
	EncryptedKey     string `db:"ENCRYPTEDKEY"`
	EncryptedCardNo  string `db:"ENCRYPTEDCARDNO"`
	WorkflowID       string `db:"WORKFLOWID"`
	Status           string `db:"STATUS"`
	PrevStatus       string `db:"PREV_STATUS"`
	ProcessingStatus string `db:"PROCESSING_STATUS"`
	CreatedDate      string `db:"CREATED_DATE"`
	LastUpdated      string `db:"LAST_UPDATED"`
}

type YugabyteRepository interface {
	Insert(ctx context.Context, entity SmartlinkRetrieveTraceEntity) error
	Get(ctx context.Context, id string) (SmartlinkRetrieveTraceEntity, error)
	Update(ctx context.Context, entity SmartlinkRetrieveTraceEntity) error
	Delete(ctx context.Context, id string) error
}
