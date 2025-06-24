package example

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/dmutex"
	"github.com/omniful/go_commons/redis"
	"sync"
	"time"
)

func ExampleDMutex() {
	key := "dummy-key"
	wg := sync.WaitGroup{}

	// Simulate process 1
	go func() {
		wg.Add(1)
		defer wg.Done()

		rc := redis.NewClient(&redis.Config{
			Hosts: []string{"127.0.0.1:6379"},
		})
		//time.Sleep(3 * time.Second)

		dmx := dmutex.New(key, time.Minute, rc)

		ok, err := dmx.TryLock(context.Background())
		if err != nil {
			panic(err)
		}
		if !ok {
			// This should not happen as this should get lock
			panic("process 1 failed to acquire lock")
		}

		fmt.Println("process 1 can now perform work")

		// Simulate work being done
		time.Sleep(10 * time.Second)

		fmt.Println("process 1 work is completed")
		if err := dmx.Unlock(context.Background()); err != nil {
			panic(err)
		}
		fmt.Println("process 1 is completed")
	}()

	// Simulate lag in different process
	time.Sleep(3 * time.Second)

	// Simulate process 2
	go func() {
		wg.Add(1)
		defer wg.Done()

		rc := redis.NewClient(&redis.Config{
			Hosts: []string{"127.0.0.1:6379"},
		})
		//time.Sleep(3 * time.Second)
		dmx := dmutex.New(key, time.Minute, rc)

		ok, err := dmx.TryLock(context.Background())
		if err != nil {
			panic(err)
		}
		if ok {
			// This should not happen as this should get lock
			panic("process 2 acquired lock")
		}
		fmt.Println("process 2 is waiting to unlock")
		isTTLExpired, err := dmx.WaitUntilUnlocked(context.Background())
		if err != nil {
			panic(err)
		}

		if isTTLExpired {
			panic("process 2 got unblocked after ttl acquired")
		}

		fmt.Println("process 2 can now perform work")
		fmt.Println("process 2 is completed")
	}()

	// Simulate process 3
	go func() {
		wg.Add(1)
		defer wg.Done()

		rc := redis.NewClient(&redis.Config{
			Hosts: []string{"127.0.0.1:6379"},
		})
		//time.Sleep(3 * time.Second)
		dmx := dmutex.New(key, time.Minute, rc)

		ok, err := dmx.TryLock(context.Background())
		if err != nil {
			panic(err)
		}
		if ok {
			// This should not happen as this should get lock
			panic("process 3 acquired lock")
		}
		fmt.Println("process 3 is waiting to unlock")
		isTTLExpired, err := dmx.WaitUntilUnlocked(context.Background())
		if err != nil {
			panic(err)
		}

		if isTTLExpired {
			panic("process 3 got unblocked after ttl acquired")
		}

		fmt.Println("process 3 can now perform work")
		fmt.Println("process 3 is completed")
	}()

	wg.Wait()
}
