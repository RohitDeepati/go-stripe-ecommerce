package beans

type Users struct{
	UserID   int    `json:"userId"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:""`
	Role 	   string  `json:"role"`
}

type Loginuser struct{
	Email			string	`json:"email"`
	Password	string	`json:"password"`
}