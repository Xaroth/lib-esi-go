package request

import (
	"net/http"

	"github.com/xaroth/lib-esi-go/request"
)

type Output struct {
	CorporationID        *int64  `json:"corporation_id"`
	Description          string  `json:"description"`
	FactionID            int64   `json:"faction_id"`
	IsUnique             bool    `json:"is_unique"`
	MilitiaCorporationID *int64  `json:"militia_corporation_id"`
	Name                 string  `json:"name"`
	SizeFactor           float64 `json:"size_factor"`
	SolarSystemID        *int64  `json:"solar_system_id"`
	StationCount         int64   `json:"station_count"`
	StationSystemCount   int64   `json:"station_system_count"`
}

var GetFactions = request.CreateStatic[[]*Output](request.RequestInfo{
	Method: http.MethodGet,
	Path:   "/universe/factions",
})
