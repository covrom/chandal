package dal

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/covrom/chandal/core"
	"github.com/covrom/chandal/libs/dbgen"
	"github.com/covrom/chandal/libs/verr"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CreateUserQuery struct {
	Ctx    context.Context
	User   core.User
	Result chan verr.ValErr[core.User]
}

func (q CreateUserQuery) Do() {
	CreateUserChan <- q
}

var CreateUserChan chan CreateUserQuery

var CreateUserTimeout = 20 * time.Second

func GoCreateUser(ctx context.Context, db *sqlx.DB, n, chanBufSize int, wg *sync.WaitGroup) {
	CreateUserChan = make(chan CreateUserQuery, chanBufSize)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go createUsersWorker(ctx, db, wg)
	}
}

func createUsersWorker(ctx context.Context, db *sqlx.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Printf("CreateUsersWorker exit: %v", ctx.Err())
			return
		case q, ok := <-CreateUserChan:
			if !ok {
				log.Printf("CreateUsersWorker exit: CreateUsersChan closed")
				return
			}
			// here you can put it in the buffer and do batch processing in the future
			createUsers(db, q)
		}
	}
}

func createUsers(db *sqlx.DB, q CreateUserQuery) {
	defer close(q.Result)

	u := q.User
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	tx := db.MustBeginTx(q.Ctx, nil)

	err := dbgen.Insert(q.Ctx, tx, "users", u)

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	if e := verr.SendTimeout(
		q.Result, verr.ValErr[core.User]{Value: &u, Err: err}, CreateUserTimeout,
	); e != nil {
		log.Printf("createUsers sending error result %s, error = %v", e, err)
	}
}
