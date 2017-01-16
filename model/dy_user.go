package model

import "time"

const DyUserTableName = "dy_user"

type DyUser struct {
	Id                  int64       `db:"id" json:"id"`
	Uid                 int64       `db:"uid" json:"uid"`
	Nickname            string      `db:"nickname" json:"nickname"`
	Level               int         `db:"level" json:"level"`
	Strength            int         `db:"strength" json:"strength"`
	Gift                int         `db:"gift" json:"gift"`
	PlatformPrivilege   DyPrivilege `db:"platform_privilege" json:"platform_privilege"`
	DeserveLevel        int         `db:"deserve_level" json:"deserve_level"`
	DeserveCount        int         `db:"deserve_count" json:"deserve_count"`
	BdeserveLevel       int         `db:"bdeserve_level" json:"bdeserve_level"`
	CreateTime          time.Time   `db:"create_time" json:"create_time"`
	last_update         time.Time   `db:"last_update" json:"last_update"`
	FirstAppearedRoomId int64       `db:"first_appeared_room_id" json:"first_appeared_room_id"`
	LastAppearedRoomId  int64       `db:"last_appeared_room_id" json:"last_appeared_room_id"`
}
