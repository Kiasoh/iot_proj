package main

import (
	"context"
	"fmt"
	"iot_proj/handlers"
	"iot_proj/repository"
	"iot_proj/services"
	"log"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectSQL() *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		"niflheim", "niflguard", "iot_postgres", "5433", "iot_db")

	poolconfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	ctx := context.Background()

	pool, err := pgxpool.NewWithConfig(ctx, poolconfig)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	if err = pool.Ping(ctx); err != nil {
		log.Fatal(err)
		panic(err)
	}
	return pool
}
func ConnectMosquitto(handler *handlers.MQTTHandler) mqtt.Client {
	client := mqtt.NewClient(mqtt.NewClientOptions().
		SetClientID("backend0").
		AddBroker(fmt.Sprintf("tcp://%s:%s", "mosquitto", "1883")).
		SetAutoReconnect(true).
		SetCleanSession(true).SetOnConnectHandler(func(c mqtt.Client) {
		log.Println("MQTT connected")
		go handler.SubscribeToCardScans(c)
	}).
		SetConnectionLostHandler(func(c mqtt.Client, err error) {
			log.Printf("MQTT connection lost: %v", err)
		}))
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}

func InitRoute(handler *handlers.MQTTHandler) http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))
	mux.Post("/login", handler.Login)

	mux.Group(func(r chi.Router) {
		r.Use(handler.AuthMiddleware)

		r.Get("/profile", handler.Profile)
		r.Get("/entry-logs", handler.GetEntryLogs)

		r.Group(func(adminRouter chi.Router) {
			adminRouter.Use(handler.AdminMiddleware)
			adminRouter.Post("/register", handler.Register)
			adminRouter.Get("/users", handler.GetUsers)
			adminRouter.Put("/users/{userID}/key-card", handler.UpdateUserKeyCard)
			adminRouter.Put("/users/{userID}/access-level", handler.UpdateUserAccessLevel)
		})
	})

	return mux
}
func main() {
	db := ConnectSQL()
	repo := &repository.Repository{DB: db}
	service := services.NewService(repo)
	handler := handlers.NewMQTTHandler(service)

	mb := ConnectMosquitto(handler)
	service.SetMqttClient(mb)

	server := &http.Server{
		Addr:    ":8081",
		Handler: InitRoute(handler),
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
