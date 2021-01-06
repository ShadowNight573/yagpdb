package info

import (
	"fmt"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/discordgo"
)

var Command = &commands.YAGCommand{
	CmdCategory: commands.CategoryGeneral,
	Name:        "Info",
	Description: "Responds with bot information",
	RunInDM:     true,
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		botUser := common.BotUser
		msg := discordgo.MessageEmbed {
			Author: &discordgo.MessageEmbedAuthor{
				Name:    botUser.Username,
				IconURL: discordgo.EndpointUserAvatar(botUser.ID, botUser.Avatar),
			},
			Description: fmt.Sprintf("This bot is a self hosted version of the bot YAGPDB which stands for Yet Another General Purpose Discord Bot. For more info please visit the yagpdb [website](https://yagpdb.xyz/) or [documentation](https://docs.yagpdb.xyz/).\n\nThis bot is hosted and updated by ShADowNIGHT#9025. You can find the control panel [here](https://%s/manage).", common.ConfHost.GetString()),
		}
		return msg, nil
	},
}
