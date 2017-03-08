package model

const DyUserTableName = "dy_user"

type DyUser struct {
	Model

	Uid                 int64       `db:"uid" json:"uid"`
	Nickname            string      `db:"nickname" json:"nickname"`
	Level               int         `db:"level" json:"level"`
	Strength            int         `db:"strength" json:"strength"`
	GiftRank            int         `db:"gift_rank" json:"gift_rank"`
	PlatformPrivilege   DyPrivilege `db:"platform_privilege" json:"platform_privilege"`
	DeserveLevel        int         `db:"deserve_level" json:"deserve_level"`
	DeserveCount        int         `db:"deserve_count" json:"deserve_count"`
	BdeserveLevel       int         `db:"bdeserve_level" json:"bdeserve_level"`
	FirstAppearedRoomId int64       `db:"first_appeared_room_id" json:"first_appeared_room_id"`
	LastAppearedRoomId  int64       `db:"last_appeared_room_id" json:"last_appeared_room_id"`
}

func (*DyUser) TableName() string {
	return DyUserTableName
}

func (a DyUser) Equals(b DyUser) bool {
	// omit some field
	a.Model = b.Model
	a.FirstAppearedRoomId = b.FirstAppearedRoomId

	return a == b
}

func (a *DyUser) Assign(b *DyUser) {
	// omit some field
	m := a.Model
	f := a.FirstAppearedRoomId

	*a = *b
	a.Model = m
	a.FirstAppearedRoomId = f
}
