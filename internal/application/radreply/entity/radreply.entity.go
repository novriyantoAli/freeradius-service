package entity

type Radreply struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string `json:"username" gorm:"index;not null;size:64"`
	Attribute string `json:"attribute" gorm:"not null;size:64"`
	Op        string `json:"op" gorm:"not null;size:2;default:'='"`
	Value     string `json:"value" gorm:"not null;size:253"`
}

func (r Radreply) TableName() string {
	return "radreply"
}
