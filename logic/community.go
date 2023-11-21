package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList() (data []*models.Community, err error) {
	return mysql.GetCommunityDetail()
}

func GetCommunityById(id int64) (data *models.CommunityDetail, err error) {

	return mysql.GetCommunityById(id)
}
