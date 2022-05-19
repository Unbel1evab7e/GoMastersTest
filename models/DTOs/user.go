package DTOs

// swagger:model User
type User struct {
	Firstname string `json:"Firstname" example:"Igor"`
	Lastname  string `json:"Lastname" example:"Kormich"`
	Email     string `json:"Email" example:"unbel1evableik@gmail.com"`
	Age       uint   `json:"Age" example:"1" format:"uint"`
}
