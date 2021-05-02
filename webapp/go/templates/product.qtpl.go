// Code generated by qtc from "product.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line templates/product.qtpl:1
package templates

//line templates/product.qtpl:1
import "github.com/Goryudyuma/ishoconkstm2021gw/webapp/go/types"

//line templates/product.qtpl:2
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line templates/product.qtpl:2
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line templates/product.qtpl:2
func StreamProductPage(qw422016 *qt422016.Writer, currentUser types.User, product types.Product, comments []types.Comment, alreadyBought bool) {
//line templates/product.qtpl:2
	qw422016.N().S(`
`)
//line templates/product.qtpl:3
	streamheader(qw422016, currentUser)
//line templates/product.qtpl:3
	qw422016.N().S(`
<div class="jumbotron">
  <div class="container">
    <h2>`)
//line templates/product.qtpl:6
	qw422016.E().S(product.Name)
//line templates/product.qtpl:6
	qw422016.N().S(`</h2>
    `)
//line templates/product.qtpl:7
	if alreadyBought {
//line templates/product.qtpl:7
		qw422016.N().S(`
      <h4>あなたはすでにこの商品を買っています</h4>
    `)
//line templates/product.qtpl:9
	}
//line templates/product.qtpl:9
	qw422016.N().S(`
  </div>
</div>
<div class="container">
  <div class="row">
    <div class="jumbotron">
      <img src="`)
//line templates/product.qtpl:15
	qw422016.E().S(product.ImagePath)
//line templates/product.qtpl:15
	qw422016.N().S(`" class="img-responsive" width="400"/>
      <h2>価格</h2>
      <p>`)
//line templates/product.qtpl:17
	qw422016.N().D(product.Price)
//line templates/product.qtpl:17
	qw422016.N().S(` 円</p>
      <h2>商品説明</h2>
      <p>`)
//line templates/product.qtpl:19
	qw422016.E().S(product.Description)
//line templates/product.qtpl:19
	qw422016.N().S(`</p>
    </div>
  </div>
</div>
`)
//line templates/product.qtpl:23
	streamfooter(qw422016)
//line templates/product.qtpl:23
	qw422016.N().S(`
`)
//line templates/product.qtpl:24
}

//line templates/product.qtpl:24
func WriteProductPage(qq422016 qtio422016.Writer, currentUser types.User, product types.Product, comments []types.Comment, alreadyBought bool) {
//line templates/product.qtpl:24
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/product.qtpl:24
	StreamProductPage(qw422016, currentUser, product, comments, alreadyBought)
//line templates/product.qtpl:24
	qt422016.ReleaseWriter(qw422016)
//line templates/product.qtpl:24
}

//line templates/product.qtpl:24
func ProductPage(currentUser types.User, product types.Product, comments []types.Comment, alreadyBought bool) string {
//line templates/product.qtpl:24
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/product.qtpl:24
	WriteProductPage(qb422016, currentUser, product, comments, alreadyBought)
//line templates/product.qtpl:24
	qs422016 := string(qb422016.B)
//line templates/product.qtpl:24
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/product.qtpl:24
	return qs422016
//line templates/product.qtpl:24
}
