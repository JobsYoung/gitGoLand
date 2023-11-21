package mysql

import (
	"bluebell/models"
	"database/sql"
	"errors"
)

func GetCommunityDetail() (CommunityList []*models.Community, err error) {
	sqlStr := `select community_id,community_name from community `
	err = db.Select(&CommunityList, sqlStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("查无数据")
		}
		return nil, errors.New("查询数据错误")
	}
	return
}

func GetCommunityById(id int64) (communityDetail *models.CommunityDetail, err error) {
	communityDetail = new(models.CommunityDetail)
	sqlStr := `select community_id,community_name,introduction,create_time from community where community_id =?`
	err = db.Get(communityDetail, sqlStr, id)
	if err == sql.ErrNoRows {
		err = errors.New("查无数据")
	}
	return
}
