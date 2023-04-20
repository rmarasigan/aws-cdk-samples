package awswrapper

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
)

// AWSLogs contains the CloudWatch Logs
// message event data that is a Base64-encoded.gzip
// file archive.
type AWSLogs struct {
	Data string `json:"data"`
}

// CloudWatchEvent is the CloudWatch Logs
// message event data.
type CloudWatchEvent struct {
	AWSLogs AWSLogs `json:"awslogs"`
}

// LogEvents contains the log event
// information.
type LogEvents struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}

// CloudWatchData is the CloudWatch Logs message
// data (decoded).
type CloudWatchData struct {
	MessageType         string      `json:"messageType"`
	Owner               string      `json:"owner"`
	LogGroup            string      `json:"logGroup"`
	LogStream           string      `json:"logStream"`
	SubscriptionFilters []string    `json:"subscriptionFilters"`
	LogEvents           []LogEvents `json:"logEvents"`
}

// DecodeData returns the decoded CloudWatch Logs Data.
func (cw *CloudWatchEvent) DecodeData() (*CloudWatchData, error) {
	var (
		data = new(CloudWatchData)
	)

	// Decode the Base64-encoded data
	decoded, err := base64.StdEncoding.DecodeString(cw.AWSLogs.Data)
	if err != nil {
		return nil, err
	}

	// Decompress the data by using the decoded data
	buffer := bytes.NewReader(decoded)
	reader, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}

	// Decode the decompressed data into the CloudWatchData struct
	err = json.NewDecoder(reader).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
