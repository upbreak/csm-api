package entity

import (
	"database/sql"
	"time"
)

type SitePos struct {
	AddressNameDepth1     string    `json:"address_name_depth1"`
	AddressNameDepth2     string    `json:"address_name_depth2"`
	AddressNameDepth3     string    `json:"address_name_depth3"`
	AddressNameDepth4     string    `json:"address_name_depth4"`
	AddressNameDepth5     string    `json:"address_name_depth5"`
	RoadAddressNameDepth1 string    `json:"road_address_name_depth1"`
	RoadAddressNameDepth2 string    `json:"road_address_name_depth2"`
	RoadAddressNameDepth3 string    `json:"road_address_name_depth3"`
	RoadAddressNameDepth4 string    `json:"road_address_name_depth4"`
	RoadAddressNameDepth5 string    `json:"road_address_name_depth5"`
	Latitude              float64   `json:"latitude"`
	Longitude             float64   `json:"longitude"`
	RegDate               time.Time `json:"reg_date"`
}

type SitePosSql struct {
	AddressNameDepth1     sql.NullString  `db:"ADDRESS_NAME_DEPTH1"`
	AddressNameDepth2     sql.NullString  `db:"ADDRESS_NAME_DEPTH2"`
	AddressNameDepth3     sql.NullString  `db:"ADDRESS_NAME_DEPTH3"`
	AddressNameDepth4     sql.NullString  `db:"ADDRESS_NAME_DEPTH4"`
	AddressNameDepth5     sql.NullString  `db:"ADDRESS_NAME_DEPTH5"`
	RoadAddressNameDepth1 sql.NullString  `db:"ROAD_ADDRESS_NAME_DEPTH1"`
	RoadAddressNameDepth2 sql.NullString  `db:"ROAD_ADDRESS_NAME_DEPTH2"`
	RoadAddressNameDepth3 sql.NullString  `db:"ROAD_ADDRESS_NAME_DEPTH3"`
	RoadAddressNameDepth4 sql.NullString  `db:"ROAD_ADDRESS_NAME_DEPTH4"`
	RoadAddressNameDepth5 sql.NullString  `db:"ROAD_ADDRESS_NAME_DEPTH5"`
	Latitude              sql.NullFloat64 `db:"LATITUDE"`
	Longitude             sql.NullFloat64 `db:"LONGITUDE"`
	RegDate               sql.NullTime    `db:"REG_DATE"`
}

func (s *SitePos) ToSitePos(sql *SitePosSql) *SitePos {
	s.AddressNameDepth1 = sql.AddressNameDepth1.String
	s.AddressNameDepth2 = sql.AddressNameDepth2.String
	s.AddressNameDepth3 = sql.AddressNameDepth3.String
	s.AddressNameDepth4 = sql.AddressNameDepth4.String
	s.AddressNameDepth5 = sql.AddressNameDepth5.String
	s.RoadAddressNameDepth1 = sql.RoadAddressNameDepth1.String
	s.RoadAddressNameDepth2 = sql.RoadAddressNameDepth2.String
	s.RoadAddressNameDepth3 = sql.RoadAddressNameDepth3.String
	s.RoadAddressNameDepth4 = sql.RoadAddressNameDepth4.String
	s.RoadAddressNameDepth5 = sql.RoadAddressNameDepth5.String
	s.Latitude = sql.Latitude.Float64
	s.Longitude = sql.Longitude.Float64
	s.RegDate = sql.RegDate.Time

	return s
}
