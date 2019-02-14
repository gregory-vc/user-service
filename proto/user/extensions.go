package user

func (model *User) BeforeCreate() {
	model.Id = 1
}
