package dao

import (
	"gin/model"
	"go/ast"
)

func SelectCelebrityById(celebrity model.Celebrity) error {
	result := DB.Find(&celebrity, celebrity.CelebrityId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func InsertCelebrityDetails(celebrity model.Celebrity) error {
	sql := "insert into celebrity(celebrityId,celebrityName,image,gender,constellation,birthDate,birthplace,jobs,nickname,family,introduction)values(?,?,?,?,?,?,?,?,?,?,?)"
	_, err := DB.Exec(sql, celebrity.CelebrityId, celebrity.CelebrityName, celebrity.Image, celebrity.Gender, celebrity.Constellation, celebrity.BirthDate, celebrity.Birthplace, celebrity.Jobs, celebrity.Nickname, celebrity.Family, celebrity.Introduction)
	if err != nil {
		return err
	}
	return nil
}

func InsertPhotos(picture string, celebrityId string) error {
	sql := "insert into photos(masterId,image)values(?,?)"
	_, err := DB.Exec(sql, celebrityId, picture)
	if err != nil {
		return err
	}
	return nil
}

func InsertRewards(rewards string, celebrityId string) error {
	sql := "insert into rewards(getId,rewardsDetails)values(?,?)"
	_, err := DB.Exec(sql, celebrityId, rewards)
	if err != nil {
		return err
	}
	return nil
}

func InsertRecentWorks(works model.RecentWorks, celebrityId string) error {
	sql := "insert into recentWorks(performerId,workId,workImage,workName,workScore)values(?,?,?,?,?)"
	_, err := DB.Exec(sql, celebrityId, works.WorkId, works.WorkImage, works.WorkName, works.WorkScores)
	if err != nil {
		return err
	}
	return nil
}

func SelectCelebrityDetails(celebrity model.Celebrity) error {
	result := DB.First(&celebrity, "celebrityId = ?", celebrity.CelebrityId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func SelectPhotos(photograph model.Photograph) (error, [5]string) {
	var photo [5]string
	result := DB.Find(&photograph, "masterId = ?", photograph.MasterId)
	if result.Error != nil {
		return result.Error, photo
	}
	rows := result.Scan(&photo)
	if rows.Error != nil {
		return rows.Error, photo
	}
	return nil, photo
}

func SelectRewards(celebrity string) (error, [3]string) {
	var rewardsArr [3]string
	var i int
	sql := "select rewardsDetails from rewards where getId = ?"
	rows, err := DB.Query(sql, celebrity)
	if err != nil {
		return err, rewardsArr
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&rewardsArr[i])
		if err != nil {
			return err, rewardsArr
		}
		i++
	}
	return nil, rewardsArr
}

func SelectRecentWorks(celebrityId string) (error, [5]model.RecentWorks) {
	var worksArr [5]model.RecentWorks
	var works model.RecentWorks
	var i int
	sql := "select workId,workImage,workName,workScore from recentWorks where performerId = ?"
	rows, err := DB.Query(sql, celebrityId)
	if err != nil {
		return err, worksArr
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&works.WorkId, &works.WorkImage, &works.WorkName, &works.WorkScores)
		if err != nil {
			return err, worksArr
		}
		worksArr[i] = works
		i++
	}
	return nil, worksArr
}
