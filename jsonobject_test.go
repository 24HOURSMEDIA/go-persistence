package persistence

import (
	"strconv"
	"testing"
)

type Obj struct {
	StringVal string `json:"string_val"`
	IntVal    int    `json:"int_val"`
	TestVal   string `json:"test_val"`
}

const path = "./tmp/test"

var config = NewJsonObjectPersisterConfig(path+"/config1", "")
var prefixConfig = NewJsonObjectPersisterConfig(path+"/confix2", "objprefix_")
var prefixConfig2 = NewJsonObjectPersisterConfig(path+"/confix2", "objprefix2_")

func TestNewJsonObjectPersister(t *testing.T) {
	persister, err := NewJsonObjectPersister(config)
	if persister != nil {
		t.Log("Got a persister")
	}
	if err != nil {
		t.Fatal("error getting persister")
	}
}

func TestJsonObjectPersister_SaveItem(t *testing.T) {
	persister, _ := NewJsonObjectPersister(config)
	key := "testkey"
	item := Obj{StringVal: "stringb", IntVal: 5}
	err := persister.SaveItem(key, &item)
	if err != nil {
		t.Fatal("error saving item")
	}

	// test deferred
	config.DeferWrites = true
	err = persister.SaveItem(key, &item)
	if err != nil {
		t.Fatal("error saving deferred item")
	}
}

func TestJsonObjectPersister_GetItem(t *testing.T) {
	persister, _ := NewJsonObjectPersister(config)
	key := "testkey"
	got := Obj{StringVal: "stringb", IntVal: 5}
	persister.SaveItem(key, &got)

	want := Obj{}
	err := persister.GetItem(key, &want)
	if err != nil {
		t.Fatal("error retrieving item", err)
	}
	if got.StringVal != want.StringVal {
		t.Fatalf("Expected %q got %q", got.StringVal, want.StringVal)
	}

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func TestJsonObjectPersister_ListKeys(t *testing.T) {
	persister, _ := NewJsonObjectPersister(prefixConfig)
	persister2, _ := NewJsonObjectPersister(prefixConfig2)
	obj := Obj{StringVal: "a string"}

	key1 := "test1"
	key2 := "test2"
	persister.SaveItem(key1, obj)
	persister.SaveItem(key2, obj)
	persister2.SaveItem(key1, obj)

	keys, _ := persister.ListKeys()
	keys2, _ := persister2.ListKeys()
	if len(keys) != 2 {
		t.Fatal("too many keys, probably prefix mixup")
	}
	if len(keys2) != 1 {
		t.Fatal("too many keys, probably prefix mixup")
	}
	if !stringInSlice(key1, keys) {
		t.Fatal("key not found")
	}
	if !stringInSlice(key2, keys) {
		t.Fatal("key not found")
	}
	if stringInSlice("does not exist", keys) {
		t.Fatal("key found")
	}
}

const readAt = 500

func BenchmarkWritePerformance(b *testing.B) {
	config := NewJsonObjectPersisterConfig(path+"/benchmark", "bench_")
	persister, _ := NewJsonObjectPersister(config)
	obj := Obj{StringVal: "a string"}
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		persister.SaveItem(key, obj)
	}
}
func BenchmarkWritePerformanceWithRead(b *testing.B) {
	config := NewJsonObjectPersisterConfig(path+"/benchmark", "bench_")
	persister, _ := NewJsonObjectPersister(config)
	obj := Obj{StringVal: "a string"}
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		persister.SaveItem(key, obj)

		if i == readAt {
			ref := Obj{}
			persister.GetItem(key, &ref)
		}
	}
}
func BenchmarkWritePerformanceDeferredWrite(b *testing.B) {
	config := NewJsonObjectPersisterConfig(path+"/benchmark-deferred", "bench_")
	config.DeferWrites = true
	persister, _ := NewJsonObjectPersister(config)
	obj := Obj{StringVal: "a string"}
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		persister.SaveItem(key, obj)
	}
}
func BenchmarkWritePerformanceDeferredWriteWithRead(b *testing.B) {
	config := NewJsonObjectPersisterConfig(path+"/benchmark-deferred", "bench_")
	config.DeferWrites = true
	persister, _ := NewJsonObjectPersister(config)
	obj := Obj{StringVal: "a string"}
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		persister.SaveItem(key, obj)
		if i == readAt {
			ref := Obj{}
			persister.GetItem(key, &ref)
		}
	}
}
