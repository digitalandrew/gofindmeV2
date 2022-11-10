package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// The below structs represent the JSON data that comes from the gofundme API and are used to unmarshal the JSON data
type donationResponse struct {
	References references `json:"references"`
	Meta       meta       `json:"meta"`
}

type references struct {
	Donations []donation `json:"donations"`
}

type meta struct {
	LastUpdated string `json:"last_updated_at"`
	HasNext     bool   `json:"has_next"`
}

type donation struct {
	DonationID  float64 `json:"donation_id"`
	Amount      float64 `json:"amount"`
	IsOffline   bool    `json:"is_offline"`
	IsAnonymous bool    `json:"is_anonymous"`
	CreatedAt   string  `json:"created_at"`
	Name        string  `json:"name"`
	ProfileUrl  string  `json:"profile_url"`
	Verified    bool    `json:"verified"`
	Currency    string  `json:"currencycode"`
	FundID      float64 `json:"fund_id"`
	CheckOutID  float64 `json:"checkout_id"`
}

//Function will take in a string of the campaign url, parse out the campaign name and return it as an API url
func createAPIUrl(url string) string {
	s := strings.Split(url, "f/")
	s = strings.Split(s[1], "?")
	return ("https://gateway.gofundme.com/web-gateway/v1/feed/" + s[0] + "/")
}

//Function takes in the API url, an offset and the channels required to return a boolean hasNext value and the donation slice
//The offset is required as the API allows for a max limit of 100 donations per request, the offset is increased from 0 by 100
//for each subsequent request until the hasNext value is false
func getDonations(url string, i int, c chan []donation, ch chan bool) {

	req, err := http.NewRequest("GET", url+"donations", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	q.Add("limit", "100")
	q.Add("offset", strconv.Itoa(i*100))
	req.URL.RawQuery = q.Encode()

	resp, err := http.Get(req.URL.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var dr donationResponse

	json.Unmarshal(body, &dr)

	ch <- dr.Meta.HasNext
	c <- dr.References.Donations

}

func printBoilerPlate() {
	fmt.Println("")
	fmt.Println("**************************************************")
	fmt.Println("               __ _           _         ")
	fmt.Println("              / _(_)         | |               ")
	fmt.Println("   __ _  ___ | |_ _ _ __   __| |_ __ ___   ___ ")
	fmt.Println("  / _` |/ _ \\|  _| | '_ \\ / _` | '_ ` _ \\ / _ \\")
	fmt.Println(" | (_| | (_) | | | | | | | (_| | | | | | |  __/")
	fmt.Println("  \\__, |\\___/|_| |_|_| |_|\\__,_|_| |_| |_|\\___|")
	fmt.Println("   __/ |                                       ")
	fmt.Println("  |___/     ")
	fmt.Println("")
	fmt.Println("**************************************************")
	fmt.Println("")
	fmt.Println("Created by: DigitalAndrew - https://github.com/digitalandrew/gofindme")
	fmt.Println("")
}

func printHelp() {
	fmt.Println("usage example: ")
	fmt.Println("go run main.go -c https://www.gofundme.com/f/digitalandrews-fundraising-campaign -f")
	fmt.Println("")
	fmt.Println("Options: ")
	fmt.Println(" -c           Set Campaign URL")
	fmt.Println(" -f           Write output to CSV file")
}

func printDonations(d []donation) {
	fmt.Printf("%-50s", "Name: ")
	fmt.Println("Amount: ")
	fmt.Println("----------------------------------------------------------------------")
	for i := range d {
		fmt.Printf("%-50s", d[i].Name)
		fmt.Printf("%.2f", d[i].Amount)
		fmt.Println(d[i].Currency)
	}
}

func main() {

	help := flag.Bool("h", false, "Shows help details")
	campaign := flag.String("c", "", "Sets the campaign to search for")
	file := flag.Bool("f", false, "If set the app will write data out to a CSV")
	url := ""
	flag.Parse()
	if *help {
		printHelp()
		os.Exit(0)
	}
	if len(*campaign) > 1 {
		url = *campaign
	} else {
		fmt.Println("Valid Campain URL Required, use -h for help")
		os.Exit(1)
	}
	printBoilerPlate()
	hasNext := true
	donationCounter := 0
	var allDonations []donation
	c := make(chan []donation)
	ch := make(chan bool)

	for hasNext {
		go getDonations(createAPIUrl(url), donationCounter, c, ch)
		donationCounter += 1
		hasNext = <-ch
		allDonations = append(allDonations, <-c...)
	}

	printDonations(allDonations)

	if *file {
		csvFile, err := os.Create("donations.csv")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		csvWriter := csv.NewWriter(csvFile)
		row := []string{"Name", "Donations", "Currency", "Date"}
		csvWriter.Write(row)
		for _, d := range allDonations {
			row[0], row[1], row[2], row[3] = d.Name, fmt.Sprintf("%.2f", d.Amount), d.Currency, d.CreatedAt
			csvWriter.Write(row)
		}
	}

}
