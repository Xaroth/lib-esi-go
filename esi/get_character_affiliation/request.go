package request

import (
	"net/http"

	"github.com/xaroth/lib-esi-go/request"
)

type Input struct {
	CharacterIDs []int64 `body:"json"`
}

type Output struct {
	AllianceID    *int64 `json:"alliance_id"`
	CharacterID   int64  `json:"character_id"`
	CorporationID int64  `json:"corporation_id"`
	FactionID     *int64 `json:"faction_id"`
}

var GetCharacterAffiliation = request.Create[Input, []*Output](request.RequestInfo{
	Method: http.MethodPost,
	Path:   "/characters/affiliation",
})
