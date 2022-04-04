package service

import (
	"gin/dao"
	"gin/model"
	"gorm.io/gorm"
)

func SelectCelebrityById(celebrity model.Celebrity) (error, bool) {
	err := dao.SelectCelebrityById(celebrity)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false //未查询到返回false，反之返回true
		}
		return err, true
	}
	return nil, true
}

func InsertCelebrityDetails(celebrity model.Celebrity) error {
	err := dao.InsertCelebrityDetails(celebrity)
	if err != nil {
		return err
	}
	return nil
}

func InsertPhotos(picture string, celebrityId string) error {
	err := dao.InsertPhotos(picture, celebrityId)
	if err != nil {
		return err
	}
	return nil
}

func InsertRewards(rewards string, celebrityId string) error {
	err := dao.InsertRewards(rewards, celebrityId)
	if err != nil {
		return err
	}
	return nil
}

func InsertRecentWorks(works model.RecentWorks, celebrityId string) error {
	err := dao.InsertRecentWorks(works, celebrityId)
	if err != nil {
		return err
	}
	return nil
}

func SelectCelebrityDetails(celebrity model.Celebrity) error {
	err := dao.SelectCelebrityDetails(celebrity)
	if err != nil {
		return err
	}
	return nil
}

func SelectPhotos(celebrityId string) (error, [5]string) {
	err, photoArr := dao.SelectPhotos(celebrityId)
	if err != nil {
		return err, photoArr
	}
	return nil, photoArr
}

func SelectRewards(celebrityId string) (error, [3]string) {
	err, rewardsArr := dao.SelectRewards(celebrityId)
	if err != nil {
		return err, rewardsArr
	}
	return nil, rewardsArr
}

func SelectRecentWorks(celebrityId string) (error, [5]model.RecentWorks) {
	err, worksArr := dao.SelectRecentWorks(celebrityId)
	if err != nil {
		return err, worksArr
	}
	return nil, worksArr
}
