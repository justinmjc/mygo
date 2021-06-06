package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//https://mp.weixin.qq.com/s/m7uMwiWK_b7oudXueUcSjA

//1 数据传递:
//有 4 个 goroutine，编号为 1、2、3、4。每秒钟会有一个 goroutine 打印出它自己的编号，要求你编写程序，
//让输出的编号总是按照 1、2、3、4、1、2、3、4……这个顺序打印出来。
func startTask() {
	n := 4
	chans := []chan struct{}{}
	for i := 0; i < n; i++ {
		chans = append(chans, make(chan struct{}))
	}
	for i := 0; i < 4; i++ {
		go func(i int) {
			for {
				token := <-chans[i] //这段代码中，token 代指“击鼓传花”中的“花”，chans 代指围坐一圈的人。每个 chan（人）都是从上一个 chan（人）手中拿到 token，放在自己手上，从而实现顺序打印 1，2，3，4。
				fmt.Printf("%d \n", i+1)
				chans[(i+1)%n] <- token
				time.Sleep(time.Second)

			}
		}(i)
	}
	chans[0] <- struct{}{}
	select {}
}

//2 信号通知
func gracefulShutdown() {
	go func() {
		//业务处理
	}()
	// 处理CTRL+C等中断信号
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	closed := make(chan struct{})

	// 执行退出之前的清理操作
	go doCleanup(closed)

	select {
	case <-closed:
	case <-time.After(time.Second):
		fmt.Println("清理超时，不等了！")
	}
	fmt.Println("优雅退出！")
}
func doCleanup(closed chan struct{}) {
	time.Sleep(time.Minute)
	close(closed)
}

// 3 锁
// 使用chan实现互斥锁
type Mutex struct {
	ch chan struct{}
}

//初始化
func NewMutex() *Mutex {
	mu := &Mutex{make(chan struct{}, 1)}
	mu.ch <- struct{}{}
	return mu
}

//lock
func (m *Mutex) Lock() {
	<-m.ch
}

//unlock
func (m *Mutex) Unlock() {
	select {
	case m.ch <- struct{}{}:
	default:
		panic("unlock of unlocked mutex")
	}
}

// 加入一个超时的设置
func (m *Mutex) LockTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	select {
	case <-m.ch:
		timer.Stop()
		return true
	case <-timer.C:
	}
	return false
}

// 锁是否已被持有
func (m *Mutex) IsLocked() bool {
	return len(m.ch) == 0
}

//4任务编排

//有一批任务需要处理，但是机器资源有限，只能承受100的并发度，该如何实现？

func task(ch chan struct{}) {
	//执行任务
	time.Sleep(time.Second)
	ch <- struct{}{}
	return
}

func concurrency100() {
	ch := make(chan struct{}, 100)
	for i := 0; i < 100; i++ {
		ch <- struct{}{}
	}

	for {
		<-ch
		go task(ch)
	}
}

//or-Done
func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0: //2
		return nil
	case 1: //3
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() { //4
		defer close(orDone)

		switch len(channels) {
		case 2: //5
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default: //6
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-or(append(channels[3:], orDone)...): //6
			}
		}
	}()
	return orDone
	//递归前，需要声明一个 orDone 变量，用来通知子函数退出。
	//len(channels) == 2 是一种特殊情况，否则会因为 append orDone 产生无限递归。

}

func main() {
	startTask()
}
