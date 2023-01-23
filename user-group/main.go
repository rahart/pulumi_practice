package main

import (
  "fmt"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
  "github.com/pulumi/pulumi-azuread/sdk/v5/go/azuread"
  "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
  "io/ioutil"
  "gopkg.in/yaml.v3"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
    err := buildUserGroup(ctx)
    if err != nil {
      return err
    }
    return nil
  })
}

type Infrastructure struct {

}

type UserGroupConfig struct {
  Configuration UserGroup `yaml:"usergroup"`
  Edate string `yaml:"date"`
}

type UserGroup struct {
  Users  map[string]User `yaml:"users"`
  Groups map[string]Group `yaml:"groups"`
}

type User struct {
  RealName string `yaml:"name"`
  Email string `yaml:"email"`
  Groups []string `yaml:"groups"`
}

type Group struct {
  Policies []string `yaml:"policies"`
  Description string `yaml:"description"`
}

func userGroupConfig(file string, ctx *pulumi.Context) (*UserGroupConfig, error) {
  ctx.Log.Info(file, nil)
  data, err := ioutil.ReadFile(file)
  if err != nil {
    return nil, err
  }

  var configUG UserGroupConfig
  err = yaml.Unmarshal(data, &configUG)
  if err != nil {
    return nil, err
  }
  return &configUG, nil 
}

func buildUserGroup(ctx *pulumi.Context) error {
  conf := config.New(ctx,"")
  configUG, err := userGroupConfig(conf.Require("user_group"), ctx)
  if err != nil {
    return err
  }
  
  current_ad, err := azuread.GetClientConfig(ctx, nil, nil)
  if err != nil {
    return err
  }
  
  for group, _ := range configUG.Configuration.Groups {
    _, err = azuread.NewGroup(ctx, group, &azuread.GroupArgs{
      DisplayName: pulumi.String(group),
      SecurityEnabled: pulumi.Bool(true),
    })
    if err != nil {
      ctx.Log.Error("Unable to add Group: " + group , nil)
      return err
    }
  }

  for user, user_data := range configUG.Configuration.Users {
    ctx.Log.Info("user: " + user_data.RealName, nil)
    _, err := azuread.NewUser(ctx, user, &azuread.UserArgs{
      DisplayName: pulumi.String(user_data.RealName),
      UserPrincipalName: pulumi.String(user_data.Email),
    })
    if err != nil {
      ctx.Log.Error("Unable to add User: " + user, nil)
      return err
    }

  }
  return nil
 
}
