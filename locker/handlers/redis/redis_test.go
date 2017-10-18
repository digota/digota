package redis

import (
	"errors"
	"testing"
	"time"

	"github.com/digota/digota/config"
	"github.com/garyburd/redigo/redis"
	uuid "github.com/satori/go.uuid"
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

type testPool struct {
	redisConn *testRedisConn
}

func (p *testPool) Get() redis.Conn {
	return p.redisConn
}
func (p *testPool) Close() error {
	return nil
}

type testRedisConn struct {
	// Do command params and return values
	doCmd      string
	doParams   []interface{}
	doReply    interface{}
	doError    error
	doBlocking time.Duration
}

func (rc *testRedisConn) Close() error {
	return nil
}
func (rc *testRedisConn) Err() error {
	return nil
}
func (rc *testRedisConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	rc.doCmd = commandName
	rc.doParams = args
	if rc.doBlocking > 0 {
		time.Sleep(rc.doBlocking)
	}
	return rc.doReply, rc.doError
}
func (rc *testRedisConn) Send(commandName string, args ...interface{}) error {
	return nil
}
func (rc *testRedisConn) Flush() error {
	return nil
}
func (rc *testRedisConn) Receive() (reply interface{}, err error) {
	return nil, nil
}

func TestNewLocker(t *testing.T) {
	if l, err := NewLocker(config.Locker{Address: []string{"localhost"}}); err != nil {
		t.Fatal(err)
	} else {
		l.Close()
	}

	if _, err := NewLocker(config.Locker{}); err == nil {
		t.Fatal("We should have an error without at least one address!")
	}
}

func TestLock_Lock(t *testing.T) {
	rc := &testRedisConn{
		doReply: "OK",
		doError: nil,
	}
	l := &locker{rp: &testPool{redisConn: rc}}

	testObj := &testObj{Id: uuid.NewV4().String()}
	unlock, err := l.Lock(testObj)
	if err != nil {
		t.Fatal(err)
	}

	if rc.doCmd != "SET" {
		t.Errorf("Wrong redis command! Expected: SET, Got: %s", rc.doCmd)
	}

	objKey, _ := getKey(testObj)
	if rc.doParams[0].(string) != objKey {
		t.Errorf("Wrong key! Expected: %s, Got: %s", objKey, rc.doParams[0].(string))
	}

	if rc.doParams[1].(string) != "NX" {
		t.Errorf("Wrong param! Expected: NX, Got: %s", rc.doParams[1].(string))
	}

	if err = unlock(); err != nil {
		t.Error(err)
	}
}

func TestLock_LockFail(t *testing.T) {
	errConnFailed := errors.New("connection failed")
	rc := &testRedisConn{
		doReply: "",
		doError: errConnFailed,
	}
	l := &locker{rp: &testPool{redisConn: rc}}

	testObj := &testObj{Id: ""}
	_, err := l.Lock(testObj)
	if err != ErrMissingInfo {
		t.Fatal(err)
	}

	testObj.Id = uuid.NewV4().String()
	_, err = l.Lock(testObj)
	if err != errConnFailed {
		t.Fatal(err)
	}
}

func TestLock_TryLockSuccess(t *testing.T) {
	rc := &testRedisConn{
		doReply:    "OK",
		doError:    nil,
		doBlocking: 10 * time.Millisecond,
	}
	l := &locker{rp: &testPool{redisConn: rc}}

	testObj := &testObj{Id: uuid.NewV4().String()}
	unlock, err := l.TryLock(testObj, 100*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	if rc.doCmd != "SET" {
		t.Errorf("Wrong redis command! Expected: SET, Got: %s", rc.doCmd)
	}

	objKey, _ := getKey(testObj)
	if rc.doParams[0].(string) != objKey {
		t.Errorf("Wrong key! Expected: %s, Got: %s", objKey, rc.doParams[0].(string))
	}

	if rc.doParams[1].(string) != "NX" {
		t.Errorf("Wrong param! Expected: NX, Got: %s", rc.doParams[1].(string))
	}

	if err = unlock(); err != nil {
		t.Error(err)
	}
}

func TestLock_TryLockTimeout(t *testing.T) {
	rc := &testRedisConn{
		doReply:    "OK",
		doError:    nil,
		doBlocking: 100 * time.Millisecond,
	}
	l := &locker{rp: &testPool{redisConn: rc}}

	testObj := &testObj{Id: uuid.NewV4().String()}
	_, err := l.TryLock(testObj, 20*time.Millisecond)
	if err != ErrTimeout {
		t.Fatal(err)
	}

	// Test GetKey error too
	testObj.Id = ""
	_, err = l.TryLock(testObj, 110*time.Millisecond)
	if err != ErrMissingInfo {
		t.Fatal(err)
	}
}

func TestLock_TryLockFailed(t *testing.T) {
	errConnFailed := errors.New("connection failed")
	rc := &testRedisConn{
		doReply:    "",
		doError:    errConnFailed,
		doBlocking: 10 * time.Millisecond,
	}
	l := &locker{rp: &testPool{redisConn: rc}}

	testObj := &testObj{Id: uuid.NewV4().String()}
	_, err := l.TryLock(testObj, 100*time.Millisecond)
	if err != errConnFailed {
		t.Fatal(err)
	}
}

func TestLock_unlock(t *testing.T) {
	rc := &testRedisConn{
		doError: nil,
	}
	l := &locker{rp: &testPool{redisConn: rc}}

	err := l.unlock("lockKey")
	if err != nil {
		t.Error(err)
	}

	if rc.doCmd != "DEL" {
		t.Errorf("Wrong redis command! Expected: DEL, Got: %s", rc.doCmd)
	}
}

func TestLock_getKey(t *testing.T) {
	_, err := getKey(&testObj{Id: ""})
	if err == nil {
		t.Fatal("getKey should return an error for missing object id")
	}
}
