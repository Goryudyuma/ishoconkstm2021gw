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

		key := commentProductIDKey{
			productID: index,
		}
		value := commentProductIDValue{}
		if valueRow, ok := commentProductID.Load(key); ok {
			value = valueRow.(commentProductIDValue)
		}

		if productRow, ok := productsID.Load(index); ok {
			product := productRow.(Product)
			products[i] = ProductWithComments{
				product.ID,
				product.Name,
				product.Description,
				product.ImagePath,
				product.Price,
				product.CreatedAt,
				value.count,
				value.commentMemo,
			}
		}
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
