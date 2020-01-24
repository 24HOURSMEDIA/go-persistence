# go-persistence

Object persistence library for GO
Allows saving objects as json files

## Example

Example:

```
import (
	"github.com/24HOURSMEDIA/go-persistence"
)

type Obj struct {
	StringVal string `json:"string_val"`
	IntVal    int    `json:"int_val"`
	TestVal   string `json:"test_val"`
}

const path = "/tmp/path"

persister, _ := persistence.NewJsonObjectPersister(JsonObjectPersisterConfig{Path: path, Prefix: "obj_"})

// write an item
item := Obj{StringVal: "A string", IntVal: 5}
_ := persister.SaveItem("test1, &item)

// reads an item by key
retrievedItem = Obj{}
_ := persister.GetItem("test1", &retrievedItem)


keys := persister.AllKeys()
// [test1, test2]
```