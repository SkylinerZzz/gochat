package queue

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

// RedisNode definition
type RedisNode struct {
	redisPool   *redis.Pool // redis connection pool
	redisClient *redis.Conn // redis connection
}

func NewRedisNode(pool *redis.Pool, client *redis.Conn) *RedisNode {
	return &RedisNode{
		redisPool:   pool,
		redisClient: client,
	}
}

func (r *RedisNode) ReceiveMessage(queueName string, timeout time.Duration) (Message, error) {
	conn := r.redisPool.Get()
	defer conn.Close()
	data, err := redis.Strings(conn.Do("BLPOP", queueName, int(timeout.Seconds())))
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return Message{}, errors.New("queue is empty")
		}
		log.Errorf("[ReceiveMessage] failed to receive message, err = %s", err)
		return Message{}, err
	}
	message := Message{}
	err = json.Unmarshal([]byte(data[1]), &message)
	if err != nil {
		log.Errorf("[ReceiveMessage] failed to unmarshal message, err = %s", err)
		return Message{}, err
	}
	log.WithFields(log.Fields{
		"queueName": queueName,
		"message":   data[1],
	}).Info("[ReceiveMessage] receive message successfully")
	return message, nil
}

func (r *RedisNode) SendMessage(queueName string, message Message) error {
	conn := r.redisPool.Get()
	defer conn.Close()
	data, err := json.Marshal(message)
	if err != nil {
		log.Errorf("[SendMessage] failed to marshal message, err = %s", err)
		return err
	}
	_, err = conn.Do("RPUSH", queueName, data)
	if err != nil {
		log.Errorf("[SendMessage] failed to send message, err = %s", err)
		return err
	}
	log.WithFields(log.Fields{
		"queueName": queueName,
		"message":   string(data),
	}).Info("[SendMessage] send message successfully")
	return nil
}

func (r *RedisNode) Subscribe(channel string) <-chan Message {
	subChan := make(chan Message, 10)
	pubSubConn := redis.PubSubConn{Conn: *r.redisClient} // pubsub connection
	err := pubSubConn.Subscribe(channel)                 // subscribe given channel
	if err != nil {
		log.Errorf("[Subscribe] failed to subscribe channel, err = %s", err)
		return nil
	}

	// constantly receive subscription message
	go func() {
		for {
			switch msg := pubSubConn.Receive().(type) {
			case redis.Subscription:
				log.WithFields(log.Fields{
					"channel": msg.Channel,
				}).Info("[Subscribe] subscribe channel successfully")
			case redis.Message:
				log.WithFields(log.Fields{
					"channel": msg.Channel,
				}).Info("[Subscribe] receive message successfully")
				message := Message{}
				err = json.Unmarshal(msg.Data, &message)
				if err != nil {
					log.Errorf("[Subscribe] failed to unmarshal message, err = %s", err)
				}
				// send subscription message to chan
				subChan <- message
			case error:
				log.Errorf("[Subscribe] err = %s", msg)
			}
		}
	}()
	return subChan
}

func (r *RedisNode) Publish(channel string, message Message) error {
	conn := r.redisPool.Get()
	defer conn.Close()
	data, err := json.Marshal(message)
	if err != nil {
		log.Errorf("[Publish] failed to marshal message, err = %s", err)
		return err
	}
	_, err = conn.Do("PUBLISH", channel, data)
	if err != nil {
		log.Errorf("[Publish] failed to publish message, err = %s", err)
		return err
	}
	log.WithFields(log.Fields{
		"channel": channel,
		"message": string(data),
	}).Info("[Publish] publish message successfully")
	return nil
}
