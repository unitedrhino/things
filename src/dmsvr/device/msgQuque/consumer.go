package main

import (
	"context"
	"fmt"
	"gitee.com/godLei6/things/shared/utils"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)
var brokers = []string{"81.68.223.176:9092"}
var topics = []string{"onPublish","onConnect","onDisconnect"}
var group = "39"


type Router struct {
	Topic string
	Handler func(msg *string)error
}

type Kafka struct {
	Brokers []string
	Routers map[string]Router//key是topic 对应的是处理函数
	Topics 	[]string
	//OffsetNewest int64 = -1
	//OffsetOldest int64 = -2
	StartOffset       int64 `json:",optional"`
	Version           string `json:",optional"`
	ready             chan bool
	Group             string `json:",optional"`
	ChannelBufferSize int `json:",default=20"`
}

func NewKafka() *Kafka {
	return &Kafka{
		ChannelBufferSize: 2,
		ready:             make(chan bool),
		Version:"1.1.1",
	}
}
func (k *Kafka)AddRouter(router Router){
	k.Routers[router.Topic] = router
	k.Topics = append(k.Topics,router.Topic)
}


func (k *Kafka) Init() func() {
	log.Printf("kafka init...")

	version, err := sarama.ParseKafkaVersion(k.Version)
	if err != nil {
		log.Fatalf("Error parsing Kafka version: %v", err)
	}
	config := sarama.NewConfig()
	config.Version = version
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange // 分区分配策略
	config.Consumer.Offsets.Initial = -2                    // 未找到组消费位移的时候从哪边开始消费
	config.ChannelBufferSize = k.ChannelBufferSize // channel长度

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(k.Brokers, k.Group, config)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			//util.HandlePanic("client.Consume panic", log.StandardLogger())
		}()
		for {
			if err := client.Consume(ctx, k.Topics, k); err != nil {
				log.Fatalf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				log.Println(ctx.Err())
				return
			}
			k.ready = make(chan bool)
		}
	}()
	<-k.ready
	log.Printf("Sarama consumer up and running!...")
	// 保证在系统退出时，通道里面的消息被消费
	return func() {
		log.Printf("kafka close")
		cancel()
		wg.Wait()
		if err = client.Close(); err != nil {
			log.Printf("Error closing client: %v", err)
		}
	}
}



// Setup is run at the beginning of a new session, before ConsumeClaim
func (k *Kafka) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(k.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (k *Kafka) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (k *Kafka) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	// 具体消费消息
	for message := range claim.Messages() {
		msg := string(message.Value)
		log.Printf("%s|%+v|msessage=%+v\n",utils.FuncName(), msg,string(message.Key))
		v,ok := k.Routers[message.Topic]
		if ok != true{
			panic(fmt.Sprintf("get msg bug topic not have hander func:%s",message.Topic))
		}
		err := v.Handler(&msg)
		if err != nil {
			
		}
		//run.Run(msg)
		// 更新位移
		session.MarkMessage(message, "")
	}
	return nil
}


func Start(){

	k := NewKafka()
	f := k.Init()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		log.Printf("terminating: via signal")
	}
	f()
}
func main()  {
	Start()
}
