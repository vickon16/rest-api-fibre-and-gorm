package models

import "time"

type Product struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `json:"name"`
	SerialNumber string    `gorm:"unique" json:"serialNumber"`
	Price        float64   `json:"price"`
	CreatedAt    time.Time `json:"createdAt"`

	// Reverse relation: one product can appear in many orders
	Orders []Order `gorm:"foreignKey:ProductID" json:"orders,omitempty"`
}

type ProductSerializer struct {
	ID           uint              `json:"id"`
	Name         string            `json:"name"`
	SerialNumber string            `json:"serialNumber"`
	Price        float64           `json:"price"`
	Orders       []OrderSerializer `json:"orders,omitempty"`
}

type CreateProductDTO struct {
	Name         string  `json:"name" validate:"required"`
	SerialNumber string  `json:"serialNumber" validate:"required"`
	Price        float64 `json:"price" validate:"required,numeric"`
}

type UpdateProductDTO struct {
	Name         string  `json:"name"`
	SerialNumber string  `json:"serialNumber"`
	Price        float64 `json:"price" validate:"omitempty,numeric"`
}

func CreateResponseProduct(product Product) ProductSerializer {
	orders := make([]OrderSerializer, len(product.Orders))
	for i, order := range product.Orders {
		orders[i] = CreateResponseOrder(order)
	}

	return ProductSerializer{
		ID:           product.ID,
		Name:         product.Name,
		SerialNumber: product.SerialNumber,
		Price:        product.Price,
		Orders:       orders,
	}
}
