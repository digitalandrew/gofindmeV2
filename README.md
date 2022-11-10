gofindmeV2 - GoFundMe OSINT tool now faster in GO

gofindme is a quick to use webscraping tool to quickly pull down a list of donators and donation amounts from gofundme campaigns.

I initially created this tool for use in TraceLabs missing person CTFs as I found GoFundMe campaigns to give lots of leads however manually scrolling through the donations list to populate the full list was time consuming.
Installation

# Installation Instructions
Clone the repo
`$ git clone https://github.com/digitalandrew/gofindmeV2`

Change the working directory to gofindme
`$ cd gofindmeV2`

# Usage Example

`$ go run main.go -c https://www.gofundme.com/f/digitalandrews-fundraising-campaign -f`

Alternatively you can compile it for your desired OS and run. 

# Flags
Required Arguments:

 -c         Specifies campaign URL
 
Optional Arguments:
 
 -f         Write Name, Donation Amount and Time of donation to csv file
 -h         Help
 
