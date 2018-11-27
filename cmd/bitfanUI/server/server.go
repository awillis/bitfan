package server

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	sessions "github.com/tommy351/gin-sessions"

	"bitfan/api/client"
	"bitfan/api/models"
)

var apiClient *client.RestClient
var apiBaseUrl string

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func Handler(baseURL string, debug bool) http.Handler {
	apiBaseUrl = baseURL
	apiClient = client.New(apiBaseUrl)

	r := gin.New()
	render := NewRender()
	render.Debug = debug
	render.TemplateFuncMap = template.FuncMap{
		"dateFormat": (*templateFunctions)(nil).dateFormat,
		"ago":        (*templateFunctions)(nil).dateAgo,
		"string":     (*templateFunctions)(nil).toString,
		"b64":        (*templateFunctions)(nil).toBase64,
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
		"hasPrefix":    (*templateFunctions)(nil).hasPrefix,
		"hasSuffix":    (*templateFunctions)(nil).hasSuffix,
		"replace":      (*templateFunctions)(nil).replace,
		"markdown":     (*templateFunctions)(nil).toMarkdown,
	}

	r.HTMLRender = render.Init()
	store := sessions.NewCookieStore([]byte("qg498+f"))
	r.Use(sessions.Middleware("my_session", store), gin.Recovery())

	if _, err := os.Stat("assets"); !os.IsNotExist(err) {
		// assets exists on disk, UseFS
		r.StaticFS("/public", http.Dir(filepath.Join("assets", "public")))
	} else {
		// assets from bindData
		r.StaticFS("/public", &assetfs.AssetFS{Asset, AssetDir, AssetInfo, "assets/public"})
	}

	r.GET("/", getPipelines)

	// list pipelines
	r.GET("/logs", getLogs)

	// list env
	r.GET("/env", getEnv)

	// list xprocessors
	r.GET("/xprocessors", getXProcessors)
	// New xprocessor
	r.GET("/xprocessors/:id/new", newXProcessor)
	// Create xprocessor
	r.POST("/xprocessors", createXProcessor)
	// Show Xprocessor
	r.GET("/xprocessors/:id", editXprocessor)
	// Save Xprocessor
	r.POST("/xprocessors/:id", updateXProcessor)
	// Delete Xprocessor
	r.GET("/xprocessors/:id/delete", deleteXProcessor)

	// list pipelines
	r.GET("/pipelines", getPipelines)
	// New pipeline
	r.GET("/pipelines/:id/new", newPipeline)
	// Create pipeline
	r.POST("/pipelines", createPipeline)
	// Show pipeline
	r.GET("/pipelines/:id", editPipeline)
	// Save pipeline
	r.POST("/pipelines/:id", updatePipeline)

	// Start pipeline
	r.GET("/pipelines/:id/start", startPipeline)
	// Restart pipeline
	r.GET("/pipelines/:id/restart", startPipeline)
	// Stop pipeline
	r.GET("/pipelines/:id/stop", stopPipeline)

	// Delete asset
	r.GET("/pipelines/:id/delete", deletePipeline)
	// Show asset
	r.GET("/pipelines/:id/assets/:assetID", showAsset)
	// Create asset
	r.POST("/pipelines/:id/assets", createAsset)
	// Update asset
	r.POST("/pipelines/:id/assets/:assetID", updateAsset)
	// Replace asset
	r.PUT("/pipelines/:id/assets/:assetID", replaceAsset)
	// Download asset
	r.GET("/pipelines/:id/assets/:assetID/download", downloadAsset)
	// Delete asset
	r.GET("/pipelines/:id/assets/:assetID/delete", deleteAsset)

	//Pipeline Playground
	r.GET("/pipelines/:id/play", playgroundPipeline)
	r.PUT("/pipelines/:id/play", playgroundsPlayDo)
	r.DELETE("/pipelines/:id/play", playgroundsPlayExit)

	// Playgrounds
	r.GET("/playgrounds/play", playgroundsPlay)
	r.PUT("/playgrounds/play", playgroundsPlayDo)
	r.DELETE("/playgrounds/play", playgroundsPlayExit)

	// Replace asset
	r.PUT("/settings/api", changeBitfanApiURL)

	return r
}

func changeBitfanApiURL(c *gin.Context) {
	var values map[string]interface{}
	err := c.BindJSON(&values)
	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
		return
	}

	newURL := values["url"].(string)

	if newURL == "" {
		c.JSON(500, "provide the bitfan api host:port")
		log.Printf("error : %v\n", err)
		return
	}

	apiBaseUrl = newURL
	apiClient = client.New(apiBaseUrl)
	c.JSON(200, values)
}

func withCommonValues(c *gin.Context, h gin.H) gin.H {
	session := sessions.Get(c)
	h["apiHost"] = apiBaseUrl
	h["flashes"] = session.Flashes()
	session.Save()
	return h
}

