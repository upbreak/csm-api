package entity

import (
	"database/sql"
	"time"
)

// 공지 기간
type NoticePeriod struct {
	PeriodCode string `json:"period_code"`
	NoticeNM   string `json:"notice_nm"`
	// Day      int64  `json:"day"`
	// Month    int64  `json:"month"`
}

type NoticePeriods []*NoticePeriod

type NoticePeriodSql struct {
	PeriodCode sql.NullString `db:"PERIOD_CODE"`
	NoticeNM   sql.NullString `db:"NOTICE_NM"`
	// Day      sql.NullInt64  `db:"DAY"`
	// Month    sql.NullInt64  `db:"MONTH"`
}

type NoticePeriodSqls []*NoticePeriodSql

// 공지사항
type NoticeID int64

type Notice struct {
	RowNum       int64     `json:"row_num"`
	Idx          NoticeID  `json:"idx"`
	Sno          int64     `json:"sno"`
	Jno          int64     `json:"jno"`
	JobName      string    `json:"job_name"`
	JobLocName   string    `json:"job_loc_name"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	ShowYN       string    `json:"show_yn"`
	RegUno       int64     `json:"reg_uno"`
	RegUser      string    `json:"reg_user"`
	RegDate      time.Time `json:"reg_date"`
	UserDutyName string    `json:"user_duty_name"`
	UserInfo     string    `json:"user_info"`
	ModUno       int64     `json:"mod_uno"`
	ModUser      string    `json:"mod_user"`
	ModDate      time.Time `json:"mod_date"`
	PeriodCode   string    `json:"period_code"`
	NoticeNm     string    `json:"notice_nm"`
	PostingDate  time.Time `json:"posting_date"`
	IsImportant  string    `json:"is_important"`
}

type Notices []*Notice

type NoticeSql struct {
	RowNum       sql.NullInt64  `db:"RNUM"`
	Idx          NoticeID       `db:"IDX"`
	Sno          sql.NullInt64  `db:"SNO"`
	Jno          sql.NullInt64  `db:"JNO"`
	JobName      sql.NullString `db:"JOB_NAME"`
	JobLocName   sql.NullString `db:"JOB_LOC_NAME"`
	Title        sql.NullString `db:"TITLE" validate:"required"`
	Content      sql.NullString `db:"CONTENT" validate:"required"`
	ShowYN       sql.NullString `db:"SHOW_YN"`
	RegUno       sql.NullInt64  `db:"REG_UNO"`
	RegUser      sql.NullString `db:"REG_USER"`
	RegDate      sql.NullTime   `db:"REG_DATE"`
	UserDutyName sql.NullString `db:"DUTY_NAME"`
	UserInfo     sql.NullString `db:"USER_INFO"`
	ModUno       sql.NullInt64  `db:"MOD_UNO"`
	ModUser      sql.NullString `db:"MOD_USER"`
	ModDate      sql.NullTime   `db:"MOD_DATE"`
	PeriodCode   sql.NullString `db:"PERIOD_CODE"`
	NoticeNm     sql.NullString `db:"NOTICE_NM"`
	PostingDate  sql.NullTime   `db:"POSTING_DATE"`
	IsImportant  sql.NullString `db:"IS_IMPORTANT"`
}

type NoticeSqls []*NoticeSql

func (n *Notice) ToNotice(noticeSql *NoticeSql) *Notice {
	n.RowNum = noticeSql.RowNum.Int64
	n.Idx = noticeSql.Idx
	n.Sno = noticeSql.Sno.Int64
	n.Jno = noticeSql.Jno.Int64
	n.JobName = noticeSql.JobName.String
	n.JobLocName = noticeSql.JobLocName.String
	n.Title = noticeSql.Title.String
	n.Content = noticeSql.Content.String
	n.ShowYN = noticeSql.ShowYN.String
	n.RegUno = noticeSql.RegUno.Int64
	n.RegUser = noticeSql.RegUser.String
	n.RegDate = noticeSql.RegDate.Time
	n.UserDutyName = noticeSql.UserDutyName.String
	n.UserInfo = noticeSql.UserInfo.String
	n.ModUno = noticeSql.RegUno.Int64
	n.ModUser = noticeSql.ModUser.String
	n.ModDate = noticeSql.ModDate.Time
	n.PeriodCode = noticeSql.PeriodCode.String
	n.NoticeNm = noticeSql.NoticeNm.String
	n.PostingDate = noticeSql.PostingDate.Time
	n.IsImportant = noticeSql.IsImportant.String

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
	n.Idx = notice.Idx

	if notice.Jno != 0 {
		n.Jno = sql.NullInt64{Valid: true, Int64: notice.Jno}
	} else {
		n.Jno = sql.NullInt64{Valid: true, Int64: 0}
	}

	if notice.Sno != 0 {
		n.Sno = sql.NullInt64{Valid: true, Int64: notice.Sno}
	} else {
		n.Sno = sql.NullInt64{Valid: false}
	}

	if notice.JobName != "" {
		n.JobName = sql.NullString{Valid: true, String: notice.JobName}
	} else {
		n.JobName = sql.NullString{Valid: false}
	}

	if notice.JobLocName != "" {
		n.JobLocName = sql.NullString{Valid: true, String: notice.JobLocName}
	} else {
		n.JobLocName = sql.NullString{Valid: false}
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

	if notice.UserDutyName != "" {
		n.UserDutyName = sql.NullString{Valid: true, String: notice.UserDutyName}
	} else {
		n.UserDutyName = sql.NullString{Valid: false}
	}

	if notice.UserInfo != "" {
		n.UserInfo = sql.NullString{Valid: true, String: notice.UserInfo}
	} else {
		n.UserInfo = sql.NullString{Valid: false}
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

	if notice.PeriodCode != "" {
		n.PeriodCode = sql.NullString{Valid: true, String: notice.PeriodCode}
	} else {
		n.PeriodCode = sql.NullString{Valid: false}
	}
	if notice.NoticeNm != "" {
		n.NoticeNm = sql.NullString{Valid: true, String: notice.NoticeNm}
	} else {
		n.NoticeNm = sql.NullString{Valid: false}
	}

	if notice.PostingDate.IsZero() != true {
		n.PostingDate = sql.NullTime{Valid: true, Time: notice.PostingDate}
	} else {
		n.PostingDate = sql.NullTime{Valid: false}
	}

	if notice.IsImportant != "" {
		n.IsImportant = sql.NullString{Valid: true, String: notice.IsImportant}
	} else {
		n.IsImportant = sql.NullString{Valid: false}
	}
	return n
}
