package kb

import (
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/antoniomralmeida/k2/lib"
	"github.com/antoniomralmeida/k2/web"
	"github.com/eiannone/keyboard"
	"github.com/madflojo/tasks"
	"gopkg.in/mgo.v2/bson"
)

func Run(wg *sync.WaitGroup) {
	defer wg.Done()

	// Start the Scheduler
	scheduler := tasks.New()
	defer scheduler.Stop()

	// Add tasks
	_, err := scheduler.Add(&tasks.Task{
		Interval: time.Duration(2 * time.Second),

		TaskFunc: func() error {
			go GKB.RunStackRules()
			return nil
		},
	})
	lib.LogFatal(err)
	_, err = scheduler.Add(&tasks.Task{
		Interval: time.Duration(60 * time.Second),
		TaskFunc: func() error {
			go GKB.RefreshRules()
			return nil
		},
	})
	lib.LogFatal(err)

	log.Println("K2 System started!")
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()
	fmt.Println("K2 System started! Press ESC to shutdown")

	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		if event.Key == keyboard.KeyEsc {
			fmt.Printf("Shutdown...")
			scheduler.Stop()
			web.Stop()
			wg.Done()
			os.Exit(0)
		}
	}

}

func (kb *KnowledgeBased) RunStackRules() error {
	log.Println("RunStackRules...")
	if len(kb.stack) > 0 {
		kb.mutex.Lock()
		localstack := kb.stack
		kb.mutex.Unlock()
		mark := len(localstack) - 1
		sort.Slice(localstack, func(i, j int) bool {
			return (localstack[i].Priority > localstack[j].Priority) || (localstack[i].Priority == localstack[j].Priority && localstack[j].lastexecution.Unix() > localstack[i].lastexecution.Unix())
		})

		runtaks := make(map[bson.ObjectId]*KBRule) //run the rule once
		for _, r := range localstack {
			if runtaks[r.Id] == nil {
				r.Run()
				runtaks[r.Id] = r
			}
		}
		kb.mutex.Lock()
		kb.stack = kb.stack[mark:]
		kb.mutex.Unlock()
	}
	for i := range kb.Rules {
		if kb.Rules[i].ExecutionInterval != 0 && time.Now().After(kb.Rules[i].lastexecution.Add(time.Duration(kb.Rules[i].ExecutionInterval)*time.Millisecond)) {
			kb.mutex.Lock()
			kb.stack = append(kb.stack, &kb.Rules[i])
			kb.mutex.Unlock()
		}
	}
	return nil
}

func (kb *KnowledgeBased) RefreshRules() error {
	log.Println("RefrehRules...")
	for i := range kb.Objects {
		if !kb.Objects[i].parsed {
			for j := range kb.Rules {
				for k := range kb.Rules[j].bkclasses {
					if kb.Rules[j].bkclasses[k] == kb.Objects[i].Bkclass {
						_, bin, err := kb.ParsingCommand(kb.Rules[j].Rule)
						lib.LogFatal(err)
						kb.linkerRule(&kb.Rules[j], bin)
					}
				}
			}
			kb.Objects[i].parsed = true
		}
	}
	return nil
}
