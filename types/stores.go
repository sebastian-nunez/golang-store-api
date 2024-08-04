package types

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int) (*User, error)
	CreateUser(user User) error
}

type ProductStore interface {
	GetProducts() ([]Product, error)
}
