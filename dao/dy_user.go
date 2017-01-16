package dao

import (
	"github.com/gocraft/dbr"
	"github.com/zwh8800/66ana/model"
)

func DyUserById(sr dbr.SessionRunner, id int64) (*model.DyUser, error) {
	dyUser := &model.DyUser{}
	if err := sr.Select("*").From(model.DyUserTableName).
		Where("id = ?", id).LoadStruct(dyUser); err != nil {
		return nil, err
	}
	return dyUser
}

func DyUserByUid(sr dbr.SessionRunner, uid int64) (*model.DyUser, error) {
	dyUser := &model.DyUser{}
	if err := sr.Select("*").From(model.DyUserTableName).
		Where("uid = ?", uid).LoadStruct(dyUser); err != nil {
		return nil, err
	}
	return dyUser
}

func InsertIncompleteDyUser(sr dbr.SessionRunner) (*model.DyUser, error) {

}
