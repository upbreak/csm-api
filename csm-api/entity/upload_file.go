package entity

import "github.com/guregu/null"

type UploadFile struct {
	FileType    null.String `json:"file_type" db:"FILE_TYPE"`
	FilePath    null.String `json:"file_path" db:"FILE_PATH"`
	FileName    null.String `json:"file_name" db:"FILE_NAME"`
	UploadRound null.Int    `json:"upload_round" db:"UPLOAD_ROUND"`
	WorkDate    null.Time   `json:"work_date" db:"WORK_DATE"`
	Jno         null.Int    `json:"jno" db:"JNO"`
	Base
}
