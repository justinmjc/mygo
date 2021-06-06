package main

import (
	"context"
	"fmt"
	"time"
)

/*
使用context如何更优雅实现协程间取消信号的同步：
	子协程本来应该要10秒才能执行完
    WithTimeout创建的ctx 5秒就会退出
    子线程监听主线程传入的ctx，一旦ctx.Done()返回空channel，子线程即可取消执行任务。
    但这个例子还无法展现context的传递取消信息的强大优势。

*/
func main() {
	messages := make(chan int, 10)

	for i := 0; i < 10; i++ {
		messages <- i
	}
	/*
		context.Background()用来做根context
		WithTimeout 创建一个定时取消的context
	*/
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	//consumer
	go func(ctx context.Context) {
		ticker := time.NewTicker(1 * time.Second)
		for _ = range ticker.C {
			select {
			case <-ctx.Done():
				fmt.Println("child process interrupt...")
				return
			default:
				fmt.Printf("send message: %d\n", <-messages)
			}
		}
	}(ctx)

	defer close(messages)
	defer cancel()

	select {
	case <-ctx.Done():
		time.Sleep(1 * time.Second)
		fmt.Println("main process exit!")
	}
}
