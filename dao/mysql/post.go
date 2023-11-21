package mysql

import (
	"bluebell/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

// CreatePost  根据信息创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(
    post_id,title,content,author_id,community_id) 
	values (?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.Id, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostById 根据id查询单个贴子数据
func GetPostDetail(pid int64) (P *models.PostDetail, err error) {
	P = new(models.PostDetail)
	sqlStr := `select post_id,title,content,author_id,community_id,create_time from post where id =?`
	err = db.Get(P, sqlStr, pid)
	return
}

func GetPostList(ids []string) (posts []*models.PostDetail, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time
				from post 
				where post_id in (?)
				order by FIND_IN_SET(post_id,?)
				`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&posts, query, args...)
	return
}
