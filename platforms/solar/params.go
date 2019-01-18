package solar

// TODO: add more paramters here that would help identify a given solar project
type SolarParams struct {
	Index int // an Index to keep quick track of how many projects exist

	PanelSize   string // size of the given panel, for diplsaying to the user who wants to bid stuff
	TotalValue  int    // the total money that we need from investors
	Location    string // where this specific solar panel is located
	MoneyRaised int    // total money that has been raised until now
	Years       int    // number of years the recipient has chosen to opt for
	Metadata    string // any other metadata can be stored here

	// once all funds have been raised, we need to set assetCodes
	InvestorAssetCode string
	DebtAssetCode     string
	PaybackAssetCode  string

	BalLeft float64 // denotes the balance left to pay by the party
	Votes   int     // the number of votes towards a proposed contract by investors

	DateInitiated string // date the project was created
	DateFunded    string // date when the project was funded
	DateLastPaid  string // date the project was last paid
	// Percentage raised is not stored in the database since that can be calculated by the UI
}
