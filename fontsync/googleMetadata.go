package fontsync

import (
	"encoding/json"
	"fmt"
	"io"
	"jinya-fonts/database"
	"net/http"
	"strings"
)

type GoogleWebfontMetadata struct {
	License     string              `json:"license"`
	Designers   []database.Designer `json:"designers"`
	Description string              `json:"description"`
	Category    string              `json:"category"`
}

func getGoogleWebfontMetadata(cpu int, family string) (*GoogleWebfontMetadata, error) {
	logWithCpu(cpu, "Download font metadata for font %s", family)
	res, err := http.Get(fmt.Sprintf("https://fonts.google.com/metadata/fonts/%s", family))
	if err != nil {
		logWithCpu(cpu, "Failed to download font metadata for font %s", cpu, family)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		logWithCpu(cpu, "Failed to download font metadata for font %s", cpu, family)
		return nil, fmt.Errorf("failed to get metadata")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		logWithCpu(cpu, "Failed to read response body for font %s", cpu, family)
		return nil, err
	}

	bodyString := string(bodyBytes)
	bodyString = strings.TrimPrefix(bodyString, ")]}'")

	var metadata GoogleWebfontMetadata
	err = json.Unmarshal([]byte(bodyString), &metadata)

	return &metadata, err
}
