package common

import (
	"time"
)

type PublicModel struct {
	ID        uint      `json:"id" gorm:"id"` // Primary Key IDs
	CreatedAt time.Time // Creation time
	UpdatedAt time.Time // update time
}

type Img struct {
	Src string `json:"src" gorm:"src"`
	Tp  string `json:"type" gorm:"type"`
}
