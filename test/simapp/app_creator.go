package simapp

// CreateAppForTest create app for test
func CreateAppForTest(w *Wallet) *SimApp {
	genAccs := NewGenesisAccounts(w.GetRootAuth())

	return SetupWithGenesisAccounts(genAccs)
}
