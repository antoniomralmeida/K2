package kb

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/lib"
	"github.com/antoniomralmeida/k2/services"

	"github.com/madflojo/tasks"
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
	initializers.Log(err, initializers.Fatal)
	_, err = scheduler.Add(&tasks.Task{
		Interval: time.Duration(60 * time.Second),
		TaskFunc: func() error {
			go GKB.RefreshRules()
			return nil
		},
	})
	initializers.Log(err, initializers.Fatal)

	initializers.Log("K2 System started!", initializers.Info)
	fmt.Println("K2 System started! Press ESC to shutdown")

	for {
		if lib.KeyPress() == 27 || GKB.halt {
			fmt.Printf("Shutdown...")
			scheduler.Stop()
			services.Stop()
			wg.Done()
			os.Exit(0)
		}
	}

}

func (kb *KnowledgeBased) RunStackRules() error {
	initializers.Log("RunStackRules...", initializers.Info)
	if len(kb.stack) > 0 {
		kb.mutex.Lock()
		localstack := kb.stack
		kb.mutex.Unlock()
		mark := len(localstack) - 1
		sort.Slice(localstack, func(i, j int) bool {
			return (localstack[i].Priority > localstack[j].Priority) || (localstack[i].Priority == localstack[j].Priority && localstack[j].lastexecution.Unix() > localstack[i].lastexecution.Unix())
		})

		runtaks := make(map[initializers.OID]*KBRule) //run the rule once
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
	initializers.Log("RefrehRules...", initializers.Info)
	for i := range kb.Objects {
		if !kb.Objects[i].parsed {
			for j := range kb.Rules {
				for k := range kb.Rules[j].bkclasses {
					if kb.Rules[j].bkclasses[k] == kb.Objects[i].Bkclass {
						_, bin, err := kb.ParsingCommand(kb.Rules[j].Rule)
						if initializers.Log(err, initializers.Error) != nil {
							kb.linkerRule(&kb.Rules[j], bin)
						}
					}
				}
			}
			kb.Objects[i].parsed = true
		}
	}
	return nil
}
