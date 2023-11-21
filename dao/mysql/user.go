package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

const secret = "bluebell"

var (
	ErrorUserNotExist = errors.New("用户不存在")
	ErrorUserExist    = errors.New("用户已存在")
	ErrorPwdIncorrect = errors.New("密码不正确")
)

// encryptPassword  对密码进行md5加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func GetUserById(uid int64) (username string, err error) {
	sqlStr := `select username from user where user_id=?`
	err = db.Get(&username, sqlStr, uid)
	return
}

// CheckUserExist  检查用户名是否存在
func CheckUserExist(username string) error {
	sqlStr := `select count(user_id) from user where username =?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return nil
}

// InsertUser  向数据库中插入一条新的用户数据
func InsertUser(user *models.User) error {
	//对密码进行加密
	user.Password = encryptPassword(user.Password)
	//执行Sql语句入库
	sqlStr := `insert into user(user_id,username,password)values(?,?,?)`
	_, err := db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return err
}

// Login  校验用户登录信息是否正确
func Login(user *models.User) error {
	oPassword := user.Password
	// 根据用户名查询信息
	sqlStr := `select user_id,username,password from user where username =?`
	err := db.Get(user, sqlStr, user.Username)

	//校验用户名是否正确
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		return err
	}

	//校验密码是否正确
	pwd := encryptPassword(oPassword)
	if pwd != user.Password {
		return ErrorPwdIncorrect
	}
	return nil
}
