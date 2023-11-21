package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func CreatePostHandler(c *gin.Context) {
	//1.获取参数及参数的校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}
	//获取用户id
	uid, err := GetCurUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
	}
	p.AuthorID = uid

	//2.创建帖子
	err = logic.CreatePost(p)
	if err != nil {
		zap.L().Error("CreatePost Error = ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}

func GetPostsHandler(c *gin.Context) {
	//1.获取参数及参数的校验
	id := c.Param("id")
	pid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
	}
	//2.获取帖子信息
	data, err := logic.GetPostDetail(pid)
	if err != nil {
		zap.L().Error("CreatePost Error = ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回数据
	ResponseSuccess(c, data)
}

func GetPostListHandler(c *gin.Context) {
	//1..获取参数及参数的校验
	p := models.ParamPosts{
		Page:  1,
		Size:  10,
		Order: models.OrderByTime,
		CId:   0,
	}
	//c.ShouldBind() 根据请求的数据类型选择相应的方法去获取数据
	err := c.ShouldBindQuery(&p)
	if err != nil {
		zap.L().Error("CreatePost Error = ", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.获取帖子信息
	data, err := logic.GetPostList(&p)
	if err != nil {
		zap.L().Error("CreatePost Error = ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回数据
	ResponseSuccess(c, data)
}

func VoteProcessHandler(c *gin.Context) {
	//1.获取参数及参数的校验
	v := new(models.ParamVote)
	err := c.ShouldBindJSON(v)
	if err != nil {
		zap.L().Error("Vote with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {
			ResponseError(c, CodeInvalidParam)
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	UserId, err := GetCurUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	//2.生成分数
	data, err := logic.VoteForPost(UserId, v)
	if err != nil {
		zap.L().Error("Vote Error", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, err.Error())
		return
	}
	//3.返回响应
	ResponseSuccess(c, data)
}
