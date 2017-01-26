package msp

import (
	"fmt"

	"github.com/hyperledger/fabric-ca/api"
	msp "github.com/hyperledger/fabric-ca/lib"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("fabric_sdk_go")

type MSPServices struct {
	mspClient *msp.Client
}

/**
 * @param {string} url The endpoint URL for msp services of the form: "http://host:port" or "https://host:port"
 */
func NewMSPServices(url string, clientConfigFile string) (*MSPServices, error) {
	if url == "" {
		return nil, fmt.Errorf("Failed to create MSPServices. Missing requirement 'url' parameter.")
	}
	mspServer := fmt.Sprintf(`{"serverURL":"%s","homeDir":"%s"}`, url, clientConfigFile)
	c, err := msp.NewClient(mspServer)
	if err != nil {
		return nil, fmt.Errorf("New mspClient failed: %s", err)
	}

	msps := &MSPServices{mspClient: c}
	logger.Infof("Constructed MSPServices instance: %v", msps)

	return msps, nil
}

func (msps *MSPServices) Enroll(enrollmentID string, enrollmentSecret string) ([]byte, []byte, error) {
	if enrollmentID == "" {
		return nil, nil, fmt.Errorf("enrollmentID is empty")
	}
	if enrollmentSecret == "" {
		return nil, nil, fmt.Errorf("enrollmentSecret is empty")
	}
	req := &api.EnrollmentRequest{
		Name:   enrollmentID,
		Secret: enrollmentSecret,
	}
	id, err := msps.mspClient.Enroll(req)
	if err != nil {
		return nil, nil, fmt.Errorf("Enroll failed: %s", err)
	}
	return id.GetECert().GetKey(), id.GetECert().GetCert(), nil
}
