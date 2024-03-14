package DAO

type Tags struct {
	Tag string `json:"Tag" gorm:"column:tag;type:varchar(15);primaryKey"`
	VID string `json:"VID" gorm:"column:VID;type:char(36);primaryKey"`
}

func (t *Tags) TableName() string {
	return "Tags"
}
