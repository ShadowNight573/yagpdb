package specialservers

import (
    "fmt"
    "errors"
    
    "github.com/jonas747/dcmd"
    "github.com/jonas747/yagpdb/commands"
    "github.com/jonas747/yagpdb/common"
    "github.com/jonas747/yagpdb/stdcommands/util"
    "github.com/mediocregopher/radix/v3"
    "github.com/volatiletech/sqlboiler/queries/qm"
    "github.com/jonas747/yagpdb/bot/models"
)

var Command = &commands.YAGCommand{
    Cooldown:             2,
    CmdCategory:          commands.CategoryDebug,
    HideFromCommandsPage: true,
    Name:                 "specialservers",
    Description:          ";))",
    HideFromHelp:         true,
    Arguments: []*dcmd.ArgDef{
      &dcmd.ArgDef{Name: "Skip", Help: "Entries to skip", Type: dcmd.Int, Default: 0},
    },
    RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
      offset := data.Args[0].Int()

      var whitelisted []interface{}
      err := common.RedisPool.Do(radix.Cmd(&whitelisted, "SMEMBERS", "special_servers"))
      if err != nil {
          return "", err
      }
      
      results, err := models.JoinedGuilds(
        qm.Where("left_at is null"), // Don't include guilds that we left
        qm.WhereIn("id in ?", whitelisted...), // Only whitelisted guilds
        qm.OrderBy("id desc"), // Needed so we have consistent output
        qm.Limit(10), // Limit to 10 results
        qm.Offset(offset),
      ).AllG(data.Context())
      if err != nil {
        errors.New("this is it")
      }

      resp := "**Whitelisted servers**\n"
      for _, v := range results {
        resp += fmt.Sprintf("`%d`: **%s**\n", v.ID, v.Name)
      }
      return resp, nil
    }),
}