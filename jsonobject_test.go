package go_persistence

import "testing"

type Obj struct {
	StringVal string `json:"string_val"`
	IntVal    int    `json:"int_val"`
	TestVal   string `json:"test_val"`
}

const path = "/tmp/test"

func TestNewJsonObjectPersister(t *testing.T) {
	persister, err := NewJsonObjectPersister(JsonObjectPersisterConfig{Path: path})
	if persister != nil {
		t.Log("Got a persister")
	}
	if err != nil {
		t.Fatal("error getting persister")
	}
}

func TestJsonObjectPersister_SaveItem(t *testing.T) {
	persister, _ := NewJsonObjectPersister(JsonObjectPersisterConfig{Path: path})
	key := "testkey"
	item := Obj{StringVal: "stringb", IntVal: 5}
	err := persister.SaveItem(key, &item)
	if err != nil {
		t.Fatal("error saving item")
	}
}

func TestJsonObjectPersister_GetItem(t *testing.T) {
	persister, _ := NewJsonObjectPersister(JsonObjectPersisterConfig{Path: path})
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
	persister, _ := NewJsonObjectPersister(JsonObjectPersisterConfig{Path: path, Prefix: "listtest_"})
	obj := Obj{StringVal: "a string"}

	key1 := "test1"
	key2 := "test2"
	persister.SaveItem(key1, obj)
	persister.SaveItem(key2, obj)

	keys, err := persister.ListKeys()
	t.Logf("keys:: %v", keys)
	if err != nil {
		t.Fatal("error retrieving item", err)
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
