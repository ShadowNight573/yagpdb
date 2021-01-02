package shadow

import (
	"fmt"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/stdcommands/util"
)

var Command = &commands.YAGCommand{
	CmdCategory: commands.CategoryGeneral,
	Name:        "shadow",
	Description: "Random command shadow made cause why not. Could do anything :)",
	RunInDM:     true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		out := fmt.Sprint("Wonder what this does...")
		return out, nil
	},
}