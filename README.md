# Hyperledger Fabric Client SDK for Go

The Hyperledger Fabric Client SDK makes it easy to use APIs to interact with a Hyperledger Fabric blockchain.

This SDK is targeted both towards the external access to a Hyperledger Fabric blockchain using a Go application, as well as being targeted at the internal library in a peer to access API functions on other parts of the network.

## Build and Test

This project must be cloned into `$GOPATH/src/github.com/hyperledger`. Package names have been chosen to match the Hyperledger project.

Execute `go test` from the project root to build the library and run the basic headless tests.

Execute `go test` in the `integration_test` to run end-to-end tests. This requires you to have:
- A working fabric set up. Refer to the Hyperledger Fabric [documentation](https://github.com/hyperledger/fabric) on how to do this.
- The `example_cc` chaincode from the Node.js SDK deployed. Refer to the fabric-sdk-node [documentation](https://github.com/hyperledger/fabric-sdk-node) on how to install it and run the `end-to-end.js` which deploys the `example_cc`
- Customized settings in the `integration_test/test_resources/config/config_test.yaml` in case your Hyperledger Fabric network is not running on `localhost` or is using different ports.

## Work in Progress

This client was last tested and found to be compatible with the following Hyperledger Fabric commit levels:
- fabric: `f7c19f88e824cbaea3c55bc218b3bbed37cc29ad`
- fabric-ca: `1ec55b2b49e9dfbfc2e28dccec0ced659ce1f246`

The following SDK features are yet to be implemented:
- Chaincode deployment
- Chain initialization
