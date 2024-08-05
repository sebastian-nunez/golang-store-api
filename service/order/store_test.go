package order

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sebastian-nunez/golang-store-api/types"
)

func TestCreateOrder(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unable to stub db %s", err)
	}
	defer db.Close()

	store := NewStore(db)

	order := types.Order{
		UserId:  1,
		Total:   100.0,
		Status:  "Pending",
		Address: "123 Main St",
	}

	mockDb.ExpectExec("INSERT INTO orders").
		WithArgs(order.UserId, order.Total, order.Status, order.Address).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := store.CreateOrder(order)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
	if id != 1 {
		t.Errorf("expected id to be 1, but got %d", id)
	}

	if err := mockDb.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestCreateOrderItem(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unable to stub db %s", err)
	}
	defer db.Close()

	store := NewStore(db)

	orderItem := types.OrderItem{
		OrderId:   1,
		ProductId: 1,
		Quantity:  2,
		Price:     50.0,
	}

	mockDb.ExpectExec("INSERT INTO order_items").
		WithArgs(orderItem.OrderId, orderItem.ProductId, orderItem.Quantity, orderItem.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.CreateOrderItem(orderItem)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}

	if err := mockDb.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
