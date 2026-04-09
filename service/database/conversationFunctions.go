package database

import (
	"strconv"
	"strings"
	"time"
)

func (db *appdbimpl) GetConversations(id UserId) ([]ConversationId, error) {
	var conversationsR string
	var conversations []ConversationId
	err := db.c.QueryRow("SELECT conversations FROM users WHERE id = ?", id.Id).Scan(&conversationsR)
	conversations = ConvertConversations(conversationsR)
	return conversations, err
}

func (db *appdbimpl) CountConversations() (int, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM conversations").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, err
}

func (db *appdbimpl) StartConversation(sender UserId, receivers []string) (ConversationId, error) {
	id, err := db.CountConversations()
	var snippet string = ""
	var groupchat bool
	var content = ""
	var members string
	var picture string
	var name string = ""
	var date = time.Now().String()
	var senderO User
	senderO, err = db.GetUser(sender)
	senderName := senderO.Name
	senderIdstr := strconv.Itoa(sender.Id)
	if len(receivers) == 1 {
		groupchat = false

		receiverId, err := db.GetUserId(receivers[0])
		if err == nil {
			receiverIdstr := strconv.Itoa(receiverId.Id)
			members = senderIdstr + "," + receiverIdstr
		}
		name = senderName + "," + receivers[0]

		picture = name
	} else {
		groupchat = true

		var size int
		size = len(receivers)
		name = senderName
		members = senderIdstr
		for i := 0; i < size; i++ {
			receiverId, err := db.GetUserId(receivers[i])
			if err == nil {
				receiverIdstr := strconv.Itoa(receiverId.Id)
				members = members + "," + receiverIdstr
			}
			name = name + "," + receivers[i]
		}

		picture = "default"
	}

	var membersId []UserId = ConvertUsers(members)
	for i := 0; i < len(membersId); i++ {
		db.AddConversations(membersId[i], id+1)

	}

	_, err = db.c.Exec("INSERT INTO conversations (id, snippet, name, picture, date, content, groupchat, members) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", id+1, snippet, name, picture, date, content, groupchat, members)
	var conversationId = ConversationId{id + 1}
	return conversationId, err
}

func (db *appdbimpl) GetConversation(id ConversationId) (Conversation, error) {
	var conversation Conversation
	var groupchat bool
	var picture string
	var name string
	var snippet string
	var content []MessageId
	var members []UserId
	var date string
	var contentR string
	var membersR string

	err := db.c.QueryRow("SELECT  snippet, name, picture, date, content, groupchat, members FROM conversations WHERE id = ?", id.Id).Scan(&snippet, &name, &picture, &date, &contentR, &groupchat, &membersR)
	content = ConvertMessages(contentR)
	members = ConvertUsers(membersR)

	if err == nil {
		conversation = Conversation{id, snippet, name, picture, date, content, groupchat, members}
	}

	return conversation, err
}

func (db *appdbimpl) ChangeConversationName(id ConversationId, NewName string) error {
	var groupchat bool
	err := db.c.QueryRow("SELECT  groupchat FROM conversations WHERE id = ?", id.Id).Scan(&groupchat)

	if groupchat {
		_, err = db.c.Exec(`UPDATE conversations
SET name = ?
WHERE id = ?;`, NewName, id.Id)
	}
	return err
}

func (db *appdbimpl) ChangeConversationPhoto(id ConversationId, NewPicture string) error {
	var groupchat bool
	err := db.c.QueryRow("SELECT  groupchat FROM conversations WHERE id = ?", id.Id).Scan(&groupchat)

	if groupchat {
		_, err = db.c.Exec(`UPDATE conversations
SET picture = ?
WHERE id = ?;`, NewPicture, id.Id)
	}
	return err
}

func (db *appdbimpl) GetAllMembers(id ConversationId) ([]string, error) {
	var membersR string
	var members []UserId
	var memberNames []string
	err := db.c.QueryRow("SELECT  members FROM conversations WHERE id = ?", id.Id).Scan(&membersR)
	if err == nil {
		members = ConvertUsers(membersR)
		for i := 0; i < len(members); i++ {
			var user User
			user, err = db.GetUser(members[i])
			memberNames = append(memberNames, user.Name)
		}
	}
	return memberNames, err
}

