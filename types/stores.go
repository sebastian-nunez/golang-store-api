package types

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(user User) (int, error)
	GetUsers() ([]User, error)
}

type ProductStore interface {
	GetProducts() ([]Product, error)
	GetProductByID(id int) (*Product, error)
	GetProductsByID(productIDs []int) ([]Product, error)
	CreateProduct(product CreateProductRequest) (int, error)
	UpdateProduct(product Product) error
}

type OrderStore interface {
	CreateOrder(Order) (int, error)
	CreateOrderItem(OrderItem) error
}
