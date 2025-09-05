package main

import (
	"context"
	"log"
	"net/http"
	"order-service/database"
	"order-service/internal/cache"
	"order-service/internal/handler"
	"order-service/internal/kafka"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка загрузки .env файла")
		log.Fatalf(err.Error())
	}
	

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "kafka:9092"
	}

	db, err := database.NewPostrgesDB()
	if err != nil {
		log.Fatalf("Не удалось подключиться к бд: %v", err)
	}
	defer db.Close()

	cache := cache.NewCache()

	log.Println("Восстанавливаем кэш из бд...")
	orders, err :=  db.GetAllOrders()
	if err != nil {
		log.Printf("Не удалось восстановить кэш: %v", err)
	} else {
		cache.Restore(orders)
	log.Printf("Восстановлено %d заказов из кэша", len(orders))
	}
	kafkaConsumer := kafka.NewConsumer(
		[]string{kafkaBrokers},
		"orders",
		cache,
		db,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go kafkaConsumer.Start(ctx)
	defer kafkaConsumer.Close()

	router := mux.NewRouter()

	orderHandler := handler.NewOrderHandler(cache)
	orderHandler.RegisterRoutes(router)


	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    	http.ServeFile(w, r, "./static/index.html")
	})

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))
	

	



	server := &http.Server{
		Addr: ":8080",
		Handler: router,
	}

	go func() {
		log.Println("Запуск HTTP сервера на :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed{
			log.Fatalf("Не удалось запустить HTTP сервер: %v",err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Завершение работы сервера")

		
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Сервер принудительно выключен: %v", err)
	}
	
	log.Println("Сервер вышел из строя")

}
