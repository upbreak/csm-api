package entity

import (
	"github.com/guregu/null"
)

type SitePos struct {
	Sno                   null.Int    `json:"sno" db:"SNO"`
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
