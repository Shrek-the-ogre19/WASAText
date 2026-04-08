package database

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite" // Add this import
)

func TestCountUsers_EmptyDatabase(t *testing.T) {
	// 1. Open in-memory SQLite database
	db, err := sql.Open("sqlite", ":memory:") // ":memory:" creates temporary DB
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close() // Automatically cleaned up when test ends

	// 2. Create your database instance
	appDB, err := New(db)
	if err != nil {
		t.Fatalf("Failed to create AppDatabase: %v", err)
	}

	// 3. Test CountUsers on empty database
	count, err := appDB.CountUsers()
	if err != nil {
		t.Errorf("CountUsers failed: %v", err)
	}

	// 4. Verify result
	if count != 0 {
		t.Errorf("Expected 0 users in empty database, got %d", count)
	}
}

func TestAddUser_Simple(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	appDB, err := New(db)
	if err != nil {
		t.Fatal(err)
	}

	// Add a user
	err = appDB.AddUser("testuser")
	if err != nil {
		t.Errorf("AddUser failed: %v", err)
	}

	// Verify count is now 1
	count, err := appDB.CountUsers()
	if err != nil {
		t.Errorf("CountUsers failed: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 user after adding, got %d", count)
	}
}

func TestMyTest2(t *testing.T) {
	//Other Tests
	a := ConvertConversations("123,1234,12222")
	t.Log(a)

	b := ConvertUsers("axlamal,sal")
	t.Log(b)

	c := ConvertMessages("")
	t.Log(c)
}

