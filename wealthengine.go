package wealthengine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Wealthengine API Client
type WealthEngine struct {
	Dev    string
	Prod   string
	Mode   string
	APIKey string
	client *http.Client
}

// A lookup by physical address
type AddressLookup struct {
	Last_name     string `json:"last_name"`
	First_name    string `json:"first_name"`
	Address_line1 string `json:"address_line1"`
	Address_line2 string `json:"address_line2"`
	City          string `json:"city"`
	State         string `json:"state"`
	Zip           string `json:"zip"`
}

// A lookup by physical address for use with scoring endpoint
type AddressLookupScore struct {
	*AddressLookup
	Model string `json:"model"`
}

// A lookup by email and name
type EmailLookup struct {
	Email      string `json:"email"`
	Last_name  string `json:"last_name"`
	first_name string `json:"first_name"`
}

// A lookup by email and name for use with scoring endpoint
type EmailLookupScore struct {
	*EmailLookup
	Model string `json:"model"`
}

// A lookup by phone number and name
type PhoneLookup struct {
	Phone      string `json:"phone"`
	Last_name  string `json:"last_name"`
	First_name string `json:"first_name"`
}

// A lookup by phone number for use with scoring endpoint
type PhoneLookupScore struct {
	*PhoneLookup
	Model string `json:"model"`
}

// A common Wealthengine object describing a range
type WealthEngineRange struct {
	Min       int    `json:"min"`
	Max       int    `json:"max"`
	Value     int    `json:"value"`
	Text      string `json:"text"`
	Text_low  string `json:"text_low"`
	Text_high string `json:"text_high"`
}

// A common Wealthengine object describing a text representation and a value
type TextAndValue struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}

// A person
type HumanBeing struct {
	First_name  string `json:"first_name"`
	Full_name   string `json:"full_name"`
	Last_name   string `json:"last_name"`
	Middle_name string `json:"middle_name"`
}

// A physical address
type Address struct {
	Street_line1 string       `json:"street_line1"`
	Street_line2 string       `json:"street_line2"`
	Street_line3 string       `json:"street_line3"`
	City         string       `json:"city"`
	State        TextAndValue `json:"state"`
	Postal_code  string       `json:"postal_code"`
}

type Location struct {
	Address        Address `json:"address"`
	Personal_phone string  `json:"personal_phone"`
}

// A single email address
type Email struct {
	Email string `json:"email"`
}

// A representation of an indvidual's wealth, including their income and assets
type Wealth struct {
	Cash_on_hand            WealthEngineRange `json:"cash_on_hand"`
	Networth                WealthEngineRange `json:"networth"`
	Total_income            WealthEngineRange `json:"total_income"`
	Business_ownership      WealthEngineRange `json:"business_ownership"`
	Business_sales_volume   WealthEngineRange `json:"business_sales_volume"`
	Accredited_investor     bool              `json:"accredited_investor"`
	Influence_rating        TextAndValue      `json:"influence_rating"`
	Total_stock             WealthEngineRange `json:"total_stock"`
	Stock_holdings_direct   WealthEngineRange `json:"stock_holdings_direct"`
	Stock_holdings_indirect WealthEngineRange `json:"stock_holdings_direct"`
	Investable_assets       WealthEngineRange `json:"investable_assets"`
	Total_assets            WealthEngineRange `json:"total_assets"`
	Total_pensions          WealthEngineRange `json:"total_pensions"`
}

type Professional struct {
	Board_member bool `json:"Board_member"`
}

// A real estate entity
type Realestate struct {
	Total_num_properties   int `json:"total_num_properties"`
	Total_realestate_value struct {
		WealthEngineRange
	}
}

// A personage
type Identity struct {
	Age    int `json:"age"`
	Gender struct {
		Value string `json:"value"`
		Text  string `json:"text"`
	}
	HumanBeing
	Marital_status TextAndValue `json:"marital_status"`
	Emails         []Email      `json:"emails"`
}

