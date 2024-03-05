package customers

type CreateTypeDTO struct {
	EntityType string `json:"entity_type"`
	ContactName string `json:"contact_name"`
	Address string `json:"address"`
	Phone string `json:"phone"`
	UserPriority string `json:"user_priority"`
	Login string `json:"login"`
	Password string `json:"password"`
	EMail string `json:"e_mail"`
}
