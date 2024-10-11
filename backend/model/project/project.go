package project

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/team"

	"github.com/pkg-id/objectid"
	"gorm.io/gorm"
)

const (
	VisibilityPrivate = "private"
	VisibilityPublic  = "public"
)

type Project struct {
	ID             string `gorm:"type:varchar(24);primarykey"`
	TeamID         string `gorm:"type:varchar(24);index;comment:team id"`
	MemberID       uint   `gorm:"type:bigint;comment:team member ID of the project manager"`
	Title          string `gorm:"type:varchar(255);not null;comment:project title"`
	Visibility     string `gorm:"type:varchar(32);not null;comment:project visibility:0-private,1-public"`
	ShareKey       string `gorm:"type:varchar(255);comment:project share key"`
	Description    string `gorm:"type:varchar(1024);comment:project description"`
	Cover          string `gorm:"type:varchar(255);comment:project cover"`
	EmbeddingModel string `gorm:"type:varchar(255);comment:project embedding model"`
	model.TimeModel
}

// Get 获取项目
func (p *Project) Get(ctx context.Context) (bool, error) {
	tx := model.DB(ctx).Take(p, "id = ?", p.ID)
	err := model.NotRecord(tx)
	return tx.Error == nil, err
}

// Create 创建项目
func (p *Project) Create(ctx context.Context, member *team.TeamMember, groupID uint) (*Project, error) {
	p.ID = objectid.New().String()
	p.MemberID = member.ID
	if err := model.DB(ctx).Transaction(
		func(tx *gorm.DB) error {
			ret := tx.Create(p)
			if ret.Error != nil {
				return ret.Error
			}

			// 添加默认管理员
			pm := &ProjectMember{
				ProjectID:  p.ID,
				MemberID:   member.ID,
				Permission: ProjectMemberManage,
				GroupID:    groupID,
			}
			_, err := pm.Create(ctx, tx)
			return err
		},
	); err != nil {
		return nil, err
	}
	return p, nil
}

// Update 更新项目
func (p *Project) Update(ctx context.Context) error {
	if p.ID == "" {
		return nil
	}
	// 只能更新Title、Visibility、Description和Cover
	return model.DB(ctx).Model(p).Updates(map[string]interface{}{
		"title":       p.Title,
		"visibility":  p.Visibility,
		"description": p.Description,
		"cover":       p.Cover,
	}).Error
}

// UpdateShareKey 更新项目分享密码
func (p *Project) UpdateShareKey(ctx context.Context) error {
	if p.ID == "" {
		return nil
	}
	return model.DB(ctx).Model(p).Update("share_key", p.ShareKey).Error
}

// Delete 删除项目
func (p *Project) Delete(ctx context.Context) error {
	if p.ID == "" {
		return nil
	}
	return model.DB(ctx).Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Where("project_id = ?", p.ID).Delete(&ProjectMember{}).Error; err != nil {
				return err
			}
			return tx.Delete(p).Error
		},
	)
}
