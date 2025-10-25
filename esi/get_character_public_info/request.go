package request

import (
	"net/http"
	"time"

	"github.com/xaroth/lib-esi-go/request"
)

type Input struct {
	CharacterID int64 `path:"character_id"`
}

type Output struct {
	AllianceID     *int64    `json:"alliance_id"`
	Birthday       time.Time `json:"birthday"`
	BloodlineID    int64     `json:"bloodline_id"`
	CorporationID  int64     `json:"corporation_id"`
	Description    string    `json:"description"`
	FactionID      *int64    `json:"faction_id"`
	Gender         string    `json:"gender"`
	Name           string    `json:"name"`
	RaceID         int64     `json:"race_id"`
	SecurityStatus float64   `json:"security_status"`
	Title          string    `json:"title"`
}

var GetCharacterPublicInfo = request.Create[Input, *Output](request.RequestInfo{
	Method: http.MethodGet,
	Path:   "/characters/{character_id}",
})
