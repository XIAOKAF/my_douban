package dao

import (
	"gin/model"
	time2 "time"
)

func IsRegister(mobile string) error {
	var userId int
	rows := DB.QueryRow("SELECT userId FROM user WHERE mobile = ? ", mobile)
	if rows.Err() != nil {
		return rows.Err()
	}
	err := rows.Scan(&userId)
	if err != nil {
		return err
	}
	return nil
}

func StoreVerifyCode(user model.User) error {
	sql := "update user set verifyCode = ?, sendTime = ? where mobile = ?"
	_, err := DB.Exec(sql, user.VerifyCode, user.SendTime, user.Mobile)
	if err != nil {
		return err
	}
	return nil
}

func SelectVerifyCodeAndSendTime(mobile string) (string, time2.Time, error) {
	var code string
	var time time2.Time
	sql := "select verifyCode,sendTime from user where mobile = ?"
	rows, err := DB.Query(sql, mobile)
	if err != nil {
		return "", time, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&code, &time)
		if err != nil {
			return code, time, err
		}
	}
	return code, time, nil
}

func Register(mobile string) error {
	sql := "insert into user(mobile)value(?)"
	_, err := DB.Exec(sql, mobile)
	if err != nil {
		return err
	}
	return nil
}

func ImprovePersonalInfo(user model.User) error {
	sql := "update user set username = ?, password = ?, selfIntroduction = ? where mobile = ?"
	_, err := DB.Exec(sql, user.Username, user.Password, user.SelfIntroduction, user.Mobile)
	if err != nil {
		return err
	}
	return nil
}

func SelectUsernameByMobile(mobile string) error {
	var username string
	rows := DB.QueryRow("SELECT username FROM user WHERE username = ? ", username)
	if rows.Err() != nil {
		return rows.Err()
	}
	err := rows.Scan(&username)
	if err != nil {
		return err
	}
	return nil
}

func SelectInfoByMobile(mobile string) (error, model.User) {
	var user model.User
	sql := "select userId,userName,password,selfIntroduction from user where mobile = ?"
	rows, err := DB.Query(sql, mobile)
	if err != nil {
		return err, user
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user.UserId, &user.Username, &user.Password, &user.SelfIntroduction)
		if err != nil {
			return err, user
		}
	}
	return nil, user
}

func SelectPasswordByMobile(user model.User) (error, string) {
	var pwd string
	sql := "select password from user where mobile = ?"
	rows, err := DB.Query(sql, user.Mobile)
	if err != nil {
		return err, pwd
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&pwd)
		if err != nil {
			return err, pwd
		}
	}
	return nil, pwd
}

func ChangePassword(user model.User) error {
	sql := "update user set password = ? where mobile = ?"
	_, err := DB.Exec(sql, user.Password, user.Mobile)
	if err != nil {
		return err
	}
	return nil
}
