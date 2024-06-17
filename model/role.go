package model

type RoleEntity struct {
	ID   	int    `gorm:"primaryKey; column:role_id;"`
	Title 	string `gorm:"not null;   column:role_title;"`
}

func (r RoleEntity) TableName() string {
	return "role"
}

type Role struct {
	ID  	int    
	Title 	string 
}