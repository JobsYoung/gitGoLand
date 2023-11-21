package redis

import (
	"bluebell/models"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const perVoteScore = 432

func CreatePost(p *models.Post) error {
	pipeline := rdb.TxPipeline()
	//缓存帖子创建时间
	pipeline.ZAdd(TotalKey(KeyPostTimeZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: p.Id,
	})

	//缓存帖子分数
	pipeline.ZAdd(TotalKey(KeyPostScoreZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: p.Id,
	})

	//把帖子id加到社区的set
	pipeline.SAdd(TotalKey(KeyCommunityPostSetPF+strconv.Itoa(int(p.CommunityID))), p.Id)
	_, err := pipeline.Exec()
	return err
}

/*
   每票加432分
   投票的限制：
   	自发布一星期之内允许用户投票
*/

func VoteForPost(Uid, Pid string, Score float64) (float64, error) {
	//1.获取帖子发布时间
	CreateTime := rdb.ZScore(TotalKey(KeyPostTimeZset), Pid).Val()
	if float64(time.Now().Unix())-CreateTime > 7*24*3600 {
		return 0, errors.New("帖子超出投票时间范围")
	}
	//2.获取用户之前投票状态，新增分数
	pervVote := rdb.ZScore(TotalKey(KeyPostVotedZsetPF+Pid), Uid).Val()
	addScore := (Score - pervVote) * perVoteScore

	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(TotalKey(KeyPostScoreZset), addScore, Pid)

	//3.更改用户投票
	if Score == 0 {
		pipeline.ZRem(TotalKey(KeyPostVotedZsetPF+Pid), Uid)
	} else {
		//有则修改，无则新增
		pipeline.ZAdd(TotalKey(KeyPostVotedZsetPF+Pid), redis.Z{
			Score:  Score, // 赞成票还是反对票
			Member: Uid,
		})
	}
	_, err := pipeline.Exec()
	return addScore, err
}
