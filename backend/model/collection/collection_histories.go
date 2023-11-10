package collection

import (
	"github.com/apicat/apicat/backend/model"
	"time"
)

type CollectionHistories struct {
	ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	CollectionId uint   `gorm:"type:bigint;index;not null;comment:集合id"`
	Title        string `gorm:"type:varchar(255);not null;comment:名称"`
	Type         string `gorm:"type:varchar(255);not null;comment:类型:category,doc,http"`
	Content      string `gorm:"type:mediumtext;comment:内容"`
	CreatedAt    time.Time
	CreatedBy    uint `gorm:"type:bigint;not null;default:0;comment:创建人id"`
}

func init() {
	model.RegMigrate(&CollectionHistories{})
}

func NewCollectionHistories(ids ...uint) (*CollectionHistories, error) {
	if len(ids) > 0 {
		ch := &CollectionHistories{ID: ids[0]}
		if err := model.Conn.Take(ch).Error; err != nil {
			return ch, err
		}
		return ch, nil
	}
	return &CollectionHistories{}, nil
}

func (ch *CollectionHistories) List(collectionIDs ...uint) ([]*CollectionHistories, error) {
	var collectionHistories []*CollectionHistories
	if len(collectionIDs) > 0 {
		return collectionHistories, model.Conn.Where("collection_id IN ?", collectionIDs).Order("created_at desc").Find(&collectionHistories).Error
	}
	return collectionHistories, model.Conn.Order("created_at desc").Find(&collectionHistories).Error
}

func (ch *CollectionHistories) Create() error {
	return model.Conn.Create(ch).Error
}

func (ch *CollectionHistories) Restore(collection *Collections, uid uint) error {
	nch, _ := NewCollectionHistories()
	nch.CollectionId = collection.ID
	nch.Title = collection.Title
	nch.Type = collection.Type
	nch.Content = collection.Content
	nch.CreatedBy = uid
	if err := nch.Create(); err != nil {
		return err
	}

	collection.Title = ch.Title
	collection.Content = ch.Content
	collection.UpdatedBy = uid
	if err := collection.Update(); err != nil {
		return err
	}

	return nil
}
