package local

type User struct {
	ID   string
	Name string
}

type UserUpdate struct {
	ID   string
	Name *string
}

type UserDTO struct {
	ID   string  `json:"id" binding:"required"`
	Name *string `json:"name,omitempty"`
}

type UserEntity struct {
	ID   string `firestore:"id"`
	Name string `firestore:"name"`
}

type UserMapper struct {
}

func (m UserMapper) ToDTO(model User) UserDTO {

	return UserDTO{
		ID:   model.ID,
		Name: &model.Name,
	}
}

func (m UserMapper) ToModel(dto UserDTO) User {

	user := User{
		ID: dto.ID,
	}

	if dto.Name != nil {
		user.Name = *dto.Name
	}

	return user
}

func (m UserMapper) ToUpdate(dto UserDTO) UserUpdate {
	return UserUpdate(dto)
}

func (m UserMapper) ToEntity(model User) UserEntity {
	return UserEntity(model)
}

func (m UserMapper) FromEntity(entity UserEntity) User {
	return User(entity)
}
