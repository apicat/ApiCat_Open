package collection

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/model/global"
	"github.com/apicat/apicat/v2/backend/model/team"
	"github.com/apicat/apicat/v2/backend/module/spec"

	"github.com/lithammer/shortuuid/v4"
	"gorm.io/gorm"
)

const (
	CategoryType = "category"
	DocType      = "doc"
	HttpType     = "http"
)

type Collection struct {
	ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	PublicID     string `gorm:"type:varchar(255);index;comment:集合公开id"`
	ProjectID    string `gorm:"type:varchar(24);index;not null;comment:项目id"`
	ParentID     uint   `gorm:"type:bigint;not null;comment:父级id"`
	Path         string `gorm:"type:varchar(255);not null;comment:请求路径"`
	Method       string `gorm:"type:varchar(255);not null;comment:请求方法"`
	Title        string `gorm:"type:varchar(255);not null;comment:名称"`
	Type         string `gorm:"type:varchar(255);not null;comment:类型:category,doc,http"`
	ShareKey     string `gorm:"type:varchar(255);comment:项目分享密码"`
	Content      string `gorm:"type:mediumtext;comment:内容"`
	DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:显示顺序"`
	CreatedBy    uint   `gorm:"type:bigint;not null;default:0;comment:创建成员id"`
	UpdatedBy    uint   `gorm:"type:bigint;not null;default:0;comment:最后更新成员id"`
	DeletedBy    uint   `gorm:"type:bigint;default:null;comment:删除成员id"`
	model.TimeModel
}

func init() {
	model.RegMigrate(&Collection{})
}

func (c *Collection) Get(ctx context.Context) (bool, error) {
	tx := model.DB(ctx)
	if c.ID != 0 {
		tx = tx.Take(c, "id = ? AND project_id = ?", c.ID, c.ProjectID)
	} else if c.PublicID != "" {
		tx = tx.Take(c, "public_id = ?", c.PublicID)
	} else if c.ProjectID != "" && c.Path != "" {
		tx = tx.First(c, "project_id = ? AND path = ? AND method = ?", c.ProjectID, c.Path, c.Method)
	} else {
		return false, errors.New("query condition error")
	}
	err := model.NotRecord(tx)
	return tx.Error == nil, err
}

func (c *Collection) HasChildren(ctx context.Context) (bool, error) {
	tx := model.DB(ctx).Model(c).Where("project_id = ? AND parent_id = ?", c.ProjectID, c.ID).Take(&Collection{})
	return tx.Error == nil, model.NotRecord(tx)
}

func (c *Collection) Create(ctx context.Context, member *team.TeamMember) error {
	if c.Type == CategoryType {
		// 创建目录时，新建的目录在目标父级集合的最上方
		if err := model.DB(ctx).Model(c).Where("project_id = ? AND parent_id = ?", c.ProjectID, c.ParentID).Update("display_order", gorm.Expr("display_order + ?", 1)).Error; err != nil {
			slog.ErrorContext(ctx, "collection.Create.UpdateOrder", "err", err)
		}
		c.DisplayOrder = 1
	} else {
		// 创建文档时，新建的文档在目标父级集合的最下方
		// 获取最大的display_order
		if c.DisplayOrder == 0 {
			var maxDisplayOrder Collection
			if err := model.DB(ctx).Model(c).Where("project_id = ? AND parent_id = ?", c.ProjectID, c.ParentID).Order("display_order desc").First(&maxDisplayOrder).Error; err != nil {
				maxDisplayOrder = Collection{DisplayOrder: 0}
			}
			c.DisplayOrder = maxDisplayOrder.DisplayOrder + 1
		}

		// 获取文档的path
		url, err := GetCollectionURLNode(ctx, c.Content)
		if err != nil {
			slog.ErrorContext(ctx, "collection.Create.GetCollectionURLNode", "err", err)
		}
		c.Path = url.Attrs.Path
		c.Method = url.Attrs.Method
	}

	c.PublicID = shortuuid.New()
	c.CreatedBy = member.ID
	c.UpdatedBy = member.ID
	return model.DB(ctx).Create(c).Error
}

