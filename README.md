# go-persistence

Object persistence library for GO
Allows saving objects as json files

## Example

Example:

```
type Obj struct {
	StringVal string `json:"string_val"`
	IntVal    int    `json:"int_val"`
	TestVal   string `json:"test_val"`
}

const path = "/tmp/path"

persister, _ := NewJsonObjectPersister(JsonObjectPersisterConfig{Path: path, Prefix: "obj_"})



item := Obj{StringVal: "A string", IntVal: 5}
_ := persister.SaveItem("test1, &item)

// reads an item with a key
retrievedItem = Obj{}
_ := persister.GetItem("test1", &item)


keys := persister.AllKeys()
// [test1, test2]
```