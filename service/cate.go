package service

import (
	"strconv"

	"github.com/zwh8800/66ana/model"
)

func InsertDyCate(cateInfo *model.CateInfo) (*model.DyCate, error) {
	committed := false
	tx := dbConn.Begin()
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	cate, err := cookModelFromCateInfo(cateInfo)
	if err != nil {
		return nil, err
	}

	updatedCate := *cate
	if err := tx.Where(model.DyCate{Cid: cate.Cid}).
		Attrs(cate).FirstOrCreate(cate).Error; err != nil {
		return nil, err
	}
	if !cate.Equals(updatedCate) {
		if err := tx.Model(cate).Update(updatedCate).
			Error; err != nil {
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	committed = true
	return cate, nil
}

func cookModelFromCateInfo(cateInfo *model.CateInfo) (*model.DyCate, error) {
	cid, err := strconv.ParseInt(cateInfo.CateID, 10, 64)
	if err != nil {
		return nil, err
	}
	cate := &model.DyCate{
		Cid:       cid,
		GameName:  cateInfo.GameName,
		ShortName: cateInfo.ShortName,
		GameUrl:   cateInfo.GameURL,
		GameSrc:   cateInfo.GameSrc,
		GameIcon:  cateInfo.GameIcon,
	}
	return cate, nil
}
