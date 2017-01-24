package cop

import (
	"fmt"

	"github.com/hyperledger/fabric-ca/api"
	cop "github.com/hyperledger/fabric-ca/lib"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("fabric_sdk_go")

type FabricCOPServices struct {
	fabricCOPClient *cop.Client
}

/**
 * @param {string} url The endpoint URL for Fabric COP services of the form: "http://host:port" or "https://host:port"
 */
func NewFabricCOPServices(url string, clientConfigFile string) (*FabricCOPServices, error) {
	if url == "" {
		return nil, fmt.Errorf("Failed to create FabricCopServices. Missing requirement 'url' parameter.")
	}
	copServer := fmt.Sprintf(`{"serverURL":"%s","homeDir":"%s"}`, url, clientConfigFile)
	c, err := cop.NewClient(copServer)
	if err != nil {
		return nil, fmt.Errorf("New fabricCOPClient failed: %s", err)
	}

	fcs := &FabricCOPServices{fabricCOPClient: c}
	logger.Infof("Constructed FabricCOPServices instance: %v", fcs)

	return fcs, nil
}

func (fcs *FabricCOPServices) Enroll(enrollmentID string, enrollmentSecret string) ([]byte, []byte, error) {
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
	id, err := fcs.fabricCOPClient.Enroll(req)
	if err != nil {
		return nil, nil, fmt.Errorf("Enroll failed: %s", err)
	}
	return id.GetECert().GetKey(), id.GetECert().GetCert(), nil
}
