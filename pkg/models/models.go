package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	CreatedAt time.Time      `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type Users struct {
	UserId     uuid.UUID   `json:"user_id,omitempty" gorm:"type:uuid;primaryKey"`
	FirstName  string      `json:"first_name,omitempty" gorm:"not null"`
	LastName   string      `json:"last_name,omitempty" gorm:"not null"`
	Email      string      `json:"email,omitempty" gorm:"unique;not null"`
	Phone      string      `json:"phone,omitempty" gorm:"unique;not null"`
	Password   string      `json:"password,omitempty" gorm:"not null"`
	Role       string      `json:"role,omitempty" gorm:"not null;check:role= 'user' or role= 'merchant' or role='admin'"`
	IsVerified bool        `json:"is_verified,omitempty" gorm:"not null"`
	Address    []Addresses `json:"address,omitempty" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Product    []Products  `json:"product,omitempty" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Order      []Orders    `json:"order,omitempty" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Base
}

type Addresses struct {
	AddressId uuid.UUID `json:"address_id,omitempty" gorm:"type:uuid;primaryKey;not null"`
	DoorNo    string    `json:"door_no,omitempty" gorm:"not null"`
	Street    string    `json:"street,omitempty" gorm:"not null"`
	City      string    `json:"city,omitempty" gorm:"not null"`
	State     string    `json:"state,omitempty" gorm:"not null"`
	ZipCode   uint      `json:"zip_code,omitempty" gorm:"not null"`
	UserId    uuid.UUID `json:"user_id,omitempty"`
	Order     Orders    `json:"-" gorm:"foreignKey:AddressId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Categories struct {
	CategoryId   uuid.UUID  `json:"category_id,omitempty" gorm:"type:uuid;primaryKey"`
	CategoryName string     `json:"category_name,omitempty" gorm:"not null"`
	Product      []Products `json:"product,omitempty" gorm:"foreignKey:CategoryId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Brands struct {
	BrandId   uuid.UUID  `json:"brand_id,omitempty" gorm:"type:uuid;primaryKey"`
	BrandName string     `json:"brand_name,omitempty" gorm:"not null"`
	Product   []Products `json:"product,omitempty" gorm:"foreignKey:BrandId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Products struct {
	ProductId   uuid.UUID `json:"product_id,omitempty" gorm:"type:uuid;primaryKey;not null"`
	ProductName string    `json:"product_name,omitempty" gorm:"not null"`
	CategoryId  uuid.UUID `json:"category_id,omitempty" gorm:"not null"`
	BrandId     uuid.UUID `json:"brand_id,omitempty" gorm:"not null"`
	UserId      uuid.UUID `json:"user_id,omitempty" gorm:"not null"`
	Price       float64   `json:"price,omitempty" gorm:"not null"`
	Rating      float32   `json:"rating,omitempty" gorm:"not null"`
	IsApproved  bool      `json:"is_Approved,omitempty" gorm:"not null"`
}

type Orders struct {
	OrderId     uuid.UUID      `json:"ordered_id,omitempty" gorm:"type:uuid;primaryKey"`
	UserId      uuid.UUID      `json:"user_id,omitempty" gorm:"not null"`
	AddressId   uuid.UUID      `json:"address_id,omitempty" gorm:"not null"`
	Name        string         `json:"name,omitempty" gorm:"not null"`
	Email       string         `json:"first_name,omitempty" gorm:"not null"`
	Phone       string         `json:"phone,omitempty" gorm:"not null"`
	Products    []OrderedItems `json:"products,omitempty" gorm:"foreignKey:OrderId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	TotalAmount float64        `json:"total_amount,omitempty" gorm:"null"`
	Status      string         `json:"status,omitempty"`
	CreatedAt   time.Time      `json:"created_at,omitempty" gorm:"autoCreateTime"`
}

type OrderedItems struct {
	OrderedItemsId uuid.UUID `json:"ordered_items_id,omitempty" gorm:"type:uuid;primaryKey"`
	ProductId      uuid.UUID `json:"product_id,omitempty" gorm:"not null"`
	ProductName    string    `json:"product_name,omitempty" gorm:"not null" `
	Quantity       uint      `json:"quantity,omitempty" gorm:"not null"`
	Price          float64   `json:"price,omitempty" gorm:"not null"`
	OrderId        uuid.UUID `json:"order_id,omitempty"`
}

func (user *Users) BeforeCreate(tx *gorm.DB) error {
	user.UserId = uuid.New()
	return nil
}

func (address *Addresses) BeforeCreate(tx *gorm.DB) error {
	address.AddressId = uuid.New()
	return nil
}

func (categories *Categories) BeforeCreate(tx *gorm.DB) error {
	categories.CategoryId = uuid.New()
	return nil
}

func (brand *Brands) BeforeCreate(tx *gorm.DB) error {
	brand.BrandId = uuid.New()
	return nil
}

func (product *Products) BeforeCreate(tx *gorm.DB) error {
	product.ProductId = uuid.New()
	return nil
}

func (order *Orders) BeforeCreate(tx *gorm.DB) error {
	order.OrderId = uuid.New()
	return nil
}

func (orderItem *OrderedItems) BeforeCreate(tx *gorm.DB) error {
	orderItem.OrderedItemsId = uuid.New()
	return nil
}
