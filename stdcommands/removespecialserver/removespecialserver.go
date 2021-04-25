package removespecialserver

import (
    "github.com/jonas747/dcmd"
    "github.com/jonas747/yagpdb/commands"
    "github.com/jonas747/yagpdb/common"
    "github.com/jonas747/yagpdb/stdcommands/util"
    "github.com/mediocregopher/radix/v3"
)

var Command = &commands.YAGCommand{
    Cooldown:             2,
    CmdCategory:          commands.CategoryDebug,
    HideFromCommandsPage: true,
    Name:                 "removespecialserver",
    Description:          ";))",
    HideFromHelp:         true,
    RequiredArgs:         1,
    Arguments: []*dcmd.ArgDef{
        {Name: "server", Type: dcmd.String},
    },
    RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {

        var whitelisted bool
        err := common.RedisPool.Do(radix.FlatCmd(&whitelisted, "SREM", "special_servers", data.Args[0].Str()))
        if err != nil {
            return nil, err
        }

        if !whitelisted {
            return "Server wasnt whitelisted", nil
        }

        return "UnWhitelisted server", nil
    }),
} 