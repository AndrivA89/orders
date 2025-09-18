package services

// CreateUserRequest объединяет параметры для создания пользователя
type CreateUserRequest struct {
	FirstName string
	LastName  string
	Age       int
	IsMarried bool
	Password  string
}

// CreateProductRequest объединяет параметры для создания товара
type CreateProductRequest struct {
	Description string
	Tags        []string
	Quantity    int
	Price       int64
}
