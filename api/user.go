package api

import "golang.org/x/crypto/bcrypt"

// User is the model struct for users
type User struct {
	Username       string `gorethink:"username" json:"username"`
	Email          string `gorethink:"email" json:"email"`
	Password       string `gorethink:"-" json:"-"`
	HashedPassword string `gorethink:"hashed_password" json:"-"`
}

func (user *User) encryptPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	user.HashedPassword = string(hashedPassword)
	return nil
}

// CreateUser inserts the user into the db
func CreateUser(user *User) error {
	err := user.encryptPassword()
	if err != nil {
		return err
	}
	_, err = UserTable.Insert(user).RunWrite(S)
	return err
}

// LoginUser returns a user from the db if the username and password matches
func LoginUser(username, password string) (*User, error) {
	res, err := UserTable.Get(username).Run(S)
	if err != nil {
		return nil, err
	}
	if res.IsNil() {
		return nil, nil
	}
	var user User
	err = res.One(&user)

	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return nil, nil
	}
	return &user, nil
}

// GetUser returns a user from the db
func GetUser(id string) (*User, error) {
	res, err := UserTable.Get(id).Run(S)
	if err != nil {
		return nil, err
	}
	if res.IsNil() {
		return nil, nil
	}
	var user User
	err = res.One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update returns a user from the db
func Update(user *User) (bool, error) {
	res, err := UserTable.Update(user).RunWrite(S)
	if err != nil {
		return false, err
	}
	return res.Updated > 0, nil
}
