// Code generated by qtc from "template.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line templates/template.qtpl:1
package templates

//line templates/template.qtpl:1
import "github.com/Goryudyuma/ishoconkstm2021gw/webapp/go/types"

//line templates/template.qtpl:3
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line templates/template.qtpl:3
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line templates/template.qtpl:3
func StreamProductPage(qw422016 *qt422016.Writer, currentUser types.User, product types.Product, alreadyBought bool) {
//line templates/template.qtpl:3
	qw422016.N().S(`
`)
//line templates/template.qtpl:4
	StreamHeader(qw422016, currentUser)
//line templates/template.qtpl:4
	qw422016.N().S(`
<div class="jumbotron">
  <div class="container">
    <h2>`)
//line templates/template.qtpl:7
	qw422016.E().S(product.Name)
//line templates/template.qtpl:7
	qw422016.N().S(`</h2>
    `)
//line templates/template.qtpl:8
	if alreadyBought {
//line templates/template.qtpl:8
		qw422016.N().S(`
      <h4>あなたはすでにこの商品を買っています</h4>
    `)
//line templates/template.qtpl:10
	}
//line templates/template.qtpl:10
	qw422016.N().S(`
  </div>
</div>
<div class="container">
  <div class="row">
    <div class="jumbotron">
      <img src="`)
//line templates/template.qtpl:16
	qw422016.E().S(product.ImagePath)
//line templates/template.qtpl:16
	qw422016.N().S(`" class="img-responsive" width="400"/>
      <h2>価格</h2>
      <p>`)
//line templates/template.qtpl:18
	qw422016.N().D(product.Price)
//line templates/template.qtpl:18
	qw422016.N().S(` 円</p>
      <h2>商品説明</h2>
      <p>`)
//line templates/template.qtpl:20
	qw422016.E().S(product.Description)
//line templates/template.qtpl:20
	qw422016.N().S(`</p>
    </div>
  </div>
</div>
`)
//line templates/template.qtpl:24
	StreamFooter(qw422016)
//line templates/template.qtpl:24
	qw422016.N().S(`
`)
//line templates/template.qtpl:25
}

//line templates/template.qtpl:25
func WriteProductPage(qq422016 qtio422016.Writer, currentUser types.User, product types.Product, alreadyBought bool) {
//line templates/template.qtpl:25
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/template.qtpl:25
	StreamProductPage(qw422016, currentUser, product, alreadyBought)
//line templates/template.qtpl:25
	qt422016.ReleaseWriter(qw422016)
//line templates/template.qtpl:25
}

//line templates/template.qtpl:25
func ProductPage(currentUser types.User, product types.Product, alreadyBought bool) string {
//line templates/template.qtpl:25
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/template.qtpl:25
	WriteProductPage(qb422016, currentUser, product, alreadyBought)
//line templates/template.qtpl:25
	qs422016 := string(qb422016.B)
//line templates/template.qtpl:25
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/template.qtpl:25
	return qs422016
//line templates/template.qtpl:25
}

