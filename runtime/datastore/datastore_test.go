package datastore

import (
	// 	"appengine/aetest"
	// 	"github.com/stretchr/testify/assert"
	"github.com/ionous/sashimi/meta"
	"testing"
)

// https://golang.org/pkg/encoding/binary/#Write

func TestDataStore(t *testing.T) {
	// if c, err := aetest.NewContext(nil); assert.NoError(t, err) {

	// 	defer c.Close()
	// 	// kind, key string, key number, parent key
	// 	key := datastore.NewKey(c, "BankAccount", "", 1, nil)
	// 	if _, err := datastore.Put(c, key, &BankAccount{100}); err != nil {
	// 		t.Fatal(err)
	// 	}

	// 	err = withdraw(c, "myid", 128, 0)
	// 	if err == nil || err.Error() != "insufficient funds" {
	// 		t.Errorf("Error: %v; want insufficient funds error", err)
	// 	}

	// 	b := BankAccount{}
	// 	if err := datastore.Get(c, key, &b); err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	if bal, want := b.Balance, 100; bal != want {
	// 		t.Errorf("Balance %d, want %d", bal, want)
	// 	}
	// }
}
