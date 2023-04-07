package core

import (
	"fmt"
	"strings"

	c "github.com/pergamenum/go-consensus-standards/constants"
	"github.com/pergamenum/go-consensus-standards/reflection"
	t "github.com/pergamenum/go-consensus-standards/types"
)

type User struct {
	ID   string `automap:"id"`
	Name string `automap:"name"`
}

type UserDTO struct {
	ID   string  `automap:"id" json:"id" binding:"required"`
	Name *string `automap:"name" json:"name,omitempty"`
}

type UserEntity struct {
	ID   string `automap:"id" firestore:"id"`
	Name string `automap:"name" firestore:"name"`
}

type UserManager struct {
	ttt map[string]string
	vos map[string]bool
}

func NewUserManager() *UserManager {

	vqk := reflection.MapTagToType("firestore", UserEntity{})
	vqk["created"] = "time"
	vqk["updated"] = "time"

	return &UserManager{
		ttt: vqk,
		vos: c.ValidRelationalOperators,
	}
}

func (m *UserManager) ToDTO(u User) UserDTO {

	ud := UserDTO{
		ID:   u.ID,
		Name: &u.Name,
	}

	return ud
}

func (m *UserManager) FromDTO(ud UserDTO) User {

	u := User{
		ID: ud.ID,
	}

	if ud.Name != nil {
		u.Name = *ud.Name
	}

	return u
}

func (m *UserManager) ToEntity(u User) UserEntity {

	return UserEntity(u)
}

func (m *UserManager) FromEntity(ue UserEntity) User {

	return User(ue)
}

// ToUpdate decides what fields get updated. Don't add a field you want protected.
func (m *UserManager) ToUpdate(ud UserDTO) t.Update {

	update := map[string]any{}

	update["id"] = ud.ID

	if ud.Name != nil {
		update["name"] = *ud.Name
	}

	return update
}

func (m *UserManager) ValidateModel(user User) error {

	var sb strings.Builder

	if user.ID == "" {
		sb.WriteString("(id: empty) ")
	}

	if len(user.Name) > 100 {
		sb.WriteString("(name: max 100 chars) ")
	}

	if len(sb.String()) > 0 {
		return fmt.Errorf("(invalid user: '%s')", strings.TrimSpace(sb.String()))
	}
	return nil
}

func (m *UserManager) ValidateUpdate(update t.Update) error {

	u := User{}

	if id, found := update["id"]; found {
		u.ID = id.(string)
	}

	if name, found := update["name"]; found {
		u.Name = name.(string)
	}

	return m.ValidateModel(u)
}

func (m *UserManager) ValidateQuery(queries []t.Query) error {

	for i := range queries {
		err := queries[i].Validate(m.ttt, m.vos)
		if err != nil {
			return err
		}
	}

	return nil
}
