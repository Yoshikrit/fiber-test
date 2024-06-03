package model

type ProductTypeEntity struct {
	ID   int    `gorm:"primaryKey; column:prodtype_code;"`
	Name string `gorm:"not null;   column:prodtype_name;"`
}

func (p ProductTypeEntity) TableName() string {
	return "producttype"
}

type ProductType struct {
	ID   int    `json:"ProdType_ID"`
	Name string `json:"ProdType_Name"`
}

type ProductTypeCreate struct {
	ID   int    `json:"ProdType_ID"      validate:"required,gte=0"`
	Name string `json:"ProdType_Name"    validate:"required,max=40"`
}

type ProductTypeUpdate struct {
	Name string `json:"ProdType_Name"    validate:"required,max=40"`
}