func (db *appdbimpl) AddMembers(id ConversationId, Name string) error {
	var groupchat bool
	var membersR string
	err := db.c.QueryRow("SELECT  members FROM conversations WHERE id = ?", id.Id).Scan(&membersR)
	err = db.c.QueryRow("SELECT  groupchat FROM conversations WHERE id = ?", id.Id).Scan(&groupchat)
	var newUserId UserId
	newUserId, err = db.GetUserId(Name)

	membersR = membersR + "," + strconv.Itoa(newUserId.Id)
	_, err = db.c.Exec(`UPDATE conversations
SET members = ?
WHERE id = ?;`, membersR, id.Id)

	var conversations string
	err = db.c.QueryRow("SELECT conversations FROM users WHERE id = ?", newUserId.Id).Scan(&conversations)
	if conversations != "" {
		conversations = conversations + ","
	}
	conversations = conversations + strconv.Itoa(id.Id)
	_, err = db.c.Exec(`UPDATE users
SET conversations = ?
WHERE id = ?;`, conversations, newUserId.Id)

	return err
}

func (db *appdbimpl) LeaveGroup(convid ConversationId, userid UserId) error {
	var membersR string
	useridstr := strconv.Itoa(userid.Id)
	err := db.c.QueryRow("SELECT  members FROM conversations WHERE id = ?", convid.Id).Scan(&membersR)
	members := strings.Split(membersR, ",")
	var newmembers string
	for i := 0; i < len(members); i++ {
		if members[i] != useridstr {
			newmembers = newmembers + "," + members[i]
		}

	}

	_, err = db.c.Exec(`UPDATE conversations
SET members = ?
WHERE id = ?;`, newmembers, convid.Id)

	var conversationsR string
	convidstr := strconv.Itoa(convid.Id)
	err = db.c.QueryRow("SELECT  conversations FROM users WHERE id = ?", userid.Id).Scan(&conversationsR)
	conversations := strings.Split(conversationsR, ",")
	var newconversations string
	for i := 0; i < len(conversations); i++ {
		if conversations[i] != convidstr {
			newconversations = newconversations + "," + conversations[i]
		}

	}
	_, err = db.c.Exec(`UPDATE users
SET conversations = ?
WHERE id = ?;`, newconversations, userid.Id)

	return err
}

func (db *appdbimpl) CheckGroupchat(id ConversationId) (bool, error) {
	var groupchat bool
	err := db.c.QueryRow("SELECT  groupchat FROM conversations WHERE id = ?", id.Id).Scan(&groupchat)
	if groupchat {
		return true, err
	} else {
		return false, err
	}
}

func (db *appdbimpl) UpdateConversationHead(id ConversationId, snippetNew string, dateNew string) error {
	_, err := db.c.Exec(`UPDATE conversations
SET snippet = ?, date = ?
WHERE id = ?;`, snippetNew, dateNew, id.Id)
	return err
}

func (db *appdbimpl) UpdateContent(id ConversationId, newMessageId int) error {
	var contentR string
	err := db.c.QueryRow("SELECT  content FROM conversations WHERE id = ?", id.Id).Scan(&contentR)
	contentR = contentR + "," + strconv.Itoa(newMessageId)

	_, err = db.c.Exec(`UPDATE conversations
SET content = ?
WHERE id = ?;`, contentR, id.Id)
	return err
}

func (db *appdbimpl) RemoveFromContent(id ConversationId, oldMessageId int) error {
	var contentR string
	err := db.c.QueryRow("SELECT  content FROM conversations WHERE id = ?", id.Id).Scan(&contentR)
	content := ConvertMessages(contentR)
	var newcontent string
	for i := 0; i < len(content); i++ {
		if content[i].Id != oldMessageId {

			newcontent = newcontent + "," + strconv.Itoa(content[i].Id)
		}
	}

	_, err = db.c.Exec(`UPDATE conversations
SET content = ?
WHERE id = ?;`, newcontent, id.Id)
	return err
}

func (db *appdbimpl) GetConversationMessages(id ConversationId) ([]MessageId, error) {
	var contentR string
	err := db.c.QueryRow("SELECT  content FROM conversations WHERE id = ?", id.Id).Scan(&contentR)
	messages := ConvertMessages(contentR)
	return messages, err
}

func (db *appdbimpl) CheckIfInConversation(userId UserId, receiver string) (ConversationId, error) {
	conversations, err := db.GetConversations(userId)

	for i := 0; i < len(conversations); i++ {
		var members []string
		members, err = db.GetAllMembers(conversations[i])
		if len(members) == 2 {
			if members[0] == receiver || members[1] == receiver {
				return conversations[i], err
			}
		}
	}

	return ConversationId{0}, err
}
