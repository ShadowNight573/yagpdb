package shadow

import (
	// "fmt"
	
	"github.com/jonas747/dcmd"
    "github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/stdcommands/util"
)

var Command = &commands.YAGCommand{
	CmdCategory: commands.CategoryDebug,
	Name:        "shadow",
	Description: "Random command shadow made cause why not. Could do anything :)",
	RunInDM:     true,
	HideFromHelp:         true,
	Arguments: []*dcmd.ArgDef{
		{Name: "content", Type: dcmd.String},
	},
	ArgSwitches: []*dcmd.ArgDef{
		{Switch: "del", Name: "delete trigger"},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		MSG := "Provide a message when..."
		if data.Args[0].Value != nil {
			MSG = data.Args[0].Value
		}
		$delTrigger := data.Switch("del").Value != nil && data.Switch("del").Value.(bool)
		if $delTrigger {
			bot.MessageDeleteQueue.DeleteMessages(c.GS.ID, c.CS.ID, c.Msg.ID)
		}
		return MSG, nil
	}),
}