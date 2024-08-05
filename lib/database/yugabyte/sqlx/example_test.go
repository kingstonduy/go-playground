package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-playground/utils"
	"log"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	_ "github.com/yugabyte/pgx/v5/stdlib"
)

type Trace struct {
	From               string `json:"frm" xml:"frm"`
	To                 string `json:"to" xml:"to"`
	Cid                string `json:"cid" xml:"cid" validate:"required"`
	Sid                string `json:"sid" xml:"sid"`
	Cts                int64  `json:"cts" xml:"cts" validate:"required"`
	Sts                int64  `json:"sts" xml:"sts"`
	Dur                int64  `json:"dur" xml:"dur"`
	Username           string `json:"userName" xml:"userName"`
	ReplyTo            string `json:"replyTo" xml:"replyTo"`
	TransactionTimeout int64  `json:"transactionTimeout"`
}

type Transaction struct {
	TransID          *string    `db:"TRANS_ID" json:"trans_id"`                   //Unique key, định danh của giao dịch phát sinh bởi phía server.
	RefID            *string    `db:"REF_ID" json:"ref_id"`                       //Khóa tham chiếu. trance.sid
	ClientTransID    *string    `db:"CLIENT_TRANS_ID" json:"client_trans_id"`     //Định danh của giao dịch phát sinh bởi phía client. trace.cid
	BranchCode       *string    `db:"BRANCH_CODE" json:"branch_code"`             //Mã chi nhánh ghi ngân hàng ghi nhận cho bút toán.
	UserID           *string    `db:"USER_ID" json:"user_id"`                     //Định danh user
	ClientChannel    *string    `db:"CLIENT_CHANNEL" json:"client_channel"`       //Kênh giao dịch phía client.(OMNIW,OMNIA,…)
	TransType        *string    `db:"TRANS_TYPE" json:"trans_type"`               //Loại giao dịch.(FastTransfer)
	CreatedAt        *time.Time `db:"CREATED_AT" json:"created_at"`               //Thời gian tạo sự kiện, quan trọng cho việc xử lý theo thứ tự và khắc phục lỗi. Giúp xác định thời điểm xảy ra sự kiện và cho phép các dịch vụ xử lý theo trình tự thời gian.
	Status           *int       `db:"STATUS" json:"status"`                       //Trạng thái của giao dịch.(0: Initiated (Khởi tạo); 1: Success (Thành công); 2: Timeout (Nghi vấn); 3: Failed (Thất bại))
	StatusDesc       *string    `db:"STATUS_DESC" json:"status_desc"`             //Mô tả cho trạng thái của giao dịch.
	State            *string    `db:"STATE" json:"state"`                         //Tình trạng giao dịch(COMMAND_TYPE cuối)
	ProcessingStatus *int       `db:"PROCESSING_STATUS" json:"processing_status"` //Cờ xác định trạng thái giao dịch có đang được xử lý hay không. Giúp hệ thống xử lý giao dịch theo cơ chế Optimistic locking.
	LastUpdated      *time.Time `db:"LAST_UPDATED" json:"last_updated"`           //Thời gian giao dịch được update sau cùng. Giúp hệ thống xử lý giao dịch theo cơ chế Optimistic locking.
	ReplyTo          *string    `db:"REPLY_TO" json:"reply_to"`                   //Topic cần trả về
	Trace            *Trace     `db:"TRACE" json:"trace"`                         //trace cua transaction
}

func NewTransaction() Transaction {
	s := utils.NewString(3)
	i := int(2)
	t := time.Now()
	return Transaction{
		TransID:          &s,
		RefID:            &s,
		ClientTransID:    &s,
		BranchCode:       &s,
		UserID:           &s,
		ClientChannel:    &s,
		TransType:        &s,
		CreatedAt:        &t,
		Status:           &i,
		StatusDesc:       &s,
		State:            &s,
		ProcessingStatus: &i,
		LastUpdated:      &t,
		ReplyTo:          &s,
		Trace: &Trace{
			From:               "",
			To:                 "",
			Cid:                "",
			Sid:                "",
			Cts:                time.Now().UnixMilli(),
			Sts:                time.Now().UnixMilli(),
			Dur:                0,
			Username:           "",
			ReplyTo:            "",
			TransactionTimeout: 0,
		},
	}
}

func insertTransaction(ctx context.Context, db *sqlx.DB, transaction Transaction) (err error) {
	sqlString :=
		`
		insert into "napas-fast-fund"."TRANSACTION" (
			"TRANS_ID",
			"REF_ID",
			"CLIENT_TRANS_ID",
			"BRANCH_CODE",
			"USER_ID",
			"CLIENT_CHANNEL",
			"TRANS_TYPE",
			"CREATED_AT",
			"STATUS",
			"STATUS_DESC",
			"STATE",
			"PROCESSING_STATUS",
			"LAST_UPDATED",
			"REPLY_TO",
			"TRACE"
			)
			values(
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
			)
		`

	if _, err = db.Exec(
		sqlString,
		transaction.TransID,
		transaction.RefID,
		transaction.ClientTransID,
		transaction.BranchCode,
		transaction.UserID,
		transaction.ClientChannel,
		transaction.TransType,
		transaction.CreatedAt,
		transaction.Status,
		transaction.StatusDesc,
		transaction.State,
		transaction.ProcessingStatus,
		transaction.LastUpdated,
		transaction.ReplyTo,
		transaction.Trace,
	); err != nil {
		err = errors.Wrap(err, "cant insert table transaction")
		return err
	}

	return nil
}

