## channel的典型的应用场景
#### 信息交流
channel 的底层是一个循环队列，当队列的长度大于 0 的时候，它会被当做线程安全队列和 buffer。利用这个特性，一个 goroutine 可以安全的往 channel 中存放数据，另一个 goroutine 可以安全的从 channel 中读取数据，这样就实现了 goroutine 之间的消息交流。


#### 数据传递
数据传递类似游戏“击鼓传花”。鼓响时，花（或者其它物件）从一个人手里传到下一个人，数据就类似这里的花。
#### 信号通知
channel 类型有这样一个特性：如果 channel 为空，那么 recevier 接收数据的时候就会阻塞，直到有新的数据进来或者 channel 被关闭。

利用这个特性，就可以实现 wait/notify 设计模式。另外还有一个经常碰到的场景，实现程序的 graceful shutdown。
#### 锁
sync.Mutex 通过修改持有锁标记位的状态达到占有锁的目的，因此 channel 可以通过转移这个标记位的所有权实现占有锁。
#### 任务编排
通过 WaitGroup，我们能很容易的实现 等待一组 goroutine 完成任务 这种任务编排需求。同样，我们也可以用 channel 实现。
###### or-Done
or-Done 模式对应的场景很好理解，n 个任务，有一个完成就算完成。
###### 扇入模式
###### 扇出模式
###### stream
###### map-reduce
