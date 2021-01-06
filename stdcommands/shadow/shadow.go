package shadow

import (
	// "fmt"
	"net/http"
	"bytes"
	"json"
	
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
		{Switch: "url", Name: "webhook url", Type: dcmd.String, Default: ""},
		{Switch: "content", Name: "message to send", Type: dcmd.String, Default: ""},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		MSG, WHL := data.Switches("content").Value, data.Switches("url").Value
		
		BR, err := json.Marshal(MSG)
		if err != nil {
			return "", err
		}
		
		resp, err := http.Post(WHL, "payload_json", bytes.NewBuffer(BR))
		if err != nil {
			return "", err
		}
		return "uh", nil
	}),
}