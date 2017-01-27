package fabric_sdk_go

import (
	"testing"
)

//
// Peer via chain setPeer/getPeer
//
// Set the Perr URL through the chain setPeer method. Verify that the
// Peer URL was set correctly through the getPeer method. Repeat the
// process by updating the Peer URL to a different address.
//
func TestPeerViaChain(t *testing.T) {
	client := NewClient()
	chain, err := client.NewChain("testChain-peer")
	if err != nil {
		t.Fatalf("error from NewChain %v", err)
	}
	peer := CreateNewPeer("localhost:7050")
	chain.AddPeer(peer)

	peers := chain.GetPeers()
	if peers == nil || len(peers) != 1 || peers[0].GetUrl() != "localhost:7050" {
		t.Fatalf("Failed to retieve the new peers URL from the chain")
	}
	chain.RemovePeer(peer)
	peer2 := CreateNewPeer("localhost:7054")
	chain.AddPeer(peer2)
	peers = chain.GetPeers()

	if peers == nil || len(peers) != 1 || peers[0].GetUrl() != "localhost:7054" {
		t.Fatalf("Failed to retieve the new peers URL from the chain")
	}
}

//
// Peer via chain missing peer
//
// Attempt to send a request to the peer with the SendTransactionProposal method
// before the peer was set. Verify that an error is reported when tying
// to send the request.
//
func TestOrdererViaChainMissingOrderer(t *testing.T) {
	client := NewClient()
	chain, err := client.NewChain("testChain-peer")
	if err != nil {
		t.Fatalf("error from NewChain %v", err)
	}
	_, err = chain.SendTransactionProposal(nil, 0)
	if err == nil {
		t.Fatalf("SendTransactionProposal didn't return error")
	}
	if err.Error() != "peers is nil" {
		t.Fatalf("SendTransactionProposal didn't return right error")
	}
}

//
// Peer via chain nil data
//
// Attempt to send a request to the peers with the SendTransactionProposal method
// with the data set to null. Verify that an error is reported when tying
// to send null data.
//
func TestPeerViaChainNilData(t *testing.T) {
	client := NewClient()
	chain, err := client.NewChain("testChain-peer")
	if err != nil {
		t.Fatalf("error from NewChain %v", err)
	}
	peer := CreateNewPeer("localhost:7050")
	chain.AddPeer(peer)
	_, err = chain.SendTransactionProposal(nil, 0)
	if err == nil {
		t.Fatalf("SendTransaction didn't return error")
	}
	if err.Error() != "signedProposal is nil" {
		t.Fatalf("SendTransactionProposal didn't return right error")
	}
}
