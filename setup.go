package gcs

import (
	"encoding/json"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

// CreateClient creates a gcloud storage client
func CreateClient(serviceAccountDetails []byte) (*GCloudStorageAgent, error) {
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	var cfg config
	err = json.Unmarshal(serviceAccountDetails, &cfg)
	if err != nil {
		return nil, err
	}

	client := &GCloudStorageAgent{
		Client:         storageClient,
		ProjectID:      cfg.ProjectID,
		GoogleAccessID: cfg.ClientEmail,
		PrivateKey:     cfg.PrivateKey,
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
