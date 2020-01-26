package persistence

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const dirPerms = 0755
const filePerms = 0600

// Configuration for the persister
type jsonObjectPersisterConfig struct {
	// private, Directory to save files in
	path string
	// private Prefix for the json files, if you have multiple persisters that save in the same directory
	prefix string
	// defers writes.. EXPERIMENTAL
	DeferWrites bool
}

func newJsonObjectPersisterConfig(path string, prefix string) *jsonObjectPersisterConfig {
	return &jsonObjectPersisterConfig{path: path, prefix: prefix}
}

// Returns a new default configuration
func NewJsonObjectPersisterConfig(path string, prefix string) *jsonObjectPersisterConfig {
	return &jsonObjectPersisterConfig{path: path, prefix: prefix, DeferWrites: false}
}

// The persister. Create with persistence.NewJsonObjectPersister
type jsonObjectPersister struct {
	config *jsonObjectPersisterConfig
}

// Internal, create a path to save the item to
func (persister jsonObjectPersister) createPath(key *string) *string {
	path := persister.config.path + "/" + persister.config.prefix + *key + ".json"
	return &path
}

// Internal, retrieves the key from a path (strips extension and prefix)
func (persister jsonObjectPersister) keyFromPath(sourceKey *string) (string, error) {
	key := strings.Replace(*sourceKey, ".json", "", 1)
	if persister.config.prefix == "" {
		return key, nil
	}
	trimmed := strings.TrimPrefix(key, persister.config.prefix)
	if trimmed == key {
		return trimmed, errors.New("prefix not found")
	}
	return trimmed, nil
}

func (persister jsonObjectPersister) saveItemNow(key string, obj interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	path := persister.createPath(&key)
	return ioutil.WriteFile(*path, b, filePerms)
}

func (persister jsonObjectPersister) SaveItem(key string, obj interface{}) error {
	if persister.config.DeferWrites {
		addToQueue(persister, key, obj)
		return nil
	} else {
		return persister.saveItemNow(key, obj)
	}
}

// Retrieves an item from a json file
func (persister jsonObjectPersister) GetItem(key string, obj interface{}) error {
	// if the persister defers writes, wait until all of the queue is written
	if persister.config.DeferWrites {
		queue.wait()
	}

	path := persister.createPath(&key)
	b, err := ioutil.ReadFile(*path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &obj)
}

// List all keys (i.e. the json files in the directory without the prefix and .json suffix)
func (persister jsonObjectPersister) ListKeys() ([]string, error) {
	keys := make([]string, 0)
	err := filepath.Walk(persister.config.path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			name := info.Name()
			key, invalid := persister.keyFromPath(&name)
			if invalid == nil {
				keys = append(keys, key)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return keys, nil
}

func NewJsonObjectPersister(config *jsonObjectPersisterConfig) (*jsonObjectPersister, error) {
	persister := jsonObjectPersister{config: config}
	err := os.MkdirAll(config.path, dirPerms)
	return &persister, err
}
