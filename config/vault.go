package config

import (
	"context"
	"encoding/json"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/spf13/viper"
	"log"
)

type VaultParameters struct {
	// connection parameters
	address             string
	approleRoleID       string
	approleSecretIDFile string

	apiKeyPath              string
	apiKeyMountPath         string
	apiKeyField             string
	databaseCredentialsPath string
}

type Vault struct {
	client     *vault.Client
	parameters VaultParameters
}

func VaultSecrets() (*Config, error) {

	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = viper.GetString("VAULT_ADDR")

	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		return nil, err
	}

	client.SetToken(viper.GetString("VAULT_AUTH_TOKEN"))

	// Read a secret from the default mount path for KV v2 in dev mode, "secret"
	secret, err := client.KVv2("secret").Get(context.Background(), viper.GetString("VAULT_SECRET_PATH"))
	if err != nil {
		log.Fatalf("unable to read secret: %v", err)
	}

	j, _ := json.Marshal(secret.Data)

	config := &Config{}

	err = json.Unmarshal(j, config)
	if err != nil {
		log.Fatalf("unable to parse secrets %v", err)
	}

	return config, nil
}

func NewVaultAppRoleClient(ctx context.Context, parameters VaultParameters) (*Vault, *vault.Secret, error) {
	log.Printf("connecting to vault @ %s", parameters.address)

	config := vault.DefaultConfig() // modify for more granular configuration
	config.Address = parameters.address

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize vault client: %w", err)
	}

	vault := &Vault{
		client:     client,
		parameters: parameters,
	}

	token, err := vault.login(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("vault login error: %w", err)
	}

	log.Println("connecting to vault: success!")

	return vault, token, nil
}

func (v *Vault) login(ctx context.Context) (*vault.Secret, error) {
	log.Printf("logging in to vault with approle auth; role id: %s", v.parameters.approleRoleID)

	approleSecretID := &approle.SecretID{
		FromFile: v.parameters.approleSecretIDFile,
	}

	appRoleAuth, err := approle.NewAppRoleAuth(
		v.parameters.approleRoleID,
		approleSecretID,
		// approle.WithWrappingToken(),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize approle authentication method: %w", err)
	}

	authInfo, err := v.client.Auth().Login(ctx, appRoleAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to login using approle auth method: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no approle info was returned after login")
	}

	log.Println("logging in to vault with approle auth: success!")

	return authInfo, nil
}

func (v *Vault) GetSecretAPIKey(ctx context.Context) (string, error) {
	log.Println("getting secret api key from vault")

	secret, err := v.client.KVv2(v.parameters.apiKeyMountPath).Get(ctx, v.parameters.apiKeyPath)
	if err != nil {
		return "", fmt.Errorf("unable to read secret: %w", err)
	}

	apiKey, ok := secret.Data[v.parameters.apiKeyField]
	if !ok {
		return "", fmt.Errorf("the secret retrieved from vault is missing %q field", v.parameters.apiKeyField)
	}

	apiKeyString, ok := apiKey.(string)
	if !ok {
		return "", fmt.Errorf("unexpected secret key type for %q field", v.parameters.apiKeyField)
	}

	log.Println("getting secret api key from vault: success!")

	return apiKeyString, nil
}

// func (v *Vault) GetDatabaseCredentials(ctx context.Context) (DatabaseCredentials, *vault.Secret, error) {
// 	log.Println("getting temporary database credentials from vault")

// 	lease, err := v.client.Logical().ReadWithContext(ctx, v.parameters.databaseCredentialsPath)
// 	if err != nil {
// 		return DatabaseCredentials{}, nil, fmt.Errorf("unable to read secret: %w", err)
// 	}

// 	b, err := json.Marshal(lease.Data)
// 	if err != nil {
// 		return DatabaseCredentials{}, nil, fmt.Errorf("malformed credentials returned: %w", err)
// 	}

// 	var credentials DatabaseCredentials

// 	if err := json.Unmarshal(b, &credentials); err != nil {
// 		return DatabaseCredentials{}, nil, fmt.Errorf("unable to unmarshal credentials: %w", err)
// 	}

// 	log.Println("getting temporary database credentials from vault: success!")

// 	return credentials, lease, nil
// }
