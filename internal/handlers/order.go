package handlers

import (
	"WB_Intern/internal/model"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"strings"
)

type Order struct {
	service orderService
	logger *logrus.Logger
}

type orderService interface {
	GetOrderById(string) (*model.Order, error)
}

func NewOrderHandler(logger *logrus.Logger, service orderService) *Order {
	return &Order{
		service: service,
		logger:  logger,
	}
}

func (o *Order) ShowForm(w http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.ParseFiles("../internal/templates/form.html"))

	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render form template", http.StatusInternalServerError)
		return
	}
}

type Response struct {
	ID string
	Data string
}

func (o *Order) ShowOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.FormValue("order_id")

	body, err := o.service.GetOrderById(orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data, err := json.Marshal(*body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
	}

	response := Response{
		ID: "Received order ID: " + body.OrderUID,
		Data: strings.ReplaceAll(string(data), ",", "\n"),
	}


	tmpl := template.Must(template.ParseFiles("../internal/templates/order.html"))
	err = tmpl.Execute(w, response)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}