package handle

type Client struct {
	Name      string          // 客户端链接的名字，这里一般为3rdsession的字符串
	in        <-chan *Message // 传进来的消息管道
	out       chan<- *Message // 发出去的消息管道
	done      <-chan *bool    // 结束的bool
	err       <-chan error    // 错误管道
	diconnect chan<- int      // 断开链接的管道
}

type Message struct {
}
