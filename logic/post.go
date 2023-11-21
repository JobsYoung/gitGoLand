package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"errors"
	"strconv"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	//1. 雪花算法生产PostId
	p.Id = snowflake.GenID()
	//2. 保存到数据库层
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	//3. 缓存到redis
	err = redis.CreatePost(p)
	return
}

func GetPostDetail(pid int64) (P *models.PostDetail, err error) {
	P, err = mysql.GetPostDetail(pid)
	P.CommunityDetail, err = mysql.GetCommunityById(P.CommunityID)
	P.AuthorName, err = mysql.GetUserById(P.AuthorID)
	return
}

func GetPostList(p *models.ParamPosts) (data []*models.PostDetail, err error) {
	//1.从redis查询id列表
	var ids []string
	if p.CId == 0 {
		ids, err = redis.GetPostListIds(p)
	} else {
		ids, err = redis.GetPostListIdsByCommunity(p)
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return nil, errors.New("查到0行数据")
	}
	//2. 根据id去MySQL数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id的顺序返回
	data, err = mysql.GetPostList(ids)
	if err != nil {
		return
	}
	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	//3.补全社区和作者、投票信息
	for idx, post := range data {
		post.CommunityDetail, err = mysql.GetCommunityById(post.CommunityID)
		post.AuthorName, err = mysql.GetUserById(post.AuthorID)
		post.VoteNum = voteData[idx]
	}
	return
}

func VoteForPost(uid int64, v *models.ParamVote) (data float64, err error) {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", uid),
		zap.String("postID", v.PostID),
		zap.Int8("direction", v.Score))
	data, err = redis.VoteForPost(strconv.Itoa(int(uid)), v.PostID, float64(v.Score))
	return
}