func TestMyTest(t *testing.T) {
	//Other Tests
	a := ConvertConversations("123,1234,12222")
	t.Log(a)

	// 1. Open in-memory SQLite database
	db, err := sql.Open("sqlite", ":memory:") // ":memory:" creates temporary DB
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close() // Automatically cleaned up when test ends

	// 2. Create your database instance
	appDB, err := New(db)
	if err != nil {
		t.Fatalf("Failed to create AppDatabase: %v", err)
	}
	t.Log("database created succesfully")

	// 3. Test CountUsers on empty database
	count, err := appDB.CountUsers()
	if err != nil {
		t.Errorf("CountUsers failed: %v", err)
	}
	t.Log("count is", count)
	// 4. Verify result
	if count != 0 {
		t.Errorf("Expected 0 users in empty database, got %d", count)
	}

	// 5. Add a User called "pirveli"
	err = appDB.AddUser("pirveli")
	if err != nil {
		t.Errorf("AddUser failed: %v", err)
	}
	t.Log("user pirveli added succesfully")

	// 6. Count Users
	count, err = appDB.CountUsers()
	if err != nil {
		t.Errorf("CountUsers failed: %v", err)
	}
	t.Log("count is", count)

	// 7. get the first user id
	var user1id UserId
	user1id, err = appDB.GetUserId("pirveli")
	if err != nil {
		t.Errorf("GetUserId failed: %v", err)
	}
	t.Log(user1id)

	// 8. get the user properties
	var user1 User
	user1, err = appDB.GetUser(user1id)
	if err != nil {
		t.Errorf("GetUser failed: %v", err)
	}
	t.Log("first enty:", user1.Id, user1.Name, user1.Picture)

	//9. add another user
	err = appDB.AddUser("meort")
	if err != nil {
		t.Errorf("AddUser failed: %v", err)
	}
	t.Log("user meort added succesfully")

	// 10. Change the user name and picture
	var user2id UserId
	user2id, err = appDB.GetUserId("meort")
	if err != nil {
		t.Errorf("GetUserId failed: %v", err)
	}

	err = appDB.ChangeUserName(user2id, "meore")
	if err != nil {
		t.Errorf("ChangeUserName failed: %v", err)
	}

	err = appDB.ChangeUserPicture(user2id, "new picture url")
	if err != nil {
		t.Errorf("ChangeUserPicture failed: %v", err)
	}

	var user2 User
	user2, err = appDB.GetUser(user2id)
	if err != nil {
		t.Errorf("GetUser failed: %v", err)
	}
	t.Log("second enty:", user2.Id, user2.Name, user2.Picture)

	// 11. List all users
	count, err = appDB.CountUsers()
	if err != nil {
		t.Errorf("CountUsers failed: %v", err)
	}
	t.Log("count is", count)
	var users []string
	users, err = appDB.ListAllUsers()
	if err != nil {
		t.Errorf("ListAllUsers failed: %v", err)
	}
	t.Log(users)

	var confirmation bool
	confirmation, err = appDB.UserLookup("pirvela")
	if err != nil {
		t.Errorf("UserLookup failed: %v", err)
	}
	t.Log(confirmation)

	// 12. Conversations
	var user2conversations []ConversationId
	user2conversations, err = appDB.GetConversations(user2id)
	if err != nil {
		t.Errorf("GetConversations failed: %v", err)
	}
	t.Log("user2conversations:", user2conversations)

	// 13. Add Conversations
	var countofconversations int
	countofconversations, err = appDB.CountConversations()
	if err != nil {
		t.Errorf("CountConversations failed: %v", err)
	}
	t.Log("Count of Conversations:", countofconversations)

	var receivers []string
	receivers = append(receivers, "pirveli")
	err = appDB.StartConversation(user2id, receivers)
	if err != nil {
		t.Errorf("StartConversation failed: %v", err)
	}

	countofconversations, err = appDB.CountConversations()
	if err != nil {
		t.Errorf("CountConversations failed: %v", err)
	}
	t.Log("Count of Conversations:", countofconversations)

	user2conversations, err = appDB.GetConversations(user2id)
	if err != nil {
		t.Errorf("GetConversations failed: %v", err)
	}
	t.Log("Conversations of meore", user2conversations)

	var conversations1 []ConversationId
	conversations1, err = appDB.GetConversations(user1id)
	if err != nil {
		t.Errorf("GetConversations failed: %v", err)
	}
	t.Log("Conversations of pirveli", conversations1)

	firstConversationId := user2conversations[0]
	t.Log(firstConversationId)

	var firstConversation Conversation
	firstConversation, err = appDB.GetConversation(firstConversationId)
	if err != nil {
		t.Errorf("GetConversation failed: %v", err)
	}
	t.Log(firstConversation)

}
func TestWorld(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:") // ":memory:" creates temporary DB
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close() // Automatically cleaned up when test ends

	appDB, err := New(db)
	if err != nil {
		t.Fatalf("Failed to create AppDatabase: %v", err)
	}
	t.Log("database created succesfully")

	//create users
	var user1 User
	var user2 User
	var user3 User
	var user4 User
	var user1Id UserId
	var user2Id UserId
	var user3Id UserId
	var user4Id UserId
	var userLookupPirveli bool
	var userLookupMeort bool
	var userLookupMesame bool
	var userLookupMeotxe bool

	//Login for "pirveli"
	userLookupPirveli, err = appDB.UserLookup("pirveli")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	if userLookupPirveli {
		user1Id, err = appDB.GetUserId("pirveli")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	} else {
		err = appDB.AddUser("pirveli")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	}

	userLookupMeort, err = appDB.UserLookup("meort")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	if userLookupMeort {
		user2Id, err = appDB.GetUserId("meort")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	} else {
		err = appDB.AddUser("meort")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	}

	userLookupMesame, err = appDB.UserLookup("mesame")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	if userLookupMesame {
		user3Id, err = appDB.GetUserId("mesame")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	} else {
		err = appDB.AddUser("mesame")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	}

	userLookupMeotxe, err = appDB.UserLookup("meotxe")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	if userLookupMeotxe {
		user4Id, err = appDB.GetUserId("meotxe")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	} else {
		err = appDB.AddUser("meotxe")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	}

	var allUserNames []string
	allUserNames, err = appDB.ListAllUsers()
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	t.Log(allUserNames)
	t.Log(user1Id, ":", user1)
	t.Log(user2Id, ":", user2)
	t.Log(user3Id, ":", user3)
	t.Log(user4Id, ":", user4)

	//login from meort and change username and picture
	var userId UserId
	var user User
	userLookupMeort, err = appDB.UserLookup("meort")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	if userLookupMeort {
		userId, err = appDB.GetUserId("meort")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	} else {
		err = appDB.AddUser("meort")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	}

	user, err = appDB.GetUser(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log()
	t.Log("now logged in", userId, ":", user)

	err = appDB.ChangeUserName(userId, "meore")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	err = appDB.ChangeUserPicture(userId, "new user2 picture")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	user, err = appDB.GetUser(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log(userId, ":", user)

	userId = UserId{}
	user = User{}

	//Login with mesame and start a conversation with meore
	userLookupMesame, err = appDB.UserLookup("mesame")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	if userLookupMesame {
		userId, err = appDB.GetUserId("mesame")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	} else {
		err = appDB.AddUser("mesame")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	}

	user, err = appDB.GetUser(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log()
	t.Log("now logged in", userId, ":", user)

	err = appDB.StartConversation(userId, []string{"pirveli"})
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	var conversations []ConversationId
	conversations, err = appDB.GetConversations(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("mesames conversations", conversations)

	//login with pirveli and check if the conversation is there
	userLookupPirveli, err = appDB.UserLookup("pirveli")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	if userLookupPirveli {
		userId, err = appDB.GetUserId("pirveli")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	} else {
		err = appDB.AddUser("pirveli")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	}

	user, err = appDB.GetUser(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log()
	t.Log("now logged in", userId, ":", user)

	conversations, err = appDB.GetConversations(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("pirvelis conversations", conversations)

	//login back in mesame and send a message to pirveli

	userLookupMesame, err = appDB.UserLookup("mesame")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	if userLookupMesame {
		userId, err = appDB.GetUserId("mesame")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	} else {
		err = appDB.AddUser("mesame")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	}

	user, err = appDB.GetUser(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log()
	t.Log("now logged in", userId, ":", user)

	conversations, err = appDB.GetConversations(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("mesames conversations", conversations)

	var conversationId ConversationId
	var conversation Conversation
	conversationId, err = appDB.CheckIfInConversation(userId, "pirveli")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	conversation, err = appDB.GetConversation(conversationId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("mesame-pirveli conversation", conversation)

	err = appDB.CreateMessage(userId, conversationId, "Hello")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	err = appDB.CreateMessage(userId, conversationId, "How are you")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	err = appDB.CreateMessage(userId, conversationId, "Can you create a groupchat with others?")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	conversation, err = appDB.GetConversation(conversationId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("mesame-pirveli conversation", conversation)

	//Log into pirveli, respond to mesame and create a groupchat with all others

	userLookupPirveli, err = appDB.UserLookup("pirveli")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	if userLookupPirveli {
		userId, err = appDB.GetUserId("pirveli")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	} else {
		err = appDB.AddUser("pirveli")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	}

	user, err = appDB.GetUser(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log()
	t.Log("now logged in", userId, ":", user)

	conversations, err = appDB.GetConversations(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("pirvelis conversations", conversations)
	conversationId, err = appDB.CheckIfInConversation(userId, "mesame")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	var messages []MessageId
	messages, err = appDB.GetConversationMessages(conversationId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("messages in pirveli-mesame conversation", messages)

	var message1 Message
	var message2 Message
	var message3 Message

	message1, err = appDB.GetMessage(messages[0])
	t.Log(messages[0])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", err)
	}
	t.Log("first message", message1)

	message2, err = appDB.GetMessage(messages[1])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("second message", message2)
	message3, err = appDB.GetMessage(messages[2])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("third message", message3)

	err = appDB.CreateMessage(userId, conversationId, "Sure")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	err = appDB.StartConversation(userId, []string{"meore", "mesame"})
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	conversations, err = appDB.GetConversations(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	conversationId = conversations[1]
	err = appDB.CreateMessage(userId, conversationId, "Welcome to our new group chat")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	err = appDB.ForwardMessage(userId, messages[2], conversationId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	conversation, err = appDB.GetConversation(conversationId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("groupchat", conversation)

	messages, err = appDB.GetConversationMessages(conversationId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("messages in groupchat", messages)

	//log into meore and mesame and react to the messages
	userLookupMeore := userLookupMeort
	userLookupMeore, err = appDB.UserLookup("meore")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	if userLookupMeore {
		userId, err = appDB.GetUserId("meore")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	} else {
		err = appDB.AddUser("meore")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	}

	user, err = appDB.GetUser(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log()
	t.Log("now logged in", userId, ":", user)

	conversations, err = appDB.GetConversations(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("meore conversations", conversations)

	conversation, err = appDB.GetConversation(conversations[0])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("groupchat", conversation)

	messages, err = appDB.GetConversationMessages(conversations[0])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("messages in groupchat", messages)
	var message5 Message
	var message6 Message
	message5, err = appDB.GetMessage(messages[0])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	message6, err = appDB.GetMessage(messages[1])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log(message5)
	t.Log(message6)

	var comments []CommentId
	comments, err = appDB.GetComments(messages[0])
	if err != nil {
		t.Log(err)
	}
	t.Log("comments", comments)

	err = appDB.CreateComment(userId, messages[0], "<3")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	comments, err = appDB.GetComments(messages[0])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("comments", comments)

	var comment Comment

	comment, err = appDB.GetComment(comments[0])
	if err != nil {
		t.Log(err)
	}
	t.Log("new comment", comment)

	userLookupMesame, err = appDB.UserLookup("mesame")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	if userLookupMesame {
		userId, err = appDB.GetUserId("mesame")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	} else {
		err = appDB.AddUser("mesame")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	}

	user, err = appDB.GetUser(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log()
	t.Log("now logged in", userId, ":", user)

	conversations, err = appDB.GetConversations(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("mesame conversations", conversations)

	conversation, err = appDB.GetConversation(conversations[1])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("groupchat", conversation)

	messages, err = appDB.GetConversationMessages(conversations[1])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("messages in groupchat", messages)

	message5, err = appDB.GetMessage(messages[0])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	message6, err = appDB.GetMessage(messages[1])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log(message5)
	t.Log(message6)

	comments, err = appDB.GetComments(messages[0])
	if err != nil {
		t.Log(err)
	}
	t.Log("comments", comments)

	err = appDB.CreateComment(userId, messages[0], ":)")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	comments, err = appDB.GetComments(messages[0])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("comments", comments)

	comment, err = appDB.GetComment(comments[1])
	if err != nil {
		t.Log(err)
	}
	t.Log("new comment", comment)

	//Log into pirveli, delete last message from the groupchat, add meotxe, change name and picture
	userLookupPirveli, err = appDB.UserLookup("pirveli")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	if userLookupPirveli {
		userId, err = appDB.GetUserId("pirveli")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	} else {
		err = appDB.AddUser("pirveli")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	}

	user, err = appDB.GetUser(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log()
	t.Log("now logged in", userId, ":", user)

	conversations, err = appDB.GetConversations(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("pirvelis conversations", conversations)

	conversationId = conversations[1]
	messages, err = appDB.GetConversationMessages(conversationId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("messages in groupchat", messages)
	message5, err = appDB.GetMessage(messages[0])
	t.Log(message5)

	var message5comments []CommentId
	message5comments, err = appDB.GetComments(messages[0])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("comments for above message", message5comments)
	var comment1 Comment
	var comment2 Comment
	comment1, err = appDB.GetComment(message5comments[0])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	comment2, err = appDB.GetComment(message5comments[1])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("these are the comments", comment1, comment2)

	err = appDB.DeleteMessage(messages[1], conversationId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	messages, err = appDB.GetConversationMessages(conversationId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("messages in groupchat", messages)

	err = appDB.ChangeConversationName(conversationId, "family")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	err = appDB.ChangeConversationPhoto(conversationId, "familypicture")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	err = appDB.AddMembers(conversationId, "meotxe")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

	//Log into mesame and leave

	userLookupMesame, err = appDB.UserLookup("mesame")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	if userLookupMesame {
		userId, err = appDB.GetUserId("mesame")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	} else {
		err = appDB.AddUser("mesame")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	}

	user, err = appDB.GetUser(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log()
	t.Log("now logged in", userId, ":", user)

	conversations, err = appDB.GetConversations(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("mesame conversations", conversations)

	conversation, err = appDB.GetConversation(conversations[1])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("groupchat", conversation)

	conversationId = conversations[1]
	err = appDB.LeaveGroup(conversationId, userId)
	if err != nil {
		t.Log(err)
	}

	//log into meore and delete comment
	userLookupMeore, err = appDB.UserLookup("meore")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	if userLookupMeore {
		userId, err = appDB.GetUserId("meore")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	} else {
		err = appDB.AddUser("meore")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	}

	user, err = appDB.GetUser(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log()
	t.Log("now logged in", userId, ":", user)

	conversations, err = appDB.GetConversations(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("meore conversations", conversations)

	conversation, err = appDB.GetConversation(conversations[0])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("groupchat", conversation)

	messages, err = appDB.GetConversationMessages(conversations[0])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("messages in groupchat", messages)

	comments, err = appDB.GetComments(messages[0])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("comments", comments)

	comment1, err = appDB.GetComment(comments[0])
	if err != nil {
		t.Log(err)
	}
	t.Log(comment1)

	err = appDB.DeleteComment(comments[0], messages[0])
	if err != nil {
		t.Log(err)
	}

	comments, err = appDB.GetComments(messages[0])
	if err != nil {
		t.Log(err)
	}
	t.Log("comments", comments)

	//log into meotxe and check everything
	userLookupMeotxe, err = appDB.UserLookup("meotxe")
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	if userLookupMeotxe {
		userId, err = appDB.GetUserId("meotxe")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	} else {
		err = appDB.AddUser("meotxe")
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
	}

	user, err = appDB.GetUser(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log()
	t.Log("now logged in", userId, ":", user)

	conversations, err = appDB.GetConversations(userId)
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("meotxe conversations:", conversations)
	for i := 0; i < len(conversations); i++ {
		conversation, err = appDB.GetConversation(conversations[i])
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
		t.Log(conversations[i], conversation)
	}

	messages, err = appDB.GetConversationMessages(conversations[0])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log("messages in groupchat:", messages)

	var message Message
	for i := 0; i < len(messages); i++ {
		message, err = appDB.GetMessage(messages[i])
		if err != nil {
			t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		}
		t.Log(messages[i], message)
	}

	comments, err = appDB.GetComments(messages[0])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log(comments)
	comment, err = appDB.GetComment(comments[0])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log(comment)

	err = appDB.UpdateStatus(messages[0])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	message, err = appDB.GetMessage(messages[0])
	if err != nil {
		t.Log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	t.Log(message)

}
