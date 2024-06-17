package model

type ProductTypeResponse struct {
	Code 	int 			`json:"code"`
	Message *ProductType 	`json:"message"`
}

type ProductTypesResponse struct {
	Code 	int 			`json:"code"`
	Message []ProductType 	`json:"message"`
}

type AuthPassportResponse struct {
	Code 	int 			`json:"code"`
	Message *UserPassport 	`json:"message"`
}

type CountResponse struct {
	Code 	int 	`json:"code"`
	Message int 	`json:"message"`
}

type StringResponse struct {
	Code 	int 	`json:"code"`
	Message string 	`json:"message"`
}


