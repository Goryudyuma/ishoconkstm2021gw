package main

import (
	"context"
	"database/sql"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/Goryudyuma/ishoconkstm2021gw/webapp/go/types"
	"github.com/gin-gonic/contrib/sessions"
)

type User types.User

func authenticate(email string, password string) (int, bool) {
	load, ok := usersEmailPassword.Load(usersEmailPasswordKey{email, password})
	if !ok {
		return 0, false
	}
	return load.(int), ok
}

func notAuthenticated(session sessions.Session) bool {
	uid := session.Get("uid")
	return uid == nil || !(uid.(int) > 0)
}

func getUser(uid int) User {
	user, ok := usersID.Load(uid)
	if !ok {
		return User{}
	}
	return user.(User)
}

func currentUser(session sessions.Session) User {

	uid := session.Get("uid")
	if uid == nil {
		return User{}
	}

	return getUser(uid.(int))
}

// BuyingHistory : products which user had bought
func (u User) BuyingHistory(c context.Context) (products []Product, totalCost int) {
	historyUserIDMutex.RLock()
	vRow, ok := historyUserID.Load(historyUserIDKey{u.ID})
	historyUserIDMutex.RUnlock()
	if !ok {
		return nil, 0
	}
	v := vRow.(historyUserIDValue)
	totalCost = v.totalPay

	beginIndex := len(v.boughtProductList) - 30
	if beginIndex < 0 {
		beginIndex = 0
	}
	boughtProduct := v.boughtProductList[beginIndex:]

	for _, p := range boughtProduct {
		productRow, _ := productsID.Load(p.productID)
		product := productRow.(Product)
		product.CreatedAt = p.createdAt
		products = append(products, product)
	}
	for i := 0; i < len(products)/2; i++ {
		products[i], products[len(products)-i-1] = products[len(products)-i-1], products[i]
	}

	return
}

var buyProductStmt *sql.Stmt

// BuyProduct : buy product
func (u *User) BuyProduct(pid string) {
	pidint, err := strconv.Atoi(pid)
	if err != nil {
		panic(err.Error())
	}
	buyProductStmt.Exec(pid, u.ID, time.Now())

	historyUserIDMutex.Lock()
	{
		v := historyUserIDValue{}
		key := historyUserIDKey{u.ID}
		vRow, ok := historyUserID.Load(key)
		if ok {
			v = vRow.(historyUserIDValue)
			v.boughtProductMap = make(map[int]struct{})
		}
		v.boughtProductMap[pidint] = struct{}{}
		boughtProduct := boughtProductListType{
			productID: pidint,
			createdAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		v.boughtProductList = append(v.boughtProductList, boughtProduct)
		product, ok := productsID.Load(pidint)
		if ok {
			v.totalPay += product.(Product).Price
		}
		historyUserID.Store(key, v)
	}
	historyUserIDMutex.Unlock()

}

var createCommentStmt *sql.Stmt

// CreateComment : create comment to the product
func (u *User) CreateComment(pid string, content string) {
	createCommentStmt.Exec(pid, u.ID, content, time.Now())

	productID, err := strconv.Atoi(pid)
	if err != nil {
		panic(err.Error())
	}
	key := commentProductIDKey{
		productID: productID,
	}
	value := commentProductIDValue{}
	if valueRow, ok := commentProductID.Load(key); ok {
		value = valueRow.(commentProductIDValue)
	}
	value.count++

	cw := types.CommentWriter{}
	load, ok := usersID.Load(u.ID)
	if ok {
		cw.Writer = load.(User).Name
	}
	if utf8.RuneCountInString(content) > 25 {
		content = string([]rune(content)[:25]) + "…"
	}
	cw.Content = content
	value.commentMemo = append(value.commentMemo, cw)
	beginIndex := len(value.commentMemo) - 5
	if beginIndex < 0 {
		beginIndex = 0
	}
	value.commentMemo = value.commentMemo[beginIndex:]
	commentProductID.Store(key, value)
}

func (u *User) UpdateLastLogin() {
	u.LastLogin = time.Now().Format("2006-01-02 15:04:05")
	usersID.Store(u.ID, *u)
}
