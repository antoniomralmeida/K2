package kb

import (
	"fmt"
	"os"
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
