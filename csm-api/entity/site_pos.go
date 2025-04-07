package entity

import (
	"github.com/guregu/null"
)

type SitePos struct {
	AddressNameDepth1     null.String `json:"address_name_depth1" db:"ADDRESS_NAME_DEPTH1"`
	AddressNameDepth2     null.String `json:"address_name_depth2" db:"ADDRESS_NAME_DEPTH2"`
	AddressNameDepth3     null.String `json:"address_name_depth3" db:"ADDRESS_NAME_DEPTH3"`
	AddressNameDepth4     null.String `json:"address_name_depth4" db:"ADDRESS_NAME_DEPTH4"`
	AddressNameDepth5     null.String `json:"address_name_depth5" db:"ADDRESS_NAME_DEPTH5"`
	RoadAddressNameDepth1 null.String `json:"road_address_name_depth1" db:"ROAD_ADDRESS_NAME_DEPTH1"`
	RoadAddressNameDepth2 null.String `json:"road_address_name_depth2" db:"ROAD_ADDRESS_NAME_DEPTH2"`
	RoadAddressNameDepth3 null.String `json:"road_address_name_depth3" db:"ROAD_ADDRESS_NAME_DEPTH3"`
	RoadAddressNameDepth4 null.String `json:"road_address_name_depth4" db:"ROAD_ADDRESS_NAME_DEPTH4"`
	RoadAddressNameDepth5 null.String `json:"road_address_name_depth5" db:"ROAD_ADDRESS_NAME_DEPTH5"`
	Latitude              null.Float  `json:"latitude" db:"LATITUDE"`
	Longitude             null.Float  `json:"longitude" db:"LONGITUDE"`
	RoadAddress           null.String `json:"road_address" db:"ROAD_ADDRESS"`
	ZoneCode              null.String `json:"zone_code" db:"ZONE_CODE"`
	BuildingName          null.String `json:"building_name" db:"BUILDING_NAME"`
	Base
}

//type SitePosSql struct {
//	AddressNameDepth1     sql.NullString  `db:"ADDRESS_NAME_DEPTH1"`
//	AddressNameDepth2     sql.NullString  `db:"ADDRESS_NAME_DEPTH2"`
//	AddressNameDepth3     sql.NullString  `db:"ADDRESS_NAME_DEPTH3"`
//	AddressNameDepth4     sql.NullString  `db:"ADDRESS_NAME_DEPTH4"`
//	AddressNameDepth5     sql.NullString  `db:"ADDRESS_NAME_DEPTH5"`
//	RoadAddressNameDepth1 sql.NullString  `db:"ROAD_ADDRESS_NAME_DEPTH1"`
//	RoadAddressNameDepth2 sql.NullString  `db:"ROAD_ADDRESS_NAME_DEPTH2"`
//	RoadAddressNameDepth3 sql.NullString  `db:"ROAD_ADDRESS_NAME_DEPTH3"`
//	RoadAddressNameDepth4 sql.NullString  `db:"ROAD_ADDRESS_NAME_DEPTH4"`
//	RoadAddressNameDepth5 sql.NullString  `db:"ROAD_ADDRESS_NAME_DEPTH5"`
//	Latitude              sql.NullFloat64 `db:"LATITUDE"`
//	Longitude             sql.NullFloat64 `db:"LONGITUDE"`
//	RegDate               sql.NullTime    `db:"REG_DATE"`
//	RoadAddress           sql.NullString  `db:"UDF_VAL_01"`
//	ZoneCode              sql.NullString  `db:"UDF_VAL_02"`
//	BuildingName          sql.NullString  `db:"UDF_VAL_03"`
//}

type MapPoint struct {
	X string `json:"x"`
	Y string `json:"y"`
}

//func (s *SitePos) ToSitePos(sql *SitePosSql) *SitePos {
//	s.AddressNameDepth1 = sql.AddressNameDepth1.String
//	s.AddressNameDepth2 = sql.AddressNameDepth2.String
//	s.AddressNameDepth3 = sql.AddressNameDepth3.String
//	s.AddressNameDepth4 = sql.AddressNameDepth4.String
//	s.AddressNameDepth5 = sql.AddressNameDepth5.String
//	s.RoadAddressNameDepth1 = sql.RoadAddressNameDepth1.String
//	s.RoadAddressNameDepth2 = sql.RoadAddressNameDepth2.String
//	s.RoadAddressNameDepth3 = sql.RoadAddressNameDepth3.String
//	s.RoadAddressNameDepth4 = sql.RoadAddressNameDepth4.String
//	s.RoadAddressNameDepth5 = sql.RoadAddressNameDepth5.String
//	s.Latitude = sql.Latitude.Float64
//	s.Longitude = sql.Longitude.Float64
//	s.RegDate = sql.RegDate.Time
//	s.RoadAddress = sql.RoadAddress.String
//	s.ZoneCode = sql.ZoneCode.String
//	s.BuildingName = sql.BuildingName.String
//
//	return s
//}
