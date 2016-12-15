package fabric_sdk_go

import (
	config "sk-git.securekey.com/vme/fabric-sdk-go/config"
)

type Chain struct {
	Name            string // Name of the chain is only meaningful to the client
	SecurityEnabled bool   // Security enabled flag
	Members         map[string]*Member
	TcertBatchSize  int // The number of tcerts to get in each batch
	Orderer         *Orderer
}

/**
 * @param {string} name to identify different chain instances. The naming of chain instances
 * is completely at the client application's discretion.
 */
func CreateNewChain(chainName string) Chain {
	m := make(map[string]*Member)
	return Chain{Name: chainName, SecurityEnabled: config.IsSecurityEnabled(), Members: m,
		TcertBatchSize: config.TcertBatchSize()}
}

func (c Chain) GetMember(memberName string) *Member {
	if val, ok := c.Members[memberName]; ok {
		return val
	}
	m := CreateNewMember(memberName, &c)
	c.Members[memberName] = m
	return m
}
