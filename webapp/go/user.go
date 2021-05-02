package main

import (
	"context"
	"time"

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
	rows, err := db.QueryContext(c,
		"SELECT p.id, p.name, p.image_path, p.price, h.created_at "+
			"FROM histories as h "+
			"LEFT OUTER JOIN products as p "+
			"ON h.product_id = p.id "+
			"WHERE h.user_id = ? "+
			"ORDER BY h.id DESC LIMIT 30", u.ID)
	if err != nil {
		return nil, 0
	}

	defer rows.Close()
	for rows.Next() {
		p := Product{}
		var cAt string
		fmt := "2006-01-02 15:04:05"
		err = rows.Scan(&p.ID, &p.Name, &p.ImagePath, &p.Price, &cAt)
		tmp, _ := time.Parse(fmt, cAt)
		p.CreatedAt = (tmp.Add(9 * time.Hour)).Format(fmt)
		if err != nil {
			panic(err.Error())
		}
		products = append(products, p)
	}

	err = db.QueryRowContext(c, "SELECT sum(p.price) FROM histories as h INNER JOIN products as p ON h.product_id = p.id WHERE h.user_id = ?", u.ID).
		Scan(&totalCost)
	if err != nil {
		return nil, 0
	}

	return
}

// BuyProduct : buy product
func (u *User) BuyProduct(pid string) {
	db.Exec(
		"INSERT INTO histories (product_id, user_id, created_at) VALUES (?, ?, ?)",
		pid, u.ID, time.Now())
}

// CreateComment : create comment to the product
func (u *User) CreateComment(pid string, content string) {
	db.Exec(
		"INSERT INTO comments (product_id, user_id, content, created_at) VALUES (?, ?, ?, ?)",
		pid, u.ID, content, time.Now())
}

func (u *User) UpdateLastLogin() {
	u.LastLogin = time.Now().Format("2006-01-02 15:04:05")
	usersID.Store(u.ID, *u)
}
