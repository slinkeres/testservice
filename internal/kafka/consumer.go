package kafka

import (
	"context"
	"encoding/json"
	"log"
	"order-service/database"
	"order-service/internal/cache"
	"order-service/internal/model"
	"time"

	"github.com/segmentio/kafka-go"
)


type Consumer struct {
	reader *kafka.Reader
	cache *cache.Cache
	db		*database.Database
}

func NewConsumer(brokers []string, topic string, cache *cache.Cache, db *database.Database) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic: topic,
		GroupID: "order-service-group",
		MaxWait: 10 * time.Second,
	})
	return &Consumer{
		reader: reader,
		cache: cache,
		db: db,
	}
}

func (c *Consumer) Start(ctx context.Context) {
    log.Print("Запуск Kafka consumer...")

    for {
        select {
        case <- ctx.Done():
            log.Println("Остановка kafka consumer...")
            c.reader.Close()
            return
        default: 
            msg, err := c.reader.ReadMessage(ctx)
            if err != nil {
                log.Printf("Ошибка чтения сообщения: %v", err)
                continue
            }
            log.Printf("Полученное сообщение: %s", string(msg.Value))

            var order model.Order
            if err := json.Unmarshal(msg.Value, &order); err != nil {
                log.Printf("Ошибка парсинга сообщения: %v", err)
                continue
            }

            if err := c.db.SaveOrder(order); err != nil {
                log.Printf("Ошибка сохранения заказа в бд: %v", err)
                continue
            }

            c.cache.Set(order)


			 if err := c.reader.CommitMessages(ctx, msg); err != nil {
                log.Printf("Ошибка коммита сообщения: %v", err)
            }

            log.Printf("Заказ %v успешно обработан", order.OrderUID)
        }
    }
}


func (c *Consumer) Close() error {
	return c.reader.Close()
}