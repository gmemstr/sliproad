package files

import "fmt"

var ProviderConfig map[string]FileProvider
var Providers map[string]*FileProviderInterface

func TranslateProvider(codename string, i *FileProviderInterface) {
	provider := ProviderConfig[codename]
	if provider.Provider == "disk" {
		*i = &DiskProvider{provider,}
		return
	}

	if provider.Provider == "backblaze" {
		bbProv := &BackblazeProvider{provider, provider.Config["bucket"], ""}
		*i = bbProv
		return
	}

	if provider.Provider == "s3" {
		s3Prov := &S3Provider{
			FileProvider: provider,
			Region:       provider.Config["region"],
			Bucket:       provider.Config["bucket"],
			Endpoint:     "",
			KeyID:        "",
			KeySecret:    "",
		}
		if _, ok := provider.Config["endpoint"]; ok {
			s3Prov.Endpoint = provider.Config["endpoint"]
		}
		if _, ok := provider.Config["keyid"]; ok {
			s3Prov.KeyID = provider.Config["keyid"]
		}
		if _, ok := provider.Config["keysecret"]; ok {
			s3Prov.KeySecret = provider.Config["keysecret"]
		}

		*i = s3Prov
		return
	}

	*i = FileProvider{}
}

func SetupProviders() {
	Providers = make(map[string]*FileProviderInterface)
	for name, provider := range ProviderConfig {
		var i FileProviderInterface
		TranslateProvider(name, &i)
		success := i.Setup(provider.Config)
		if !success {
			fmt.Printf("%s failed to initialize\n", name)
		} else {
			Providers[name] = &i
			fmt.Printf("%s initialized successfully\n", name)
		}
	}
}
