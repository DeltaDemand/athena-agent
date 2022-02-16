package main

import (
	"context"
	"flag"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func main() {
	n := flag.Int("n", 100, "创造死循环个数")
	t := flag.Int("t", 90, "死循环时间(s)")
	sl := flag.Int("s", 0, "每次循环内睡眠时间(ns)")
	flag.Parse()
	println(*n, "个goroutine，开始跑cpu: ", *t, "s")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(*t))
	wg.Add(*n)
	for i := 1; i < *n; i++ {
		go cycle(ctx, *sl)
	}
	cycle(ctx, *sl)
	cancel()
	wg.Wait()
}

//死循环
func cycle(ctx context.Context, sle int) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if sle > 0 {
				time.Sleep(time.Duration(sle))
			}
		}
	}

}
