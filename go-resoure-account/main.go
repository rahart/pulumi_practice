package main

import (
  "fmt"
	"github.com/pulumi/pulumi-azure-native-sdk/resources"
	"github.com/pulumi/pulumi-azure-native-sdk/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
  "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
    infra, err := buildInfra(ctx)
    if err != nil {
      return err
    }
    if infra == nil {
      return nil
    }
    // Export the primary key of the Storage Account
    ctx.Export("primaryStorageKey", pulumi.All(infra.resourceGroup.Name, infra.storageAccount.Name).ApplyT(
      func(args []interface{}) (string, error) {
        resourceGroupName := args[0].(string)
        accountName := args[1].(string)
        accountKeys, err := storage.ListStorageAccountKeys(ctx, &storage.ListStorageAccountKeysArgs{
          ResourceGroupName: resourceGroupName,
          AccountName:       accountName,
        })
        if err != nil {
          return "", err
        }

        return accountKeys.Keys[0].Value, nil
      },
      ))
    ctx.Export("resourceGroupName", pulumi.Sprintf("%s", infra.resourceGroup.Name))
    return nil
	})
}

type infrastructure struct {
  resourceGroup *resources.ResourceGroup
  storageAccount *storage.StorageAccount
}

// func tags(ctx *pulumi.Context) (pulumi.StringMap){
//
// }

func buildInfra(ctx *pulumi.Context) (*infrastructure, error) {
  // Create an Azure Resource Group
  // tags := tags(ctx)
  conf := config.New(ctx,"")
  tags := pulumi.StringMap{
    "Project": pulumi.String(ctx.Project()),
    "Stack": pulumi.String(ctx.Stack()),
    "Environment": pulumi.String(conf.Get("shortname")),
    "Name": pulumi.String(fmt.Sprintf("%s-%s", ctx.Project(), ctx.Stack())),
  }
  resourceGroup, err := resources.NewResourceGroup(ctx, "rg", &resources.ResourceGroupArgs{
    Tags: tags,
  })
  if err != nil {
    return nil, err
  }

  // Create an Azure resource (Storage Account)
  account, err := storage.NewStorageAccount(ctx, "sa", &storage.StorageAccountArgs{
    ResourceGroupName: resourceGroup.Name,
    Sku: &storage.SkuArgs{
      Name: pulumi.String("Standard_LRS"),
    },
    Kind: pulumi.String("StorageV2"),
    AllowBlobPublicAccess: pulumi.Bool(false),
    Tags: tags,
  })
  if err != nil {
    return nil, err
  }

  return &infrastructure{
    resourceGroup: resourceGroup,
    storageAccount: account,
  }, nil
}
