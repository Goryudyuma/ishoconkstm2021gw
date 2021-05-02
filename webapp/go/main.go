//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=templates
package main

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"sync"
	"unicode/utf8"

	"github.com/Goryudyuma/ishoconkstm2021gw/webapp/go/templates"
	"github.com/Goryudyuma/ishoconkstm2021gw/webapp/go/types"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

var usersEmailPassword sync.Map

type usersEmailPasswordKey struct {
	email    string
	password string
}

var usersID sync.Map

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	// database setting
	user := getEnv("ISHOCON1_DB_USER", "ishocon")
	pass := getEnv("ISHOCON1_DB_PASSWORD", "ishocon")
	dbname := getEnv("ISHOCON1_DB_NAME", "ishocon1")
	db, _ = sql.Open("mysql", user+":"+pass+"@/"+dbname)
	db.SetMaxIdleConns(5)

	r := gin.Default()
	pprof.Register(r)
	// load templates
	r.Use(static.Serve("/css", static.LocalFile("public/css", true)))
	r.Use(static.Serve("/images", static.LocalFile("public/images", true)))

	// session store
	store := sessions.NewCookieStore([]byte("mysession"))
	store.Options(sessions.Options{HttpOnly: true})
	r.Use(sessions.Sessions("showwin_happy", store))

	// GET /login
	r.GET("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(templates.Login("ECサイトで爆買いしよう！！！！")))
	})

	// POST /login
	r.POST("/login", func(c *gin.Context) {
		email := c.PostForm("email")
		pass := c.PostForm("password")

		session := sessions.Default(c)
		userID, result := authenticate(email, pass)
		if result {
			// 認証成功
			session.Set("uid", userID)
			session.Save()

			// user.UpdateLastLogin()

			c.Redirect(http.StatusSeeOther, "/")
		} else {
			// 認証失敗

			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(templates.Login("ログインに失敗しました")))
		}
	})

	// GET /logout
	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()

		c.Redirect(http.StatusFound, "/login")
	})

	// GET /
	r.GET("/", func(c *gin.Context) {
		cUser := currentUser(sessions.Default(c))

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 0
		}
		products := getProductsWithCommentsAt(page)
		// shorten description and comment
		var sProducts []types.ProductWithComments
		for _, p := range products {
			if utf8.RuneCountInString(p.Description) > 70 {
				p.Description = string([]rune(p.Description)[:70]) + "…"
			}

			var newCW []types.CommentWriter
			for _, c := range p.Comments {
				if utf8.RuneCountInString(c.Content) > 25 {
					c.Content = string([]rune(c.Content)[:25]) + "…"
				}
				newCW = append(newCW, c)
			}
			p.Comments = newCW
			sProducts = append(sProducts, types.ProductWithComments(p))
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(templates.Index(types.User(cUser), sProducts)))
	})

	// GET /users/:userId
	r.GET("/users/:userId", func(c *gin.Context) {
		cUser := currentUser(sessions.Default(c))

		uid, _ := strconv.Atoi(c.Param("userId"))
		user := getUser(uid)

		products, totalPay := user.BuyingHistory(c)

		// shorten description
		var sdProducts []types.Product
		for _, p := range products {
			if utf8.RuneCountInString(p.Description) > 70 {
				p.Description = string([]rune(p.Description)[:70]) + "…"
			}
			sdProducts = append(sdProducts, types.Product(p))

			if len(sdProducts) > 30 {
				break
			}
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8",
			[]byte(templates.MyPage(types.User(cUser), types.User(user), sdProducts, totalPay)))
	})

	// GET /products/:productId
	r.GET("/products/:productId", func(c *gin.Context) {
		pid, _ := strconv.Atoi(c.Param("productId"))
		product := getProduct(pid)
		comments := getComments(pid)

		cUser := currentUser(sessions.Default(c))
		bought := product.isBought(cUser.ID)

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(templates.ProductPage(types.User(cUser), types.Product(product), comments, bought)))
	})

	// POST /products/buy/:productId
	r.POST("/products/buy/:productId", func(c *gin.Context) {
		// need authenticated
		if notAuthenticated(sessions.Default(c)) {
			c.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte(templates.Login("先にログインをしてください")))
		} else {
			// buy product
			cUser := currentUser(sessions.Default(c))
			cUser.BuyProduct(c.Param("productId"))

			// redirect to user page
			c.Redirect(http.StatusFound, "/users/"+strconv.Itoa(cUser.ID))
		}
	})

	// POST /comments/:productId
	r.POST("/comments/:productId", func(c *gin.Context) {
		// need authenticated
		if notAuthenticated(sessions.Default(c)) {
			c.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte(templates.Login("先にログインをしてください")))
		} else {
			// create comment
			cUser := currentUser(sessions.Default(c))
			cUser.CreateComment(c.Param("productId"), c.PostForm("content"))

			// redirect to user page
			c.Redirect(http.StatusFound, "/users/"+strconv.Itoa(cUser.ID))
		}
	})

	// GET /initialize
	r.GET("/initialize", func(c *gin.Context) {
		db.Exec("DELETE FROM users WHERE id > 5000")
		db.Exec("DELETE FROM products WHERE id > 10000")
		db.Exec("DELETE FROM comments WHERE id > 200000")
		db.Exec("DELETE FROM histories WHERE id > 500000")

		usersEmailPassword = sync.Map{}
		usersID = sync.Map{}

		rows, err := db.Query("SELECT id, name, email, password, last_login FROM users")
		if err != nil {
			c.String(http.StatusServiceUnavailable, err.Error())
			return
		}
		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.LastLogin)
			if err != nil {
				c.String(http.StatusServiceUnavailable, err.Error())
				return
			}

			usersEmailPassword.Store(usersEmailPasswordKey{user.Email, user.Password}, user.ID)
			usersID.Store(user.ID, user)
		}

		c.String(http.StatusOK, "Finish")
	})

	r.RunUnix("/var/run/go/go.sock")
}
