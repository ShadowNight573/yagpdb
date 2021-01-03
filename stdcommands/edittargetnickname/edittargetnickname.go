package edittargetnickname

import (
    "fmt"
    "strings"

    "github.com/jonas747/dcmd"
    "github.com/jonas747/yagpdb/bot"
    "github.com/jonas747/yagpdb/commands"
    "github.com/jonas747/yagpdb/common"
)

var Command = &commands.YAGCommand{
    CmdCategory:  commands.CategoryTool,
    Name:         "EditTargetNickname",
    Aliases:      []string{"etn"},
    Description:  "Edits the nickname of the specified user",
    RequiredArgs: 1,
    Arguments: []*dcmd.ArgDef{
        {Name: "User", Type: dcmd.UserID},
        {Name: "Nick", Type: dcmd.String},
    },
    RunFunc: func(data *dcmd.Data) (interface{}, error) {
        ms, err := bot.GetMember(data.GS.ID, data.Args[0].Int64())
        if err != nil {
            return "Member not found.", nil
        }

        nick := SafeArgString(data, 1)
        if strings.Compare(ms.Nick, nick) == 0 {
            return "This is that user's nickname already.", nil
        }

        err = common.BotSession.GuildMemberNickname(data.GS.ID, ms.ID, nick)
        if err != nil {
            return "", err
        }

        if nick == "" {
            return fmt.Sprintf("The nickname of user <@%d> was removed.", ms.ID), nil
        }

        return fmt.Sprintf("The nickname of user <@%d> was updated to `%s`.", ms.ID, nick), nil
    },
}

func SafeArgString(data *dcmd.Data, arg int) string {
    if arg >= len(data.Args) || data.Args[arg].Value == nil {
        return ""
    }

    return data.Args[arg].Str()
}