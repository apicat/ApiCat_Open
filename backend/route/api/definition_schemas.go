package api

import (
	"encoding/json"
	"fmt"
	"github.com/apicat/apicat/backend/model/collection"
	"github.com/apicat/apicat/backend/model/definition"
	"github.com/apicat/apicat/backend/model/project"
	"net/http"
	"strconv"
	"time"

	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/enum"
	"github.com/gin-gonic/gin"
)

type DefinitionSchemaCreate struct {
	ParentId    uint                   `json:"parent_id" binding:"gte=0"`
	Name        string                 `json:"name" binding:"required,lte=255"`
	Description string                 `json:"description" binding:"lte=255"`
	Type        string                 `json:"type" binding:"required,oneof=category schema"`
	Schema      map[string]interface{} `json:"schema"`
}

type DefinitionSchemaUpdate struct {
	Name        string                 `json:"name" binding:"required,lte=255"`
	Description string                 `json:"description" binding:"lte=255"`
	Schema      map[string]interface{} `json:"schema"`
}

type DefinitionSchemaSearch struct {
	ParentId uint   `form:"parent_id" binding:"gte=0"`
	Name     string `form:"name" binding:"lte=255"`
	Type     string `form:"type" binding:"omitempty,oneof=category schema"`
}

type DefinitionSchemaID struct {
	ID uint `uri:"schemas-id" binding:"required,gte=0"`
}

type DefinitionSchemaMove struct {
	Target OrderContent `json:"target" binding:"required"`
	Origin OrderContent `json:"origin" binding:"required"`
}

type OrderContent struct {
	Pid uint   `json:"pid" binding:"gte=0"`
	Ids []uint `json:"ids" binding:"required,dive,gte=0"`
}

func DefinitionSchemasList(ctx *gin.Context) {
	var data DefinitionSchemaSearch

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	p, _ := ctx.Get("CurrentProject")

	d, _ := definition.NewDefinitionSchemas()
	d.ProjectId = p.(*project.Projects).ID
	d.ParentId = data.ParentId
	d.Name = data.Name
	d.Type = data.Type

	definitions, err := d.List()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    enum.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.NotFound"}),
		})
		return
	}

	result := make([]gin.H, 0)
	for _, d := range definitions {
		schema := make(map[string]interface{})
		if err := json.Unmarshal([]byte(d.Schema), &schema); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
			})
			return
		}

		result = append(result, gin.H{
			"id":          d.ID,
			"parent_id":   d.ParentId,
			"name":        d.Name,
			"description": d.Description,
			"type":        d.Type,
			"schema":      schema,
		})
	}
	ctx.JSON(http.StatusOK, result)
}

func DefinitionSchemasCreate(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var data DefinitionSchemaCreate

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	schemaJson, err := json.Marshal(data.Schema)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
		})
		return
	}

	p, _ := ctx.Get("CurrentProject")
	d, _ := definition.NewDefinitionSchemas()
	d.ProjectId = p.(*project.Projects).ID
	d.Name = data.Name
	definitions, err := d.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.QueryFailed"}),
		})
		return
	}
	if len(definitions) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.NameExists"}),
		})
		return
	}

	d.Description = data.Description
	d.Type = data.Type
	d.Schema = string(schemaJson)
	d.CreatedBy = currentProjectMember.(*project.ProjectMembers).UserID
	d.UpdatedBy = currentProjectMember.(*project.ProjectMembers).UserID
	if err := d.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.CreateFail"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":          d.ID,
		"parent_id":   d.ParentId,
		"name":        d.Name,
		"description": d.Description,
		"type":        d.Type,
		"schema":      data.Schema,
		"created_at":  d.CreatedAt.Format("2006-01-02 15:04:05"),
		"created_by":  d.Creator(),
		"updated_at":  d.UpdatedAt.Format("2006-01-02 15:04:05"),
		"updated_by":  d.Updater(),
	})
}

