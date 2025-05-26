package utility

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type UploadImageRequest struct {
	Filename string `json:"filename"`
	Filedata string `json:"filedata"`
}

type UploadImageResponse struct {
	URL string `json:"url"`
}

func UploadBase64Image(filename, filedata string) (string, error) {
	payload := UploadImageRequest{
		Filename: filename,
		Filedata: filedata,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api-k4g7wd2cpq-uc.a.run.app/upload", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-bb", "apithxlne")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload failed with status: %v", resp.Status)
	}

	var uploadResp UploadImageResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return "", err
	}

	return uploadResp.URL, nil
}
