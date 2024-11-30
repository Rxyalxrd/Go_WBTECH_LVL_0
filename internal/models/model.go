package models

// Delivery - таблица для информации о доставке
type Delivery struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	OrderID uint   `gorm:"index"`
	Name    string `gorm:"type:varchar(255)"`
	Phone   string `gorm:"type:varchar(255)"`
	Zip     string `gorm:"type:varchar(255)"`
	City    string `gorm:"type:varchar(255)"`
	Address string `gorm:"type:varchar(255)"`
	Region  string `gorm:"type:varchar(255)"`
	Email   string `gorm:"type:varchar(255)"`
}

// Payment - таблица для информации об оплате
type Payment struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	OrderID      uint   `gorm:"index"`
	Transaction  string `gorm:"type:varchar(255)"`
	RequestID    string `gorm:"type:varchar(255)"`
	Currency     string `gorm:"type:varchar(255)"`
	Provider     string `gorm:"type:varchar(255)"`
	Amount       int    `gorm:"type:int"`
	PaymentDT    int64  `gorm:"type:int"`
	Bank         string `gorm:"type:varchar(255)"`
	DeliveryCost int    `gorm:"type:int"`
	GoodsTotal   int    `gorm:"type:int"`
	CustomFee    int    `gorm:"type:int"`
}

// Item - таблица для информации о товарах
type Item struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	OrderID     uint   `gorm:"index"`
	ChrtID      int    `gorm:"type:int"`
	TrackNumber string `gorm:"type:varchar(255)"`
	Price       int    `gorm:"type:int"`
	Rid         string `gorm:"type:varchar(255)"`
	Name        string `gorm:"type:varchar(255)"`
	Sale        int    `gorm:"type:int"`
	Size        string `gorm:"type:varchar(255)"`
	TotalPrice  int    `gorm:"type:int"`
	NmID        int    `gorm:"type:int"`
	Brand       string `gorm:"type:varchar(255)"`
	Status      int    `gorm:"type:int"`
}

// Order - таблица для заказов
type Order struct {
	ID                uint     `json:"id" gorm:"primaryKey"`
	OrderUID          string   `json:"order_uid" gorm:"unique;not null"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	DeliveryID        uint     `json:"-"`
	PaymentID         uint     `json:"-"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	Shardkey          string   `json:"shardkey"`
	SmID              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
	Delivery          Delivery `json:"delivery" gorm:"foreignKey:DeliveryID"`
	Payment           Payment  `json:"payment" gorm:"foreignKey:PaymentID"`
	Items             []Item   `json:"items" gorm:"foreignKey:OrderID"`
}
