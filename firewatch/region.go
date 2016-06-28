package firewatch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const metadataUrl string = "http://169.254.169.254/latest/dynamic/instance-identity/document"

type EC2Metadata struct {
	Region string
}

func getRegion() (string, error) {
	if os.Getenv("AWS_REGION") != "" {
		return os.Getenv("AWS_REGION"), nil
	}

	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(metadataUrl)
	if err != nil {
		return "", fmt.Errorf("Timeout while trying to auto-detect AWS_REGION: %v", err)
	}
	defer resp.Body.Close()

	md := &EC2Metadata{}
	err = json.NewDecoder(resp.Body).Decode(md)
	if err != nil {
		return "", err
	}
	return md.Region, nil
}
