package model

type Admin struct {
	AdminID string `json:"admin_id"`
	Name    string `json:"name"`
}

func NewAdmin(adminID, name string) *Admin {
	return &Admin{
		AdminID: adminID,
		Name:    name,
	}
}