func (c *Collection) Update(ctx context.Context, title, content string, memberID uint) error {
	if c.Type != CategoryType {
		h := &CollectionHistory{
			CollectionID: c.ID,
			Title:        c.Title,
			Content:      c.Content,
		}
		h.Create(ctx, memberID)
	}

	// 获取文档的path
	url, err := GetCollectionURLNode(ctx, content)
	if err != nil {
		slog.ErrorContext(ctx, "collection.Update.GetCollectionURLNode", "err", err)
	}

	return model.DB(ctx).Model(c).Updates(map[string]interface{}{
		"path":       url.Attrs.Path,
		"method":     url.Attrs.Method,
		"title":      title,
		"content":    content,
		"updated_by": memberID,
	}).Error
}

// UpdateShareKey 更新项目分享密码
func (c *Collection) UpdateShareKey(ctx context.Context) error {
	if c.ID == 0 {
		return nil
	}
	return model.DB(ctx).Model(c).Update("share_key", c.ShareKey).Error
}

func (c *Collection) Sort(ctx context.Context, parentID uint, displayOrder int) error {
	return model.DB(ctx).Model(c).UpdateColumns(map[string]interface{}{
		"parent_id":     parentID,
		"display_order": displayOrder,
	}).Error
}

func (c *Collection) ToSpec() (*spec.Collection, error) {
	sc := &spec.Collection{
		ID:       c.ID,
		ParentID: c.ParentID,
		Title:    c.Title,
		Type:     spec.CollectionType(c.Type),
	}

	if c.Content != "" {
		if err := json.Unmarshal([]byte(c.Content), &sc.Content); err != nil {
			return nil, err
		}
	}

	return sc, nil
}

// DelRefSchema 删除公共模型引用
// deref: 是否解引用，true: 展开引用自身(collectiuon.$ref to schema detail)，false: 清除引用自身(delete $ref)
func (c *Collection) DelRefSchema(ctx context.Context, refSchema *definition.DefinitionSchema, deref bool) error {
	collectionSpec, err := c.ToSpec()
	if err != nil {
		return err
	}

	refSchemaSpec, err := refSchema.ToSpec()
	if err != nil {
		return err
	}

	if deref {
		if err := collectionSpec.DerefSchema(refSchemaSpec); err != nil {
			return err
		}
	} else {
		if err := collectionSpec.DelRefSchema(refSchemaSpec); err != nil {
			return err
		}
	}

	content, err := json.Marshal(collectionSpec.Content)
	if err != nil {
		return err
	}

	return model.DB(ctx).Model(c).Select("content").UpdateColumn("content", string(content)).Error
}

// DelRefResponse 删除公共响应引用
// deref: 是否解引用，true: 展开引用自身(collectiuon.$ref to response detail)，false: 清除引用自身(delete $ref)
func (c *Collection) DelRefResponse(ctx context.Context, refResponse *definition.DefinitionResponse, deref bool) error {
	collectionSpec, err := c.ToSpec()
	if err != nil {
		return err
	}

	refResponseSpec, err := refResponse.ToSpec()
	if err != nil {
		return err
	}

	if deref {
		if err := collectionSpec.DerefResponse(refResponseSpec); err != nil {
			return err
		}
	} else {
		if err := collectionSpec.DelResponseByRefId(int64(refResponseSpec.ID)); err != nil {
			return err
		}
	}

	content, err := json.Marshal(collectionSpec.Content)
	if err != nil {
		return err
	}

	return model.DB(ctx).Model(c).Select("content").UpdateColumn("content", string(content)).Error
}

// DelExceptParam 删除全局参数排除关系
// unpack: 是否展开，true: 将globalParam详情添加到parameters中，false: 在glabalExcept中删除globalParamID
func (c *Collection) DelExceptParam(ctx context.Context, exceptParam *global.GlobalParameter, unpack bool) error {
	collectionSpec, err := c.ToSpec()
	if err != nil {
		return err
	}

	exceptParamSpec, err := exceptParam.ToSpec()
	if err != nil {
		return err
	}

	if unpack {
		if err := collectionSpec.AddParameter(exceptParam.In, exceptParamSpec); err != nil {
			return err
		}
	} else {
		if err := collectionSpec.DelGlobalExceptID(exceptParam.In, int64(exceptParam.ID)); err != nil {
			return err
		}
	}

	content, err := json.Marshal(collectionSpec.Content)
	if err != nil {
		return err
	}

	return model.DB(ctx).Model(c).Select("content").UpdateColumn("content", string(content)).Error
}
