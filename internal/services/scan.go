package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kai-sec/internal/daos"
	"kai-sec/internal/daos/models"
	"kai-sec/internal/dtos"
	"kai-sec/internal/logger"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

func fetchFileFromGitHub(repo string, file string) ([]byte, error) {
	l := logger.GetLogger()
	// Ensures if the git repo URL is in raw format
	if !strings.HasPrefix(repo, "https://raw.githubusercontent.com/") {
		err := errors.New("invalid repo URL. Must be a raw.githubusercontent.com link")
		l.Error("Invalid repository URL", zap.String("repo", repo), zap.Error(err))
		return nil, err
	}

	// Constructing the full URL
	url := fmt.Sprintf("%s/%s", strings.TrimRight(repo, "/"), file)
	l.Info("Fetching file from GitHub", zap.String("url", url))

	var resp *http.Response
	var err error

	for attempt := 1; attempt <= 2; attempt++ {
		resp, err = http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			l.Info("Successfully fetched file", zap.String("file", file), zap.Int("status_code", resp.StatusCode))
			return io.ReadAll(resp.Body)
		}

		l.Warn("Failed to fetch file", zap.Int("attempt", attempt), zap.String("file", file), zap.Int("status_code", resp.StatusCode), zap.Error(err))
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
	l := logger.GetLogger()
	data, err := fetchFileFromGitHub(repo, file)
	if err != nil {
		l.Error("Failed to fetch file", zap.String("file", file), zap.Error(err))
		return err
	}

	var scans []dtos.ScanReport
	err = json.Unmarshal(data, &scans)
	if err != nil {
		l.Error("Failed to parse JSON", zap.String("file", file), zap.Error(err))
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	for _, scan := range scans {
		//Ensure scan metadata fields have defaults if missing
		if scan.ScanResults.ScanMetadata.ScannerVersion == "" {
			scan.ScanResults.ScanMetadata.ScannerVersion = "unknown"
		}
		if scan.ScanResults.ScanMetadata.PoliciesVersion == "" {
			scan.ScanResults.ScanMetadata.PoliciesVersion = "unknown"
		}
		if scan.ScanResults.ScanMetadata.ScanningRules == nil {
			scan.ScanResults.ScanMetadata.ScanningRules = []string{}
		}
		if scan.ScanResults.ScanMetadata.ExcludedPaths == nil {
			scan.ScanResults.ScanMetadata.ExcludedPaths = []string{}
		}
		scanDAO := convertToDAO(scan.ScanResults)
		err = daos.UpsertScan(scanDAO)
		if err != nil {
			l.Error("Failed to store scan in DB", zap.String("scan_id", scanDAO.ScanId), zap.Error(err))
			return fmt.Errorf("failed to store scan in DB: %w", err)
		}
	}

	l.Info("Successfully processed file", zap.String("file", file))
	return err
}

// Process JSON files concurrently
func ProcessScan(repo string, files []string) error {
	l := logger.GetLogger()
	if repo == "" || len(files) == 0 {
		err := errors.New("repo and files are required")
		l.Error("Missing repo or files", zap.Error(err))
		return errors.New("repo and files are required")
	}

	l.Info("Starting scan processing", zap.String("repo", repo), zap.Int("file_count", len(files)))

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
			workerID := 0
			for file := range fileChan {
				workerID++
				err := processFile(repo, file)
				if err != nil {
					log.Printf("Error processing file %s: %v\n", file, err)
					l.Error("error processing file", zap.String("file", file), zap.Int("worker", workerID), zap.Error(err))
				}
			}
		}()
	}

	wg.Wait()
	l.Info("Completed processing all files")
	return nil
}
