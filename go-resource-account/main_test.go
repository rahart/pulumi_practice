package main

import (
  "sync"
  "testing"
  "github.com/pulumi/pulumi/sdk/v3/go/common/resource"
  "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
  "github.com/stretchr/testify/assert"
)
// "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
type mocks int

func (mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return args.Name + "_id", args.Inputs, nil
}

func (mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return args.Args, nil
}

func TestBuildInfra(t *testing.T){
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		infra, err := buildInfra(ctx)
		assert.NoError(t, err)
		var wg sync.WaitGroup
		wg.Add(2)
		// Test if the service has tags and a name tag.
		pulumi.All(infra.resourceGroup.URN(), infra.resourceGroup.Tags).ApplyT(func(all []interface{}) error {
			urn := all[0].(pulumi.URN)
			tags := all[1].(map[string]string)
      assert.Containsf(t, tags, "Name", "missing a Name tag on resourceGroup %v", urn)
			wg.Done()
			return nil
		})

    pulumi.All(infra.storageAccount.URN(), infra.storageAccount.Tags).ApplyT(func(all []interface{}) error {
			urn := all[0].(pulumi.URN)
			tags := all[1].(map[string]string)
      assert.Containsf(t, tags, "Name", "missing a Name tag on storageAccount %v", urn)
			wg.Done()
			return nil
		})

		wg.Wait()
		return nil
  }, pulumi.WithMocks("project", "stack", mocks(0)))
	assert.NoError(t, err)
}
