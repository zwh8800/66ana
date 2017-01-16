package model

type DyPrivilege int

func (p DyPrivilege) IsGuest() bool {
	return p == 0
}

func (p DyPrivilege) IsUser() bool {
	return p == 1
}

func (p DyPrivilege) IsPrivileged() bool {
	return p != 0 && p != 1
}
