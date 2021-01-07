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
	Arguments: []*dcmd.ArgDef{
		{Name: "content", Type: dcmd.String},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		MSG := "Provide a message when..."
		if data.Args[0].Value != nil {
			MSG = data.Args[0].Value.String()
		}
		return MSG, nil
	}),
}