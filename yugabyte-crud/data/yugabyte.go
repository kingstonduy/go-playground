package data_yugabyte

import (
	"context"
	"yugabyte-crud/domain"

	"10.96.24.141/UDTN/integration/microservices/mcs-go/mcs-go-modules/mcs-go-core.git/database"
	"10.96.24.141/UDTN/integration/microservices/mcs-go/mcs-go-modules/mcs-go-core.git/logger"
)

type CRUDRepository struct {
	db *database.Gdbc
}

type YugabyteDatabase struct {
	DB *database.Gdbc
}

func NewCRUDRepository(db *YugabyteDatabase) domain.YugabyteRepository {
	return &CRUDRepository{
		db: db.DB,
	}
}

// Insert implements domain.YugabyteRepository.
func (c *CRUDRepository) Insert(ctx context.Context, entity domain.SmartlinkRetrieveTraceEntity) error {
	sql := `
	INSERT INTO "SMARTLINK_RETRIEVE_TRACE"(
		"FIELD_0", "FIELD_2", "FIELD_3", "FIELD_4", 
		"TRANS_ID", "CLIENT_ID"
		) VALUES
		($1,$2,$3,$4,$5,$6)
	`
	_, err := c.db.Exec(ctx, sql, entity.Field0, entity.Field2, entity.Field3, entity.Field4, entity.TransID, entity.ClientID)
	if err != nil {
		logger.Errorf(ctx, "errorf when insert db: %s", err)
	}
	return err
}

// Delete implements domain.YugabyteRepository.
func (c *CRUDRepository) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// Get implements domain.YugabyteRepository.
func (c *CRUDRepository) Get(ctx context.Context, id string) (domain.SmartlinkRetrieveTraceEntity, error) {
	panic("unimplemented")
}

// Update implements domain.YugabyteRepository.
func (c *CRUDRepository) Update(ctx context.Context, entity domain.SmartlinkRetrieveTraceEntity) error {
	panic("unimplemented")
}
