package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	assets "github.com/YaleOpenLab/openx/assets"
	xlm "github.com/YaleOpenLab/openx/xlm"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
	wallet "github.com/YaleOpenLab/openx/wallet"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func parseYamlProject(fileName string, feJson string, projIndex int) error {
	viper.SetConfigType("yaml")
	viper.SetConfigName(fileName)
	viper.AddConfigPath("./data-sandbox")
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "error while reading values from config file")
	}

	project, err := opensolar.RetrieveProject(8)
	if err != nil {
		return err
	}

	termsHelper := viper.Get("Terms").(map[string]interface{})
	if projIndex == 8 {
		terms := make([]opensolar.TermsHelper, 8)
		i := 0
		for _, elem := range termsHelper {
			// elem inside here is a map of "variable": values.
			newMap := elem.(map[string]interface{})
			terms[i].Variable = newMap["variable"].(string)
			terms[i].Value = newMap["value"].(string)
			terms[i].RelevantParty = newMap["relevantparty"].(string)
			terms[i].Note = newMap["note"].(string)
			terms[i].Status = newMap["status"].(string)
			terms[i].SupportDoc = newMap["supportdoc"].(string)
			i += 1
		}

		project.Terms = terms
	} else {
		terms := make([]opensolar.TermsHelper, 6)
		i := 0
		for _, elem := range termsHelper {
			// elem inside here is a map of "variable": values.
			newMap := elem.(map[string]interface{})
			terms[i].Variable = newMap["variable"].(string)
			terms[i].Value = newMap["value"].(string)
			terms[i].RelevantParty = newMap["relevantparty"].(string)
			terms[i].Note = newMap["note"].(string)
			terms[i].Status = newMap["status"].(string)
			terms[i].SupportDoc = newMap["supportdoc"].(string)
			i += 1
		}

		project.Terms = terms
	}

	var executiveSummary opensolar.ExecutiveSummaryHelper

	execSummaryReader := viper.Get("ExecutiveSummary.Investment").(map[string]interface{})
	execSummaryWriter := make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.Investment = execSummaryWriter

	execSummaryReader = viper.Get("ExecutiveSummary.Financials").(map[string]interface{})
	execSummaryWriter = make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.Financials = execSummaryWriter

	execSummaryReader = viper.Get("ExecutiveSummary.ProjectSize").(map[string]interface{})
	execSummaryWriter = make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.ProjectSize = execSummaryWriter

	execSummaryReader = viper.Get("ExecutiveSummary.SustainabilityMetrics").(map[string]interface{})
	execSummaryWriter = make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.SustainabilityMetrics = execSummaryWriter

	project.ExecutiveSummary = executiveSummary

	var bullets opensolar.BulletHelper
	bullets.Bullet1 = viper.Get("Bullets.Bullet1").(string)
	bullets.Bullet2 = viper.Get("Bullets.Bullet2").(string)
	bullets.Bullet3 = viper.Get("Bullets.Bullet3").(string)

	project.Bullets = bullets

	var architecture opensolar.ArchitectureHelper

	architecture.SolarArray = viper.Get("Architecture.SolarArray").(string)
	architecture.DailyAvgGeneration = viper.Get("Architecture.DailyAvgGeneration").(string)
	architecture.System = viper.Get("Architecture.System").(string)
	architecture.InverterSize = viper.Get("Architecture.InverterSize").(string)

	project.Architecture = architecture

	project.Index = viper.Get("Index").(int)
	project.Name = viper.Get("Name").(string)
	project.State = viper.Get("State").(string)
	project.Country = viper.Get("Country").(string)
	project.TotalValue = viper.Get("TotalValue").(float64)
	project.Metadata = viper.Get("Metadata").(string)
	project.PanelSize = viper.Get("PanelSize").(string)
	project.PanelTechnicalDescription = viper.Get("PanelTechnicalDescription").(string)
	project.Inverter = viper.Get("Inverter").(string)
	project.ChargeRegulator = viper.Get("ChargeRegulator").(string)
	project.ControlPanel = viper.Get("ControlPanel").(string)
	project.CommBox = viper.Get("CommBox").(string)
	project.ACTransfer = viper.Get("ACTransfer").(string)
	project.SolarCombiner = viper.Get("SolarCombiner").(string)
	project.Batteries = viper.Get("Batteries").(string)
	project.IoTHub = viper.Get("IoTHub").(string)
	project.Rating = viper.Get("Rating").(string)
	project.EstimatedAcquisition = viper.Get("EstimatedAcquisition").(int)
	project.BalLeft = viper.Get("BalLeft").(float64)
	project.InterestRate = viper.Get("InterestRate").(float64)
	project.Tax = viper.Get("Tax").(string)
	project.DateInitiated = viper.Get("DateInitiated").(string)
	project.DateFunded = viper.Get("DateFunded").(string)
	project.AuctionType = viper.Get("AuctionType").(string)
	project.InvestmentType = viper.Get("InvestmentType").(string)
	project.PaybackPeriod = viper.Get("PaybackPeriod").(int)
	project.Stage = viper.Get("Stage").(int)
	project.SeedInvestmentFactor = viper.Get("SeedInvestmentFactor").(float64)
	project.SeedInvestmentCap = viper.Get("SeedInvestmentCap").(float64)
	project.ProposedInvestmentCap = viper.Get("ProposedInvestmentCap").(float64)
	project.SelfFund = viper.Get("SelfFund").(float64)
	project.SecurityIssuer = viper.Get("SecurityIssuer").(string)
	project.BrokerDealer = viper.Get("BrokerDealer").(string)
	project.EngineeringLayoutType = viper.Get("EngineeringLayoutType").(string)
	project.MapLink = viper.Get("MapLink").(string)

	project.FEText, err = parseJsonText(feJson)
	if err != nil {
		log.Fatal(err)
	}

	return project.Save()
}

