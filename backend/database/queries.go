package database

// Insert statements to add data to the database
const (
	AddUser    = `INSERT INTO users(username, email, firstname, lastname, gender, dob, password) values(?, ?, ?, ?, ?, ?, ?)`
	AddPost    = `INSERT INTO posts(user_id, category, title, content, date) values(?, ?, ?, ?, ?)`
	AddComment = `INSERT INTO comments(post_id, user_id, content, date) values(?, ?, ?, ?)`
	AddMessage = `INSERT INTO messages(sender_id, receiver_id, content, date) values(?, ?, ?, ?)`
	AddSession = `INSERT INTO sessions(session_uuid, user_id) values(?, ?)`
	AddChat    = `INSERT INTO chats(id_one, id_two, time) values(? ,?, ?)`
)

// Query statements to filter data from the database
const (
	GetUserById          = `SELECT * FROM users WHERE id = ?`
	GetUserByUsername    = `SELECT * FROM users WHERE username = ?`
	GetUserByEmail       = `SELECT * FROM users WHERE email = ?`
	GetAllUser           = `SELECT * FROM users ORDER BY username ASC`
	GetPostById          = `SELECT * FROM posts WHERE id = ? ORDER BY id DESC`
	GetAllPost           = `SELECT * FROM posts ORDER BY id DESC`
	GetAllPostByCategory = `SELECT * FROM posts WHERE category = ? ORDER BY id DESC`
	GetAllPostByUser     = `SELECT * FROM posts WHERE user_id = ? ORDER BY id DESC`
	GetCommentById       = `SELECT * FROM comments WHERE id = ?`
	GetAllPostComment    = `SELECT * FROM comments WHERE post_id = ?`
	GetAllUserComment    = `SELECT * FROM comments WHERE user_id = ?`
	GetMessage           = `SELECT * FROM messages WHERE id = ?`
	GetAllChatMessage    = `SELECT * FROM messages WHERE sender_id = ? AND receiver_id = ? OR sender_id = ? AND receiver_id = ?`
	GetSessionUser       = `SELECT users.* FROM sessions INNER JOIN users ON sessions.user_id = users.id WHERE sessions.session_uuid = ?`
	GetUserChats         = `SELECT * FROM chats WHERE id_one = ? OR id_two = ? ORDER BY time DESC`
	GetChatBetween       = `SELECT * FROM chats WHERE id_one = ? AND id_two = ? OR id_one = ? AND id_two = ?`
)

// Query statements to remove data from database
const (
	RemoveCookie = `DELETE FROM sessions WHERE user_id = ?`
)

// Query statements to update data in database
const (
	UpdateChat = `UPDATE chats SET time = ? WHERE id_one = ? AND id_two = ?`
)
