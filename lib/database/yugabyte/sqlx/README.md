# yugate-crud-sqlx
gồm 2 file example là file example của nhà nsx và example_test là 1 file dùng sqlx để CRUD table TRANSACTION có NULL và JSONB

demonstrate how to select, update, delete an entity in db.
This entity has NULL field and jsonb type

sql to generate TRANSACTION
```
CREATE TABLE "napas-fast-fund"."TRANSACTION" (
	"TRANS_ID" varchar(128) NOT NULL,
	"REF_ID" varchar(128) NULL,
	"CLIENT_TRANS_ID" varchar(128) NOT NULL,
	"BRANCH_CODE" varchar(16) NOT NULL,
	"USER_ID" varchar(128) NOT NULL,
	"CLIENT_CHANNEL" varchar(128) NOT NULL,
	"TRANS_TYPE" varchar(128) NOT NULL,
	"CREATED_AT" timestamp NOT NULL,
	"STATUS" numeric(3) NOT NULL,
	"STATUS_DESC" varchar(256) NOT NULL,
	"STATE" varchar(128) NOT NULL,
	"PROCESSING_STATUS" numeric(1) NULL,
	"LAST_UPDATED" timestamp NULL,
	"REPLY_TO" varchar(256) NOT NULL,
	"TRACE" jsonb NULL,
	CONSTRAINT transaction_pk PRIMARY KEY ("TRANS_ID"),
	CONSTRAINT transaction_un UNIQUE ("CLIENT_TRANS_ID")
);
```
## How to run
```
go test -v
```