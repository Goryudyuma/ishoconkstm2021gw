package main

import "github.com/Goryudyuma/ishoconkstm2021gw/webapp/go/types"

type Product types.Product
type ProductWithComments types.ProductWithComments
type CommentWriter types.CommentWriter

func getProduct(pid int) Product {
	product, _ := productsID.Load(pid)
	return product.(Product)
}

func getProductsWithCommentsAt(page int) []ProductWithComments {
	// select 50 products with offset page*50
	products := make([]ProductWithComments, 50)
	page50 := page * 50

	for i := 0; i < 50; i++ {
		index := 10000 - page50 - i
		if productRow, ok := productsID.Load(index); ok {
			product := productRow.(Product)
			products[i] = ProductWithComments{
				product.ID,
				product.Name,
				product.Description,
				product.ImagePath,
				product.Price,
				product.CreatedAt,
				0,
				[]types.CommentWriter{},
			}
		}
	}

	rows, err := db.Query("SELECT product_id, count(1) as c FROM comments WHERE ? >= product_id AND product_id > ? GROUP BY product_id", 10000-page50, 10000-page50-50)
	if err != nil {
		panic(err)
		return nil
	}
	m := make(map[int]int)
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
	comment := make(map[int][]types.CommentWriter)
	for rows.Next() {
		var content string
		var userId, productId int
		err := rows.Scan(&content, &userId, &productId)
		if err != nil {
			panic(err)
			return nil
		}
		if len(comment[productId]) < 5 {
			var cw types.CommentWriter
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
		products[i].Comments = (comment[products[i].ID])
	}

	return products
}

func (p *Product) isBought(uid int) bool {
	historyUserIDMutex.RLock()
	vRow, ok := historyUserID.Load(historyUserIDKey{uid})
	historyUserIDMutex.RUnlock()
	if !ok {
		return false
	}
	v := vRow.(historyUserIDValue)
	_, ok = v.boughtProductMap[p.ID]
	return ok
}
