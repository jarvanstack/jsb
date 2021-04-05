package snowflake

import (
	"fmt"
	"sync"
	"testing"
)

func TestWorker_NextID(t *testing.T) {
	worker := NewWorker(0, 0)
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {
				_, err := worker.nextID()
				if err != nil {
					//fmt.Printf("string=%s\n", "[错误]")
				}
				//fmt.Printf("id=%d\n", id)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("string=%s\n", "1000万用时")
}
