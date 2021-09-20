package model

type UserAndAccounts struct {
	User     UserHTTP  `json:"user,omitempty"`
	Accounts []Account `json:"accounts,omitempty"`
}