func parseYaml(fileName string, feJson string) error {
	viper.SetConfigType("yaml")
	viper.SetConfigName(fileName)
	viper.AddConfigPath("./data-sandbox")
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "error while reading values from config file")
	}

	var project opensolar.Project
	terms := make([]opensolar.TermsHelper, 6)
	termsHelper := viper.Get("Terms").(map[string]interface{})

	i := 0
	for _, elem := range termsHelper {
		// elem inside here is a map of "variable": values.
		newMap := elem.(map[string]interface{})
		terms[i].Variable = newMap["variable"].(string)
		terms[i].Value = newMap["value"].(string)
		terms[i].RelevantParty = newMap["relevantparty"].(string)
		terms[i].Note = newMap["note"].(string)
		terms[i].Status = newMap["status"].(string)
		terms[i].SupportDoc = newMap["supportdoc"].(string)
		i += 1
	}

	project.Terms = terms
	var executiveSummary opensolar.ExecutiveSummaryHelper

	execSummaryReader := viper.Get("ExecutiveSummary.Investment").(map[string]interface{})
	execSummaryWriter := make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.Investment = execSummaryWriter

	execSummaryReader = viper.Get("ExecutiveSummary.Financials").(map[string]interface{})
	execSummaryWriter = make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.Financials = execSummaryWriter

	execSummaryReader = viper.Get("ExecutiveSummary.ProjectSize").(map[string]interface{})
	execSummaryWriter = make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.ProjectSize = execSummaryWriter

	execSummaryReader = viper.Get("ExecutiveSummary.SustainabilityMetrics").(map[string]interface{})
	execSummaryWriter = make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.SustainabilityMetrics = execSummaryWriter

	project.ExecutiveSummary = executiveSummary

	var bullets opensolar.BulletHelper
	bullets.Bullet1 = viper.Get("Bullets.Bullet1").(string)
	bullets.Bullet2 = viper.Get("Bullets.Bullet2").(string)
	bullets.Bullet3 = viper.Get("Bullets.Bullet3").(string)

	project.Bullets = bullets

	var architecture opensolar.ArchitectureHelper

	architecture.SolarArray = viper.Get("Architecture.SolarArray").(string)
	architecture.DailyAvgGeneration = viper.Get("Architecture.DailyAvgGeneration").(string)
	architecture.System = viper.Get("Architecture.System").(string)
	architecture.InverterSize = viper.Get("Architecture.InverterSize").(string)

	project.Architecture = architecture

	project.Index = viper.Get("Index").(int)
	project.Name = viper.Get("Name").(string)
	project.State = viper.Get("State").(string)
	project.Country = viper.Get("Country").(string)
	project.TotalValue = viper.Get("TotalValue").(float64)
	project.Metadata = viper.Get("Metadata").(string)
	project.PanelSize = viper.Get("PanelSize").(string)
	project.PanelTechnicalDescription = viper.Get("PanelTechnicalDescription").(string)
	project.Inverter = viper.Get("Inverter").(string)
	project.ChargeRegulator = viper.Get("ChargeRegulator").(string)
	project.ControlPanel = viper.Get("ControlPanel").(string)
	project.CommBox = viper.Get("CommBox").(string)
	project.ACTransfer = viper.Get("ACTransfer").(string)
	project.SolarCombiner = viper.Get("SolarCombiner").(string)
	project.Batteries = viper.Get("Batteries").(string)
	project.IoTHub = viper.Get("IoTHub").(string)
	project.Rating = viper.Get("Rating").(string)
	project.EstimatedAcquisition = viper.Get("EstimatedAcquisition").(int)
	project.BalLeft = viper.Get("BalLeft").(float64)
	project.InterestRate = viper.Get("InterestRate").(float64)
	project.Tax = viper.Get("Tax").(string)
	project.DateInitiated = viper.Get("DateInitiated").(string)
	project.DateFunded = viper.Get("DateFunded").(string)
	project.AuctionType = viper.Get("AuctionType").(string)
	project.InvestmentType = viper.Get("InvestmentType").(string)
	project.PaybackPeriod = viper.Get("PaybackPeriod").(int)
	project.Stage = viper.Get("Stage").(int)
	project.SeedInvestmentFactor = viper.Get("SeedInvestmentFactor").(float64)
	project.SeedInvestmentCap = viper.Get("SeedInvestmentCap").(float64)
	project.ProposedInvestmentCap = viper.Get("ProposedInvestmentCap").(float64)
	project.SelfFund = viper.Get("SelfFund").(float64)
	project.SecurityIssuer = viper.Get("SecurityIssuer").(string)
	project.BrokerDealer = viper.Get("BrokerDealer").(string)
	project.EngineeringLayoutType = viper.Get("EngineeringLayoutType").(string)
	project.MapLink = viper.Get("MapLink").(string)

	project.FEText, err = parseJsonText(feJson)
	if err != nil {
		log.Fatal(err)
	}

	return project.Save()
}

