package entity

import "github.com/guregu/null"

type Device struct {
	RowNum   null.Int    `json:"rnum" db:"RNUM"`
	Dno      null.Int    `json:"dno" db:"DNO"`             // 홍채인식기 고유번호
	Sno      null.Int    `json:"sno" db:"SNO"`             // 현장 고유번호
	DeviceSn null.String `json:"device_sn" db:"DEVICE_SN"` // 홍채인식기 시리얼번호
	DeviceNm null.String `json:"device_nm" db:"DEVICE_NM"` // 홍채인식기 장치명
	Etc      null.String `json:"etc" db:"ETC"`             // 비고
	IsUse    null.String `json:"is_use" db:"IS_USE"`       // 사용여부
	SiteNm   null.String `json:"site_nm" db:"SITE_NM"`
	Base
}

type Devices []*Device
