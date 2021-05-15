package device

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"yl/shared/utils"
)



func NewDevice(){
	defer func() {
		if p := recover(); p != nil {
			utils.HandleThrow(p)
		}
	}()
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = true
	client,err := sarama.NewClient([]string{"localhost:9092"}, config)
	if err != nil {
		log.Printf("%s|err1=%+v\n",utils.FuncName(),err)
		panic(err)
	}
	defer client.Close()

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Printf("%s|err2=%+v\n",utils.FuncName(),err)
		panic(err)
	}
	defer consumer.Close()

	// get partitionId list

	partitions,err := consumer.Partitions("onConnect")
	if err != nil {
		log.Printf("%s|err3=%+v\n",utils.FuncName(),err)
		panic(err)
	}
	var wg sync.WaitGroup

	for k, partitionId := range partitions{
		log.Printf("key=%d|id=%d\n",k,partitionId)
		// create partitionConsumer for every partitionId
		partitionConsumer, err := consumer.ConsumePartition("onConnect", partitionId, sarama.OffsetNewest)
		if err != nil {
			log.Printf("%s|err4=%+v\n",utils.FuncName(),err)
			panic(err)
		}
		wg.Add(1)
		go func(pc *sarama.PartitionConsumer) {
			defer func() {
				wg.Done()
				if p := recover(); p != nil {
					utils.HandleThrow(p)
				}
			}()
			defer (*pc).AsyncClose()
			// block
			for message := range (*pc).Messages(){
				value := string(message.Value)
				key := string(message.Key)
				log.Printf("topic:%s,Partitionid: %d; offset:%d,key:%s, value: %s\n",message.Topic, message.Partition,message.Offset,key, value)
			}

		}(&partitionConsumer)
	}
	wg.Wait()
	log.Printf("NewDevice end\n")
}

type Kafka struct {
	brokers []string
	topics  []string
	//OffsetNewest int64 = -1
	//OffsetOldest int64 = -2
	startOffset       int64
	version           string
	ready             chan bool
	group             string
	channelBufferSize int
}

func NewKafka() *Kafka {
	return &Kafka{
		brokers:           brokers,
		topics:            []string{
			topics,
		},
		group:             group,
		channelBufferSize: 2,
		ready:             make(chan bool),
		version:"1.1.1",
	}
}

var brokers = []string{"127.0.0.1:9092"}
var topics = "onConnect"
var group = "39"

func (p *Kafka) Init() func() {
	log.Printf("kafka init...")

	version, err := sarama.ParseKafkaVersion(p.version)
	if err != nil {
		log.Fatalf("Error parsing Kafka version: %v", err)
	}
	config := sarama.NewConfig()
	config.Version = version
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange // 分区分配策略
	config.Consumer.Offsets.Initial = -2                    // 未找到组消费位移的时候从哪边开始消费
	config.ChannelBufferSize = p.channelBufferSize // channel长度

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(p.brokers, p.group, config)
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
			if err := client.Consume(ctx, p.topics, p); err != nil {
				log.Fatalf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				log.Println(ctx.Err())
				return
			}
			p.ready = make(chan bool)
		}
	}()
	<-p.ready
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
func (p *Kafka) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(p.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (p *Kafka) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (p *Kafka) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	// 具体消费消息
	for message := range claim.Messages() {
		msg := string(message.Value)
		log.Printf("%s|%+v",utils.FuncName(), msg)
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