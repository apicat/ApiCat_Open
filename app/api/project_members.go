package api

import (
	"math"
	"net/http"
	"strings"

	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/enum"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

type ProjectMembersListData struct {
	Page     int `form:"page" binding:"omitempty,gte=1"`
	PageSize int `form:"page_size" binding:"omitempty,gte=1,lte=100"`
}

type GetPathUserID struct {
	UserID uint `uri:"user-id" binding:"required"`
}

type CreateProjectMemberData struct {
	UserIDs   []uint `json:"user_ids" binding:"required,gt=0,dive,required"`
	Authority string `json:"authority" binding:"required,oneof=manage write read"`
}

type UpdateProjectMemberAuthData struct {
	Authority string `json:"authority" binding:"required,oneof=manage write read"`
}

type ProjectMemberData struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Authority string `json:"authority"`
	CreatedAt string `json:"created_at"`
}

// MembersList handles GET requests to retrieve a list of members in the current project.
func ProjectMembersList(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")

	var data GetMembersData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data.Page <= 0 {
		data.Page = 1
	}
	if data.PageSize <= 0 {
		data.PageSize = 15
	}

	member, _ := models.NewProjectMembers()
	member.ProjectID = currentProject.(*models.Projects).ID
	totalMember, err := member.Count()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	member.ProjectID = currentProject.(*models.Projects).ID
	members, err := member.List(data.Page, data.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	userIDs := []uint{}
	for _, v := range members {
		userIDs = append(userIDs, v.UserID)
	}

	users, err := models.UserListByIDs(userIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	userIDToNameMap := map[uint]*models.Users{}
	for _, v := range users {
		userIDToNameMap[v.ID] = v
	}

	membersList := []any{}
	for _, v := range members {
		email := userIDToNameMap[v.UserID].Email
		parts := strings.Split(email, "@")
		membersList = append(membersList, map[string]any{
			"id":         v.ID,
			"user_id":    v.UserID,
			"username":   userIDToNameMap[v.UserID].Username,
			"authority":  v.Authority,
			"is_enabled": userIDToNameMap[v.UserID].IsEnabled,
			"email":      parts[0][0:1] + "***" + parts[0][len(parts[0])-1:] + "@" + parts[len(parts)-1],
			"created_at": v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"current_page": data.Page,
		"total_page":   int(math.Ceil(float64(totalMember) / float64(data.PageSize))),
		"total":        totalMember,
		"records":      membersList,
	})
}

// (Abandoned) MemberGet retrieves the project member data for a given user and project ID.
func MemberGetByUserID(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")

	data := GetPathUserID{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := models.NewUsers(data.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	pm, _ := models.NewProjectMembers()
	pm.UserID = user.ID
	pm.ProjectID = currentProject.(*models.Projects).ID
	if err := pm.GetByUserIDAndProjectID(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusOK, ProjectMemberData{
		ID:        pm.ID,
		UserID:    user.ID,
		Username:  user.Username,
		Authority: pm.Authority,
		CreatedAt: pm.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}

// ProjectMembersCreate projects the creation of a new member.
func ProjectMembersCreate(ctx *gin.Context) {
	// 项目管理员才添加成员
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberIsManage() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := CreateProjectMemberData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	result := []gin.H{}
	for _, v := range data.UserIDs {
		user, err := models.NewUsers(v)
		if err != nil {
			continue
		}

		pm, _ := models.NewProjectMembers()
		pm.UserID = user.ID
		pm.ProjectID = currentProjectMember.(*models.ProjectMembers).ProjectID
		if err := pm.GetByUserIDAndProjectID(); err == nil {
			continue
		}

		pm.Authority = data.Authority
		if err := pm.Create(); err != nil {
			continue
		}

		result = append(result, gin.H{
			"id":         pm.ID,
			"user_id":    user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"is_enabled": user.IsEnabled,
			"authority":  pm.Authority,
			"created_at": pm.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	ctx.JSON(http.StatusOK, result)
}

// DeleteMember deletes a project member by checking if the given member exists in the project.
func ProjectMembersDelete(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberIsManage() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := GetPathUserID{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data.UserID == currentProjectMember.(*models.ProjectMembers).UserID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.DeleteFailed"}),
		})
		return
	}

	pm, _ := models.NewProjectMembers()
	pm.UserID = data.UserID
	pm.ProjectID = currentProjectMember.(*models.ProjectMembers).ProjectID
	if err := pm.GetByUserIDAndProjectID(); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.NotFound"}),
		})
		return
	}

	if err := pm.Delete(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.DeleteFailed"}),
		})
	}

	ctx.Status(http.StatusNoContent)
}

// UpdateMember updates the authority of a project member in the database.
func ProjectMembersAuthUpdate(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberIsManage() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := GetPathUserID{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	pm, _ := models.NewProjectMembers()
	pm.UserID = data.UserID
	pm.ProjectID = currentProjectMember.(*models.ProjectMembers).ProjectID
	if err := pm.GetByUserIDAndProjectID(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.NotFound"}),
		})
		return
	}

	bodyData := UpdateProjectMemberAuthData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&bodyData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	pm.Authority = bodyData.Authority
	if err := pm.Update(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.UpdateFailed"}),
		})
	}

	ctx.Status(http.StatusCreated)
}

func ProjectMembersWithout(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*models.ProjectMembers).MemberIsManage() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	user, _ := models.NewUsers()
	users, err := user.List(0, 0)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	projectMember, _ := models.NewProjectMembers()
	projectMember.ProjectID = currentProjectMember.(*models.ProjectMembers).ProjectID
	projectMembers, err := projectMember.List(0, 0)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	pmMap := map[uint]models.ProjectMembers{}
	for _, v := range projectMembers {
		pmMap[v.UserID] = v
	}

	result := []map[string]any{}
	for _, u := range users {
		if _, ok := pmMap[u.ID]; !ok {
			result = append(result, map[string]any{
				"user_id":  u.ID,
				"username": u.Username,
				"email":    u.Email,
			})
		}
	}

	ctx.JSON(http.StatusOK, result)
}
