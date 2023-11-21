package redis

import (
	"bluebell/models"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func getPostIdsByKey(key string, p *models.ParamPosts) (ids []string, err error) {
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	//按分数从大到小查询指定数量的元素
	ids, err = rdb.ZRevRange(key, start, end).Result()
	//3.返回ids
	return ids, nil
}

func GetPostListIds(p *models.ParamPosts) (ids []string, err error) {
	//1.根据p.order获取查询的key，默认按时间排序
	key := TotalKey(KeyPostTimeZset)
	if p.Order == models.OrderByScore {
		key = TotalKey(KeyPostScoreZset)
	}

	//2.根据key获取redis中的PostIds
	return getPostIdsByKey(key, p)

}

func GetPostListIdsByCommunity(p *models.ParamPosts) (ids []string, err error) {
	//1.判断排序依据Time 或 Score
	var orderKey string
	if p.Order == models.OrderByTime {
		orderKey = TotalKey(KeyPostTimeZset)
	} else if p.Order == models.OrderByScore {
		orderKey = TotalKey(KeyPostScoreZset)
	} else {
		return nil, errors.New("错误的Order参数")
	}
	cKey := TotalKey(KeyCommunityPostSetPF + strconv.Itoa(int(p.CId)))
	//2.缓存 社区Set与分数ZSet 的ZInterStore生成的ZSet
	key := orderKey + strconv.Itoa(int(p.CId))
	//不存在需要计算，缓存减少计算次数
	if rdb.Exists(key).Val() < 1 {
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey)
		pipeline.Expire(key, 60*time.Second) //设置超时时间
		_, err = pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	//3.根据key获取redis中的PostIds
	return getPostIdsByKey(key, p)
}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPF + id)
	//	// 查找key中分数是1的元素的数量->统计每篇帖子的赞成票的数量
	//	v := client.ZCount(key, "1", "1").Val()
	//	data = append(data, v)
	//}
	// 使用pipeline一次发送多条命令,减少RTT
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := TotalKey(KeyPostVotedZsetPF + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}
