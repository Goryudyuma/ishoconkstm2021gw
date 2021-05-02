package types

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	LastLogin string
}

// Product Model
type Product struct {
	ID          int
	Name        string
	Description string
	ImagePath   string
	Price       int
	CreatedAt   string
}

// ProductWithComments Model
type ProductWithComments struct {
	ID           int
	Name         string
	Description  string
	ImagePath    string
	Price        int
	CreatedAt    string
	CommentCount int
	Comments     []CommentWriter
}

// CommentWriter Model
type CommentWriter struct {
	Content string
	Writer  string
}

// Comment Model
type Comment struct {
	ID        int
	ProductID int
	UserID    int
	Content   string
	CreatedAt string
}