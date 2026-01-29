package database

func (db *appdbimpl) CountUsers()(int, error){
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, err
}
func (db *appdbimpl) AddUser(name string)(error){
	id, err := db.CountUsers()
	var picture = "this will be the default picture"
	_, err = db.c.Exec("INSERT INTO users (id, name, picture) VALUES (?, ?, ?)", id+1,name, picture)
	return err
}


func (db *appdbimpl) ChangeUserName(id UserId, name string)(error){
	_, err := db.c.Exec(`UPDATE users
SET name = ?
WHERE id = ?;`, name, id.Id)
	return err
}

func (db *appdbimpl) ChangeUserPicture(id UserId, picture string)(error){
	_, err := db.c.Exec(`UPDATE users
SET picture = ?
WHERE id = ?;`, picture, id.Id)
	return err
}

func (db *appdbimpl) GetUserId(name string) (UserId, error) {
	var id int
	err := db.c.QueryRow("SELECT id FROM users WHERE name = ?", name).Scan(&id)
	var UserId = UserId{id}
	return UserId, err
}


func (db *appdbimpl) GetUser(id UserId) (User, error) {
	var user User
	var name string
	var picture string
	err := db.c.QueryRow("SELECT name, picture FROM users WHERE id = ?", id.Id).Scan(&name, &picture)
	if err == nil {
		user = User{id, name, picture}
	}
	return user, err
}

func (db *appdbimpl) ListAllUsers() ([]string, error){
	var users []string
	count, err := db.CountUsers()
	var id int = 1
	var name string
	for count > 0{
		err = db.c.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&name)
		users = append(users, name)
		id = id+1
		count = count-1
	}
	return users, err
}

func (db *appdbimpl) UserLookup(name string) (bool, error){
	var bool bool
	err := db.c.QueryRow("SELECT 1 From users WHERE name = ?", name).Scan(&bool)
	if err == nil{
		return bool, err
	}
	return false, nil
}
