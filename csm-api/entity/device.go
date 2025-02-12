package entity

import (
	"database/sql"
	"time"
)

type Device struct {
	RowNum   int64     `json:"row_num"`
	Dno      int64     `json:"dno"`       // 홍채인식기 고유번호
	Sno      int64     `json:"sno"`       // 현장 고유번호
	DeviceSn string    `json:"device_sn"` // 홍채인식기 시리얼번호
	DeviceNm string    `json:"device_nm"` // 홍채인식기 장치명
	Etc      string    `json:"etc"`       // 비고
	IsUse    string    `json:"is_use"`    // 사용여부
	SiteNm   string    `json:"site_nm"`
	RegDate  time.Time `json:"reg_date"`  // 최초 생성일시
	RegAgent string    `json:"reg_agent"` // 최초 생성정보
	RegUser  string    `json:"reg_user"`  // 최초 생성자
	RegUno   int64     `json:"reg_uno"`   // 최초 생성 UNO
	ModDate  time.Time `json:"mod_date"`  // 최종 수정일시
	ModAgent string    `json:"mod_agent"` // 최종 수정정보
	ModUser  string    `json:"mod_user"`  // 최종 수정자
	ModUno   int64     `json:"mod_uno"`   // 최종 수정 UNO
}

type Devices []*Device

type DeviceSql struct {
	RowNum   sql.NullInt64  `db:"RNUM"`
	Dno      sql.NullInt64  `db:"DNO"`
	Sno      sql.NullInt64  `db:"SNO"`
	DeviceSn sql.NullString `db:"DEVICE_SN"`
	DeviceNm sql.NullString `db:"DEVICE_NM"`
	Etc      sql.NullString `db:"ETC"`
	IsUse    sql.NullString `db:"IS_USE"`
	SiteNm   sql.NullString `db:"SITE_NM"`
	RegDate  sql.NullTime   `db:"REG_DATE"`
	RegAgent sql.NullString `db:"REG_AGENT"`
	RegUser  sql.NullString `db:"REG_USER"`
	RegUno   sql.NullInt64  `db:"REG_UNO"`
	ModDate  sql.NullTime   `db:"MOD_DATE"`
	ModAgent sql.NullString `db:"MOD_AGENT"`
	ModUser  sql.NullString `db:"MOD_USER"`
	ModUno   sql.NullInt64  `db:"MOD_UNO"`
}

type DeviceSqls []*DeviceSql

func (d *Device) ToDevice(sql *DeviceSql) *Device {
	d.RowNum = sql.RowNum.Int64
	d.Dno = sql.Dno.Int64
	d.Sno = sql.Sno.Int64
	d.DeviceSn = sql.DeviceSn.String
	d.DeviceNm = sql.DeviceNm.String
	d.Etc = sql.Etc.String
	d.IsUse = sql.IsUse.String
	d.SiteNm = sql.SiteNm.String
	d.RegDate = sql.RegDate.Time
	d.RegAgent = sql.RegAgent.String
	d.RegUser = sql.RegUser.String
	d.RegUno = sql.RegUno.Int64
	d.ModDate = sql.ModDate.Time
	d.ModAgent = sql.ModAgent.String
	d.ModUser = sql.ModUser.String
	d.ModUno = sql.ModUno.Int64

	return d
}

func (ds *Devices) ToDevices(sqls *DeviceSqls) *Devices {
	for _, sql := range *sqls {
		d := &Device{}
		d.ToDevice(sql)
		*ds = append(*ds, d)
	}
	return ds
}

// - device entity.DeviceSql: SNO, DEVICE_SN, DEVICE_NM, ETC, IS_USE, REG_USER
func (d *DeviceSql) OfDeviceSql(device Device) *DeviceSql {
	if device.Sno != 0 {
		d.Sno = sql.NullInt64{Valid: true, Int64: device.Sno}
	} else {
		d.Sno = sql.NullInt64{Valid: false}
	}
	if device.DeviceSn != "" {
		d.DeviceSn = sql.NullString{Valid: true, String: device.DeviceSn}
	} else {
		d.DeviceSn = sql.NullString{Valid: false} // NULL로 설정
	}
	if device.DeviceNm != "" {
		d.DeviceNm = sql.NullString{Valid: true, String: device.DeviceNm}
	} else {
		d.DeviceNm = sql.NullString{Valid: false} // NULL로 설정
	}
	if device.Etc != "" {
		d.Etc = sql.NullString{Valid: true, String: device.Etc}
	} else {
		d.Etc = sql.NullString{Valid: false} // NULL로 설정
	}
	if device.IsUse != "" {
		d.IsUse = sql.NullString{Valid: true, String: device.IsUse}
	} else {
		d.IsUse = sql.NullString{Valid: false} // NULL로 설정
	}
	if device.RegUser != "" {
		d.RegUser = sql.NullString{Valid: true, String: device.RegUser}
	} else {
		d.RegUser = sql.NullString{Valid: false} // NULL로 설정
	}
	if device.ModUser != "" {
		d.ModUser = sql.NullString{Valid: true, String: device.ModUser}
	} else {
		d.ModUser = sql.NullString{Valid: false} // NULL로 설정
	}

	return d
}
