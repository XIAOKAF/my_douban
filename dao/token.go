package dao

func CheckRefreshToken(mobile string) error {
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
