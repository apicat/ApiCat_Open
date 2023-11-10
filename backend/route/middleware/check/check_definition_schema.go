package check

import (
	"github.com/apicat/apicat/backend/model/definition"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/module/translator"
	"net/http"

	"github.com/apicat/apicat/backend/enum"
	"github.com/gin-gonic/gin"
)

type SchemaUriData struct {
	ProjectID string `uri:"project-id" binding:"required,gt=0"`
	SchemaID  uint   `uri:"schemas-id" binding:"required,gt=0"`
}

// 需要先通过CheckProject中间件。检验模型是否存在，是否所属请求对应的项目
func CheckDefinitionSchema() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentProject, _ := ctx.Get("CurrentProject")

		var data SchemaUriData

		responseCode := enum.Display404ErrorMessage
		if ctx.Request.Method == "GET" {
			responseCode = enum.Redirect404Page
		}

		if err := ctx.ShouldBindUri(&data); err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    responseCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.NotFound"}),
			})
			ctx.Abort()
			return
		}

		ds, err := definition.NewDefinitionSchemas(data.SchemaID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    responseCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.NotFound"}),
			})
			ctx.Abort()
			return
		}

		if ds.ProjectId != currentProject.(*project.Projects).ID {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    responseCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.NotFound"}),
			})
			ctx.Abort()
			return
		}

		ctx.Set("CurrentDefinitionSchema", ds)
	}
}
