package example

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/dchannel"
	"github.com/omniful/go_commons/redis"
	"time"
)

func ExampleDChannel() {
	rc := redis.NewClient(&redis.Config{
		Hosts: []string{"127.0.0.1:6379"},
	})

	ch1 := dchannel.New("test", rc)
	ch2 := dchannel.New("test", rc)
	ch3 := dchannel.New("test", rc)

	go func() {
		for {
			fmt.Println("CH1: Waiting for message")
			d, err := ch1.Listen(context.Background())
			if err != nil {
				fmt.Printf("CH1: Error listening: %s\n", err)
				return
			}
			fmt.Printf("CH1: Got: %s\n", d)
		}
	}()

	go func() {
		for {
			fmt.Println("CH2: Waiting for message")
			d, err := ch2.Listen(context.Background())
			if err != nil {
				fmt.Printf("CH2: Error listening: %s\n", err)
				return
			}
			fmt.Printf("CH2: Got: %s\n", d)
		}
	}()

	go func() {
		for {
			fmt.Println("CH3: Waiting for message")
			d, err := ch3.Listen(context.Background())
			if err != nil {
				fmt.Printf("CH3: Error listening: %s\n", err)
				return
			}
			fmt.Printf("CH3: Got: %s\n", d)
		}
	}()

	// To ensure all listeners are ready
	time.Sleep(1 * time.Second)

	if err := ch1.Push(context.Background(), "dummy-message-1"); err != nil {
		fmt.Printf("CH1: Error pushing: %s\n", err)
	}

	ch1.Close()
	ch2.Close()
	ch3.Close()

	if err := ch1.Push(context.Background(), "dummy-message-2"); err != nil {
		fmt.Printf("CH1-1: Error pushing: %s\n", err)
	}

	fmt.Println("Completed")
}
