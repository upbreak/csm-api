package entity

import "github.com/guregu/null"

type Menu struct {
	MenuId   null.String `json:"menu_id" db:"MENU_ID"`
	MenuNm   null.String `json:"menu_nm" db:"MENU_NM"`
	HasChild null.String `json:"has_child" db:"HAS_CHILD"`
	ParentId null.String `json:"parent_id" db:"PARENT_ID"`
	SvgName  null.String `json:"svg_name" db:"SVG_NAME"`
	IsTemp   null.String `json:"is_temp" db:"IS_TEMP"`
	RoleCode null.String `json:"role_code" db:"ROLE_CODE"`
}

type MenuRes struct {
	Parent []Menu `json:"parent"`
	Child  []Menu `json:"child"`
}
