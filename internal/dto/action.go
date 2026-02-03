package dto

import "github.com/hcd233/aris-mem-api/internal/common/enum"

// ActionReq action request
//
//	author centonhuang
//	update 2026-01-30 21:00:00
type ActionReq struct {
	Body *ActionReqBody `json:"body" doc:"Request body containing entity info"`
}

// ActionReqBody action request body
//
//	author centonhuang
//	update 2026-02-03 22:30:00
type ActionReqBody struct {
	EntityType enum.ActionEntityType `json:"entityType" enum:"article,comment" required:"true" doc:"Entity type to interact with"`
	EntityID   uint                  `json:"entityID" required:"true" minimum:"1" doc:"Entity ID to interact with"`
	ActionType enum.ActionType       `json:"actionType" enum:"like,save" required:"true" doc:"Action type to perform"`
}
