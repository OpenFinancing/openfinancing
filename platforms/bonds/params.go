package bonds

// this param is shared by both bonds and coops, could add other stuff in here as well
// the idea is to ahve a common set of params for each platform and then each model
// could boorow this base params and build upon it as desired.
type BondCoopParams struct {
	Index             int
	MaturationDate    string
	MemberRights      string
	SecurityType      string
	InterestRate      float64
	Rating            string
	BondIssuer        string
	Underwriter       string
	DateInitiated     string // date the project was created
	InvestorAssetCode string
	Title             string
	Description       string
	Location          string
}