// An employment
type Job struct {
	Org_name string       `json:"org_name"`
	Org_type TextAndValue `json:"org_type"`
	Phone    string       `json:"phone"`
	Email    string       `json:"email"`
	Title    string       `json:"title"`
	Address  Address      `json:"address"`
}

// A spousal relationship
type Relationship struct {
	Spouse struct {
		HumanBeing
	}
}

type Demographics struct {
	Has_children bool `json:"has_children"`
}

// Vehicles owned
type Vehicles struct {
	Ownership []TextAndValue `json:"ownership"`
}

// An indvidual's donation activity
type Giving struct {
	Affiliation_inclination    TextAndValue      `json:"affiliation_inclination"`
	P2G_score                  TextAndValue      `json:"p2g_score"`
	Planned_giving             []TextAndValue    `json:"planned_giving"`
	Gift_capacity              WealthEngineRange `json:"gift_capacity"`
	Charitable_donations       WealthEngineRange `json:"charitable_donations"`
	Total_political_donations  WealthEngineRange `json:"total_political_donations"`
	Total_donations            WealthEngineRange `json:"total_donations"`
	Estimated_annual_donations WealthEngineRange `json:"estimated_annual_donations"`
}

type Input struct {
		Output_type string `json:"output_type"`
		Endpoint string `json:"endpoint"`
		Environment string `json:"environment"`
		Query struct {
			Email string `json:"email"`
			First_name string `json:"first_name"`
			Last_name string 	`json:"last_name"`
			Address_line1 string `json:"address_line1"`
			Address_line2 string `json:"address_line2"`
			City string `json:"city"`
			State string `json:"state"`
			Zip string `json:"zip"`
			Phone string `json:"phone"`
			Models string `json:"model"`
		}
	}

type Score struct {
	Model string `json:"model"`
	Score int `json:"score"`
}

type ScoreProfile struct {
	Id int `json:"id"`
	Scores []Score `json:"scores"`
}

type ScoredProfile struct {
	Input Input `json:"input"`
	Profile ScoreProfile `json:"profile"`
}

// A complete WealthEngine profile as described here:
// https://dev.wealthengine.com/api
type Profile struct {
	Id           int          `json:"id"`
	Identity     Identity     `json:"identity"`
	Demographics Demographics `json:"demographics"`
	Relationship Relationship `json:"relationship"`
	Wealth       Wealth       `json:"wealth"`
	Giving       Giving       `json:"giving"`
	Locations    []Location   `json:"locations"`
	Realestate   Realestate   `json:"realestate"`
	Professional Professional `json:"professional"`
	Vehicles     Vehicles     `json:"vehicles"`
	Jobs         []Job        `json:"jobs"`
}

