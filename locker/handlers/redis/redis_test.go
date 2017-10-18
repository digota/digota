package redis

import (
	"testing"

	"github.com/digota/digota/config"
)

type testObj struct {
	Id string `bson:"_id"`
}

func (o *testObj) GetNamespace() string {
	return "mongo_test"
}

func (o *testObj) GetId() string {
	return o.Id
}

func (o *testObj) SetId(id string) {
	o.Id = id
}

func TestNewLocker(t *testing.T) {
	if l, err := NewLocker(config.Locker{Address: []string{"localhost"}}); err != nil {
		t.Fatal(err)
	} else {
		l.Close()
	}
}
