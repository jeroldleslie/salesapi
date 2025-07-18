package importer

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	CSVPath  string
	Interval int
}

func InitConfig() (*Config, error) {
	config := &Config{
		CSVPath:  "sales.csv",
		Interval: 24,
	}

	intervalStr := strings.TrimSpace(os.Getenv("IMPORTER_INTERVAL_HOURS"))
	if intervalStr != "" {
		invl, err := strconv.Atoi(intervalStr)
		if err != nil {
			return nil, err
		}

		config.Interval = invl
	}

	csvPath := strings.TrimSpace(os.Getenv("IMPORTER_CSV_PATH"))
	if csvPath != "" {
		config.CSVPath = csvPath
	}

	return config, nil
}
