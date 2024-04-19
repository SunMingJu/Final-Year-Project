package common

import (
	"time"
)

type PublicModel struct {
	ID uint `json:"id" gorm:"column:id"` // primary key id
	CreatedAt time.Time // Creation time
	UpdatedAt time.Time // update time
}

type Img struct {
	Src string `json:"src" gorm:"column:src"`
	Tp  string `json:"type" gorm:"column:type"`
}
