package logic

import (
	"VideoWeb/DAO"
	EntitySets "VideoWeb/DAO/EntitySets"
	"VideoWeb/define"
)

func CreateMessage(msg *define.Message) error {
	err := DAO.DB.Model(&EntitySets.Message{}).Create(&msg).Error
	return err
}
