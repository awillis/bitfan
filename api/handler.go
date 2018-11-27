package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"bitfan/core"
)

// var plugins map[string]map[string]core.ProcessorFactory
func init() {
	gin.SetMode(gin.ReleaseMode)
}

var apiLogger *core.Logger

func Handler(path string) http.Handler {

	apiLogger = core.NewLogger("api", nil)

	logs, _ := newHook(hookConfig{Size: 100})
	logrus.AddHook(logs)

	r := gin.New()
	r.Use(
		gin.Recovery(),
		func(c *gin.Context) {
			c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, PATCH")
			c.Writer.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, PATCH")
			if c.Request.Method == "OPTIONS" {
				c.Writer.Header().Add("Access-Control-Allow-Headers", "*")
				c.JSON(http.StatusOK, struct{}{})
				return
			}
			c.Next()
		},
	)
	v2 := r.Group(path)
	{

		pipelineCtrl := &PipelineApiController{
			path: path,
		}

		xprocessorCtrl := &XProcessorApiController{
			path: path,
		}

		assetCtrl := &AssetApiController{
			path: path,
		}

		envvariablesCtrl := &EnvApiController{
			path: path,
		}

		dbCtrl := &DatabaseController{}

		logsCtrl := &LogApiController{
			Hub: newHub(logs.String),
		}
		logs.AddChan(logsCtrl.Hub.Broadcast)
		go logsCtrl.Hub.run()

		docsCtrl := &DocsController{}

		v2.GET("/logs", logsCtrl.Stream) // Websocket

		// curl -i -X POST http://localhost:5123/api/v2/pipelines
		v2.POST("/pipelines", pipelineCtrl.Create) // créer pipeline

		v2.GET("/xprocessors", xprocessorCtrl.Find)                  // list xprocessors
		v2.POST("/xprocessors", xprocessorCtrl.Create)               // list xprocessors
		v2.GET("/xprocessors/:uuid", xprocessorCtrl.FindOneByUUID)   // show xprocessors
		v2.PATCH("/xprocessors/:uuid", xprocessorCtrl.UpdateByUUID)  // update xprocessors
		v2.DELETE("/xprocessors/:uuid", xprocessorCtrl.DeleteByUUID) // delete xprocessors

		// curl -i -X GET http://localhost:5123/api/v2/pipelines
		v2.GET("/pipelines", pipelineCtrl.Find)            // list pipelines
		v2.GET("/pipelines.zip", pipelineCtrl.DownloadAll) // backup
		// curl -i -X GET http://localhost:5123/api/v2/pipelines/408b9a7b-933e-4d3d-6df1-65324a0a5315
		v2.GET("/pipelines/:uuid", pipelineCtrl.FindOneByUUID) // show pipeline

		// curl -i -X PATCH http://localhost:5123/api/v2/pipelines/408b9a7b-933e-4d3d-6df1-65324a0a5315
		v2.PATCH("/pipelines/:uuid", pipelineCtrl.UpdateByUUID) // update pipeline / stop / start / restart

		// curl -i -X DELETE http://localhost:5123/api/v2/pipelines/408b9a7b-933e-4d3d-6df1-65324a0a5315
		v2.DELETE("/pipelines/:uuid", pipelineCtrl.DeleteByUUID) // delete pipeline

		v2.POST("/assets", assetCtrl.Create)                         // create asset
		v2.GET("/assets/:uuid", assetCtrl.FindOneByUUID)             // show asset
		v2.GET("/assets/:uuid/content", assetCtrl.DownloadOneByUUID) // dl asset
		v2.PUT("/assets/:uuid", assetCtrl.ReplaceByUUID)             // replace asset
		v2.PATCH("/assets/:uuid", assetCtrl.UpdateByUUID)            // update asset
		v2.DELETE("/assets/:uuid", assetCtrl.DeleteByUUID)           // delete asset

		v2.POST("/assets/:uuid/syntax-check", assetCtrl.CheckSyntax) // check syntax

		v2.GET("/docs/processors", docsCtrl.FindAllProcessors)
		v2.GET("/docs/processors/:code", docsCtrl.FindOneProcessorByCode)
		// v1.GET("/docs/inputs", getDocsInputs)
		// v1.GET("/docs/inputs/:name", getDocsInputsByName)
		// v1.GET("/docs/filters", getDocsFilters)
		// v1.GET("/docs/filters/:name", getDocsFiltersByName)
		// v1.GET("/docs/outputs", getDocsOutputs)
		// v1.GET("/docs/outputs/:name", getDocsOutputsByName)

		v2.GET("/env", envvariablesCtrl.Find)
		v2.POST("/env", envvariablesCtrl.Create)
		v2.GET("/env/:uuid", envvariablesCtrl.FindOneByUUID)
		v2.DELETE("/env/:uuid", envvariablesCtrl.DeleteByUUID)

		v2.GET("/db.zip", dbCtrl.Download)
	}

	apiLogger.Debugf("Serving API on /%s/ ", path)

	return r
}
