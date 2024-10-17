package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	startTime := time.Now()
	ctx := context.Background()
	deadline := time.Now().Add(400 * time.Millisecond)
	ctx1, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()
	ctx2 := context.WithValue(ctx1, "key", 1)
	userID := 10
	val, err := fetchUserData(ctx2, userID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("response :", val)
	fmt.Println("time taken :", time.Since(startTime))
}

type Response struct {
	value int
	err   error
}

func fetchUserData(ctx context.Context, userID int) (int, error) {
	value := ctx.Value("key")
	fmt.Println("key:", value)
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	respCh := make(chan Response)

	go func() {
		val, err := fetchThirdPartyStuffWhichCanBeSlow()
		respCh <- Response{
			value: val,
			err:   err,
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return 0, fmt.Errorf("fetching data from third party took too long")
		case resp := <-respCh:
			return resp.value, resp.err
		}
	}
}

func fetchThirdPartyStuffWhichCanBeSlow() (int, error) {
	time.Sleep(400 * time.Millisecond)

	return 66, nil
}
