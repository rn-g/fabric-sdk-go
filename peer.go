package fabric_sdk_go

import (
	"time"

	pb "github.com/hyperledger/fabric/protos/peer"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	config "sk-git.securekey.com/vme/fabric-sdk-go/config"
)

type Peer struct {
	Url            string
	GrpcDialOption []grpc.DialOption
}

func CreateNewPeer(url string) *Peer {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTimeout(time.Second*3))
	if config.IsTlsEnabled() {
		creds := credentials.NewClientTLSFromCert(config.GetTlsCACertPool(), config.GetTlsServerHostOverride())
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	return &Peer{Url: url, GrpcDialOption: opts}
}

func (p *Peer) sendProposal(signedProposal *pb.SignedProposal) (*pb.ProposalResponse, error) {
	conn, err := grpc.Dial(p.Url, p.GrpcDialOption...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	endorserClient := pb.NewEndorserClient(conn)
	proposalResponse, err := endorserClient.ProcessProposal(context.Background(), signedProposal)
	if err != nil {
		return nil, err
	}
	return proposalResponse, nil
}