//line templates/template.qtpl:27
func StreamMyPage(qw422016 *qt422016.Writer, currentUser bool, user types.User, products []types.Product, totalPay int) {
//line templates/template.qtpl:27
	qw422016.N().S(`
<div class="jumbotron">
  <div class="container">
    <h2>`)
//line templates/template.qtpl:30
	qw422016.E().S(user.Name)
//line templates/template.qtpl:30
	qw422016.N().S(` さんの購入履歴</h2>
    <h4>合計金額: `)
//line templates/template.qtpl:31
	qw422016.N().D(totalPay)
//line templates/template.qtpl:31
	qw422016.N().S(`円</h4>
  </div>
</div>
<div class="container">
  <div class="row">
    `)
//line templates/template.qtpl:36
	for _, product := range products {
//line templates/template.qtpl:36
		qw422016.N().S(`
      <div class="col-md-4">
        <div class="panel panel-default">
          <div class="panel-heading">
            <a href="/products/`)
//line templates/template.qtpl:40
		qw422016.N().D(product.ID)
//line templates/template.qtpl:40
		qw422016.N().S(`">`)
//line templates/template.qtpl:40
		qw422016.E().S(product.Name)
//line templates/template.qtpl:40
		qw422016.N().S(`</a>
          </div>
          <div class="panel-body">
            <a href="/products/`)
//line templates/template.qtpl:43
		qw422016.N().D(product.ID)
//line templates/template.qtpl:43
		qw422016.N().S(`"><img src="`)
//line templates/template.qtpl:43
		qw422016.E().S(product.ImagePath)
//line templates/template.qtpl:43
		qw422016.N().S(`" class="img-responsive" /></a>
            <h4>価格</h4>
            <p>`)
//line templates/template.qtpl:45
		qw422016.N().D(product.Price)
//line templates/template.qtpl:45
		qw422016.N().S(`円</p>
            <h4>商品説明</h4>
            <p>`)
//line templates/template.qtpl:47
		qw422016.E().S(product.Description)
//line templates/template.qtpl:47
		qw422016.N().S(`</p>
            <h4>購入日時</h4>
            <p>`)
//line templates/template.qtpl:49
		qw422016.E().S(product.CreatedAt)
//line templates/template.qtpl:49
		qw422016.N().S(`</p>
          </div>
          `)
//line templates/template.qtpl:51
		if currentUser {
//line templates/template.qtpl:51
			qw422016.N().S(`
            <div class="panel-footer">
              <form method="POST" action="/comments/`)
//line templates/template.qtpl:53
			qw422016.N().D(product.ID)
//line templates/template.qtpl:53
			qw422016.N().S(`">
                <fieldset>
                  <div class="form-group">
                    <input class="form-control" placeholder="Comment Here" name="content" value="">
                  </div>
                  <input class="btn btn-success btn-block" type="submit" name="send_comment" value="コメントを送信" />
                </fieldset>
              </form>
            </div>
          `)
//line templates/template.qtpl:62
		}
//line templates/template.qtpl:62
		qw422016.N().S(`
        </div>
      </div>
    `)
//line templates/template.qtpl:65
	}
//line templates/template.qtpl:65
	qw422016.N().S(`
  </div>
</div>
`)
//line templates/template.qtpl:68
}

//line templates/template.qtpl:68
func WriteMyPage(qq422016 qtio422016.Writer, currentUser bool, user types.User, products []types.Product, totalPay int) {
//line templates/template.qtpl:68
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/template.qtpl:68
	StreamMyPage(qw422016, currentUser, user, products, totalPay)
//line templates/template.qtpl:68
	qt422016.ReleaseWriter(qw422016)
//line templates/template.qtpl:68
}

//line templates/template.qtpl:68
func MyPage(currentUser bool, user types.User, products []types.Product, totalPay int) string {
//line templates/template.qtpl:68
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/template.qtpl:68
	WriteMyPage(qb422016, currentUser, user, products, totalPay)
//line templates/template.qtpl:68
	qs422016 := string(qb422016.B)
//line templates/template.qtpl:68
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/template.qtpl:68
	return qs422016
//line templates/template.qtpl:68
}

//line templates/template.qtpl:70
func StreamLogin(qw422016 *qt422016.Writer, message string) {
//line templates/template.qtpl:70
	qw422016.N().S(`
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html" charset="utf-8">
<link rel="stylesheet" href="/css/bootstrap.min.css">
<title>すごいECサイト</title>
</head>
<body class="container">
  <h1 class="jumbotron"><a href="/">すごいECサイト</a></h1>
  <div class="container">
    <div class="row">
      <div class="col-md-4 col-md-offset-4">
        <div class="text-danger" id="logout-message">`)
//line templates/template.qtpl:83
	qw422016.E().S(message)
//line templates/template.qtpl:83
	qw422016.N().S(`</div>
        <div class="login-panel panel panel-default">
          <div class="panel-heading">
            <h3 class="panel-title">Sign In</h3>
          </div>
          <div class="panel-body">
            <form method="POST" action="/login">
              <fieldset>
                <div class="form-group">
                  <input class="form-control" placeholder="E-mail" name="email" type="email" autofocus>
                </div>
                <div class="form-group">
                  <input class="form-control" placeholder="Password" name="password" type="password" value="">
                </div>
                <!-- Change this to a button or input when using this as a form -->
                <input class="btn btn-lg btn-success btn-block" type="submit" name="Login" value="Login" />
              </fieldset>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</body>
</html>
`)
//line templates/template.qtpl:108
}

