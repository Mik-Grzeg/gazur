package common

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"log"
)

type Identity struct {
	*azidentity.ClientSecretCredential
}

func NewIdentity(tenantId string, clientId string, clientSecret string) Identity {
	creds, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	if err != nil {
		log.Fatalf("Unable to parse azure credentials: %v", err)
	}

	return Identity{
		creds,
	}
}