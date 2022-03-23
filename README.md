# MarsMileCalculator
The MarsMileCalculator is used by the Mars Initiative to calculate the issuance of Mars Miles based on donations. For any individual donor, donations totalling $100 leads to the issuance of 1 Mars Mile certificate. 

The preferred version is Go, because it compiles to an executable and has no installation prerequisites. The Go version reads MarsMileCalculator.json from the same folder to determine the locations of the input donations CSV file, which consists of the following 3 columns:
- Date (string)
- Donor Name (string)
- Donation Amount (float32) 

The Date field is treated as a string and is not processed. The calculator assumes the file is already in order by date. Usually the CSV file is output from Excel or Sheets, so this is the default usage and is not an issue. 

## Useful make commands
- make clean (removes executables) 
- make (compile for current platform to the file MarsMileCalculator, appropriate for Mac or Linux)
- make compile (compile for Windows x64)
