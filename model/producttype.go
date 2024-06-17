package model

type ProductTypeEntity struct {
	ID   int    `gorm:"primaryKey; column:prodtype_code;"`
	Name string `gorm:"not null;   column:prodtype_name;"`
}

func (p ProductTypeEntity) TableName() string {
	return "producttype"
}

type ProductType struct {
	ID   int    `json:"prodtype_id"`
	Name string `json:"prodtype_name"`
}

type ProductTypeCreate struct {
	ID   int    `json:"prodtype_id"      validate:"required,gte=0"`
	Name string `json:"prodtype_name"    validate:"required,max=40"`
}

type ProductTypeUpdate struct {
	Name string `json:"prodtype_name"    validate:"required,max=40"`
}