// A single personage lookup meant to be passed as one of a batch of lookups
type BatchLookup struct {
	Last_name     string `json:"last_name"`
	First_name    string `json:"first_name"`
	Address_line1 string `json:"address_line1"`
	Address_line2 string `json:"address_line2"`
	City          string `json:"city"`
	State         string `json:"state"`
	Zip           string `json:"zip"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
}

type BatchProfile struct {
	Input Input `json:"input"`
	Profile Profile `json:"profile"`
}

// A batch of lookups
type Batch struct {
	Profiles []BatchLookup `json:"profiles"`
}

// A batch processing job ID returned by the Wealthengine API
type BatchID struct {
	ID string `json:"batchID"`
}

// A status entry for a batch processing job returned by the Wealthengine API
type BatchJobStatus struct {
	Status string `json:"Status"`
}

// A batch processing result
type BatchJobResults struct {
	Profiles []BatchProfile `json:"profiles"`
}

// Create and return a new WealthEngine client
// Panics if any of the required parameters is missing
func New(apiKey, mode string) *WealthEngine {

	if apiKey == "" {
		panic(errors.New("WealthEngine must be initialized with a valid APIKey from https://dev.wealthengine.com"))
	}

	if mode == "" && (mode != "Dev" || mode != "Prod") {
		panic(errors.New("WealthEngine client must be initialized in either Dev or Prod mode"))
	}

	return &WealthEngine{
		Dev:    "https://api-sandbox.wealthengine.com/v1",
		Prod:   "https://api.wealthengine.com/v1",
		Mode:   mode,
		APIKey: apiKey,
		client: &http.Client{},
	}
}

// Returns a properly formatted resource URL for WealthEngine REST calls
func (w *WealthEngine) FormatWealthEngineUrl(ResourcePath string) string {

	var root = ""
	if w.Mode == "Dev" {
		root = w.Dev
	} else {
		root = w.Prod
	}

	return fmt.Sprintf("%s/%s", root, ResourcePath)
}

// Builds a properly formatted and authorized request to WealthEngine's API
func (w *WealthEngine) FormatRequest(data []byte, ResourceType string) (*http.Request, error) {

	req, buildErr := http.NewRequest("POST", w.FormatWealthEngineUrl(ResourceType), bytes.NewBuffer(data))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("APIKey %s", w.APIKey))

	return req, buildErr
}

// Builds a properly formatted and authorized GET request to WealthEngine's API
func (w *WealthEngine) FormatGetRequest(ResourceType string) (*http.Request, error) {
	req, buildErr := http.NewRequest("GET", w.FormatWealthEngineUrl(ResourceType), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("APIKey %s", w.APIKey))

	return req, buildErr
}

// Look up a profile using an address and attempt to retrieve a score
func (w *WealthEngine) ScoreOneByAddress(last_name, first_name, address_line1, address_line2, city, state, zip, model string) (*ScoredProfile, error) {

	a := &AddressLookup{
		last_name,
		first_name,
		address_line1,
		address_line2,
		city,
		state,
		zip,
	}

	l := &AddressLookupScore{
		a,
		model,
	}

	return w.MakeScoreRequest("score/score_one/by_address/", l)
}

// Look up a profile using an address and attempt to retrieve a score
func (w *WealthEngine) ScoreOneByEmail(email, last_name, first_name, model string) (*ScoredProfile, error) {

	e := &EmailLookup{
		email,
		last_name,
		first_name,
	}

	l := &EmailLookupScore{
		e,
		model,
	}

	return w.MakeScoreRequest("score/score_one/by_email/", l)
}

// Look up a profile using a phone number and attempt to retrieve a score
func (w *WealthEngine) ScoreOneByPhone(phone, last_name, first_name, model string) (*ScoredProfile, error) {

	p := &PhoneLookup{
		phone,
		last_name,
		first_name,
	}

	l := &PhoneLookupScore{
		p,
		model,
	}

	return w.MakeScoreRequest("score/score_one/by_phone/", l)
}

// Match a profile by address
func (w *WealthEngine) MatchOneByAddress(last_name, first_name, address_line1, address_line2, city, state, zip, mode string) (*Profile, error) {

	if mode == "" || (mode != "full" && mode != "basic") {
		mode = "full"
	}

	l := &AddressLookup{
		last_name,
		first_name,
		address_line1,
		address_line2,
		city,
		state,
		zip,
	}

	return w.MakeRequest(fmt.Sprintf("profile/find_one/by_address/%s", mode), l)
}

// Match a profile by email
func (w *WealthEngine) MatchOneByEmail(email, last_name, first_name, mode string) (*Profile, error) {

	if mode == "" || (mode != "full" && mode != "basic") {
		mode = "full"
	}

	l := &EmailLookup{
		email,
		last_name,
		first_name,
	}

	return w.MakeRequest(fmt.Sprintf("profile/find_one/by_email/%s", mode), l)
}

// Match a profile by phone number
func (w *WealthEngine) MatchOneByPhone(phone, last_name, first_name, mode string) (*Profile, error) {

	if mode == "" || (mode != "full" && mode != "basic") {
		mode = "full"
	}

	l := &PhoneLookup{
		phone,
		last_name,
		first_name,
	}

	return w.MakeRequest(fmt.Sprintf("profile/find_one/by_phone/%s", mode), l)
}

// Batch processing endpoint. Accepts a batch of BatchLookups
func (w *WealthEngine) FindMany(b *Batch, mode string) (*BatchID, error) {

	if mode == "" || (mode != "full" && mode != "basic") {
		mode = "full"
	}

	return w.MakeBatchRequest(b, fmt.Sprintf("profile/find_many/%s", mode))
}

// Get final processed results of a batch processing job
func (w *WealthEngine) GetBatchJobResults(b *BatchID) (*BatchJobResults, error) {

	return w.MakeGetRequest(fmt.Sprintf("profile/find_many/results/%s", b.ID))
}

// Check on the status of an existing batch processing job
func (w *WealthEngine) GetBatchJobStatus(b *BatchID) (*BatchJobStatus, error) {

	req, fErr := w.FormatGetRequest(fmt.Sprintf("job/status/%s", b.ID))

	if fErr != nil {
		return nil, fErr
	}

	resp, err := w.client.Do(req)

	if err != nil {
		return nil, err
	}

	body, readErr := ioutil.ReadAll(resp.Body)

	if readErr != nil {
		return nil, readErr
	}

	s := &BatchJobStatus{}

	uErr := json.Unmarshal(body, &s)

	if uErr != nil {
		return nil, uErr
	}

	return s, nil
}

func (w *WealthEngine) MakeScoreRequest(path string, v interface{}) (*ScoredProfile, error) {

	j, mErr := json.Marshal(v)

	if mErr != nil {
		return nil, mErr
	}

	req, buildErr := w.FormatRequest(j, path)

	if buildErr != nil {
		return nil, buildErr
	}

	resp, reqErr := w.client.Do(req)

	if reqErr != nil {
		return nil, reqErr
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return DecodeScoreResponse(body)
}

// Accepts a WealthEngine API path as a string and any Lookup struct
// Builds an executes a properly formatted request to the WealthEngine API
// Returns a WealthEngine Profile struct
func (w *WealthEngine) MakeRequest(path string, v interface{}) (*Profile, error) {

	j, mErr := json.Marshal(v)

	if mErr != nil {
		return nil, mErr
	}

	req, buildErr := w.FormatRequest(j, path)

	if buildErr != nil {
		return nil, buildErr
	}

	resp, reqErr := w.client.Do(req)

	if reqErr != nil {
		return nil, reqErr
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return DecodeProfile(body)
}

// Format a request to the batch processing results endpoint
func (w *WealthEngine) MakeGetRequest(path string) (*BatchJobResults, error) {

	req, buildErr := w.FormatGetRequest(path)

	if buildErr != nil {
		return nil, buildErr
	}

	resp, reqErr := w.client.Do(req)

	if reqErr != nil {
		return nil, reqErr
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return DecodeBatchResultsResponse(body)
}

// Accepts a batch, comprised of BatchLookups
// Returns a batch processing job ID that can be used for checking batch processing job status
func (w *WealthEngine) MakeBatchRequest(b *Batch, path string) (*BatchID, error) {

	j, mErr := json.Marshal(b)

	if mErr != nil {
		return nil, mErr
	}

	req, buildErr :=w.FormatRequest(j, path)

	if buildErr != nil {
		return nil, buildErr
	}

	resp, reqErr := w.client.Do(req)

	if reqErr != nil {
		return nil, reqErr
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return DecodeBatchResponse(body)
}

// Decode a score profile response
func DecodeScoreResponse(b []byte) (*ScoredProfile, error) {
	s := &ScoredProfile{}
	uErr := json.Unmarshal(b, &s)
	return s, uErr
}

// Decode the results of a batch job
func DecodeBatchResultsResponse(b []byte) (*BatchJobResults, error) {
	r := &BatchJobResults{}
	uErr := json.Unmarshal(b, &r.Profiles)
	return r, uErr
}

// Decode the status of a batch job
func DecodeBatchResponse(b []byte) (*BatchID, error) {
	i := &BatchID{}
	uErr := json.Unmarshal(b, &i)
	return i, uErr
}

// Decode a profile returned by the matching endpoints
func DecodeProfile(b []byte) (*Profile, error) {
	p := &Profile{}
	uErr := json.Unmarshal(b, &p)
	return p, uErr
}
