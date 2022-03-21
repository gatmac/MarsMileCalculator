package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

const (
	configFileName = "MarsMileCalculator.json"
)

type ConfigFile struct {
	Marsmile      float32 `json:"Marsmile"`
	DonationsIn   string  `json:"Donations_In"`
	DonorsOut     string  `json:"Donors_Out"`
	MilesOut      string  `json:"Miles_Out"`
	DuplicatesOut string  `json:"Duplicates_Out"`
}

type Donation struct {
	Date   string
	Name   string
	Amount float32
}

type MarsMile struct {
	Mile  int
	Date  string
	Donor string
}

func (mm MarsMile) String() []string {
	mString := make([]string, 3)
	mString[0] = fmt.Sprintf("%d", mm.Mile)
	mString[1] = mm.Date
	mString[2] = mm.Donor
	return mString
}

func ReadConfigJson(fn string) ConfigFile {
	var config ConfigFile
	yfile, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(yfile, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func ReadDonations(fn string) []Donation {
	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	c := csv.NewReader(f)
	donations, err := c.ReadAll()
	if err != nil {
		panic(err)
	}
	donArray := make([]Donation, 0, len(donations))
	for i, d := range donations {
		if i > 0 {
			if len(d[0]) > 0 {
				dAmount, err := strconv.ParseFloat(d[2], 32)
				if err != nil {
					fmt.Printf("Unable to convert '%s' to float.", d[2])
					panic(err)
				}
				don := Donation{
					Date:   d[0],
					Name:   d[1],
					Amount: float32(dAmount),
				}
				donArray = append(donArray, don)
			}
		}
	}
	return donArray
}

func WriteMarsMiles(fn string, mm []MarsMile) {
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Unable to open file '%s' for writing.\n", fn)
		panic(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	for _, m := range mm {
		w.Write(m.String())
	}
}

func WriteDonors(fn string, donors map[string]float32) {
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Unable to open file '%s' for writing.\n", fn)
		panic(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	for d, v := range donors {
		donor := make([]string, 2)
		donor[0] = d
		donor[1] = fmt.Sprintf("%.2f", v)
		w.Write(donor)
	}
}

func main() {
	var mm []MarsMile
	donorTotal := map[string]float32{}
	donorRunning := map[string]float32{}
	mmCount := 1

	cfg := ReadConfigJson(configFileName)
	donations := ReadDonations(cfg.DonationsIn)
	for _, d := range donations {
		donorTotal[d.Name] += d.Amount
		donorRunning[d.Name] += d.Amount
		for donorRunning[d.Name] >= cfg.Marsmile {
			newMm := MarsMile{
				Mile:  mmCount,
				Donor: d.Name,
				Date:  d.Date,
			}
			mm = append(mm, newMm)
			mmCount += 1
			donorRunning[d.Name] -= cfg.Marsmile
		}
	}
	fmt.Println("Writing file " + cfg.MilesOut)
	WriteMarsMiles(cfg.MilesOut, mm)
	fmt.Println("Writing file " + cfg.DonorsOut)
	WriteDonors(cfg.DonorsOut, donorTotal)
	/* for _, m := range mm {
		fmt.Println(m)
	}
	for k, v := range donorTotal {
		fmt.Printf("%s: %f\n", k, v)
	} */
}
