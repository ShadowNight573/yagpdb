package shadow

import (
	// "fmt"
	
	"github.com/jonas747/dcmd"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/stdcommands/util"
)

var Command = &commands.YAGCommand{
	CmdCategory: commands.CategoryDebug,
	Name:        "shadow",
	Description: "Random command shadow made cause why not. Could do anything :)",
	RunInDM:     true,
	HideFromHelp:         true,
	ArgSwitches: []*dcmd.ArgDef{
		{Switch: "content", Name: "message to send", Type: dcmd.String, Default: ""},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		MSG := data.Switch("content").Value || "Uh"
		return MSG, nil
	}),
}