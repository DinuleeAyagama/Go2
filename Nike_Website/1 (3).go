// User model
type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
	IsAdmin  bool   `json:"is_admin" gorm:"default:false"`
}

// Product model
type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	Category    string  `json:"category"`
	Stock       int     `json:"stock"`
}

// Order model
type Order struct {
	gorm.Model
	UserID     uint        `json:"user_id"`
	User       User        `json:"-" gorm:"foreignKey:UserID"`
	Products   []Product   `json:"products" gorm:"many2many:order_products;"`
	Status     string      `json:"status"` // "pending", "completed", "cancelled"
	TotalPrice float64     `json:"total_price"`
}

// Cart is a temporary order with status "cart"
type Cart struct {
	Order
}