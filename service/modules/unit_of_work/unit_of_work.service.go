package unit_of_work

import (
	"context"
	"fmt"
	"hireforwork-server/db"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type UnitOfWork struct {
	DB      *db.DB
	changes []func(ctx mongo.SessionContext) error
}

func NewUnitOfWork(dbInstance *db.DB) *UnitOfWork {
	return &UnitOfWork{DB: dbInstance}
}

func (uow *UnitOfWork) RegisterChange(change func(ctx mongo.SessionContext) error) {
	uow.changes = append(uow.changes, change)
}

func (uow *UnitOfWork) Commit() error {
	wc := writeconcern.New(writeconcern.WMajority())
	opts := options.Transaction().SetWriteConcern(wc)

	session, err := uow.DB.StartSession()
	if err != nil {
		return fmt.Errorf("Error starting session: %v", err)
	}
	defer session.EndSession(context.Background())

	return mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(opts); err != nil {
			return fmt.Errorf("Error starting transaction: %v", err)
		}

		for _, change := range uow.changes {
			if err := change(sc); err != nil {
				_ = session.AbortTransaction(sc)
				return fmt.Errorf("Error executing change: %v", err)
			}
		}

		if err := session.CommitTransaction(sc); err != nil {
			_ = session.AbortTransaction(sc)
			return fmt.Errorf("Error committing transaction: %v", err)
		}

		return nil
	})
}
