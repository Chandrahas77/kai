package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kai-sec/internal/daos"
	"kai-sec/internal/daos/models"
	"kai-sec/internal/dtos"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func fetchFileFromGitHub(repo string, file string) ([]byte, error) {
	// Ensures if the repo URL is in raw format
	if !strings.HasPrefix(repo, "https://raw.githubusercontent.com/") {
		return nil, errors.New("invalid repo URL. Must be a raw.githubusercontent.com link")
	}

	// Constructing the full URL
	url := fmt.Sprintf("%s/%s", strings.TrimRight(repo, "/"), file)
	log.Printf("Fetching: %s\n", url)

	var resp *http.Response
	var err error

	for attempt := 1; attempt <= 2; attempt++ {
		resp, err = http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			return io.ReadAll(resp.Body)
		}

		log.Printf("Attempt %d: Failed to fetch %s, retrying...\n", attempt, url)
		//waiting for 2 secs before retrying
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to fetch file from GitHub: %w", err)
}

// Convert dtos to daos before inserting into DB
func convertToDAO(scan dtos.ScanResults) models.ScanResultsDAO {
	dao := models.ScanResultsDAO{
		ScanId:       scan.ScanId,
		Timestamp:    scan.Timetamp,
		ScanStatus:   scan.ScanStatus,
		ResourceType: scan.ResourceType,
		ResourceName: scan.ResourceName,
	}

	for _, v := range scan.Vulnerabilities {
		dao.Vulnerabilities = append(dao.Vulnerabilities, models.VulnerabilityDAO{
			ID:             v.ID,
			Severity:       v.Severity,
			CVSS:           v.CVSS,
			Status:         v.Status,
			PackageName:    v.PackageName,
			CurrentVersion: v.CurrentVersion,
			FixedVersion:   v.FixedVersion,
			Description:    v.Description,
			PublishedDate:  v.PublishedDate,
			Link:           v.Link,
			RiskFactors:    v.RiskFactors,
		})
	}

	dao.Summary = models.ScanSummaryDAO{
		TotalVulnerabilities: scan.Summary.TotalVulnerabilities,
		SeverityCounts: models.SeverityCountsDAO{
			Critical: scan.Summary.SeverityCounts.Critical,
			High:     scan.Summary.SeverityCounts.High,
			Medium:   scan.Summary.SeverityCounts.Medium,
			Low:      scan.Summary.SeverityCounts.Low,
		},
		FixableCount: scan.Summary.FixableCount,
		Compliant:    scan.Summary.Complaint,
	}

	dao.ScanMetadata = models.ScanMetadataDAO{
		ScannerVersion:  scan.ScanMetadata.ScannerVersion,
		PoliciesVersion: scan.ScanMetadata.PoliciesVersion,
		ScanningRules:   scan.ScanMetadata.ScanningRules,
		ExcludedPaths:   scan.ScanMetadata.ExcludedPaths,
	}

	return dao
}

// process a single file to fetch,parse and store into DB
func processFile(repo string, file string) error {
	data, err := fetchFileFromGitHub(repo, file)
	if err != nil {
		return err
	}

	var scans []dtos.ScanReport
	err = json.Unmarshal(data, &scans)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	for _, scan := range scans {
		scanDAO := convertToDAO(scan.ScanResults)
		err = daos.InsertScan(scanDAO)
		if err != nil {
			return fmt.Errorf("failed to store scan in DB: %w", err)
		}
	}

	log.Printf("Successfully processed %s\n", file)
	return err
}

// Process JSON files concurrently
func ProcessScan(repo string, files []string) error {
	if repo == "" || len(files) == 0 {
		return errors.New("repo and files are required")
	}

	var wg sync.WaitGroup
	fileChan := make(chan string, len(files))

	// Add files to queue
	for _, file := range files {
		fileChan <- file
	}
	close(fileChan)

	// Process files in parallel (3 workers)
	//TODO: Add the maxWorkers const to env or make it constant
	maxWorkers := 3
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileChan {
				err := processFile(repo, file)
				if err != nil {
					log.Printf("Error processing file %s: %v\n", file, err)
				}
			}
		}()
	}

	wg.Wait()
	return nil
}
