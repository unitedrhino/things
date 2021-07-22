package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"strings"
)

func main() {
	ip := "192.168.1.2:803"
	addr := string([]byte(ip)[0:strings.LastIndex(ip, ":")])
	fmt.Printf("ip=%s|addr=%s\n", ip, addr)
}
func subscripe() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	pubsub := client.Subscribe("chat")
	defer pubsub.Close()
	for msg := range pubsub.Channel() {
		fmt.Printf("channel=%s message=%s\n", msg.Channel, msg.Payload)
	}
}
func subscribe() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	pubsub := client.PSubscribe("*")
	defer pubsub.Close()
	for msg := range pubsub.Channel() {
		fmt.Printf("channel=%s message=%s\n", msg.Channel, msg.Payload)
	}

}