func updateTransaction(ctx context.Context, db *sqlx.DB, transaction Transaction) (err error) {
	sqlString := `
		UPDATE "napas-fast-fund"."TRANSACTION"
		SET
			"REF_ID" = COALESCE($1, "REF_ID"),
			"CLIENT_TRANS_ID" = COALESCE($2, "CLIENT_TRANS_ID"),
			"BRANCH_CODE" = COALESCE($3, "BRANCH_CODE"),
			"USER_ID" = COALESCE($4, "USER_ID"),
			"CLIENT_CHANNEL" = COALESCE($5, "CLIENT_CHANNEL"),
			"TRANS_TYPE" = COALESCE($6, "TRANS_TYPE"),
			"CREATED_AT" = COALESCE($7, "CREATED_AT"),
			"STATUS" = COALESCE($8, "STATUS"),
			"STATUS_DESC" = COALESCE($9, "STATUS_DESC"),
			"STATE" = COALESCE($10, "STATE"),
			"PROCESSING_STATUS" = COALESCE($11, "PROCESSING_STATUS"),
			"LAST_UPDATED" = COALESCE($12, "LAST_UPDATED"),
			"REPLY_TO" = COALESCE($13, "REPLY_TO"),
			"TRACE" = COALESCE($14, "TRACE")
		WHERE "TRANS_ID" = $15
	`

	_, err = db.Exec(
		sqlString,
		transaction.RefID,
		transaction.ClientTransID,
		transaction.BranchCode,
		transaction.UserID,
		transaction.ClientChannel,
		transaction.TransType,
		transaction.CreatedAt,
		transaction.Status, // Make sure this is a pointer to an int
		transaction.StatusDesc,
		transaction.State,
		transaction.ProcessingStatus, // Make sure this is a pointer to an int
		transaction.LastUpdated,
		transaction.ReplyTo,
		transaction.Trace,
		transaction.TransID,
	)

	if err != nil {
		return errors.Wrap(err, "can not update table TRANSACTION")
	}

	return nil
}

func deleteTransaction(ctx context.Context, db *sqlx.DB, refID string) (err error) {
	sqlString :=
		`
		DELETE FROM "napas-fast-fund"."TRANSACTION" WHERE "TRANS_ID" =$1	
	`
	_, err = db.Exec(
		sqlString,
		refID,
	)

	if err != nil {
		return errors.Wrap(err, "can not delete TRANSACTION_DETAIL by REF_ID")
	}

	return nil
}

func SelectTransactionByTransID(ctx context.Context, db *sqlx.DB, transID string) (transaction Transaction, err error) {
	sqlString :=
		`
		SELECT "TRANS_ID", "REF_ID", "CLIENT_TRANS_ID", "BRANCH_CODE", "USER_ID", "CLIENT_CHANNEL", "TRANS_TYPE", "CREATED_AT", "STATUS", "STATUS_DESC", "STATE", "PROCESSING_STATUS", "LAST_UPDATED", "REPLY_TO", "TRACE"
FROM "napas-fast-fund"."TRANSACTION" where "TRANS_ID" = $1;
	`

	transaction = Transaction{}
	var trace json.RawMessage

	row := db.QueryRow(sqlString, transID)

	if err = row.Scan(
		&transaction.TransID,
		&transaction.RefID,
		&transaction.ClientTransID,
		&transaction.BranchCode,
		&transaction.UserID,
		&transaction.ClientChannel,
		&transaction.TransType,
		&transaction.CreatedAt,
		&transaction.Status,
		&transaction.StatusDesc,
		&transaction.State,
		&transaction.ProcessingStatus,
		&transaction.LastUpdated,
		&transaction.ReplyTo,
		&trace,
	); err != nil {
		err = errors.Wrap(err, "can not get TRANSACTION by transID")
		return transaction, err
	}

	if err = json.Unmarshal(trace, &transaction.Trace); err != nil {
		err = errors.Wrap(err, "can not marshal trace")
		return transaction, err
	}

	return transaction, nil
}

func InitDb() *sqlx.DB {
	const (
		host              = "10.96.20.152"
		port              = "5433"
		user              = "postgres"
		password          = "postgres"
		dbname            = "postgres"
		maxConn           = 30
		healthCheckPeriod = 1 * time.Minute
		maxConnIdleTime   = 1 * time.Minute
		maxConnLifetime   = 3 * time.Minute
		minConns          = 10
		lazyConnect       = false
	)

	db, err := sqlx.Connect("pgx", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname))
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

func Test(t *testing.T) {
	ctx := context.Background()
	db := InitDb()
	_ = db

	transaction := NewTransaction()

	if err := insertTransaction(ctx, db, transaction); err != nil {
		panic(err)
	}

	transaction, err := SelectTransactionByTransID(ctx, db, *transaction.TransID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("before update: %+v\n", transaction)

	newTransaction := NewTransaction()
	newTransaction.TransID = transaction.TransID

	if err := updateTransaction(ctx, db, newTransaction); err != nil {
		panic(err)
	}

	transaction1, err := SelectTransactionByTransID(ctx, db, *transaction.TransID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("after update: %+v\n", transaction1)

	if err := deleteTransaction(ctx, db, *transaction.RefID); err != nil {
		panic(err)
	}

	transaction2, err := SelectTransactionByTransID(ctx, db, *transaction.TransID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			t.Log("Transaction not found after deletion, as expected.")
		} else {
			panic(err) // Panic on unexpected error
		}
	}
	fmt.Printf("after delete: %+v\n", transaction2)
}
