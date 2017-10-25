package ui

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/k0kubun/pp"
	uuid "github.com/nu7hatch/gouuid"

	eztemplate "github.com/vjeantet/ez-gin-template"
)

var db *gorm.DB

func Handler(assetsPath, path string) http.Handler {
	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema

	db.AutoMigrate(&Pipeline{}, &Asset{})

	// database, _ := sql.Open("sqlite3", "./nraboy.db")
	// statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	// statement.Exec()
	// statement, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	// statement.Exec("Nic", "Raboy")
	// rows, _ := database.Query("SELECT id, firstname, lastname FROM people")

	// pp.Println("rows-->", rows)

	r := gin.New()
	render := eztemplate.New()
	render.TemplatesDir = assetsPath + "/views/" // default
	render.Ext = ".html"                         // default
	render.Debug = true                          // default
	render.TemplateFuncMap = template.FuncMap{
		"dateFormat": (*templateFunctions)(nil).dateFormat,
		"ago":        (*templateFunctions)(nil).dateAgo,
		"string":     (*templateFunctions)(nil).toString,
		"int":        (*templateFunctions)(nil).toInt,
		"time":       (*templateFunctions)(nil).asTime,
		"now":        (*templateFunctions)(nil).now,
		"isset":      (*templateFunctions)(nil).isSet,

		"numFmt": (*templateFunctions)(nil).numFmt,

		"safeHTML":     (*templateFunctions)(nil).safeHtml,
		"hTMLUnescape": (*templateFunctions)(nil).htmlUnescape,
		"hTMLEscape":   (*templateFunctions)(nil).htmlEscape,
		"lower":        (*templateFunctions)(nil).lower,
		"upper":        (*templateFunctions)(nil).upper,
		"trim":         (*templateFunctions)(nil).trim,
		"trimPrefix":   (*templateFunctions)(nil).trimPrefix,
		"replace":      (*templateFunctions)(nil).replace,
		"markdown":     (*templateFunctions)(nil).toMarkdown,
	}

	r.HTMLRender = render.Init()
	store := sessions.NewCookieStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store), gin.Recovery())
	// r.Use(gin.Recovery())
	g := r.Group(path)
	{
		g.StaticFS("/public", http.Dir(assetsPath+"/public"))
		g.GET("/", index)

		g.GET("/pipelines", getPipelines)
		g.POST("/pipelines", createPipeline)

		g.GET("/pipelines/:id", editPipeline)
		g.POST("/pipelines/:id", updatePipeline)

		g.GET("/pipelines/:id/new", newPipeline)

		g.POST("/pipelines/:id/assets", createAsset)
		g.GET("/pipelines/:id/assets/:assetID", showAsset)
		g.POST("/pipelines/:id/assets/:assetID", updateAsset)
		g.GET("/pipelines/:id/assets/:assetID/delete", deleteAsset)

		g.GET("/pipelines/:id/start", index)
		g.GET("/pipelines/:id/stop", index)
	}

	return r
}

func newPipeline(c *gin.Context) {

	c.HTML(200, "pipelines/new", gin.H{})
}

func deleteAsset(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		pp.Println("c-->", c.Param("id"))
		c.Abort()
		return
	}

	assetID, err := strconv.Atoi(c.Param("assetID"))
	if err != nil {
		pp.Println("c-->", c.Param("assetID"))
		c.Abort()
		return
	}

	var p = Pipeline{ID: id}
	var a = Asset{ID: assetID}
	db.First(&p)
	db.Model(&p).Association("Assets").Delete(&a)
	db.Delete(&a)

	c.Redirect(302, fmt.Sprintf("/ui/pipelines/%d", p.ID))
}

func updateAsset(c *gin.Context) {
	c.Request.ParseForm()
	pp.Println("c.Request.PostForm-->", c.Request.PostForm)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		pp.Println("c-->", c.Param("id"))
		c.Abort()
		return
	}

	assetID, err := strconv.Atoi(c.Param("assetID"))
	if err != nil {
		pp.Println("c-->", c.Param("assetID"))
		c.Abort()
		return
	}

	var p = Pipeline{ID: id}
	db.First(&p)
	var a = Asset{ID: assetID}
	db.First(&a)

	a.Name = c.Request.PostFormValue("name")
	a.Value = []byte(c.Request.PostFormValue("content"))
	db.Save(&a)

	session := sessions.Default(c)
	session.AddFlash(fmt.Sprintf("Asset %s saved", a.Name))
	session.Save()

	c.Redirect(302, fmt.Sprintf("/ui/pipelines/%d/assets/%d", p.ID, a.ID))
}

