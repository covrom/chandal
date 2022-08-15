package api

import (
	"encoding/json"
	"net/http"

	"github.com/covrom/chandal/core"
	"github.com/covrom/chandal/dal"
	"github.com/covrom/chandal/libs/verr"
)

func (a *Api) GetUsers(w http.ResponseWriter, r *http.Request) {
	t := make(chan verr.ValErr[core.User], 100)

	q := dal.GetUsersQuery{
		Ctx:    r.Context(),
		Result: t,
	}
	q.Do()

	enc := json.NewEncoder(w)

	for v := range t {
		enc.Encode(v)
		// chunked transfer without [ ] ,
		w.(http.Flusher).Flush()
		if v.Err != nil {
			// drain
			for range t {
			}
			return
		}
	}
}

func (a *Api) CreateUser(w http.ResponseWriter, r *http.Request) {
	t := make(chan verr.ValErr[core.User], 1)

	q := dal.CreateUserQuery{
		Ctx:    r.Context(),
		Result: t,
	}

	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := dec.Decode(&q.User); err != nil {
		http.Error(w, "bad user data", http.StatusBadRequest)
	}

	q.Do()

	enc := json.NewEncoder(w)

	for v := range t {
		enc.Encode(v)
		// chunked transfer without [ ] ,
		w.(http.Flusher).Flush()
		if v.Err != nil {
			// drain
			for range t {
			}
			return
		}
	}
}
