package collection

import (
	"encoding/json"
	"fmt"
	"github.com/apicat/apicat/backend/i18n"
	"github.com/apicat/apicat/backend/model"
	"github.com/apicat/apicat/backend/model/collection"
	"github.com/apicat/apicat/backend/model/iteration"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/module/exportresp"
	"github.com/apicat/apicat/backend/module/spec"
	"github.com/apicat/apicat/backend/module/spec/plugin/export"
	"github.com/apicat/apicat/backend/module/spec/plugin/openapi"
	"github.com/apicat/apicat/backend/route/proto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

func CollectionsList(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	p, _ := currentProject.(*project.Projects)

	var data proto.CollectionsListData
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c, _ := collection.NewCollections()
	c.ProjectId = p.ID
	collections, err := c.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.QueryFailed"}),
		})
	}

	if data.IterationID == "" {
		ctx.JSON(http.StatusOK, buildProjectTree(0, collections))
	} else {
		i, err := iteration.NewIterations(data.IterationID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.QueryFailed"}),
			})
			return
		}

		ia, _ := iteration.NewIterationApis()
		cIDs, err := ia.GetCollectionIDByIterationID(i.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.QueryFailed"}),
			})
			return
		}

		ctx.JSON(http.StatusOK, buildIterationTree(0, collections, cIDs))
	}
}

func buildProjectTree(parentID uint, collections []*collection.Collections) []*proto.CollectionList {
	return buildTree(parentID, collections, false)
}

func buildIterationTree(parentID uint, collections []*collection.Collections, selectCIDs []uint) []*proto.CollectionList {
	return buildTree(parentID, collections, true, selectCIDs...)
}

func buildTree(parentID uint, collections []*collection.Collections, isIteration bool, selectCIDs ...uint) []*proto.CollectionList {
	result := make([]*proto.CollectionList, 0)

	for _, c := range collections {
		if c.ParentId == parentID {
			children := buildTree(c.ID, collections, isIteration, selectCIDs...)

			cl := proto.CollectionList{
				ID:       c.ID,
				ParentID: c.ParentId,
				Title:    c.Title,
				Type:     c.Type,
				Items:    children,
			}

			isSelected := false
			if isIteration {
				for _, cid := range selectCIDs {
					if cid == cl.ID {
						isSelected = true
						break
					}
					if !isSelected {
						for _, v := range cl.Items {
							if *v.Selected {
								isSelected = true
								break
							}
						}
					}
				}
				cl.Selected = &isSelected
			}

			result = append(result, &cl)
		}
	}

	return result
}

func CollectionsGet(ctx *gin.Context) {
	currentCollection, _ := ctx.Get("CurrentCollection")
	c := currentCollection.(*collection.Collections)

	ctx.JSON(http.StatusOK, gin.H{
		"id":         c.ID,
		"parent_id":  c.ParentId,
		"title":      c.Title,
		"type":       c.Type,
		"content":    c.Content,
		"created_at": c.CreatedAt.Format("2006-01-02 15:04:05"),
		"created_by": c.Creator(),
		"updated_at": c.UpdatedAt.Format("2006-01-02 15:04:05"),
		"updated_by": c.Updater(),
	})
}

func CollectionsCreate(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := proto.CollectionCreate{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	currentProject, _ := ctx.Get("CurrentProject")
	p, _ := currentProject.(*project.Projects)

	if data.IterationID != "" {
		_, err := iteration.NewIterations(data.IterationID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.CreateFailed"}),
			})
			return
		}
	}

	c, _ := collection.NewCollections()
	c.ProjectId = p.ID
	c.ParentId = data.ParentID
	c.Title = data.Title
	c.Type = data.Type
	c.Content = data.Content
	c.CreatedBy = currentProjectMember.(*project.ProjectMembers).UserID
	c.UpdatedBy = currentProjectMember.(*project.ProjectMembers).UserID

	var err error
	if c.Type == "category" {
		err = c.CreateCategory()
	} else {
		err = c.CreateDoc()
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.CreateFailed"}),
		})
		return
	}

	if data.IterationID != "" {
		i, err := iteration.NewIterations(data.IterationID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.CreateFailed"}),
			})
			return
		}

		ia, _ := iteration.NewIterationApis()
		ia.IterationID = i.ID
		ia.CollectionID = c.ID
		ia.CollectionType = c.Type
		if err := ia.Create(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.CreateFailed"}),
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         c.ID,
		"parent_id":  c.ParentId,
		"title":      c.Title,
		"type":       c.Type,
		"content":    c.Content,
		"created_at": c.CreatedAt.Format("2006-01-02 15:04:05"),
		"created_by": c.Creator(),
		"updated_at": c.UpdatedAt.Format("2006-01-02 15:04:05"),
		"updated_by": c.Updater(),
	})
}

