package main

import "log"

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

func getProduct(pid int) Product {
	p := Product{}
	row := db.QueryRow("SELECT * FROM products WHERE id = ? LIMIT 1", pid)
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.ImagePath, &p.Price, &p.CreatedAt)
	if err != nil {
		panic(err.Error())
	}

	return p
}

func getProductsWithCommentsAt(page int) []ProductWithComments {
	// select 50 products with offset page*50
	products := []ProductWithComments{}
	page50 := page * 50
	rows, err := db.Query("SELECT * FROM products WHERE ? >= id AND id > ? ORDER BY id DESC", 10000-page50, 10000-page50-50)
	if err != nil {
		return nil
	}

	defer rows.Close()
	for rows.Next() {
		p := ProductWithComments{}
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.ImagePath, &p.Price, &p.CreatedAt)
		products = append(products, p)
	}

	rows, err = db.Query("SELECT product_id, count(1) as c FROM comments WHERE ? >= product_id AND product_id > ? GROUP BY product_id ORDER BY id DESC", 10000-page50, 10000-page50-50)
	if err != nil {
		panic(err)
		return nil
	}
	var m map[int]int
	for rows.Next() {
		var productId, c int
		err := rows.Scan(&productId, &c)
		if err != nil {
			panic(err)
			return nil
		}
		m[productId] = c
	}

	rows, err = db.Query("SELECT c.content, c.user_id, c.product_id FROM comments as c WHERE ? >= c.product_id AND c.product_id > ? ORDER BY c.created_at", 10000-page50, 10000-page50-50)
	var comment map[int][]CommentWriter
	for rows.Next() {
		var content string
		var userId, productId int
		err := rows.Scan(&content, &userId, &productId)
		if err != nil {
			panic(err)
			return nil
		}
		if len(comment[productId]) < 5 {
			var cw CommentWriter
			load, ok := usersID.Load(userId)
			if ok {
				cw.Writer = load.(User).Name
			}
			cw.Content = content
			comment[productId] = append(comment[productId], cw)
		}
	}

	for i := 0; i < len(products); i++ {
		products[i].CommentCount = m[products[i].ID]
		products[i].Comments = comment[products[i].ID]
	}

	return products
}

func (p *Product) isBought(uid int) bool {
	var count int
	log.Print(uid)
	log.Print(p.ID)
	err := db.QueryRow(
		"SELECT count(*) as count FROM histories WHERE product_id = ? AND user_id = ?",
		p.ID, uid,
	).Scan(&count)
	if err != nil {
		panic(err.Error())
	}

	return count > 0
}
