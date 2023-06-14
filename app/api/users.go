package api

import (
	"net/http"

	"github.com/apicat/apicat/common/auth"
	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

type SetUserInfoData struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Username string `json:"username" binding:"required,lte=255"`
}

type ChangePasswordData struct {
	Password           string `json:"password" binding:"required,gte=6,lte=255"`
	NewPassword        string `json:"new_password" binding:"required,gte=6,lte=255"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required,gte=6,lte=255,eqfield=NewPassword"`
}

func GetUserInfo(ctx *gin.Context) {
	CurrentUser, _ := ctx.Get("CurrentUser")
	user, _ := CurrentUser.(*models.Users)

	ctx.JSON(200, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"role":       user.Role,
		"is_enabled": user.IsEnabled,
		"created_at": user.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at": user.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

func SetUserInfo(ctx *gin.Context) {
	CurrentUser, _ := ctx.Get("CurrentUser")
	currentUser, _ := CurrentUser.(*models.Users)

	var data SetUserInfoData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, _ := models.NewUsers()
	if err := user.GetByEmail(data.Email); err == nil && user.ID != currentUser.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.MailboxAlreadyExists"}),
		})
		return
	}

	currentUser.Email = data.Email
	currentUser.Username = data.Username
	if err := currentUser.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.UpdateFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         currentUser.ID,
		"username":   currentUser.Username,
		"email":      currentUser.Email,
		"role":       currentUser.Role,
		"is_enabled": currentUser.IsEnabled,
		"created_at": currentUser.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at": currentUser.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

func ChangePassword(ctx *gin.Context) {
	CurrentUser, _ := ctx.Get("CurrentUser")
	currentUser, _ := CurrentUser.(*models.Users)

	var data ChangePasswordData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if !auth.CheckPasswordHash(data.Password, currentUser.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.WrongPassword"}),
		})
		return
	}

	hashedPassword, err := auth.HashPassword(data.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.UpdateFailed"}),
		})
		return
	}

	currentUser.Password = hashedPassword
	if err := currentUser.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.UpdateFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}
