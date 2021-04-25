package listroles

import (
	"fmt"
	"sort"
	"bytes"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/dutil"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
)

var Command = &commands.YAGCommand{
	CmdCategory: commands.CategoryTool,
	Name:        "ListRoles",
	Description: "List roles, their id's, color hex code, and 'mention everyone' perms (useful if you wanna double check to make sure you didn't give anyone mention everyone perms that shouldn't have it)",
	ArgSwitches: []*dcmd.ArgDef{
		{Switch: "nomanaged", Name: "Don't list managed/bot roles"},
		{Switch: "file", Name: "Send the roles in a file"},
	},

	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		var out, outFinal string
		var noMana bool
		var file bool
		
		if data.Switches["nomanaged"].Value != nil && data.Switches["nomanaged"].Value.(bool) {
			noMana = true
		}
		
		if data.Switches["file"].Value != nil && data.Switches["file"].Value.(bool) {
			file = true
		}
		
		data.GS.Lock()
		defer data.GS.Unlock()

		sort.Sort(dutil.Roles(data.GS.Guild.Roles))
		
		counter := 0
		for _, r := range data.GS.Guild.Roles {
			if noMana && r.Managed {
				continue
			} else {
				counter++
				me := r.Permissions&discordgo.PermissionAdministrator != 0 || r.Permissions&discordgo.PermissionMentionEveryone != 0
				if !file {
					out += fmt.Sprintf("`%-25s: %-19d #%-6x  ME:%5t`\n", r.Name, r.ID, r.Color, me)
				} else {
					out += fmt.Sprintf("%-25s: %-19d #%-6x  ME:%5t\n\n", r.Name, r.ID, r.Color, me)
				}
			}
		}
		outFinal = fmt.Sprintf("Total role count: %d\n", counter)
		outFinal += fmt.Sprintf("%s", "(ME = mention everyone perms)\n")
		outFinal += out
		if !file {
			return outFinal, nil
		}
		
		var buf bytes.Buffer
		buf.WriteString(outFinal)
		msg := &discordgo.MessageSend{}
		msg.Content = "Open the file to see all the roles in this server."
		msg.File = &discordgo.File{
			Name:        "Roles.txt",
			ContentType: "text/plain",
			Reader:      &buf,
		}
		
		_, err := common.BotSession.ChannelMessageSendComplex(data.CS.ID, msg)
		if err != nil {
			return "", err
		}

		return "", nil
	},
}
