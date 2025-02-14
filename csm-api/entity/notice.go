package entity

import (
	"database/sql"
	"time"
)

type NoticeID int64

type Notice struct {
	RowNum  int64     `json:"row_num"`
	Idx     NoticeID  `json:"idx"`
	Sno     int64     `json:"sno"`
	SiteNm  string    `json:"site_nm"`
	LocCode string    `json:"loc_code"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	ShowYN  string    `json:"show_yn"`
	RegUno  int64     `json:"reg_uno"`
	RegUser string    `json:"reg_user"`
	RegDate time.Time `json:"reg_date"`
	ModUno  int64     `json:"mod_uno"`
	ModUser string    `json:"mod_user"`
	ModDate time.Time `json:"mod_date"`
}

type Notices []*Notice

type NoticeSql struct {
	RowNum  sql.NullInt64  `db:"RNUM"`
	Idx     NoticeID       `db:"IDX"`
	Sno     sql.NullInt64  `db:"SNO"`
	SiteNm  sql.NullString `db:"SITE_NM"`
	LocCode sql.NullString `db:"LOC_CODE"`
	Title   sql.NullString `db:"TITLE" validate:"required"`
	Content sql.NullString `db:"CONTENT" validate:"required"`
	ShowYN  sql.NullString `db:"SHOW_YN"`
	RegUno  sql.NullInt64  `db:"REG_UNO" validate:"required"`
	RegUser sql.NullString `db:"REG_USER" validate:"required"`
	RegDate sql.NullTime   `db:"REG_DATE"`
	ModUno  sql.NullInt64  `db:"MOD_UNO" validate:"required"`
	ModUser sql.NullString `db:"MOD_USER" validate:"required"`
	ModDate sql.NullTime   `db:"MOD_DATE"`
}

type NoticeSqls []*NoticeSql

func (n *Notice) ToNotice(noticeSql *NoticeSql) *Notice {
	n.RowNum = noticeSql.RowNum.Int64
	n.Idx = noticeSql.Idx
	n.Sno = noticeSql.Sno.Int64
	n.SiteNm = noticeSql.SiteNm.String
	n.Title = noticeSql.Title.String
	n.Content = noticeSql.Content.String
	n.ShowYN = noticeSql.ShowYN.String
	n.RegUno = noticeSql.RegUno.Int64
	n.RegUser = noticeSql.RegUser.String
	n.RegDate = noticeSql.RegDate.Time
	n.ModUno = noticeSql.RegUno.Int64
	n.ModUser = noticeSql.ModUser.String
	n.ModDate = noticeSql.ModDate.Time

	return n
}

func (n *Notices) ToNotices(noticeSqls *NoticeSqls) *Notices {
	for _, noticeSql := range *noticeSqls {
		notice := Notice{}
		notice.ToNotice(noticeSql)
		*n = append(*n, &notice)
	}
	return n
}

func (n *NoticeSql) OfNoticeSql(notice Notice) *NoticeSql {
	if notice.Sno != 0 {
		n.Sno = sql.NullInt64{Valid: true, Int64: notice.Sno}
	} else {
		n.Sno = sql.NullInt64{Valid: false}
	}

	if notice.SiteNm != "" {
		n.SiteNm = sql.NullString{Valid: true, String: notice.SiteNm}
	} else {
		n.SiteNm = sql.NullString{Valid: false}
	}

	if notice.LocCode != "" {
		n.LocCode = sql.NullString{Valid: true, String: notice.LocCode}
	} else {
		n.LocCode = sql.NullString{Valid: false}
	}

	if notice.Title != "" {
		n.Title = sql.NullString{Valid: true, String: notice.Title}
	} else {
		n.Title = sql.NullString{Valid: false}
	}

	if notice.Content != "" {
		n.Content = sql.NullString{Valid: true, String: notice.Content}
	} else {
		n.Content = sql.NullString{Valid: false}
	}

	if notice.ShowYN != "" {
		n.ShowYN = sql.NullString{Valid: true, String: notice.ShowYN}
	} else {
		n.ShowYN = sql.NullString{Valid: false}
	}

	if notice.RegUno != 0 {
		n.RegUno = sql.NullInt64{Valid: true, Int64: notice.RegUno}
	} else {
		n.RegUno = sql.NullInt64{Valid: false}
	}

	if notice.RegUser != "" {
		n.RegUser = sql.NullString{Valid: true, String: notice.RegUser}
	} else {
		n.RegUser = sql.NullString{Valid: false}
	}

	if notice.RegDate.IsZero() != true {
		n.RegDate = sql.NullTime{Valid: true, Time: notice.RegDate}
	} else {
		n.RegDate = sql.NullTime{Valid: false}
	}

	if notice.ModUno != 0 {
		n.ModUno = sql.NullInt64{Valid: true, Int64: notice.ModUno}
	} else {
		n.ModUno = sql.NullInt64{Valid: false}
	}

	if notice.ModUser != "" {
		n.ModUser = sql.NullString{Valid: true, String: notice.ModUser}
	} else {
		n.ModUser = sql.NullString{Valid: false}
	}

	if notice.ModDate.IsZero() != true {
		n.ModDate = sql.NullTime{Valid: true, Time: notice.ModDate}
	} else {
		n.ModDate = sql.NullTime{Valid: false}
	}

	return n
}
