package backend

type CommentDetails struct {
	CommentID   int    `json:"comment_id"`
	PostID      int    `json:"post_id"`
	Content     string `json:"comment_content"`
	UserID      int    `json:"user_id"`
	Username    string `json:"username"`
	CreatedAt   string `json:"created_at"`
	Likes       int    `json:"likes"`
	Dislikes    int    `json:"dislikes"`
	LikedNow    bool   `json:"liked_now"`
	DislikedNow bool   `json:"disliked_now"`
}

type PostDetails struct {
	PostID      int              `json:"post_id"`
	UserID      int              `json:"user_id"`
	Username    string           `json:"username"`
	PostTitle   string           `json:"post_title"`
	PostContent string           `json:"post_content"`
	Comments    []CommentDetails `json:"comments"`
	Categories  []string         `json:"categories"`
	CreatedAt   string           `json:"created_at"`
	Likes       int              `json:"likes"`
	Dislikes    int              `json:"dislikes"`
	LikedNow    bool             `json:"liked_now"`
	DislikedNow bool             `json:"disliked_now"`
}

type PageDetails struct {
	LoggedIn         bool
	Username         string
	Categories       []CategoryDetails
	Posts            []PostDetails
	SelectedCategory string
	SelectedFilter   string
	ValidationError  string
}

type CategoryDetails struct {
	CategoryID   int
	CategoryName string
}
