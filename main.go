package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dubyte/dmarcr/feedback"
	"github.com/dubyte/dmarcr/vcs"
)

var version = vcs.Version()

func main() {
	reports := feedback.New()

	switch {
	case len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "-v"):
		// If -version is provided, print the version and exit.
		fmt.Println(version)
		os.Exit(0)

	case len(os.Args) == 2 && (os.Args[1] == "--help" || os.Args[1] == "-h"):
		// Display the help message.
		displayHelp()
		os.Exit(0)

	case len(os.Args) == 2:
		// If a filename is provided, read the DMARC report from the file.
		err := reports.ReadFromFile(os.Args[1])
		if err != nil {
			log.Fatalf("Error reading from file: %v", err)
		}

	case len(os.Args) == 1:
		// If no arguments are provided, read the DMARC reports from standard input.
		err := reports.ReadFromStdin()
		if err != nil {
			log.Fatalf("Error reading from stdin: %v", err)
		}
	}

	fmt.Printf("%s\n", reports)

}

// displayHelp prints usage instructions and explanations for each term in the table.
func displayHelp() {
	fmt.Print(`
Usage: dmarcr <dmarc_report.xml>
       cat dmarc_report.xml | dmarcr

This tool reads a DMARC report in XML format and displays the relevant information in a table format. It can read from a file specified as an argument or from standard input if no file is provided.

Table columns:
  Source IP: The IP address from which the emails originated.
  Count: The number of emails received from this IP address.
  Disposition: The action taken on the emails based on the DMARC policy (none, quarantine, or reject).
  DKIM: The result of the DomainKeys Identified Mail (DKIM) authentication check (pass or fail).
  SPF: The result of the Sender Policy Framework (SPF) authentication check (pass or fail).
`)
}
