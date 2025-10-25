package request

import (
	"net/http"
	"time"

	"github.com/xaroth/lib-esi-go/request"
)

type Input struct {
	CorporationID int64 `path:"corporation_id"`
}

type Output struct {
	AllianceID    *int64    `json:"alliance_id"`
	CeoID         int64     `json:"ceo_id"`
	CreatorID     int64     `json:"creator_id"`
	DateFounded   time.Time `json:"date_founded"`
	Description   string    `json:"description"`
	FactionID     *int64    `json:"faction_id"`
	HomeStationID *int64    `json:"home_station_id"`
	MemberCount   int64     `json:"member_count"`
	Name          string    `json:"name"`
	Shares        int64     `json:"shares"`
	TaxRate       float64   `json:"tax_rate"`
	Ticker        string    `json:"ticker"`
	URL           string    `json:"url"`
	WarEligible   bool      `json:"war_eligible"`
}

var GetCorporationPublicInfo = request.Create[Input, *Output](request.RequestInfo{
	Method: http.MethodGet,
	Path:   "/corporations/{corporation_id}",
})
