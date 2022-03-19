package dao

import "gin/model"

func CheckRegisterByUserName(user model.User) error {
	var userId int
	rows := DB.QueryRow("SELECT userId FROM user WHERE username = ? ", user.Username)
	if rows.Err() != nil {
		return rows.Err()
	}
	err := rows.Scan(&userId)
	if err != nil {
		return err
	}
	return nil
}

func RegisterByOAuth(user model.User) error {
	sql := "insert into user(username,password)values(?,?)"
	_, err := DB.Exec(sql, user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}
