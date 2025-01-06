package models

import (
	"shopping-site/utils/constants"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// base model details
type Base struct {
	CreatedAt time.Time      `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

// user details
type Users struct {
	UserID     uuid.UUID   `json:"user_id,omitempty" gorm:"type:uuid;primaryKey"`
	FirstName  string      `json:"first_name,omitempty" gorm:"not null"`
	LastName   string      `json:"last_name,omitempty" gorm:"not null"`
	Email      string      `json:"email,omitempty" gorm:"unique;not null"`
	Phone      string      `json:"phone,omitempty" gorm:"unique;not null"`
	Password   string      `json:"password,omitempty" gorm:"not null"`
	Role       string      `json:"role,omitempty" gorm:"not null;check:role= 'user' or role= 'merchant' or role='admin'"`
	IsVerified bool        `json:"is_verified,omitempty" gorm:"not null"`
	Address    []Addresses `json:"address,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Product    []Products  `json:"product,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Order      []Orders    `json:"order,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Base
}

// address details
type Addresses struct {
	AddressID uuid.UUID `json:"address_id,omitempty" gorm:"type:uuid;primaryKey;not null"`
	DoorNo    string    `json:"door_no,omitempty" gorm:"not null"`
	Street    string    `json:"street,omitempty" gorm:"not null"`
	City      string    `json:"city,omitempty" gorm:"not null"`
	State     string    `json:"state,omitempty" gorm:"not null"`
	ZipCode   uint      `json:"zip_code,omitempty" gorm:"not null"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
	Order     Orders    `json:"-" gorm:"foreignKey:AddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// category details
type Categories struct {
	CategoryID   uuid.UUID  `json:"category_id,omitempty" gorm:"type:uuid;primaryKey"`
	CategoryName string     `json:"category_name,omitempty" gorm:"not null"`
	Product      []Products `json:"-" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// brand details
type Brands struct {
	BrandID   uuid.UUID  `json:"brand_id,omitempty" gorm:"type:uuid;primaryKey"`
	BrandName string     `json:"brand_name,omitempty" gorm:"not null"`
	Product   []Products `json:"-" gorm:"foreignKey:BrandID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// product details
type Products struct {
	ProductID   uuid.UUID `json:"product_id,omitempty" gorm:"type:uuid;primaryKey;not null"`
	ProductName string    `json:"product_name,omitempty" gorm:"not null"`
	CategoryID  uuid.UUID `json:"category_id,omitempty" gorm:"not null"`
	BrandID     uuid.UUID `json:"brand_id,omitempty" gorm:"not null"`
	UserID      uuid.UUID `json:"user_id,omitempty" gorm:"not null"`
	Price       float64   `json:"price,omitempty" gorm:"not null"`
	Rating      float32   `json:"rating,omitempty" gorm:"not null"`
	IsApproved  bool      `json:"is_Approved,omitempty" gorm:"not null"`
}

// order details
type Orders struct {
	OrderID     uuid.UUID      `json:"ordered_id,omitempty" gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID      `json:"user_id,omitempty" gorm:"not null"`
	AddressID   uuid.UUID      `json:"address_id,omitempty" gorm:"not null"`
	Name        string         `json:"name,omitempty" gorm:"not null"`
	Email       string         `json:"first_name,omitempty" gorm:"not null"`
	Phone       string         `json:"phone,omitempty" gorm:"not null"`
	Products    []OrderedItems `json:"products,omitempty" gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	TotalAmount float64        `json:"total_amount,omitempty" gorm:"null"`
	Status      string         `json:"status,omitempty"`
	CreatedAt   time.Time      `json:"created_at,omitempty" gorm:"autoCreateTime"`
}

// order items details
type OrderedItems struct {
	OrderedItemsID uuid.UUID `json:"ordered_items_id,omitempty" gorm:"type:uuid;primaryKey"`
	ProductID      uuid.UUID `json:"product_id,omitempty" gorm:"not null"`
	MerchantId     uuid.UUID `json:"merchant_id,omitempty" gorm:"not null"`
	ProductName    string    `json:"product_name,omitempty" gorm:"not null" `
	Quantity       uint      `json:"quantity,omitempty" gorm:"not null"`
	Price          float64   `json:"price,omitempty" gorm:"not null"`
	Status         string    `json:"status,omitempty"`
	OrderID        uuid.UUID `json:"order_id,omitempty"`
	CustomerID     uuid.UUID `json:"customer_id,omitempty"`
	AddressID      uuid.UUID `json:"address_id,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}

// hooks for uuid generation
func (user *Users) BeforeCreate(tx *gorm.DB) error {
	user.UserID = uuid.New()
	return nil
}

// hooks for uuid generation
func (address *Addresses) BeforeCreate(tx *gorm.DB) error {
	address.AddressID = uuid.New()
	return nil
}

// hooks for uuid generation
func (categories *Categories) BeforeCreate(tx *gorm.DB) error {
	categories.CategoryID = uuid.New()
	return nil
}

// hooks for uuid generation
func (brand *Brands) BeforeCreate(tx *gorm.DB) error {
	brand.BrandID = uuid.New()
	return nil
}

// hooks for uuid generation
func (product *Products) BeforeCreate(tx *gorm.DB) error {
	product.ProductID = uuid.New()
	return nil
}

// hooks for uuid generation
func (order *Orders) BeforeCreate(tx *gorm.DB) error {
	order.OrderID = uuid.New()
	return nil
}

// hooks for uuid generation
func (orderItem *OrderedItems) BeforeCreate(tx *gorm.DB) error {
	orderItem.OrderedItemsID = uuid.New()
	return nil
}

// hooks to update order items status and created time  whent status updated on orders table
func (order *Orders) AfterCreate(tx *gorm.DB) error {
	if order.Status == constants.Inprogress {
		tx.Model(&OrderedItems{}).Where("order_id = ?", order.OrderID).Update("status", constants.Placed)
	}

	if !order.CreatedAt.IsZero() {
		tx.Model(&OrderedItems{}).Where("order_id = ?", order.OrderID).Update("created_at", order.CreatedAt)
	}
	return nil
}

// hooks to update order items status whent status updates on orders table
func (order *Orders) AfterUpdate(tx *gorm.DB) error {
	if order.Status == constants.Cancelled {
		tx.Model(&OrderedItems{}).Where("order_id = ?", order.OrderID).Update("status", constants.Cancelled)
	}
	return nil
}
