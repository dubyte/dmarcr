package feedback

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// Feedback represents the structure of a DMARC report.
type Feedback struct {
	XMLName xml.Name `xml:"feedback"`
	Reports []Record `xml:"record"`
}

// Record represents a single record in the DMARC report.
type Record struct {
	Row         Row         `xml:"row"`
	Identifiers Identifiers `xml:"identifiers"`
	AuthResults AuthResults `xml:"auth_results"`
}

// Row contains information about the source IP, email count, and policy evaluation.
type Row struct {
	SourceIP        string          `xml:"source_ip"`
	Count           int             `xml:"count"`
	PolicyEvaluated PolicyEvaluated `xml:"policy_evaluated"`
}

// PolicyEvaluated contains the results of the DMARC policy evaluation.
type PolicyEvaluated struct {
	Disposition string `xml:"disposition"`
	DKIM        string `xml:"dkim"`
	SPF         string `xml:"spf"`
}

// Identifiers contains the identifiers for the email (e.g., the "From" header).
type Identifiers struct {
	HeaderFrom string `xml:"header_from"`
}

// AuthResults contains the authentication results for DKIM and SPF.
type AuthResults struct {
	DKIM AuthResult `xml:"dkim"`
	SPF  AuthResult `xml:"spf"`
}

// AuthResult represents the result of an authentication check.
type AuthResult struct {
	Domain string `xml:"domain"`
	Result string `xml:"result"`
}

type AggregatedReports struct {
	Feedbacks []Feedback
}

func New() *AggregatedReports {
	return &AggregatedReports{}
}

func (r *AggregatedReports) ReadFromFile(filename string) error {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	var feedback Feedback
	if err := xml.Unmarshal(fileData, &feedback); err != nil {
		return fmt.Errorf("error unmarshalling XML: %v", err)
	}

	r.Feedbacks = append(r.Feedbacks, feedback)
	return nil
}

func (r *AggregatedReports) ReadFromStdin() error {
	stdinData, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("error reading from stdin: %v", err)
	}

	// Split the input data into individual XML documents using a regular expression.
	xmlDeclPattern := `\<\?xml version="1.0"( encoding="UTF-8")?\s*\?\>`
	re := regexp.MustCompile(xmlDeclPattern)
	xmlDocs := re.Split(string(stdinData), -1)

	for _, xmlDoc := range xmlDocs {
		if strings.TrimSpace(xmlDoc) == "" {
			continue
		}
		xmlDoc = `<?xml version="1.0" encoding="UTF-8" ?>` + xmlDoc

		var feedback Feedback
		if err := xml.Unmarshal([]byte(xmlDoc), &feedback); err != nil {
			return fmt.Errorf("error unmarshalling XML: %v", err)
		}
		r.Feedbacks = append(r.Feedbacks, feedback)
	}

	return nil
}

func (r *AggregatedReports) String() string {
	// Define a struct to hold the key fields for aggregation.
	type key struct {
		SourceIP    string
		Disposition string
		DKIM        string
		SPF         string
	}

	// Create a map to hold the aggregated counts.
	aggregatedCounts := make(map[key]int)

	// Populate the map with aggregated counts.
	for _, feedback := range r.Feedbacks {
		for _, report := range feedback.Reports {
			k := key{
				SourceIP:    report.Row.SourceIP,
				Disposition: report.Row.PolicyEvaluated.Disposition,
				DKIM:        report.Row.PolicyEvaluated.DKIM,
				SPF:         report.Row.PolicyEvaluated.SPF,
			}
			aggregatedCounts[k] += report.Row.Count
		}
	}

	// Build the output string in table format.
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%-20s %-5s %-10s %-5s %-5s\n", "Source IP", "Count", "Disposition", "DKIM", "SPF"))
	builder.WriteString(strings.Repeat("-", 50) + "\n")
	for k, count := range aggregatedCounts {
		builder.WriteString(fmt.Sprintf("%-20s %-5d %-10s %-5s %-5s\n",
			k.SourceIP,
			count,
			k.Disposition,
			k.DKIM,
			k.SPF))
	}

	return builder.String()
}
