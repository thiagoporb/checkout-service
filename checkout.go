package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
	"github.com/thiagoporb/checkout-service/queue"
	"net/http"
	"os"
	"time"
)

var connection *amqp.Channel

func init() {
	_ = os.Setenv("RABBITMQ_DEFAULT_USER", "rabbitmq")
	_ = os.Setenv("RABBITMQ_DEFAULT_PASS", "rabbitmq")
	_ = os.Setenv("RABBITMQ_DEFAULT_HOST", "localhost")
	_ = os.Setenv("RABBITMQ_DEFAULT_PORT", "5672")
	_ = os.Setenv("RABBITMQ_DEFAULT_VHOST", "/")
}

// Order Model
type Order struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	ProductId string    `json:"product_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at,string"`
}

func create(c echo.Context) error {
	order := new(Order)
	if err := c.Bind(order); err != nil {
		return err
	}
	order.ID = uuid.New()
	order.CreatedAt = time.Now()
	data, _ := json.Marshal(order)
	queue.Notify(data, "checkout_ex", "", connection)
	return c.JSON(http.StatusOK, order)
}

func getOrders(c echo.Context) error {
	var orders []Order
	return c.JSON(http.StatusOK, orders)
}

func main() {

	connection = queue.Connect()
	connection.ExchangeDeclare("checkout_ex", "direct", true, false, false, false, nil)
	connection.QueueBind("checkout_queue", "", "checkout_ex", false, nil)
	e := echo.New()
	e.GET("/checkouts", func(c echo.Context) error {
		return getOrders(c)
	})

	e.POST("/checkouts", func(c echo.Context) error {
		return create(c)
	})
	e.Logger.Fatal(e.Start(":8083"))

}
