package check

import (
	"github.com/apicat/apicat/backend/i18n"
	"github.com/apicat/apicat/backend/model/collection"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/route/proto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CollectionUriData struct {
	ProjectID    string `uri:"project-id" binding:"required,gt=0"`
	CollectionID uint   `uri:"collection-id" binding:"required,gt=0"`
}

// 需要先通过CheckProject中间件。检验集合是否存在，是否所属请求对应的项目
func CheckCollection() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentProject, _ := ctx.Get("CurrentProject")

		var data CollectionUriData

		responseCode := proto.Display404ErrorMessage
		if ctx.Request.Method == "GET" {
			responseCode = proto.Redirect404Page
		}

		if err := ctx.ShouldBindUri(&data); err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    responseCode,
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.NotFound"}),
			})
			ctx.Abort()
			return
		}

		c, err := collection.NewCollections(data.CollectionID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    responseCode,
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.NotFound"}),
			})
			ctx.Abort()
			return
		}

		if c.ProjectId != currentProject.(*project.Projects).ID {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    responseCode,
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.NotFound"}),
			})
			ctx.Abort()
			return
		}

		ctx.Set("CurrentCollection", c)
	}
}
