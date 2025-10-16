package models

import "time"

type Order struct {
	ID uint `gorm:"primaryKey" json:"id"`

	ProductID uint    `json:"productId"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`

	UserID uint `json:"userId"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
}

type OrderSerializer struct {
	ID        uint              `json:"id"`
	Product   *ProductSerializer `json:"product"`
	User      *UserSerializer    `json:"user"`
	CreatedAt time.Time         `json:"createdAt"`
}

type CreateOrderDTO struct {
	ProductId uint `json:"productId" validate:"required"`
	UserId    uint `json:"userId" validate:"required"`
}

type UpdateOrderDTO struct {
	ProductId uint `json:"productId"`
	UserId    uint `json:"userId"`
}

func CreateResponseOrder(order Order) OrderSerializer {
	var product *ProductSerializer
	if order.Product.ID != 0 {
		p := CreateResponseProduct(order.Product)
		product = &p
	}

	var user *UserSerializer
	if order.User.ID != 0 {
		u := CreateResponseUser(order.User)
		user = &u
	}

	return OrderSerializer{
		ID:        order.ID,
		Product:   product,
		User:      user,
		CreatedAt: order.CreatedAt,
	}
}
