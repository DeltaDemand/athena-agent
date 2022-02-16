package main

import (
	"context"
	"flag"
	"sync"
	"time"
)

var (
	wg  = sync.WaitGroup{}
	len int
	sle int
	a   []string
)

func main() {
	n := flag.Int("goN", 10, "创造goroutine跑死循环个数")
	t := flag.Int("times", 90, "死循环时间(s)")
	flag.IntVar(&sle, "sleep", 0, "每次循环睡眠时间(ns)")
	flag.IntVar(&len, "append", 0, "每个goroutine内append字符串的次数")
	flag.Parse()
	println(*n, "个goroutine，开始跑cpu: ", *t, "s")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(*t))
	wg.Add(*n)
	for i := 1; i < *n; i++ {
		go cycle(ctx)
	}
	cycle(ctx)
	cancel()
	wg.Wait()
}

//循环体
func cycle(ctx context.Context) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for i := 0; i < len; i++ {
				a = append(a, "testMemory")
			}
			if sle > 0 {
				time.Sleep(time.Duration(sle))
			}
		}
	}

}
