package openfga

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/openfga/go-sdk/client"
	openfga "github.com/openfga/go-sdk"
)

type FGAClient struct {
	Client  *client.OpenFgaClient
	Context context.Context
}

func NewFGAClient() (*FGAClient, error) {
	fgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiUrl:  "http://openfga:8080",
		StoreId: os.Getenv("FGA_STORE_ID"),
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	// Check if a store named "auth_app_store" exists
	listStoresResponse, err := fgaClient.ListStores(ctx).Execute()
	if err != nil {
		return nil, err
	}

	var store *openfga.Store
	for _, s := range listStoresResponse.GetStores() {
		if s.GetName() == "auth_app_store" {
			store = &s
			break
		}
	}

	if store != nil {
		log.Printf("OpenFGA Store ID: %s, Name: %s", store.GetId(), store.GetName())
		fgaClient.SetStoreId(store.GetId())
	} else {
		// If not, create a new store
		createStoreResponse, err := fgaClient.CreateStore(ctx).Body(client.ClientCreateStoreRequest{Name: "auth_app_store"}).Execute()
		if err != nil {
			return nil, err
		}
		fgaClient.SetStoreId(createStoreResponse.GetId())
		log.Printf("Created OpenFGA Store. ID: %s, Name: %s", createStoreResponse.GetId(), createStoreResponse.GetName())
	}

	// Read the authorization model from model.json
	jsonModelBytes, err := os.ReadFile("model.json")
	if err != nil {
		return nil, err
	}

	var authModel openfga.AuthorizationModel
	if err := json.Unmarshal(jsonModelBytes, &authModel); err != nil {
		return nil, err
	}

	writeModelResponse, err := fgaClient.WriteAuthorizationModel(ctx).Body(client.ClientWriteAuthorizationModelRequest{
		SchemaVersion:   authModel.SchemaVersion,
		TypeDefinitions: authModel.TypeDefinitions,
	}).Execute()
	if err != nil {
		return nil, err
	}

	fgaClient.SetAuthorizationModelId(writeModelResponse.GetAuthorizationModelId())
	log.Printf("Created OpenFGA Authorization Model with ID: %s", writeModelResponse.GetAuthorizationModelId())

	return &FGAClient{Client: fgaClient, Context: ctx}, nil
}
