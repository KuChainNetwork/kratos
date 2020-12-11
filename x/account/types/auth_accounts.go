package types

type Accounts []string
type AuthAccounts struct {
	Auth     string   `json:"auth"`
	Accounts Accounts `json:"accounts"`
}

func NewAuthAccount(auth, acc string) AuthAccounts {
	return AuthAccounts{
		Auth:     auth,
		Accounts: []string{acc},
	}
}

func (a *AuthAccounts) AddAccount(acc string) {
	a.Accounts = append(a.Accounts, acc)
}

func (a *AuthAccounts) DeleteAccount(acc string) {
	for i := len(a.Accounts) - 1; i >= 0; i-- {
		if a.Accounts[i] == acc {
			a.Accounts = append(a.Accounts[:i], a.Accounts[i+1:]...)
		}
	}
}

func (a AuthAccounts) GetAccounts() Accounts {
	return a.Accounts
}
