package dal

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/covrom/chandal/core"
	"github.com/covrom/chandal/libs/dbgen"
	"github.com/covrom/chandal/libs/verr"
	"github.com/jmoiron/sqlx"
)

type GetUsersQuery struct {
	Ctx    context.Context
	Result chan verr.ValErr[core.User]
}

func (q GetUsersQuery) Do() {
	GetUsersChan <- q
}

var GetUsersChan chan GetUsersQuery

var GetUsersTimeout = 20 * time.Second

func GoGetUsers(ctx context.Context, db *sqlx.DB, n, chanBufSize int, wg *sync.WaitGroup) {
	GetUsersChan = make(chan GetUsersQuery, chanBufSize)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go getUsersWorker(ctx, db, wg)
	}
}

func getUsersWorker(ctx context.Context, db *sqlx.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Printf("getUsersWorker exit: %v", ctx.Err())
			return
		case q, ok := <-GetUsersChan:
			if !ok {
				log.Printf("getUsersWorker exit: GetUsersChan closed")
				return
			}
			// here you can put it in the buffer and do batch processing in the future
			getUsers(db, q)
		}
	}
}

func getUsers(db *sqlx.DB, q GetUsersQuery) {
	defer close(q.Result)

	tx := db.MustBeginTx(q.Ctx, nil)
	defer tx.Rollback()

	err := dbgen.SelectAll(q.Ctx, tx, "users", func(ctx context.Context, t *core.User) error {
		return verr.SendTimeout(
			q.Result, verr.ValErr[core.User]{Value: t}, GetUsersTimeout,
		)
	})

	if err != nil && err != sql.ErrNoRows {
		if e := verr.SendTimeout(
			q.Result, verr.ValErr[core.User]{Err: err}, GetUsersTimeout,
		); e != nil {
			log.Printf("getUsers sending error result %s, error = %v", e, err)
		}
		return
	}
}
