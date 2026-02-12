package database

func (db *appdbimpl) CountComments()(int, error){
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM comments").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, err
}

func (db *appdbimpl) CreateComment(userId UserId, messageId MessageId, content string)(CommentId, error){
	id, err := db.CountComments()
	_, err = db.c.Exec("INSERT INTO comments (id, content, user, message) VALUES (?, ?, ?, ?)", id+1, content, userId.Id, messageId.Id)
	if err ==nil{db.AddCommentToMessage(messageId, CommentId{id+1})}
	return CommentId{id+1}, err
}

func (db *appdbimpl) GetComment(id CommentId)(Comment, error){
	var content string
	var userId int
	var messageId int
	var comment Comment
	err := db.c.QueryRow("SELECT content, user, message FROM comments WHERE id = ?", id.Id).Scan(&content, &userId, &messageId)

	comment = Comment{id, content, UserId{userId}, MessageId{messageId}}
	return comment, err
}

func (db *appdbimpl) DeleteComment(id CommentId, messageId MessageId)(error){
	_, err := db.c.Exec(`UPDATE comments
SET  content = ?
WHERE id = ?;`, "", id.Id)

	if err == nil {
		err = db.RemoveCommentFromMessage(messageId, id)
	}
	return err
}


func (db *appdbimpl) CheckIfUserCommented(id CommentId, userId UserId)(bool, error){
	var userIdR int
	err := db.c.QueryRow("SELECT  user FROM comments WHERE id = ?", id.Id).Scan(&userIdR)
	if userIdR == userId.Id{
		return true, err
	}
	return false,err
}