func getLogs(c *gin.Context) {
	// TODO : proxy WS:// github.com/koding/websocketproxy
	c.HTML(200, "logs/logs", withCommonValues(c, gin.H{
		"bitfanHost": apiBaseUrl,
	}))
}

func getEnv(c *gin.Context) {
	c.HTML(200, "env/index", withCommonValues(c, gin.H{
		"error": nil,
	}))
}

func getXProcessors(c *gin.Context) {
	producers, err := apiClient.XProcessors("producer")
	transformers, err := apiClient.XProcessors("transformer")
	consumers, err := apiClient.XProcessors("consumer")
	c.HTML(200, "xprocessors/index", withCommonValues(c, gin.H{
		"producers":    producers,
		"transformers": transformers,
		"consumers":    consumers,
		"error":        err,
	}))
}

func newXProcessor(c *gin.Context) {
	c.HTML(200, "xprocessors/new", withCommonValues(c, gin.H{
		"xprocessor": models.XProcessor{
			Behavior: c.Query("behavior"),
		},
	}))
}

func createXProcessor(c *gin.Context) {
	c.Request.ParseForm()

	var p = models.XProcessor{
		Label:       c.Request.PostFormValue("label"),
		Behavior:    c.Request.PostFormValue("behavior"),
		Kind:        c.Request.PostFormValue("kind"),
		Description: c.Request.PostFormValue("description"),
	}

	if p.Kind == "php" {
		p.Command = "php"
	}
	if p.Kind == "python" {
		p.Command = "python"
	}
	if p.Kind == "golang" {
		p.Command = "go"
	}

	pnew, _ := apiClient.NewXProcessor(&p)
	flash(c, fmt.Sprintf("XProcessor %s created", p.Label))
	c.Redirect(302, fmt.Sprintf("/xprocessors/%s", pnew.Uuid))
}
func updateXProcessor(c *gin.Context) {
	c.Request.ParseForm()
	xprocessorUUID := c.Param("id")

	var data = map[string]interface{}{
		"stream":  false,
		"has_doc": false,
	}
	if _, ok := c.Request.PostForm["behavior"]; ok {
		data["behavior"] = c.Request.PostFormValue("behavior")
	}
	if _, ok := c.Request.PostForm["label"]; ok {
		data["label"] = c.Request.PostFormValue("label")
	}
	if _, ok := c.Request.PostForm["description"]; ok {
		data["description"] = c.Request.PostFormValue("description")
	}
	if _, ok := c.Request.PostForm["stream"]; ok {
		data["stream"] = true
	}
	if _, ok := c.Request.PostForm["has_doc"]; ok {
		data["has_doc"] = true
	}
	if _, ok := c.Request.PostForm["stdin_as"]; ok {
		data["stdin_as"] = c.Request.PostFormValue("stdin_as")
	}
	if _, ok := c.Request.PostForm["stdout_as"]; ok {
		data["stdout_as"] = c.Request.PostFormValue("stdout_as")
	}
	if _, ok := c.Request.PostForm["code"]; ok {
		data["code"] = c.Request.PostFormValue("code")
	}
	if _, ok := c.Request.PostForm["command"]; ok {
		data["command"] = c.Request.PostFormValue("command")
	}

	if _, ok := c.Request.PostForm["args"]; ok {
		data["args"] = strings.Split(c.Request.PostFormValue("args"), "\n")
	}
	if _, ok := c.Request.PostForm["options_composition_tpl"]; ok {
		data["options_composition_tpl"] = c.Request.PostFormValue("options_composition_tpl")
	}

	pnew, _ := apiClient.UpdateXProcessor(xprocessorUUID, &data)
	c.JSON(200, pnew)
}
func editXprocessor(c *gin.Context) {
	id := c.Param("id")

	p, _ := apiClient.XProcessor(id)
	c.HTML(200, "xprocessors/edit", withCommonValues(c, gin.H{
		"xprocessor": p,
	}))
}

func deleteXProcessor(c *gin.Context) {
	uuid := c.Param("id")

	err := apiClient.DeleteXProcessor(uuid)

	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
		return
	}

	c.Redirect(302, "/xprocessors")
}

func getPipelines(c *gin.Context) {
	pipelines, err := apiClient.Pipelines()
	c.HTML(200, "pipelines/index", withCommonValues(c, gin.H{
		"pipelines": pipelines,
		"error":     err,
	}))
}

func flash(c *gin.Context, message string) {
	session := sessions.Get(c)
	session.AddFlash(message)
	session.Save()
}

func editPipeline(c *gin.Context) {
	id := c.Param("id")

	p, _ := apiClient.Pipeline(id)

	c.HTML(200, "pipelines/edit", withCommonValues(c, gin.H{
		"pipeline": p,
	}))

}

func newPipeline(c *gin.Context) {
	c.HTML(200, "pipelines/new", withCommonValues(c, gin.H{}))
}

