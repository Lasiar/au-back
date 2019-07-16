package auth

// ChangeUser изменяет данные пользователя если он существует иначе добавляет нового
func (a *Auth) ChangeUser(user *User) (int, error) {
	var id int
	var err error
	if _, err = a.db.selectUserByID(user.ID); err != nil {
		user.PermID = 1
		id, err = a.db.insertUser(user)
		if err != nil {
			return -1, err
		}
	} else {
		err = a.db.updateUser(user)
		if err != nil {
			return -1, err
		}
		id = user.ID
	}
	return id, nil
}

// AddUser Регистрация пользователя
func (a *Auth) AddUser(user *User) (int, error) {
	user.ID = -1
	return a.ChangeUser(user)
}
