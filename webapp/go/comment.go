package main

import "github.com/Goryudyuma/ishoconkstm2021gw/webapp/go/types"

// Comment Model
type Comment struct {
	ID        int
	ProductID int
	UserID    int
	Content   string
	CreatedAt string
}

func getComments(pid int) []types.Comment {
	rows, err := db.Query("SELECT * FROM comments WHERE product_id = ? ", pid)
	if err != nil {
		return nil
	}

	defer rows.Close()
	var comments []types.Comment
	for rows.Next() {
		c := Comment{}
		err = rows.Scan(&c.ID, &c.ProductID, &c.UserID, &c.Content, &c.CreatedAt)
		comments = append(comments, types.Comment(c))
	}

	return comments
}