func createAsset(c *gin.Context) {

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}

	filename := header.Filename
	contentType := header.Header.Get("Content-Type")
	size := header.Size

	pp.Println(header)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Fatal(err)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Abort()
		return
	}

	var p = Pipeline{ID: id}
	db.First(&p)

	var a = Asset{
		Name:        filename,
		ContentType: contentType,
		Value:       buf.Bytes(),
		Size:        int(size),
	}
	p.Assets = append(p.Assets, a)
	db.Save(&p)
	c.String(200, "ok")
}

func showAsset(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		pp.Println("c-->", c.Param("id"))
		c.Abort()
		return
	}

	assetID, err := strconv.Atoi(c.Param("assetID"))
	if err != nil {
		pp.Println("c-->", c.Param("assetID"))
		c.Abort()
		return
	}

	var a = Asset{ID: (assetID)}
	db.First(&a)

	var p = Pipeline{ID: id}
	db.Preload("Assets").First(&p)

	flashes := []string{}
	for _, m := range sessions.Default(c).Flashes() {
		flashes = append(flashes, m.(string))
	}
	sessions.Default(c).Save()

	c.HTML(200, "assets/edit", gin.H{
		"asset":    a,
		"pipeline": p,
		"flashes":  flashes,
	})
}

func createPipeline(c *gin.Context) {
	c.Request.ParseForm()

	uid, _ := uuid.NewV4()
	var p = Pipeline{
		Version:     1,
		Label:       c.Request.PostFormValue("label"),
		Description: c.Request.PostFormValue("description"),
		Uuid:        uid.String(),
		Content:     "input{ }\n\nfilter{ }\n\noutput{ }",
	}

	db.Create(&p)
	session := sessions.Default(c)
	session.AddFlash(fmt.Sprintf("Pipeline %s created", p.Label))
	session.Save()

	c.Redirect(302, fmt.Sprintf("/ui/pipelines/%d", p.ID))
}

func updatePipeline(c *gin.Context) {
	c.Request.ParseForm()
	pp.Println("c.Request.PostForm-->", c.Request.PostForm)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Abort()
		return
	}

	var p = Pipeline{ID: id}
	db.First(&p)

	p.ID = 0
	p.Version = p.Version + 1
	p.Label = c.Request.PostFormValue("label")
	p.Content = c.Request.PostFormValue("content")
	p.Description = c.Request.PostFormValue("description")
	p.UpdatedAt = time.Now()
	db.Create(&p)

	session := sessions.Default(c)
	session.AddFlash(fmt.Sprintf("Pipeline %s saved (version %d)", p.Label, p.Version))
	session.Save()

	c.Redirect(302, fmt.Sprintf("/ui/pipelines/%d", p.ID))
}

func editPipeline(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Abort()
		return
	}

	var p = Pipeline{ID: (idInt)}
	db.Preload("Assets").First(&p)

	flashes := []string{}
	for _, m := range sessions.Default(c).Flashes() {
		flashes = append(flashes, m.(string))
	}
	sessions.Default(c).Save()

	c.HTML(200, "pipelines/edit", gin.H{
		"pipeline": p,
		"flashes":  flashes,
	})

}

func index(c *gin.Context) {

	var pipelines []Pipeline

	rows, err := db.Model(&Pipeline{}).Select("*, MAX(version)").Group("uuid").Rows()
	if err != nil {
		c.AbortWithError(500, err)
	}

	defer rows.Close()

	for rows.Next() {
		var p Pipeline
		db.ScanRows(rows, &p)
		pipelines = append(pipelines, p)
	}

	c.HTML(200, "pipelines/index", gin.H{
		"pipelines": pipelines,
	})
}

func getPipelines(c *gin.Context) {

	// var err error

	// ppls := core.Pipelines()

	// if err == nil {
	// 	c.JSON(200, ppls)
	// } else {
	// 	c.JSON(404, gin.H{"error": "no pipelines(s) running"})
	// }
	// curl -i http://localhost:8080/api/v1/pipelines
}
