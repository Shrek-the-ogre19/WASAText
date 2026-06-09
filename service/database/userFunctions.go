package database

import "strconv"

func (db *appdbimpl) CountUsers() (int, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, err
}

func (db *appdbimpl) AddUser(name string) error {
	id, err := db.CountUsers()
	if err != nil {
		return err
	}
	_, err = db.c.Exec("INSERT INTO users (id, name, picture, conversations) VALUES (?, ?, ?, ?)", id+1, name, "default", "")
	return err
}

func (db *appdbimpl) ChangeUserName(id UserId, name string) error {
	_, err := db.c.Exec(`UPDATE users
SET name = ?
WHERE id = ?;`, name, id.Id)
	return err
}

func (db *appdbimpl) ChangeUserPicture(id UserId, picture string) error {
	_, err := db.c.Exec(`UPDATE users
SET picture = ?
WHERE id = ?;`, picture, id.Id)
	return err
}

func (db *appdbimpl) GetUserId(name string) (UserId, error) {
	var id int
	err := db.c.QueryRow("SELECT id FROM users WHERE name = ?", name).Scan(&id)
	return UserId{id}, err
}

func (db *appdbimpl) GetUser(id UserId) (User, error) {
	var user User
	var name string
	var picture string
	var conversationsR string
	var conversations []ConversationId
	err := db.c.QueryRow("SELECT name, picture, conversations FROM users WHERE id = ?", id.Id).Scan(&name, &picture, &conversationsR)
	conversations = ConvertConversations(conversationsR)
	if err == nil {
		user = User{id, name, picture, conversations}
	}
	return user, err
}

func (db *appdbimpl) ListAllUsers() ([]string, error) {
	var users []string
	count, err := db.CountUsers()
	var id int = 1
	var name string
	for count > 0 {
		err = db.c.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&name)
		users = append(users, name)
		id = id + 1
		count = count - 1
	}
	return users, err
}

func (db *appdbimpl) UserLookup(name string) (bool, error) {
	var bool bool
	err := db.c.QueryRow("SELECT 1 From users WHERE name = ?", name).Scan(&bool)
	if err == nil {
		return bool, err
	}
	return false, nil
}

func (db *appdbimpl) AddConversations(id UserId, newConvId int) error {
	var conversations string
	err := db.c.QueryRow("SELECT conversations FROM users WHERE id = ?", id.Id).Scan(&conversations)
	if err != nil {
		return err
	}
	if conversations != "" {
		conversations = conversations + ","
	}
	conversations = conversations + strconv.Itoa(newConvId)
	_, err = db.c.Exec(`UPDATE users
SET conversations = ?
WHERE id = ?;`, conversations, id.Id)
	return err
}