func createPipeline(c *gin.Context) {
	c.Request.ParseForm()

	defaultValue := []byte("input{ }\n\nfilter{ }\n\noutput{ }")
	var p = models.Pipeline{
		Label:       c.Request.PostFormValue("label"),
		Description: c.Request.PostFormValue("description"),
		Assets: []models.Asset{{
			Name:        "entrypoint.conf",
			Type:        "entrypoint",
			ContentType: "text/plain",
			Value:       defaultValue,
			Size:        len(defaultValue),
		}},
	}

	pnew, _ := apiClient.NewPipeline(&p)
	flash(c, fmt.Sprintf("Pipeline %s created", p.Label))
	c.Redirect(302, fmt.Sprintf("/pipelines/%s", pnew.Uuid))
}

func updatePipeline(c *gin.Context) {
	c.Request.ParseForm()
	pipelineUUID := c.Param("id")

	var data = map[string]interface{}{
		"label":       c.Request.PostFormValue("label"),
		"description": c.Request.PostFormValue("description"),
		"auto_start":  false,
	}

	if _, ok := c.Request.PostForm["auto_start"]; ok {
		data["auto_start"] = true
	}

	pnew, _ := apiClient.UpdatePipeline(pipelineUUID, &data)
	flash(c, fmt.Sprintf("Pipeline %s saved", pnew.Label))
	c.Redirect(302, fmt.Sprintf("/pipelines/%s", pnew.Uuid))
}

func startPipeline(c *gin.Context) {
	pipelineUUID := c.Param("id")

	pipeline, err := apiClient.StartPipeline(pipelineUUID)
	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
	} else {
		c.JSON(200, pipeline)
	}
}

func stopPipeline(c *gin.Context) {
	pipelineUUID := c.Param("id")

	pipeline, err := apiClient.StopPipeline(pipelineUUID)
	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
	} else {
		c.JSON(200, pipeline)
	}
}

func deletePipeline(c *gin.Context) {
	pipelineUUID := c.Param("id")

	err := apiClient.DeletePipeline(pipelineUUID)

	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
		return
	}

	c.Redirect(302, "/pipelines")
}

func createAsset(c *gin.Context) {
	pipelineUUID := c.Param("id")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}

	filename := header.Filename

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Fatal(err)
	}

	nasset := &models.Asset{
		Name:         filename,
		Value:        buf.Bytes(),
		PipelineUUID: pipelineUUID,
	}

	asset, err := apiClient.NewAsset(nasset)

	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
		return
	}

	c.JSON(201, asset)
}

func showAsset(c *gin.Context) {
	pipelineUUID := c.Param("id")
	assetUUID := c.Param("assetID")

	p, _ := apiClient.Pipeline(pipelineUUID)
	a, _ := apiClient.Asset(assetUUID)

	c.HTML(200, "pipelines/assets/edit", withCommonValues(c, gin.H{
		"asset":    a,
		"pipeline": p,
	}))
}

func deleteAsset(c *gin.Context) {
	pipelineUUID := c.Param("id")
	assetUUID := c.Param("assetID")

	err := apiClient.DeleteAsset(assetUUID)

	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
		return
	}

	c.Redirect(302, fmt.Sprintf("/pipelines/%s", pipelineUUID))
}

func downloadAsset(c *gin.Context) {
	assetUUID := c.Param("assetID")

	asset, err := apiClient.Asset(assetUUID)

	if err != nil {
		c.JSON(404, err.Error())
		log.Printf("error : %v\n", err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=\""+asset.Name+"\"")
	c.Data(200, asset.ContentType, asset.Value)
}

func updateAsset(c *gin.Context) {
	c.Request.ParseForm()
	assetUUID := c.Param("assetID")

	var data = map[string]interface{}{
		"name": c.Request.PostFormValue("name"),
	}
	if _, ok := c.Request.PostForm["content"]; ok {
		data["value"] = []byte(c.Request.PostFormValue("content"))
	}

	asset, err := apiClient.UpdateAsset(assetUUID, &data)
	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
		return
	}

	flash(c, fmt.Sprintf("Asset %s saved", asset.Name))

	c.Redirect(302, fmt.Sprintf("/pipelines/%s/assets/%s", asset.PipelineUUID, asset.Uuid))
}

func replaceAsset(c *gin.Context) {
	pipelineUUID := c.Param("id")
	assetUUID := c.Param("assetID")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}

	filename := header.Filename

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Fatal(err)
	}

	nasset := &models.Asset{
		Uuid:         assetUUID,
		Name:         filename,
		Value:        buf.Bytes(),
		PipelineUUID: pipelineUUID,
	}

	asset, err := apiClient.ReplaceAsset(assetUUID, nasset)
	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
		return
	}

	flash(c, fmt.Sprintf("Asset %s saved", asset.Name))

	c.JSON(200, asset)
}
