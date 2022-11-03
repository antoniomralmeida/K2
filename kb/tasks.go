package kb

import (
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/antoniomralmeida/k2/lib"
	"github.com/eiannone/keyboard"
	"github.com/gofiber/fiber/v2"
	"github.com/madflojo/tasks"
)

func (kb *KnowledgeBased) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	if lib.IsMainThread() {

		// Start the Scheduler
		scheduler := tasks.New()
		defer scheduler.Stop()

		// Add tasks
		_, err := scheduler.Add(&tasks.Task{
			Interval: time.Duration(2 * time.Second),

			TaskFunc: func() error {
				go kb.Scan()
				return nil
			},
		})
		lib.LogFatal(err)
		_, err = scheduler.Add(&tasks.Task{
			Interval: time.Duration(60 * time.Second),
			TaskFunc: func() error {
				go kb.ReLink()
				return nil
			},
		})
		lib.LogFatal(err)

		_, err = scheduler.Add(&tasks.Task{
			Interval: time.Duration(2 * time.Second),
			TaskFunc: func() error {
				go kb.IOTParsing()
				return nil
			},
		})
		lib.LogFatal(err)

		//TODO: Criar Clean History Task

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
				wg.Done()
				os.Exit(0)
			}
		}
	}

}

func (kb *KnowledgeBased) Scan() error {
	log.Println("Scaning...")
	if len(kb.stack) > 0 {
		kb.mutex.Lock()
		localstack := kb.stack
		kb.mutex.Unlock()
		mark := len(localstack) - 1
		sort.Slice(localstack, func(i, j int) bool {
			return (localstack[i].Priority > localstack[j].Priority) || (localstack[i].Priority == localstack[j].Priority && localstack[j].lastexecution.Unix() > localstack[i].lastexecution.Unix())
		})

		for len(localstack) > 0 {
			r := localstack[0]
			r.Run()
			localstack = localstack[1:]
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

func (kb *KnowledgeBased) IOTParsing() error {
	log.Println("IOTParsing...")
	//TODO: https://nabto.com/guide-iot-protocols-standards/, definir protocolo para IOT SET and SET

	for i := range kb.Objects {
		for j := range kb.Objects[i].Attributes {
			a := &kb.Objects[i].Attributes[j]
			if !a.Validity() {
				if a.KbAttribute.isSource(KBSource(User)) && kb.IOTApi != "" {
					iotapi := kb.IOTApi + "?" + a.getFullName()
					api := fiber.AcquireAgent()
					req := api.Request()
					req.Header.SetMethod("post")
					req.SetRequestURI(iotapi)
					if err := api.Parse(); err != nil {
						log.Println(err)
					} else {
						_, body, errs := api.Bytes()
						if errs != nil {
							a.SetValue(string(body), IOT, 100.0)
						}
					}
				}
			}
		}
	}
	return nil
}

func (kb *KnowledgeBased) ReLink() error {
	log.Println("ReLink...")
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
