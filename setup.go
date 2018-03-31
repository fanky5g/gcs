package gcs

import (
	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

// CreateClient creates a gcloud storage client
func CreateClient(serviceAccountPath string) (*GCloudStorageAgent, error) {
	ctx := context.Background()
	var storageClient *storage.Client

	if serviceAccountPath != "" {
		c, err := storage.NewClient(ctx, option.WithServiceAccountFile(serviceAccountPath))

		if err != nil {
			return nil, err
		}

		storageClient = c
	} else {
		c, err := storage.NewClient(ctx) //application default credentials

		if err != nil {
			return nil, err
		}

		storageClient = c
	}

	client := &GCloudStorageAgent{
		Client: storageClient,
	}

	return client, nil
}

/*
	eg serviceAccountDetails
	{
	  "type": "",
	  "project_id": "",
	  "private_key_id": "",
	  "private_key": "",
	  "auth_uri": "",
	  "token_uri": "",
		"client_email": "silverbird@silverbird-192810.iam.gserviceaccount.com",
	  "client_id": "115772240233797321241",
	  "auth_provider_x509_cert_url": "",
	  "client_x509_cert_url": "",
	}
*/

type config struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}
