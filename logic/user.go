package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

// SignUp  进行用户注册业务逻辑
func SignUp(p *models.ParamSignUp) error {
	//1.判断用户是否已存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	//2.生成userID
	userID := snowflake.GenID()
	u := models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//3.存进数据库
	return mysql.InsertUser(&u)

}

// Login  进行用户登录业务逻辑
func Login(u *models.ParamLogin) (User *models.User, err error) {
	//1.根据输入信息构建用户信息
	User = &models.User{
		Username: u.Username,
		Password: u.Password,
	}
	//2.进行用户登录校验
	err = mysql.Login(User)
	if err != nil {
		return nil, err
	}
	//生产jwt
	User.Token, err = jwt.GenToken(User.UserID, User.Username)
	return
}