func populateStaticData() error {
	var err error
	log.Println("populating db with static data")
	err = createAllStaticEntities()
	if err != nil {
		return err
	}
	err = parseYaml("1kwy", "data-sandbox/1kw.json")
	if err != nil {
		return err
	}
	// project: One Kilowatt Project
	// STAGE 7 - Puerto Rico
	err = populateStaticData1kw()
	if err != nil {
		return err
	}
	err = parseYaml("1mwy", "data-sandbox/1mw.json")
	if err != nil {
		return err
	}
	// project: One Megawatt Project
	// STAGE 4 - New Hampshire
	err = populateStaticData1mw()
	if err != nil {
		return err
	}
	err = parseYaml("10kwy", "data-sandbox/10kw.json")
	if err != nil {
		return err
	}
	// project: Ten Kilowatt Project
	// STAGE 8 - Connecticut Homeless Shelter
	err = populateStaticData10kw()
	if err != nil {
		return err
	}
	err = parseYaml("10mwy", "data-sandbox/10mw.json")
	if err != nil {
		return err
	}
	// project: Ten Megawatt Project
	// STAGE 2 - Puerto Rico Public School Bond
	err = populateStaticData10MW()
	if err != nil {
		return err
	}
	err = parseYaml("100kwy", "data-sandbox/100kw.json")
	if err != nil {
		return err
	}
	// project: One HUndred Kilowatt Project
	// STAGE 1 - Rwanda Project
	err = populateStaticData100KW()
	if err != nil {
		return err
	}
	return nil
}

func populateDynamicData() error {
	var err error
	// we ignore errors here since they are b ound to happen (guarantor related errors)
	err = populateDynamicData1kw()
	if err != nil {
		log.Println("error while populating 1kw project", err)
	}
	err = populateDynamicData1mw()
	if err != nil {
		log.Println("error while populating 1mw project", err)
	}
	err = populateDynamicData10kw()
	if err != nil {
		log.Println("error while populating 10kw project", err)
	}
	return nil
}

// CreateSandbox is the main function that controls data insertion as part of the sandbox environment
func CreateSandbox() error {
	// project: Puerto Rico Project
	// STAGE 7 - Puerto Rico
	var err error
	err = populateStaticData()
	if err != nil {
		return err
	}
	err = populateDynamicData()
	if err != nil {
		return err
	}
	err = populateAdditionalData()
	if err != nil {
		return err
	}
	return nil
}

func parseJsonText(fileName string) (map[string]interface{}, error) {

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}

