# Hyperledger Fabric Client SDK for golang

The Hyperledger Fabric Client SDK makes it easy to use APIs to interact with a Hyperledger Fabric blockchain.

This SDK is targeted both towards the external access to a Hyperledger Fabric blockchain using a golang application and the use as an internal library in a peer to access API functions on other parts of the network.

## Build and Test
Execute `go test` from the project root to build the library and run the basic headless tests.

Execute `go test` in the `integration_test` to run end-to-end tests. This requires you to have:
- a working fabric set up. Refer to the Hyperledger Fabric documentation on how to do this.
- the `example_cc` chaincode from the nodejs SDK deployed. Refer to the fabric-sdk-node documentation on how to install it and run the `end-to-end.js` which deploys the `example_cc`
- adjust settings in the `integration_test/test_resources/config/config_test.yaml` in case your Hyperledger Fabric network is not running on `localhost` or is using different ports.

