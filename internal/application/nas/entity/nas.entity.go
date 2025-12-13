package entity

import (
	"time"

	"gorm.io/gorm"
)

type NAS struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	NASName          string    `json:"nasname" gorm:"index;uniqueIndex;not null;size:128"`
	ShortName        string    `json:"shortname" gorm:"size:32"`
	Type             string    `json:"type" gorm:"size:30;default:'other'"`
	Ports            int       `json:"ports"`
	Secret           string    `json:"secret" gorm:"not null;default:'secret'"`
	Server           string    `json:"server" gorm:"size:64"`
	Community        string    `json:"community" gorm:"size:50"`
	Description      string    `json:"description" gorm:"size:200;default:'RADIUS Client'"`
	RequireMa        string    `json:"require_ma" gorm:"size:4;default:'auto'"`
	LimitProxyState  string    `json:"limit_proxy_state" gorm:"size:4;default:'auto'"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (n NAS) TableName() string {
	return "nas"
}
