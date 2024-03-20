package response

import (
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	"github.com/apicat/apicat/v2/backend/route/proto/collection/base"
)

type CollectionHistory struct {
	protobase.IdCreateTimeInfo
	CollectionHistoryData
	CreatedBy string `json:"createdBy"`
}

type CollectionHistoryData struct {
	base.CollectionIDOption
	base.CollectionTypeOption
	base.CollectionData
}

type CollectionHistoryList []*CollectionHistoryItem

type CollectionHistoryItem struct {
	protobase.IdCreateTimeInfo
	base.CollectionTypeOption
	CreatedBy string `json:"createdBy"`
}

type DiffCollectionHistories struct {
	Doc1 *CollectionHistory `json:"doc1"`
	Doc2 *CollectionHistory `json:"doc2"`
}
