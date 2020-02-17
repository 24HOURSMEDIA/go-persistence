package persistence

import (
	"fmt"
	"sync"
	"time"
)

// Saves an item to a json file
type queueItem struct {
	persister ObjectPersisterInterface
	key       string
	obj       interface{}
}
type queueType struct {
	items           []*queueItem
	mux             sync.Mutex
	maxSize         int
	enableLogging   bool
	loggingInterval time.Duration
}

// waits until queue is not busy
func (queue *queueType) wait() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func(queue *queueType, wg *sync.WaitGroup) {
		defer wg.Done()
		waiting := true
		for waiting {
			if len(queue.items) == 0 {
				waiting = false
			} else {
			}
			time.Sleep(time.Microsecond)
		}
	}(queue, &wg)
	wg.Wait()
}
func (queue *queueType) Size() int {
	return len(queue.items)
}

func addToQueue(persister ObjectPersisterInterface, key string, obj interface{}) {
	if queue.Size() > queue.maxSize {
		queue.wait()
	}
	queue.mux.Lock()
	queue.items = append(queue.items, &queueItem{
		persister: persister,
		key:       key,
		obj:       obj,
	})
	queue.mux.Unlock()
}

func handleQueue() {
	if len(queue.items) > 0 {
		var process = false
		var item = new(queueItem)
		queue.mux.Lock()
		if len(queue.items) > 1 {
			item, queue.items = queue.items[0], queue.items[1:]
			process = true
		} else if len(queue.items) == 1 {
			item, queue.items = queue.items[0], make([]*queueItem, 0)
			process = true
		}
		queue.mux.Unlock()
		if process {
			go func(item *queueItem) {
				queue.mux.Lock()
				item.persister.saveItemNow(item.key, item.obj)
				queue.mux.Unlock()
			}(item)
		}
	}
}

var queue = queueType{maxSize: 1000, enableLogging: false, loggingInterval: time.Second * 5}

func init() {
	// get status logger

	go func() {
		for true {
			if queue.enableLogging {
				fmt.Println("Items in queue", queue.Size())
			}
			time.Sleep(queue.loggingInterval)
		}
	}()

	// launch async queue handlers
	go func() {
		for true {
			handleQueue()
			time.Sleep(time.Millisecond)
		}
	}()
}
