//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=templates
package main

import (
	"bytes"
	"database/sql"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/Goryudyuma/ishoconkstm2021gw/webapp/go/templates"
	"github.com/Goryudyuma/ishoconkstm2021gw/webapp/go/types"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
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

var productsID sync.Map

var productDescriptionMemo sync.Map

var historyUserID sync.Map
var historyUserIDMutex sync.RWMutex

type historyUserIDKey struct {
	userID int
}

type historyUserIDValue struct {
	totalPay          int
	boughtProductMap  map[int]struct{}
	boughtProductList []boughtProductListType
}

type boughtProductListType struct {
	productID int
	createdAt string
}

var commentProductID sync.Map

type commentProductIDKey struct {
	productID int
}

type commentProductIDValue struct {
	count       int
	commentMemo []types.CommentWriter
}

var usersHTMLCache sync.Map

type usersHTMLCacheKey struct {
	userID      int
	currentUser bool
}

type usersHTMLCacheValue struct {
	html []byte
}

var headerHTMLCache sync.Map

type headerHTMLCacheKey struct {
	userID int
}

type headerHTMLCacheValue struct {
	html []byte
}

var indexHTMLCache sync.Map

type indexHTMLCacheKey struct {
	page        int
	isLoginUser bool
}

type indexHTMLCacheValue struct {
	html []byte
}

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
	db, _ = sql.Open("mysql", user+":"+pass+"@unix(/var/run/mysqld/mysqld.sock)/"+dbname)
	db.SetMaxIdleConns(5)

	var err error
	createCommentStmt, err = db.Prepare("INSERT INTO comments (product_id, user_id, content, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	buyProductStmt, err = db.Prepare("INSERT INTO histories (product_id, user_id, created_at) VALUES (?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	r := gin.Default()
	pprof.Register(r)
	// load templates
	r.Use(static.Serve("/css", static.LocalFile("public/css", true)))
	r.Use(static.Serve("/images", static.LocalFile("public/images", true)))

	// session store
	store := memstore.NewStore([]byte("mysession"))
	store.Options(sessions.Options{HttpOnly: true})
	r.Use(sessions.Sessions("showwin_happy", store))

	// GET /login
	r.GET("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(templates.Login("EC??????????????????????????????????????????")))
	})

	// POST /login
	r.POST("/login", func(c *gin.Context) {
		email := c.PostForm("email")
		pass := c.PostForm("password")

		session := sessions.Default(c)
		userID, result := authenticate(email, pass)
		if result {
			// ????????????
			session.Set("uid", userID)
			session.Save()

			// user.UpdateLastLogin()

			c.Redirect(http.StatusSeeOther, "/")
		} else {
			// ????????????

			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(templates.Login("?????????????????????????????????")))
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
		var ret bytes.Buffer
		cUser := currentUser(sessions.Default(c))

		if valueRow, ok := headerHTMLCache.Load(headerHTMLCacheKey{userID: cUser.ID}); ok {
			ret.Write(valueRow.(headerHTMLCacheValue).html)
		} else {
			templates.WriteHeader(&ret, types.User(cUser))
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 0
		}

		isLoginUser := cUser.ID > 0
		key := indexHTMLCacheKey{
			page:        page,
			isLoginUser: isLoginUser,
		}
		if valueRow, ok := indexHTMLCache.Load(key); ok {
			ret.Write(valueRow.(indexHTMLCacheValue).html)
		} else {
			html := indexPageHTML(page, isLoginUser)
			ret.Write(html)

			indexHTMLCache.Store(key, indexHTMLCacheValue{html: html})
		}
		templates.WriteFooter(&ret)

		c.Data(http.StatusOK, "text/html; charset=utf-8", ret.Bytes())
	})

	// GET /users/:userId
	r.GET("/users/:userId", func(c *gin.Context) {
		var ret bytes.Buffer
		cUser := currentUser(sessions.Default(c))
		cUserUser := types.User(cUser)

		if valueRow, ok := headerHTMLCache.Load(headerHTMLCacheKey{userID: cUserUser.ID}); ok {
			ret.Write(valueRow.(headerHTMLCacheValue).html)
		} else {
			templates.WriteHeader(&ret, cUserUser)
		}

		uid, _ := strconv.Atoi(c.Param("userId"))
		user := getUser(uid)

		key := usersHTMLCacheKey{
			userID:      uid,
			currentUser: uid == cUserUser.ID,
		}
		if valueRow, ok := usersHTMLCache.Load(key); ok {
			value := valueRow.(usersHTMLCacheValue)
			ret.Write(value.html)
		} else {
			products, totalPay := user.BuyingHistory(c)

			// shorten description
			var sdProducts []types.Product
			for _, p := range products {
				var productDescription string
				productDescriptionRow, ok := productDescriptionMemo.Load(p.ID)
				if !ok {
					if utf8.RuneCountInString(p.Description) > 70 {
						productDescription = string([]rune(p.Description)[:70]) + "???"
					} else {
						productDescription = p.Description
					}
					productDescriptionMemo.Store(p.ID, productDescription)
				} else {
					productDescription = productDescriptionRow.(string)
				}
				p.Description = productDescription
				sdProducts = append(sdProducts, types.Product(p))

				if len(sdProducts) > 30 {
					break
				}
			}

			var html bytes.Buffer
			templates.WriteMyPage(&html, cUserUser.ID == user.ID, types.User(user), sdProducts, totalPay)
			value := usersHTMLCacheValue{
				html: html.Bytes(),
			}
			usersHTMLCache.Store(key, value)

			if _, err := html.WriteTo(&ret); err != nil {
				panic(err)
			}
		}
		templates.WriteFooter(&ret)
		c.Data(http.StatusOK, "text/html; charset=utf-8", ret.Bytes())
	})

	// GET /products/:productId
	r.GET("/products/:productId", func(c *gin.Context) {
		pid, _ := strconv.Atoi(c.Param("productId"))
		product := getProduct(pid)

		cUser := currentUser(sessions.Default(c))
		bought := product.isBought(cUser.ID)

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(templates.ProductPage(types.User(cUser), types.Product(product), bought)))
	})

	// POST /products/buy/:productId
	r.POST("/products/buy/:productId", func(c *gin.Context) {
		// need authenticated
		if notAuthenticated(sessions.Default(c)) {
			c.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte(templates.Login("???????????????????????????????????????")))
		} else {
			// buy product
			cUser := currentUser(sessions.Default(c))
			cUser.BuyProduct(c.Param("productId"))

			{
				key := usersHTMLCacheKey{
					userID:      cUser.ID,
					currentUser: false,
				}
				usersHTMLCache.Delete(key)
				key.currentUser = true
				usersHTMLCache.Delete(key)
			}
			// redirect to user page
			c.Redirect(http.StatusFound, "/users/"+strconv.Itoa(cUser.ID))
		}
	})

	// POST /comments/:productId
	r.POST("/comments/:productId", func(c *gin.Context) {
		// need authenticated
		if notAuthenticated(sessions.Default(c)) {
			c.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte(templates.Login("???????????????????????????????????????")))
		} else {
			// create comment
			cUser := currentUser(sessions.Default(c))
			cUser.CreateComment(c.Param("productId"), c.PostForm("content"))

			productId, err := strconv.Atoi(c.Param("productId"))
			if err != nil {
				panic(err)
			}
			key := indexHTMLCacheKey{
				page:        (10000 - productId) / 50,
				isLoginUser: false,
			}
			indexHTMLCache.Delete(key)
			key.isLoginUser = true
			indexHTMLCache.Delete(key)

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
		productsID = sync.Map{}
		productDescriptionMemo = sync.Map{}
		historyUserID = sync.Map{}
		commentProductID = sync.Map{}
		usersHTMLCache = sync.Map{}
		headerHTMLCache = sync.Map{}
		indexHTMLCache = sync.Map{}

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

			key := headerHTMLCacheKey{
				userID: user.ID,
			}
			var html bytes.Buffer
			templates.WriteHeader(&html, types.User(user))
			headerHTMLCache.Store(key, headerHTMLCacheValue{
				html: html.Bytes(),
			})
		}
		err = rows.Close()
		if err != nil {
			c.String(http.StatusServiceUnavailable, err.Error())
			return
		}
		var user User
		key := headerHTMLCacheKey{
			userID: 0,
		}
		var html bytes.Buffer
		templates.WriteHeader(&html, types.User(user))
		headerHTMLCache.Store(key, headerHTMLCacheValue{
			html: html.Bytes(),
		})

		rows, err = db.Query("SELECT * FROM products")
		if err != nil {
			c.String(http.StatusServiceUnavailable, err.Error())
			return
		}

		for rows.Next() {
			p := Product{}
			err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.ImagePath, &p.Price, &p.CreatedAt)
			productsID.Store(p.ID, p)

			var productDescription string
			if utf8.RuneCountInString(p.Description) > 70 {
				productDescription = string([]rune(p.Description)[:70]) + "???"
			} else {
				productDescription = p.Description
			}
			productDescriptionMemo.Store(p.ID, productDescription)
		}
		err = rows.Close()
		if err != nil {
			c.String(http.StatusServiceUnavailable, err.Error())
			return
		}

		rows, err = db.Query("SELECT id, product_id, user_id, created_at FROM histories ORDER BY id")
		if err != nil {
			c.String(http.StatusServiceUnavailable, err.Error())
			return
		}

		for rows.Next() {
			var id, productID, userID int
			var createdAt string
			err = rows.Scan(&id, &productID, &userID, &createdAt)

			key := historyUserIDKey{
				userID: userID,
			}
			value := historyUserIDValue{}
			historyUserIDMutex.Lock()
			if v, ok := historyUserID.Load(key); ok {
				value = v.(historyUserIDValue)
			} else {
				value.boughtProductMap = make(map[int]struct{})
			}
			value.boughtProductMap[productID] = struct{}{}

			fmt := "2006-01-02 15:04:05"
			tmp, _ := time.Parse(fmt, createdAt)
			value.boughtProductList = append(value.boughtProductList,
				boughtProductListType{
					productID: productID,
					createdAt: (tmp.Add(9 * time.Hour)).Format(fmt),
				})

			product, _ := productsID.Load(productID)
			value.totalPay += product.(Product).Price
			historyUserID.Store(key, value)
			historyUserIDMutex.Unlock()
		}
		err = rows.Close()
		if err != nil {
			c.String(http.StatusServiceUnavailable, err.Error())
			return
		}

		rows, err = db.Query("SELECT c.content, c.user_id, c.product_id FROM comments as c ORDER BY c.created_at, c.id")
		if err != nil {
			c.String(http.StatusServiceUnavailable, err.Error())
			return
		}

		for rows.Next() {
			var content string
			var userID, productID int
			err = rows.Scan(&content, &userID, &productID)

			key := commentProductIDKey{
				productID: productID,
			}
			value := commentProductIDValue{}
			if valueRow, ok := commentProductID.Load(key); ok {
				value = valueRow.(commentProductIDValue)
			}
			value.count++

			cw := types.CommentWriter{}
			load, ok := usersID.Load(userID)
			if ok {
				cw.Writer = load.(User).Name
			}
			if utf8.RuneCountInString(content) > 25 {
				content = string([]rune(content)[:25]) + "???"
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
		err = rows.Close()
		if err != nil {
			c.String(http.StatusServiceUnavailable, err.Error())
			return
		}

		for page := 0; page < 200; page++ {
			isLoginUser := true
			indexHTMLCache.Store(indexHTMLCacheKey{
				page:        page,
				isLoginUser: isLoginUser,
			}, indexHTMLCacheValue{
				html: indexPageHTML(page, isLoginUser),
			})
			isLoginUser = false
			indexHTMLCache.Store(indexHTMLCacheKey{
				page:        page,
				isLoginUser: isLoginUser,
			}, indexHTMLCacheValue{
				html: indexPageHTML(page, isLoginUser),
			})
		}

		runtime.GC()

		c.String(http.StatusOK, "Finish")
	})

	r.RunUnix("/var/run/go/go.sock")
}

func indexPageHTML(page int, isLoginUser bool) []byte {
	products := getProductsWithCommentsAt(page)
	// shorten description and comment
	var sProducts []types.ProductWithComments
	for _, p := range products {
		var productDescription string
		productDescriptionRow, ok := productDescriptionMemo.Load(p.ID)
		if !ok {
			if utf8.RuneCountInString(p.Description) > 70 {
				productDescription = string([]rune(p.Description)[:70]) + "???"
			} else {
				productDescription = p.Description
			}
			productDescriptionMemo.Store(p.ID, productDescription)
		} else {
			productDescription = productDescriptionRow.(string)
		}
		p.Description = productDescription

		var newCW []types.CommentWriter
		for _, c := range p.Comments {
			newCW = append(newCW, c)
		}
		p.Comments = newCW
		sProducts = append(sProducts, types.ProductWithComments(p))
	}
	return []byte(templates.Index(isLoginUser, sProducts))
}
