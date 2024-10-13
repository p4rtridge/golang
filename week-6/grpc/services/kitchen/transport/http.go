package transport

import (
	"context"
	"kitchen/services/kitchen/repository/rpc"
	"kitchen/services/orders/entity"
	"log"
	"net/http"
	"text/template"
	"time"
)

type httpTransport struct {
	repo rpc.OrderRepo
}

func NewHTTPTransport(repo rpc.OrderRepo) *httpTransport {
	return &httpTransport{
		repo,
	}
}

func (api *httpTransport) ServeHomepage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err := api.repo.CreateOrder(ctx, &entity.Order{
		OrderID:    1,
		CustomerID: 1,
		ProductID:  1,
		Quantity:   10,
	})
	if err != nil {
		log.Fatalf("Client error: %v", err)
	}

	res, err := api.repo.GetOrders(ctx, 1)
	if err != nil {
		log.Fatalf("Client error: %v", err)
	}

	t := template.Must(template.New("orders").Parse(ordersTemplate))

	if err := t.Execute(w, res); err != nil {
		log.Fatalf("template error: %v", err)
	}
}

var ordersTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Kitchen Orders</title>
</head>
<body>
    <h1>Orders List</h1>
    <table border="1">
        <tr>
            <th>Order ID</th>
            <th>Customer ID</th>
            <th>Quantity</th>
        </tr>
        {{range .}}
        <tr>
            <td>{{.OrderID}}</td>
            <td>{{.CustomerID}}</td>
            <td>{{.Quantity}}</td>
        </tr>
        {{end}}
    </table>
</body>
</html>`
