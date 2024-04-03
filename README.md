# DMARCr

DMARCr is a command-line tool for reading and aggregating DMARC (Domain-based Message Authentication, Reporting, and Conformance) reports. It provides a simple way to parse DMARC XML reports and display aggregated data in a table format.

## Features

- Read DMARC reports from files or standard input.
- Aggregate data based on source IP, disposition, DKIM, and SPF results.
- Display aggregated data in a table format.

## Installation

To install DMARCr, use the following command:

```bash
go install github.com/dubyte/dmarcr@latest
```

## Usage

To read a DMARC report from a file:

```bash
dmarcr path/to/report.xml
```

To read DMARC reports from standard input:

```bash
cat *.xml | dmarcr
```

## Example Outputs

Assuming you have DMARC reports with data similar to the following:

- Source IP: `192.0.2.1`, Count: `2`, Disposition: `none`, DKIM: `pass`, SPF: `pass`
- Source IP: `203.0.113.5`, Count: `1`, Disposition: `reject`, DKIM: `fail`, SPF: `fail`

The output of `dmarcr` would look like this:

```bash
Source IP            Count Disposition DKIM  SPF  
--------------------------------------------------
192.0.2.1            2     none        pass  pass 
203.0.113.5          1     reject      fail  fail 
```

## Contributing

Contributions to DMARCr are welcome! Please feel free to open an issue or submit a pull request.

## License

DMARCr is released under the [GNU General Public License v3.0](LICENSE).
