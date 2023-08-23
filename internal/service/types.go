package service

type User struct {
	ID int
}

type Balance struct {
	User
	Balance int
}

type BalanceChange struct {
	User
	Changing int
}
