package models

// 定义请求的参数结构体

const (
	OrderByTime  = "Time"
	OrderByScore = "Score"
)

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamVote struct {
	PostID string `json:"post_id" binding:"required"`
	Score  int8   `json:"score,string" binding:"oneof=-1 0 1"`
}

type ParamPosts struct {
	Page  int64  `json:"page" form:"page"`
	Size  int64  `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
	CId   int64  `json:"community_id" form:"cid"`
}
