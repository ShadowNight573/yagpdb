package specialservers

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
	Name:                 "specialservers",
	Description:          ";))",
	HideFromHelp:         true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
			var whitelisted []int64
			err := common.RedisPool.Do(radix.Cmd(&whitelisted, "SMEMBERS", "special_servers"))
			if err != nil {
				return "", err
			}
			
			return "**List Of Whitelisted Servers:**\n" + whitelisted, nil
	}),
}