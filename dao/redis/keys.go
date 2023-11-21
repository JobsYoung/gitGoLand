package redis

const (
	KeyPrefix             = "bluebell:"
	KeyPostTimeZset       = "post:time"  //zset；帖子及发帖时间
	KeyPostScoreZset      = "post:score" //zset；帖子及帖子分数
	KeyPostVotedZsetPF    = "post:voted"
	KeyCommunityPostSetPF = "community:post"
)

func TotalKey(key string) string {
	return KeyPrefix + key
}
