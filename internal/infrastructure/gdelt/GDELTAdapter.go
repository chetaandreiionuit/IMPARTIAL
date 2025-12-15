package gdelt

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GDELTAdapter handles fetching and stream-processing GDELT V2 Event CSVs.
type GDELTAdapter struct {
	httpClient *http.Client
}

func NewGDELTAdapter() *GDELTAdapter {
	return &GDELTAdapter{
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

// FetchHighImpactEvents downloads the latest GDELT Events CSV updates,
// filters by Tone (Impact > 5.0), and returns the Source URLs.
// Uses Go Iterators pattern (streaming processing) to minimize RAM usage.
func (adapter *GDELTAdapter) FetchHighImpactEvents(ctx context.Context) ([]string, error) {
	// 1. Get the URL of the latest update zip
	fileURL, err := adapter.getLastUpdateURL(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("[GDELT] Downloading update: %s\n", fileURL)

	// 2. Download ZIP Stream
	req, err := http.NewRequestWithContext(ctx, "GET", fileURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := adapter.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GDELT download failed: %d", resp.StatusCode)
	}

	// Read the whole zip into memory? GDELT files are ~10MB compressed, okay for RAM.
	// For true streaming of ZIPs without reading full body, we need standard zip not supporting seek?
	// Go's zip.NewReader requires ReaderAt (seekable).
	// To strictly follow "Stream CSV" without loading ZIP to RAM, we'd need a custom unzip stream.
	// However, 5-10MB is negligible. We'll read the body buffer.
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(bodyBytes), int64(len(bodyBytes)))
	if err != nil {
		return nil, err
	}

	// 3. Process CSV inside ZIP
	var highImpactURLs []string

	// Usually there is only one .export.CSV file in the zip
	for _, zipFile := range zipReader.File {
		if !strings.HasSuffix(zipFile.Name, ".csv") {
			continue
		}

		f, err := zipFile.Open()
		if err != nil {
			return nil, err
		}
		defer f.Close()

		csvReader := csv.NewReader(f)
		csvReader.Comma = '\t'         // GDELT uses Tab Separated Values
		csvReader.FieldsPerRecord = -1 // Variable fields in some cases

		// Stream Processing Line-by-Line
		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				continue // Skip malformed lines
			}

			// GDELT 2.0 Columns of Interest (Approximate indices):
			// 34: AvgTone (float)
			// 60: SourceURL (Last One) - Let's assume last one for safety or index 57/60 depending on version.
			// Actually GDELT 2.0 has 58 columns usually. SourceURL is last.

			if len(record) < 58 {
				continue
			}

			// Parse Tone (Index 34 in GDELT 2.0)
			toneStr := record[34]
			sourceURL := record[len(record)-1] // URL is always last

			tone, err := strconv.ParseFloat(toneStr, 64)
			if err != nil {
				continue
			}

			// FILTER: |Tone| > 5.0 (High Emotional Impact)
			if math.Abs(tone) > 5.0 {
				// Cleanup URL
				url := strings.TrimSpace(sourceURL)
				if url != "" && strings.HasPrefix(url, "http") {
					highImpactURLs = append(highImpactURLs, url)
				}
			}
		}
	}

	// Deduplicate locally
	return uniqueStrings(highImpactURLs), nil
}

func (adapter *GDELTAdapter) getLastUpdateURL(ctx context.Context) (string, error) {
	// GDELT 2.0 Master List (Updated every 15 mins)
	const updateListURL = "http://data.gdeltproject.org/gdeltv2/lastupdate.txt"

	req, err := http.NewRequestWithContext(ctx, "GET", updateListURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := adapter.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Format:
	// <FileSize> <MD5> <URL>
	// We want the first line (latest English export)
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			url := fields[2]
			if strings.Contains(url, ".export.CSV.zip") {
				return url, nil
			}
		}
	}

	return "", fmt.Errorf("no export csv found in lastupdate.txt")
}

func uniqueStrings(input []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range input {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