//line templates/template.qtpl:108
func WriteLogin(qq422016 qtio422016.Writer, message string) {
//line templates/template.qtpl:108
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/template.qtpl:108
	StreamLogin(qw422016, message)
//line templates/template.qtpl:108
	qt422016.ReleaseWriter(qw422016)
//line templates/template.qtpl:108
}

//line templates/template.qtpl:108
func Login(message string) string {
//line templates/template.qtpl:108
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/template.qtpl:108
	WriteLogin(qb422016, message)
//line templates/template.qtpl:108
	qs422016 := string(qb422016.B)
//line templates/template.qtpl:108
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/template.qtpl:108
	return qs422016
//line templates/template.qtpl:108
}

//line templates/template.qtpl:110
func StreamHeader(qw422016 *qt422016.Writer, currentUser types.User) {
//line templates/template.qtpl:110
	qw422016.N().S(`
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html" charset="utf-8">
<link rel="stylesheet" href="/css/bootstrap.min.css">
<title>すごいECサイト</title>
</head>

<body>
<nav class="navbar navbar-inverse navbar-fixed-top">
  <div class="container">
    <div class="navbar-header">
      <a class="navbar-brand" href="/">すごいECサイトで爆買いしよう!</a>
    </div>
    <div class="header clearfix">
    `)
//line templates/template.qtpl:126
	if currentUser.ID != 0 {
//line templates/template.qtpl:126
		qw422016.N().S(`
      <nav>
        <ul class="nav nav-pills pull-right">
          <li role="presentation"><a href="/users/`)
//line templates/template.qtpl:129
		qw422016.N().D(currentUser.ID)
//line templates/template.qtpl:129
		qw422016.N().S(`">`)
//line templates/template.qtpl:129
		qw422016.E().S(currentUser.Name)
//line templates/template.qtpl:129
		qw422016.N().S(`さんの購入履歴</a></li>
          <li role="presentation"><a href="/logout">Logout</a></li>
        </ul>
      </nav>
    `)
//line templates/template.qtpl:133
	} else {
//line templates/template.qtpl:133
		qw422016.N().S(`
    <nav>
      <ul class="nav nav-pills pull-right">
        <li role="presentation"><a href="/login">Login</a></li>
      </ul>
    </nav>
    `)
//line templates/template.qtpl:139
	}
//line templates/template.qtpl:139
	qw422016.N().S(`
  </div>
</nav>

`)
//line templates/template.qtpl:143
}

//line templates/template.qtpl:143
func WriteHeader(qq422016 qtio422016.Writer, currentUser types.User) {
//line templates/template.qtpl:143
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/template.qtpl:143
	StreamHeader(qw422016, currentUser)
//line templates/template.qtpl:143
	qt422016.ReleaseWriter(qw422016)
//line templates/template.qtpl:143
}

//line templates/template.qtpl:143
func Header(currentUser types.User) string {
//line templates/template.qtpl:143
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/template.qtpl:143
	WriteHeader(qb422016, currentUser)
//line templates/template.qtpl:143
	qs422016 := string(qb422016.B)
//line templates/template.qtpl:143
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/template.qtpl:143
	return qs422016
//line templates/template.qtpl:143
}

//line templates/template.qtpl:144
func StreamFooter(qw422016 *qt422016.Writer) {
//line templates/template.qtpl:144
	qw422016.N().S(`

</body>
</html>
`)
//line templates/template.qtpl:148
}

//line templates/template.qtpl:148
func WriteFooter(qq422016 qtio422016.Writer) {
//line templates/template.qtpl:148
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/template.qtpl:148
	StreamFooter(qw422016)
//line templates/template.qtpl:148
	qt422016.ReleaseWriter(qw422016)
//line templates/template.qtpl:148
}

//line templates/template.qtpl:148
func Footer() string {
//line templates/template.qtpl:148
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/template.qtpl:148
	WriteFooter(qb422016)
//line templates/template.qtpl:148
	qs422016 := string(qb422016.B)
//line templates/template.qtpl:148
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/template.qtpl:148
	return qs422016
//line templates/template.qtpl:148
}