// seed additional data for a few specific investors that are useful for showing in demos
func populateAdditionalData() error {
	openlab, err := database.RetrieveInvestor(46)
	if err != nil {
		return err
	}
	openlab.U.Email = "martin.wainstein@yale.edu"
	openlab.U.Address = "254 Elm Street"
	openlab.U.Country = "US"
	openlab.U.City = "New Haven"
	openlab.U.ZipCode = "06511"
	openlab.U.RecoveryPhone = "1800SECRETS"
	openlab.U.Description = "The Yale OPen Lab is the open innovation lab at the Tsai Centre for Innovative Thinking at Yale"
	err = openlab.U.Save()
	if err != nil {
		return err
	}
	err = openlab.Save()
	if err != nil {
		return err
	}

	// insert data for one specific recipient
	pasto, err := database.RetrieveRecipient(47)
	if err != nil {
		return err
	}
	pasto.U.Email = "supasto2018@gmail.com"
	pasto.U.Address = "Puerto Rico, PR"
	pasto.U.Country = "US"
	pasto.U.City = "Puerto Rico"
	pasto.U.ZipCode = "00909"
	pasto.U.RecoveryPhone = "1800SECRETS"
	pasto.U.Description = "S.U. Pasto School is a school in Puerto Rico"
	err = pasto.U.Save()
	if err != nil {
		return err
	}
	err = pasto.Save()
	if err != nil {
		return err
	}
	dci, err := opensolar.RetrieveEntity(1)
	if err != nil {
		return err
	}

	dci.U.Email = "dci@mit.edu"
	dci.U.Address = "MIT Media Lab"
	dci.U.Country = "US"
	dci.U.City = "Cambridge"
	dci.U.ZipCode = "02142"
	dci.U.RecoveryPhone = "1800SECRETS"
	dci.U.Description = "The Digital Currency Initiative at the MIT Media Lab"

	err = dci.U.Save()
	if err != nil {
		return err
	}

	// we now need to register the dci as an investor as well
	var inv database.Investor
	inv.U = dci.U
	err = inv.Save()
	if err != nil {
		return err
	}
	var recp database.Recipient
	recp.U = dci.U
	err = recp.Save()
	if err != nil {
		return err
	}

	err = xlm.GetXLM(dci.U.PublicKey)
	if err != nil {
		return err
	}

	seed, err := wallet.DecryptSeed(recp.U.EncryptedSeed, "x")
	if err != nil {
		return err
	}

	// send the pasto school account some money so we can demo using it on the frontend
	txhash, err := assets.TrustAsset(consts.Code, consts.StablecoinPublicKey, "10000000000", recp.U.PublicKey, seed)
	if err != nil {
		return err
	}
	log.Println("TX HASH for dci trusting stableUSD: ", txhash)

	_, txhash, err = assets.SendAssetFromIssuer(consts.Code, recp.U.PublicKey, "600", consts.StablecoinSeed, consts.StablecoinPublicKey)
	if err != nil {
		log.Println("SEED: ", consts.StablecoinSeed)
		return err
	}
	log.Println("TX HASH for dci getting stableUSD: ", txhash)

	recp, err = database.RetrieveRecipient(47)
	if err != nil {
		return err
	}

	seed, err = wallet.DecryptSeed(recp.U.EncryptedSeed, "x")
	if err != nil {
		return err
	}

	// send the pasto school account some money so we can demo using it on the frontend
	txhash, err = assets.TrustAsset(consts.Code, consts.StablecoinPublicKey, "10000000000", recp.U.PublicKey, seed)
	if err != nil {
		return err
	}
	log.Println("TX HASH for pasto school trusting stableUSD: ", txhash)

	_, txhash, err = assets.SendAssetFromIssuer(consts.Code, recp.U.PublicKey, "600", consts.StablecoinSeed, consts.StablecoinPublicKey)
	if err != nil {
		log.Println("SEED: ", consts.StablecoinSeed)
		return err
	}
	log.Println("TX HASH for pasto school getting stableUSD: ", txhash)

	err = xlm.GetXLM(recp.U.SecondaryWallet.PublicKey)
	if err != nil {
		return err
	}

	seed, err = wallet.DecryptSeed(recp.U.SecondaryWallet.EncryptedSeed, "x")
	if err != nil {
		return err
	}

	// send the pasto school account some money so we can demo using it on the frontend

	txhash, err = assets.TrustAsset(consts.Code, consts.StablecoinPublicKey, "10000000000", recp.U.SecondaryWallet.PublicKey, seed)
	if err != nil {
		return err
	}
	log.Println("TX HASH for pasto school sec wallet trusting stableUSD: ", txhash)

	_, txhash, err = assets.SendAssetFromIssuer(consts.Code, recp.U.SecondaryWallet.PublicKey, "10000", consts.StablecoinSeed, consts.StablecoinPublicKey)
	if err != nil {
		log.Println("SEED: ", consts.StablecoinSeed)
		return err
	}
	log.Println("TX HASH for pasto school sec wallet getting stableUSD: ", txhash)

	return nil
}
