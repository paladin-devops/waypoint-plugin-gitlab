package registry

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	gitlab "github.com/xanzy/go-gitlab"
)

type RegistryConfig struct {
	Name      string `hcl:"name,attr"`
	Version   string `hcl:"version,attr"`
	ProjectID int    `hcl:"project_id,attr"`
	FileName  string `hcl:"file_name,attr"`
	Token     string `hcl:"token,attr"`
}

type Registry struct {
	config RegistryConfig
}

// Implement Configurable
func (r *Registry) Config() (interface{}, error) {
	return &r.config, nil
}

// Implement ConfigurableNotify
func (r *Registry) ConfigSet(config interface{}) error {
	c, ok := config.(*RegistryConfig)
	if !ok {
		// The Waypoint SDK should ensure this never gets hit
		return fmt.Errorf("Expected *RegisterConfig as parameter")
	}

	// validate the config
	if c.Name == "" {
		return fmt.Errorf("Name must be set to a valid directory")
	}

	return nil
}

// Implement Registry
func (r *Registry) PushFunc() interface{} {
	// return a function which will be called by Waypoint
	return r.push
}

// Implement Registry
func (r *Registry) AccessInfoFunc() interface{} {
	return r.accessInfo
}

func (r *Registry) accessInfo() (*AccessInfo, error) {
	return &AccessInfo{}, nil
}

// A PushFunc does not have a strict signature, you can define the parameters
// you need based on the Available parameters that the Waypoint SDK provides.
// Waypoint will automatically inject parameters as specified
// in the signature at run time.
//
// Available input parameters:
// - context.Context
// - *component.Source
// - *component.JobInfo
// - *component.DeploymentConfig
// - hclog.Logger
// - terminal.UI
// - *component.LabelSet
//
// In addition to default input parameters the builder.Binary from the Build step
// can also be injected.
//
// The output parameters for PushFunc must be a Struct which can
// be serialzied to Protocol Buffers binary format and an error.
// This Output Value will be made available for other functions
// as an input parameter.
// If an error is returned, Waypoint stops the execution flow and
// returns an error to the user.
func (r *Registry) push(ctx context.Context, ui terminal.UI) (*Artifact, error) {
	u := ui.Status()
	defer u.Close()
	u.Update("Pushing binary to registry")

	git, err := gitlab.NewClient(r.config.Token)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Publishing artifact returns GenericPackagesFile
	gpf, _, err := git.GenericPackages.PublishPackageFile(r.config.ProjectID, r.config.Name, r.config.Version, r.config.FileName, strings.NewReader(""), &gitlab.PublishPackageFileOptions{})
	fmt.Println(gpf)
	if err != nil {
		log.Fatalf("Failed to push package: %v", err)
	}

	return &Artifact{}, nil
}

// Implement Authenticator
