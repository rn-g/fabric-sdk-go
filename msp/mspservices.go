/*
Copyright SecureKey Technologies Inc. All Rights Reserved.


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at


      http://www.apache.org/licenses/LICENSE-2.0


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package msp

import (
	"fmt"

	"github.com/hyperledger/fabric-ca/api"
	msp "github.com/hyperledger/fabric-ca/lib"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("fabric_sdk_go")

// Services ...
type Services struct {
	mspClient *msp.Client
}

// NewMSPServices ...
/**
 * @param {string} url The endpoint URL for msp services of the form: "http://host:port" or "https://host:port"
 */
func NewMSPServices(url string, clientConfigFile string) (*Services, error) {
	if url == "" {
		return nil, fmt.Errorf("Failed to create MSPServices. Missing requirement 'url' parameter.")
	}
	mspServer := fmt.Sprintf(`{"serverURL":"%s","homeDir":"%s"}`, url, clientConfigFile)
	c, err := msp.NewClient(mspServer)
	if err != nil {
		return nil, fmt.Errorf("New mspClient failed: %s", err)
	}

	msps := &Services{mspClient: c}
	logger.Infof("Constructed MSPServices instance: %v", msps)

	return msps, nil
}

// Enroll ...
/**
 * Enroll a registered user in order to receive a signed X509 certificate
 * @param {string} enrollmentID The registered ID to use for enrollment
 * @param {string} enrollmentSecret The secret associated with the enrollment ID
 * @returns {[]byte} X509 certificate
 * @returns {[]byte} private key
 */
func (msps *Services) Enroll(enrollmentID string, enrollmentSecret string) ([]byte, []byte, error) {
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