func CollectionsUpdate(ctx *gin.Context) {
	currentCollection, _ := ctx.Get("CurrentCollection")
	c := currentCollection.(*collection.Collections)

	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := proto.CollectionUpdate{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ch, _ := collection.NewCollectionHistories()
	ch.CollectionId = c.ID
	ch.Title = c.Title
	ch.Type = c.Type
	ch.Content = c.Content
	ch.CreatedBy = currentProjectMember.(*project.ProjectMembers).UserID

	// 不是同一个人编辑的文档或5分钟后编辑文档内容，保存历史记录
	if c.UpdatedBy != currentProjectMember.(*project.ProjectMembers).UserID || c.UpdatedAt.Add(5*time.Minute).Before(time.Now()) {
		if err := ch.Create(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.UpdateFailed"}),
			})
			return
		}
	}

	c.Title = data.Title
	c.Content = data.Content
	c.UpdatedBy = currentProjectMember.(*project.ProjectMembers).UserID
	if err := c.Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.UpdateFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func CollectionsCopy(ctx *gin.Context) {
	currentCollection, _ := ctx.Get("CurrentCollection")
	c := currentCollection.(*collection.Collections)

	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := proto.CollectionCopyData{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data.IterationID != "" {
		_, err := iteration.NewIterations(data.IterationID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.CreateFailed"}),
			})
			return
		}
	}

	newCollection := collection.Collections{
		ProjectId:    c.ProjectId,
		ParentId:     c.ParentId,
		Title:        fmt.Sprintf("%s (copy)", c.Title),
		Type:         c.Type,
		Content:      c.Content,
		DisplayOrder: c.DisplayOrder,
		CreatedBy:    currentProjectMember.(*project.ProjectMembers).UserID,
	}

	if err := newCollection.CreateDoc(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.CreateFailed"}),
		})
		return
	}

	if data.IterationID != "" {
		i, err := iteration.NewIterations(data.IterationID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.CreateFailed"}),
			})
			return
		}

		ia, _ := iteration.NewIterationApis()
		ia.IterationID = i.ID
		ia.CollectionID = newCollection.ID
		ia.CollectionType = newCollection.Type
		if err := ia.Create(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.CreateFailed"}),
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         newCollection.ID,
		"parent_id":  newCollection.ParentId,
		"title":      newCollection.Title,
		"type":       newCollection.Type,
		"content":    newCollection.Content,
		"created_at": newCollection.CreatedAt.Format("2006-01-02 15:04:05"),
		"created_by": newCollection.Creator(),
		"updated_at": newCollection.UpdatedAt.Format("2006-01-02 15:04:05"),
		"updated_by": newCollection.Updater(),
	})
}

func CollectionsMovement(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := proto.CollectionMovement{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	for i, id := range data.Target.Ids {
		if c, err := collection.NewCollections(id); err == nil {
			c.ParentId = data.Target.Pid
			c.DisplayOrder = i
			c.Update()
		}
	}

	if data.Target.Pid != data.Origin.Pid {
		for i, id := range data.Origin.Ids {
			if c, err := collection.NewCollections(id); err == nil {
				c.ParentId = data.Origin.Pid
				c.DisplayOrder = i
				c.Update()
			}
		}
	}

	ctx.Status(http.StatusCreated)
}

func CollectionsDelete(ctx *gin.Context) {
	currentCollection, _ := ctx.Get("CurrentCollection")
	c := currentCollection.(*collection.Collections)
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := proto.CollectionDeleteData{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data.IterationID != "" {
		_, err := iteration.NewIterations(data.IterationID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.DeleteFailed"}),
			})
			return
		}

		collections, err := c.GetSubCollectionsContainsSelf()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.DeleteFailed"}),
			})
			return
		}
		var cIDs []uint
		for _, v := range collections {
			cIDs = append(cIDs, v.ID)
		}

		if err := iteration.DeleteIterationApisByCollectionID(cIDs...); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.DeleteFailed"}),
			})
			return
		}
	}

	if err := collection.Deletes(c.ID, model.Conn, currentProjectMember.(*project.ProjectMembers).UserID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.DeleteFailed"}),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func CollectionDataGet(ctx *gin.Context) {
	uriData := proto.CollectionDataGetData{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	data := proto.ExportCollection{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	p, err := project.NewProjects(uriData.ProjectID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Projects.NotFound"}),
		})
		return
	}
	c, err := collection.NewCollections(uriData.CollectionID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.NotFound"}),
		})
		return
	}

	apicatData := collection.CollectionExport(p, c)
	if apicatDataContent, err := json.Marshal(apicatData); err == nil {
		slog.InfoCtx(ctx, "Export", slog.String("apicat", string(apicatDataContent)))
	}

	var content []byte
	switch data.Type {
	case "swagger":
		content, err = openapi.Encode(apicatData, "2.0")
	case "openapi3.0.0":
		content, err = openapi.Encode(apicatData, "3.0.0")
	case "openapi3.0.1":
		content, err = openapi.Encode(apicatData, "3.0.1")
	case "openapi3.0.2":
		content, err = openapi.Encode(apicatData, "3.0.2")
	case "openapi3.1.0":
		content, err = openapi.Encode(apicatData, "3.1.0")
	case "HTML":
		content, err = export.HTML(apicatData)
	case "md":
		content, err = export.Markdown(apicatData)
	default:
		content, err = apicatData.ToJSON(spec.JSONOption{Indent: "  "})
	}

	slog.InfoCtx(ctx, "Export", slog.String(data.Type, string(content)))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Collections.ExportFail"}),
		})
		return
	}

	exportresp.ExportResponse(data.Type, data.Download, p.Title+"-"+data.Type, content, ctx)
}