func DefinitionSchemasUpdate(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	currentDefinitionSchema, _ := ctx.Get("CurrentDefinitionSchema")
	d := currentDefinitionSchema.(*definition.DefinitionSchemas)

	var (
		data DefinitionSchemaUpdate
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	dsh, _ := definition.NewDefinitionSchemaHistories()
	dsh.SchemaID = d.ID
	dsh.Name = d.Name
	dsh.Description = d.Description
	dsh.Type = d.Type
	dsh.Schema = d.Schema
	dsh.CreatedBy = currentProjectMember.(*project.ProjectMembers).UserID
	if d.UpdatedBy != currentProjectMember.(*project.ProjectMembers).UserID || d.UpdatedAt.Add(5*time.Minute).Before(time.Now()) {
		if err := dsh.Create(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.UpdateFail"}),
			})
			return
		}
	}

	schemaJson, err := json.Marshal(data.Schema)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
		})
		return
	}

	d.Name = data.Name
	d.Description = data.Description
	definitions, err := d.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.QueryFailed"}),
		})
		return
	}

	if len(definitions) > 0 && definitions[0].ID != d.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.NameExists"}),
		})
		return
	}

	d.Schema = string(schemaJson)
	d.UpdatedBy = currentProjectMember.(*project.ProjectMembers).UserID
	if err := d.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.UpdateFail"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func DefinitionSchemasDelete(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	currentDefinitionSchema, _ := ctx.Get("CurrentDefinitionSchema")
	d := currentDefinitionSchema.(*definition.DefinitionSchemas)

	// 模型解引用
	isUnRefData := IsUnRefData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&isUnRefData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := collection.DefinitionsSchemaUnRefByCollections(d, isUnRefData.IsUnRef); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err := definition.DefinitionsSchemaUnRefByDefinitionsResponse(d, isUnRefData.IsUnRef); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err := definition.DefinitionsSchemaUnRefByDefinitionsSchema(d, isUnRefData.IsUnRef); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := d.Delete(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.DeleteFail"}),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func DefinitionSchemasGet(ctx *gin.Context) {
	currentDefinitionSchema, _ := ctx.Get("CurrentDefinitionSchema")
	d := currentDefinitionSchema.(*definition.DefinitionSchemas)
	var data DefinitionSchemaID

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	schema := make(map[string]interface{})
	if err := json.Unmarshal([]byte(currentDefinitionSchema.(*definition.DefinitionSchemas).Schema), &schema); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.ContentParsingFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":          d.ID,
		"parent_id":   d.ParentId,
		"name":        d.Name,
		"description": d.Description,
		"type":        d.Type,
		"schema":      schema,
		"created_at":  d.CreatedAt.Format("2006-01-02 15:04:05"),
		"created_by":  d.Creator(),
		"updated_at":  d.UpdatedAt.Format("2006-01-02 15:04:05"),
		"updated_by":  d.Updater(),
	})
}

func DefinitionSchemasCopy(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	currentDefinitionSchema, _ := ctx.Get("CurrentDefinitionSchema")
	oldDefinition := currentDefinitionSchema.(*definition.DefinitionSchemas)

	schema := map[string]interface{}{}
	if err := json.Unmarshal([]byte(oldDefinition.Schema), &schema); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.CopyFail"}),
		})
		return
	}

	newDefinition, _ := definition.NewDefinitionSchemas()
	newDefinition.ProjectId = oldDefinition.ProjectId
	newDefinition.Name = oldDefinition.Name
	newDefinition.Description = oldDefinition.Description
	newDefinition.Type = oldDefinition.Type
	newDefinition.Schema = oldDefinition.Schema
	if err := newDefinition.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.CopyFail"}),
		})
		return
	}

	newDefinition.Name = fmt.Sprintf("%s_%s", newDefinition.Name, strconv.Itoa(int(newDefinition.ID)))
	if err := newDefinition.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "DefinitionSchemas.CopyFail"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":          newDefinition.ID,
		"parent_id":   newDefinition.ParentId,
		"name":        newDefinition.Name,
		"description": newDefinition.Description,
		"type":        newDefinition.Type,
		"schema":      schema,
		"created_at":  newDefinition.CreatedAt.Format("2006-01-02 15:04:05"),
		"created_by":  newDefinition.Creator(),
		"updated_at":  newDefinition.UpdatedAt.Format("2006-01-02 15:04:05"),
		"updated_by":  newDefinition.Updater(),
	})
}

func DefinitionSchemasMove(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var data DefinitionSchemaMove

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	for i, id := range data.Target.Ids {
		if d, err := definition.NewDefinitionSchemas(id); err == nil {
			d.ParentId = data.Target.Pid
			d.DisplayOrder = i
			d.Save()
		}
	}

	if data.Target.Pid != data.Origin.Pid {
		for i, id := range data.Origin.Ids {
			if d, err := definition.NewDefinitionSchemas(id); err == nil {
				d.ParentId = data.Origin.Pid
				d.DisplayOrder = i
				d.Save()
			}
		}
	}

	ctx.Status(http.StatusCreated)
}
