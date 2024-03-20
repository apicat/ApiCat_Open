package request

import (
	protobase "github.com/apicat/apicat/backend/route/proto/base"
)

type ProjectShareSwitchOption struct {
	protobase.ProjectIdOption
	Status bool `json:"status" binding:"boolean"`
}

type CheckProjectShareSecretKeyOption struct {
	protobase.ProjectIdOption
	protobase.SecretKeyOption
}
