package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	NameLc string
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

type Duplicate struct {
	Date   string
	Name   string
	Amount float32
}

func (d Duplicate) String() []string {
	dString := make([]string, 3)
	dString[0] = d.Date
	dString[1] = d.Name
	dString[2] = fmt.Sprintf("%.2f", d.Amount)
	return dString
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
					NameLc: strings.ToLower(d[1]),
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

func WriteDuplicates(fn string, dups []Duplicate) {
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Unable to open file '%s' for writing.\n", fn)
		panic(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	for _, d := range dups {
		w.Write(d.String())
	}
}

func main() {
	var mm []MarsMile
	var dups []Duplicate
	donorTotal := map[string]float32{}
	donorRunning := map[string]float32{}
	mmCount := 1

	cfg := ReadConfigJson(configFileName)
	donations := ReadDonations(cfg.DonationsIn)

	// iterate donations to determine mars miles and donors
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

	// First scan the duplicates, record current if match found, skip the donations if found in duplicates.
	// If we make it to donations, scan from start to current - 1, record both duplicates, break inner loop
	// iterate donations (twice) to determine duplicates
	ld := len(donations)
	if ld > 2 {
		for i := 0; i < ld; i++ {
			// first compare against existing duplicates
			skipDupsCheck := false
			for j := 0; j < len(dups); j++ {
				if donations[i].Date == dups[j].Date && donations[i].NameLc == strings.ToLower(dups[j].Name) {
					newDup := Duplicate{
						Date:   donations[i].Date,
						Name:   donations[i].Name,
						Amount: donations[i].Amount,
					}
					dups = append(dups, newDup)
					skipDupsCheck = true
					break
				}
			}
			if !skipDupsCheck {
				for j := i + 1; j < ld; j++ {
					if donations[i].Date == donations[j].Date && donations[i].NameLc == donations[j].NameLc {
						newDup := Duplicate{
							Date:   donations[i].Date,
							Name:   donations[i].Name,
							Amount: donations[i].Amount,
						}
						dups = append(dups, newDup)
						break
					}
				}
			}
		}
	}
	/* for i := 0; i < len(donations)-1; i++ {
		for j := i + 1; j < len(donations); j++ {
			if donations[i].Date == donations[j].Date && donations[i].NameLc == donations[j].NameLc {
				newDup := Duplicate{
					Date: donations[i].Date,
					Name: donations[i].Name,
				}
				dups = append(dups, newDup)
			}
		}
	} */

	//fmt.Println("Writing file " + cfg.MilesOut)
	WriteMarsMiles(cfg.MilesOut, mm)
	//fmt.Println("Writing file " + cfg.DonorsOut)
	WriteDonors(cfg.DonorsOut, donorTotal)
	//fmt.Println("Writing file " + cfg.DuplicatesOut)
	WriteDuplicates(cfg.DuplicatesOut, dups)

	/* for _, m := range mm {
		fmt.Println(m)
	}
	for k, v := range donorTotal {
		fmt.Printf("%s: %f\n", k, v)
	} */
}