//line templates/template.qtpl:150
func StreamIndex(qw422016 *qt422016.Writer, isLoginUser bool, products []types.ProductWithComments) {
//line templates/template.qtpl:150
	qw422016.N().S(`
<div class="jumbotron">
  <div class="container">
    <h1>今日は大安売りの日です！</h1>
  </div>
</div>
<div class="container">
  <div class="row">
    `)
//line templates/template.qtpl:158
	for _, product := range products {
//line templates/template.qtpl:158
		qw422016.N().S(`
      <div class="col-md-4">
        <div class="panel panel-default">
          <div class="panel-heading">
            <a href="/products/`)
//line templates/template.qtpl:162
		qw422016.N().D(product.ID)
//line templates/template.qtpl:162
		qw422016.N().S(`">`)
//line templates/template.qtpl:162
		qw422016.E().S(product.Name)
//line templates/template.qtpl:162
		qw422016.N().S(`</a>
          </div>
          <div class="panel-body">
            <a href="/products/`)
//line templates/template.qtpl:165
		qw422016.N().D(product.ID)
//line templates/template.qtpl:165
		qw422016.N().S(`"><img src="`)
//line templates/template.qtpl:165
		qw422016.E().S(product.ImagePath)
//line templates/template.qtpl:165
		qw422016.N().S(`" class="img-responsive" /></a>
            <h4>価格</h4>
            <p>`)
//line templates/template.qtpl:167
		qw422016.N().D(product.Price)
//line templates/template.qtpl:167
		qw422016.N().S(`円</p>
            <h4>商品説明</h4>
            <p>`)
//line templates/template.qtpl:169
		qw422016.E().S(product.Description)
//line templates/template.qtpl:169
		qw422016.N().S(`</p>
            <h4>`)
//line templates/template.qtpl:170
		qw422016.N().D(product.CommentCount)
//line templates/template.qtpl:170
		qw422016.N().S(`件のレビュー</h4>
            <ul>
              `)
//line templates/template.qtpl:172
		for _, cw := range product.Comments {
//line templates/template.qtpl:172
			qw422016.N().S(`
                <li>`)
//line templates/template.qtpl:173
			qw422016.E().S(cw.Content)
//line templates/template.qtpl:173
			qw422016.N().S(` by `)
//line templates/template.qtpl:173
			qw422016.E().S(cw.Writer)
//line templates/template.qtpl:173
			qw422016.N().S(`</li>
              `)
//line templates/template.qtpl:174
		}
//line templates/template.qtpl:174
		qw422016.N().S(`
            </ul>
          </div>
          `)
//line templates/template.qtpl:177
		if isLoginUser {
//line templates/template.qtpl:177
			qw422016.N().S(`
            <div class="panel-footer">
              <form method="POST" action="/products/buy/`)
//line templates/template.qtpl:179
			qw422016.N().D(product.ID)
//line templates/template.qtpl:179
			qw422016.N().S(`">
                <fieldset>
                  <input class="btn btn-success btn-block" type="submit" name="buy" value="購入" />
                </fieldset>
              </form>
            </div>
          `)
//line templates/template.qtpl:185
		}
//line templates/template.qtpl:185
		qw422016.N().S(`
        </div>
      </div>
    `)
//line templates/template.qtpl:188
	}
//line templates/template.qtpl:188
	qw422016.N().S(`
  </div>
</div>
`)
//line templates/template.qtpl:191
}

//line templates/template.qtpl:191
func WriteIndex(qq422016 qtio422016.Writer, isLoginUser bool, products []types.ProductWithComments) {
//line templates/template.qtpl:191
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/template.qtpl:191
	StreamIndex(qw422016, isLoginUser, products)
//line templates/template.qtpl:191
	qt422016.ReleaseWriter(qw422016)
//line templates/template.qtpl:191
}

//line templates/template.qtpl:191
func Index(isLoginUser bool, products []types.ProductWithComments) string {
//line templates/template.qtpl:191
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/template.qtpl:191
	WriteIndex(qb422016, isLoginUser, products)
//line templates/template.qtpl:191
	qs422016 := string(qb422016.B)
//line templates/template.qtpl:191
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/template.qtpl:191
	return qs422016
//line templates/template.qtpl:191
}
