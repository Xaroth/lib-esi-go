package request

import (
	"net/http"
	"time"

	"github.com/xaroth/lib-esi-go/request"
)

type Input struct {
	AllianceID int64 `path:"alliance_id"`
}

type Output struct {
	CreatorCorporationID  int64     `json:"creator_corporation_id"`
	CreatorID             int64     `json:"creator_id"`
	DateFounded           time.Time `json:"date_founded"`
	ExecutorCorporationID *int64    `json:"executor_corporation_id"`
	FactionID             *int64    `json:"faction_id"`
	Name                  string    `json:"name"`
	Ticker                string    `json:"ticker"`
}

var GetAlliancePublicInfo = request.Create[Input, *Output](request.RequestInfo{
	Method: http.MethodGet,
	Path:   "/alliances/{alliance_id}",
})
