package templates

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jonas747/discordgo"
	"github.com/jonas747/dstate/v2"
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/common/scheduledevents2"
)

var ErrTooManyCalls = errors.New("Too many calls to this function")
var ErrTooManyAPICalls = errors.New("Too many potential discord api calls function")

func (c *Context) buildDM(gName string, s ...interface{}) *discordgo.MessageSend {
	msgSend := &discordgo.MessageSend{
		AllowedMentions: discordgo.AllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
		},
	}

	switch t := s[0].(type) {
	case *discordgo.MessageEmbed:
		msgSend.Embed = t
	case *discordgo.MessageSend:
		msgSend = t
		if (strings.TrimSpace(msgSend.Content) == "") && (msgSend.File == nil) {
			return nil
		}
	default:
		msgSend.Content = fmt.Sprintf("%s", fmt.Sprint(s...))
	}

	if !bot.IsSpecialGuild(c.GS.Guild.ID) {
		info := fmt.Sprintf("DM from server: %s", gName)
		if msgSend.Embed != nil {
			msgSend.Embed.Footer = &discordgo.MessageEmbedFooter{
				Text: info,
			}
		} else {
			info := fmt.Sprintf("DM from server: **%s**", gName)
			msgSend.Content = info + "\n" + msgSend.Content
		}
	}

	return msgSend
}

func (c *Context) tmplSendDM(s ...interface{}) string {
	if len(s) < 1 || c.IncreaseCheckCallCounter("send_dm", 3) || c.MS == nil {
		return ""
	}

	c.GS.RLock()
	memberID, gName := c.MS.ID, c.GS.Guild.Name
	c.GS.RUnlock()
	
	msgSend := c.buildDM(gName, s...)
	if msgSend == nil {
		return ""
	}

	channel, err := common.BotSession.UserChannelCreate(memberID)
	if err != nil {
		return ""
	}
	_, _ = common.BotSession.ChannelMessageSendComplex(channel.ID, msgSend)
	return ""
}

func (c *Context) tmplSendTargetDM(target interface{}, s ...interface{}) string {
	if bot.IsSpecialGuild(c.GS.ID) {
		if len(s) < 1 || c.IncreaseCheckCallCounter("send_dm", 3) {
			return ""
		}

		targetID := targetUserID(target)
		if targetID == 0 {
			return ""
		}

		ts, err := bot.GetMember(c.GS.ID, targetID)
		if err != nil {
			return ""
		}

		msgSend := c.buildDM("", s...)
		if msgSend == nil {
			return ""
		}

		channel, err := common.BotSession.UserChannelCreate(ts.ID)
		if err != nil {
			return ""
		}
		_, _ = common.BotSession.ChannelMessageSendComplex(channel.ID, msgSend)
	}

	return ""
}


// ChannelArg converts a verity of types of argument into a channel, verifying that it exists
func (c *Context) ChannelArg(v interface{}) int64 {

	c.GS.RLock()
	defer c.GS.RUnlock()

	// Look for the channel
	if v == nil && c.CurrentFrame.CS != nil {
		// No channel passed, assume current channel
		return c.CurrentFrame.CS.ID
	}

	verifiedExistence := false
	var cid int64
	if v != nil {
		switch t := v.(type) {
		case int, int64:
			// Channel id passed
			cid = ToInt64(t)
		case string:
			parsed, err := strconv.ParseInt(t, 10, 64)
			if err == nil {
				// Channel id passed in string format
				cid = parsed
			} else {
				// Channel name, look for it
				for _, v := range c.GS.Channels {
					if strings.EqualFold(t, v.Name) && v.Type == discordgo.ChannelTypeGuildText {
						cid = v.ID
						verifiedExistence = true
						break
					}
				}
			}
		}
	}

	if !verifiedExistence {
		// Make sure the channel is part of the guild
		for k, _ := range c.GS.Channels {
			if k == cid {
				verifiedExistence = true
				break
			}
		}
	}

	if !verifiedExistence {
		return 0
	}

	return cid
}



// ChannelArgNoDM is the same as ChannelArg but will not accept DM channels
func (c *Context) ChannelArgNoDM(v interface{}) int64 {

	c.GS.RLock()
	defer c.GS.RUnlock()

	// Look for the channel
	if v == nil && c.CurrentFrame.CS != nil {
		// No channel passed, assume current channel
		v = c.CurrentFrame.CS.ID
	}

	verifiedExistence := false
	var cid int64
	if v != nil {
		switch t := v.(type) {
		case int, int64:
			// Channel id passed
			cid = ToInt64(t)
		case string:
			parsed, err := strconv.ParseInt(t, 10, 64)
			if err == nil {
				// Channel id passed in string format
				cid = parsed
			} else {
				// Channel name, look for it
				for _, v := range c.GS.Channels {
					if strings.EqualFold(t, v.Name) && v.Type == discordgo.ChannelTypeGuildText {
						cid = v.ID
						verifiedExistence = true
						break
					}
				}
			}
		}
	}

	if !verifiedExistence {
		// Make sure the channel is part of the guild
		for k, _ := range c.GS.Channels {
			if k == cid {
				verifiedExistence = true
				break
			}
		}
	}

	if !verifiedExistence {
		return 0
	}

	return cid
}

func (c *Context) tmplSendTemplateDM(name string, data ...interface{}) (interface{}, error) {
	return c.sendNestedTemplate(nil, true, false, name, data...)
}

func (c *Context) tmplSendTemplate(channel interface{}, name string, data ...interface{}) (interface{}, error) {
	return c.sendNestedTemplate(channel, false, false, name, data...)
}

func (c *Context) tmplExecTemplate(channel interface{}, name string, data ...interface{}) (interface{}, error) {
 	return c.sendNestedTemplate(channel, false, true, name, data...)
 }

func (c *Context) sendNestedTemplate(channel interface{}, dm , exec bool, name string, data ...interface{}) (interface{}, error) {
	if c.IncreaseCheckCallCounter("exec_child", 3) {
		return "", ErrTooManyCalls
	}
	if name == "" {
		return "", errors.New("No template name passed")
	}
	if c.CurrentFrame.isNestedTemplate {
		return "", errors.New("Can't call this in a nested template")
	}

	t := c.CurrentFrame.parsedTemplate.Lookup(name)
	if t == nil {
		return "", errors.New("Unknown template")
	}

	var cs *dstate.ChannelState
	// find the new context channel
	if !dm {
		if channel == nil {
			cs = c.CurrentFrame.CS
		} else {
			cID := c.ChannelArg(channel)
			if cID == 0 {
				return "", errors.New("Unknown channel")
			}

			cs = c.GS.ChannelCopy(true, cID)
			if cs == nil {
				return "", errors.New("Unknown channel")
			}
		}
	} else {
		if c.CurrentFrame.SendResponseInDM {
			cs = c.CurrentFrame.CS
		} else {
			ch, err := common.BotSession.UserChannelCreate(c.MS.ID)
			if err != nil {
				return "", err
			}

			cs = &dstate.ChannelState{
				Owner: c.GS,
				Guild: c.GS,
				ID:    ch.ID,
				Name:  c.MS.Username,
				Type:  discordgo.ChannelTypeDM,
			}
		}
	}

	oldFrame := c.newContextFrame(cs)
	defer func() {
		c.CurrentFrame = oldFrame
	}()

	if dm {
		c.CurrentFrame.SendResponseInDM = oldFrame.SendResponseInDM
	} else if channel == nil {
		// inherit
		c.CurrentFrame.SendResponseInDM = oldFrame.SendResponseInDM
	}
	c.CurrentFrame.execMode = exec
	// pass some data
	if len(data) > 1 {
		dict, _ := Dictionary(data...)
		c.Data["TemplateArgs"] = dict
		if !c.checkSafeDictNoRecursion(dict, 0) {
			return nil, errors.New("trying to pass the entire current context data in as templateargs, this is not needed, just use nil and access all other data normally")
		}
	} else if len(data) == 1 {
		if cast, ok := data[0].(map[string]interface{}); ok && reflect.DeepEqual(cast, c.Data) {
			return nil, errors.New("trying to pass the entire current context data in as templateargs, this is not needed, just use nil and access all other data normally")
		}
		c.Data["TemplateArgs"] = data[0]
	}

	// and finally execute the child template
	c.CurrentFrame.parsedTemplate = t
	resp, err := c.executeParsed()
	if err != nil {
		return "", err
	}
	
	if exec {
 		var execReturnStruct CtxExecReturn
 		execReturnStruct.Response = c.MessageSend(resp)
 		execReturnStruct.Return = c.CurrentFrame.execReturn
 		return execReturnStruct, err
 	}

	m, err := c.SendResponse(resp)
	if err != nil {
		return "", err
	}

	if m != nil {
		return m.ID, err
	}
	return "", err
}

func (c* Context) tmplAddReturn(data...interface{}) (interface{} ,error) {
 	if !c.CurrentFrame.isNestedTemplate || !c.CurrentFrame.execMode {
 		return "", errors.New("Can only be used in nested templates in exec mode.")
 	}

 	if len(c.CurrentFrame.execReturn) + len (data) > 10 {
 		return "", errors.New("Return length cannot exceed 10")
 	}

 	c.CurrentFrame.execReturn = append(c.CurrentFrame.execReturn, data...)
 	return "", nil
}

func (c *Context) checkSafeStringDictNoRecursion(d SDict, n int) bool {
	if n > 1000 {
		return false
	}

	for _, v := range d {
		if cast, ok := v.(Dict); ok {
			if !c.checkSafeDictNoRecursion(cast, n+1) {
				return false
			}
		}

		if cast, ok := v.(SDict); ok {
			if !c.checkSafeStringDictNoRecursion(cast, n+1) {
				return false
			}
		}

		if reflect.DeepEqual(v, c.Data) {
			return false
		}
	}

	return true
}

func (c *Context) checkSafeDictNoRecursion(d Dict, n int) bool {
	if n > 1000 {
		return false
	}

	for _, v := range d {
		if cast, ok := v.(Dict); ok {
			if !c.checkSafeDictNoRecursion(cast, n+1) {
				return false
			}
		}

		if cast, ok := v.(SDict); ok {
			if !c.checkSafeStringDictNoRecursion(cast, n+1) {
				return false
			}
		}

		if reflect.DeepEqual(v, c.Data) {
			return false
		}
	}

	return true
}

func (c *Context) tmplSendMessage(filterSpecialMentions bool, returnID bool) func(channel interface{}, msg interface{}) interface{} {
	parseMentions := []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers}
	if !filterSpecialMentions {
		parseMentions = append(parseMentions, discordgo.AllowedMentionTypeRoles, discordgo.AllowedMentionTypeEveryone)
	}

	return func(channel interface{}, msg interface{}) interface{} {
		if c.IncreaseCheckGenericAPICall() {
			return ""
		}

		cid := c.ChannelArg(channel)
		if cid == 0 {
			return ""
		}

		isDM := cid != c.ChannelArgNoDM(channel)
 		c.GS.RLock()
 		info := fmt.Sprintf("DM from server: %s", c.GS.Guild.Name)
 		c.GS.RUnlock()
 		WL := bot.IsSpecialGuild(c.GS.ID)
		
		var m *discordgo.Message
		msgSend := &discordgo.MessageSend{
			AllowedMentions: discordgo.AllowedMentions{
				Parse: parseMentions,
			},
		}
		var err error

		switch typedMsg := msg.(type) {

		case *discordgo.MessageEmbed:
			if isDM && !WL {
 				typedMsg.Footer = &discordgo.MessageEmbedFooter{
 					Text: info,
 				}
 			}
			msgSend.Embed = typedMsg
		case *discordgo.MessageSend:
			msgSend = typedMsg
			if !filterSpecialMentions {
				msgSend.AllowedMentions = discordgo.AllowedMentions{Parse: parseMentions}
			}
			
			if isDM && !WL {
				if typedMsg.Embed != nil {
					typedMsg.Embed.Footer = &discordgo.MessageEmbedFooter{
						Text: info,
					}
				} else {
					info = fmt.Sprintf("DM from server **%s**", c.GS.Guild.Name)
					typedMsg.Content = info + "\n" + typedMsg.Content
 				}
 			}
		default:
			if isDM && !WL {
				info = fmt.Sprintf("DM from server **%s**", c.GS.Guild.Name)
 				msgSend.Content = info + "\n" + fmt.Sprint(msg)
 			} else {
 				msgSend.Content = fmt.Sprint(msg)
 			}
		}

		m, err = common.BotSession.ChannelMessageSendComplex(cid, msgSend)

		if err == nil && returnID {
			return m.ID
		}

		return ""
	}
}

func (c *Context) tmplEditMessage(filterSpecialMentions bool) func(channel interface{}, msgID interface{}, msg interface{}) (interface{}, error) {
	return func(channel interface{}, msgID interface{}, msg interface{}) (interface{}, error) {
		if c.IncreaseCheckGenericAPICall() {
			return "", ErrTooManyAPICalls
		}

		cid := c.ChannelArgNoDM(channel)
		if cid == 0 {
			return "", errors.New("Unknown channel")
		}

		mID := ToInt64(msgID)
		msgEdit := &discordgo.MessageEdit{
			ID:      mID,
			Channel: cid,
		}
		var err error

		switch typedMsg := msg.(type) {

		case *discordgo.MessageEmbed:
			msgEdit.Embed = typedMsg
		case *discordgo.MessageEdit:
			//If both Embed and string are explicitly set as null, give an error message.
			if typedMsg.Content != nil && strings.TrimSpace(*typedMsg.Content) == "" && typedMsg.Embed != nil && typedMsg.Embed.GetMarshalNil() {
				return "", errors.New("both content and embed cannot be null")
			}
			msgEdit.Content = typedMsg.Content
			msgEdit.Embed = typedMsg.Embed
		default:
			temp := fmt.Sprint(msg)
			msgEdit.Content = &temp
		}

		_, err = common.BotSession.ChannelMessageEditComplex(msgEdit)

		if err != nil {
			return "", err
		}

		return "", nil
	}
}

func (c *Context) tmplMentionEveryone() string {
	c.CurrentFrame.MentionEveryone = true
	return "@everyone"
}

func (c *Context) tmplMentionHere() string {
	c.CurrentFrame.MentionHere = true
	return "@here"
}

func (c *Context) tmplMentionRoleID(roleID interface{}) string {
	if c.IncreaseCheckStateLock() {
		return ""
	}

	var role int64
	switch r := roleID.(type) {
	case int64:
		role = r
	case int:
		role = int64(r)
	case string:
		role, _ = strconv.ParseInt(r, 10, 64)
	default:
		return ""
	}

	r := c.GS.RoleCopy(true, role)
	if r == nil {
		return "(role not found)"
	}

	if common.ContainsInt64Slice(c.CurrentFrame.MentionRoles, role) {
		return "<@&" + discordgo.StrID(role) + ">"
	}

	c.CurrentFrame.MentionRoles = append(c.CurrentFrame.MentionRoles, role)
	return "<@&" + discordgo.StrID(role) + ">"
}

func (c *Context) tmplMentionRoleName(role string) string {
	if c.IncreaseCheckStateLock() {
		return ""
	}

	var found *discordgo.Role
	c.GS.RLock()
	for _, r := range c.GS.Guild.Roles {
		if r.Name == role {
			if !common.ContainsInt64Slice(c.CurrentFrame.MentionRoles, r.ID) {
				c.CurrentFrame.MentionRoles = append(c.CurrentFrame.MentionRoles, r.ID)
				found = r
			}
		}
	}
	c.GS.RUnlock()
	if found == nil {
		return "(role not found)"
	}

	return "<@&" + discordgo.StrID(found.ID) + ">"
}

func (c *Context) tmplHasRoleID(roleID interface{}) bool {
	role := ToInt64(roleID)
	if role == 0 {
		return false
	}

	contains := common.ContainsInt64Slice(c.MS.Roles, role)
	return contains
}

func (c *Context) tmplHasRoleName(name string) (bool, error) {
	if c.IncreaseCheckStateLock() {
		return false, ErrTooManyCalls
	}

	role := c.findRoleByName(name)
	if role == nil {
		return false, nil
	}

	if common.ContainsInt64Slice(c.MS.Roles, role.ID) {
		return true, nil
	}

	// Role not found, default to false
	return false, nil
}

func targetUserID(input interface{}) int64 {
	switch t := input.(type) {
	case *discordgo.User:
		return t.ID
	case string:
		str := strings.TrimSpace(t)
		if strings.HasPrefix(str, "<@") && strings.HasSuffix(str, ">") && (len(str) > 4) {
			trimmed := str[2 : len(str)-1]
			if trimmed[0] == '!' {
				trimmed = trimmed[1:]
			}
			str = trimmed
		}

		return ToInt64(str)
	default:
		return ToInt64(input)
	}
}

func (c *Context) tmplTargetHasRoleID(target interface{}, roleID interface{}) bool {
	if c.IncreaseCheckStateLock() {
		return false
	}

	targetID := targetUserID(target)
	if targetID == 0 {
		return false
	}

	ts, err := bot.GetMember(c.GS.ID, targetID)
	if err != nil {
		return false
	}

	role := ToInt64(roleID)
	if role == 0 {
		return false
	}

	contains := common.ContainsInt64Slice(ts.Roles, role)

	return contains

}

	func (c *Context) tmplSetRoles(target interface{}, roleSlice interface{}) (string, error) {
 	if c.IncreaseCheckGenericAPICall() {
 		return "", ErrTooManyAPICalls
 	}

 	targetID := targetUserID(target)
 	if targetID == 0 {
 		return "", nil
 	}

 	if c.IncreaseCheckCallCounter("set_roles"+discordgo.StrID(targetID), 1) {
 		return "", errors.New("Too many calls to setRoles for specific user ID (max 1 / user)")
 	}

 	rSlice := reflect.ValueOf(roleSlice)
 	switch rSlice.Kind() {
 	case reflect.Slice, reflect.Array:
 		// ok
 	default:
 		return "", errors.New("Value passed was not an array or slice")
 	}

 	if rSlice.Len() > 250 {
 		return "", errors.New("Length of slice passed was > 250 (Discord role limit)")
 	}

 	roles := make([]string, 0, rSlice.Len())
 	for i := 0; i < rSlice.Len(); i++ {
 		switch v := rSlice.Index(i).Interface().(type) {
 		case string:
 			roles = append(roles, v)
 		case int, int64:
 			roles = append(roles, discordgo.StrID(reflect.ValueOf(v).Int()))
 		default:
 			return "", errors.New("Could not convert slice to string slice")
 		}
 	}

 	err := common.BotSession.GuildMemberEdit(c.GS.ID, targetID, roles)
 	if err != nil {
 		return "", err
 	}
 	return "", nil
}

func (c *Context) tmplTargetHasRoleName(target interface{}, name string) bool {
	if c.IncreaseCheckStateLock() {
		return false
	}

	targetID := targetUserID(target)
	if targetID == 0 {
		return false
	}

	ts, err := bot.GetMember(c.GS.ID, targetID)
	if err != nil {
		return false
	}

	role := c.findRoleByName(name)
	if role == nil {
		return false
	}

	if common.ContainsInt64Slice(ts.Roles, role.ID) {			
		return true
	}
	return false

}

func (c *Context) tmplGiveRoleID(target interface{}, roleID interface{}, optionalArgs ...interface{}) string {
	if c.IncreaseCheckGenericAPICall() {
		return ""
	}

	delay := 0
	if len(optionalArgs) > 0 {
		delay = tmplToInt(optionalArgs[0])
	}

	targetID := targetUserID(target)
	if targetID == 0 {
		return ""
	}

	role := ToInt64(roleID)
	if role == 0 {
		return ""
	}

	// Check to see if we can save a API request here, if this isn't delayed
	if delay <= 0 {
		c.GS.RLock()
		ms := c.GS.Member(false, targetID)
		hasRole := true
		if ms != nil && ms.MemberSet {
			hasRole = common.ContainsInt64Slice(ms.Roles, role)
		}
		c.GS.RUnlock()

		if !hasRole {
			return ""
		}
	}

	if delay > 0 {
		scheduledevents2.ScheduleAddRole(context.Background(), c.GS.ID, targetID, role, time.Now().Add(time.Second*time.Duration(delay)))
	} else {
		common.BotSession.GuildMemberRoleAdd(c.GS.ID, targetID, role)
	}

	return ""
}

func (c *Context) tmplGiveRoleName(target interface{}, name string, optionalArgs ...interface{}) string {
	if c.IncreaseCheckGenericAPICall() {
		return ""
	}

	delay := 0
	if len(optionalArgs) > 0 {
		delay = tmplToInt(optionalArgs[0])
	}

	targetID := targetUserID(target)
	if targetID == 0 {
		return ""
	}

	role := c.findRoleByName(name)
	if role == nil {
		return "no role by the name of " + name + " found"
	}

	// Maybe save an api request
	if delay <= 0 {
		ms := c.GS.Member(false, targetID)
		if ms != nil {
			if common.ContainsInt64Slice(ms.Roles, role.ID) {
				return ""
			}
		}
	}

	if delay > 0 {
		scheduledevents2.ScheduleAddRole(context.Background(), c.GS.ID, targetID, role.ID, time.Now().Add(time.Second*time.Duration(delay)))
	} else {
		common.BotSession.GuildMemberRoleAdd(c.GS.ID, targetID, role.ID)
	}

	return ""
}

func (c *Context) tmplTakeRoleID(target interface{}, roleID interface{}, optionalArgs ...interface{}) string {
	if c.IncreaseCheckGenericAPICall() {
		return ""
	}

	delay := 0
	if len(optionalArgs) > 0 {
		delay = tmplToInt(optionalArgs[0])
	}

	targetID := targetUserID(target)
	if targetID == 0 {
		return ""
	}

	role := ToInt64(roleID)
	if role == 0 {
		return ""
	}

	// Check to see if we can save a API request here, if this isn't delayed
	if delay <= 0 {
		c.GS.RLock()
		ms := c.GS.Member(false, targetID)
		hasRole := true
		if ms != nil && ms.MemberSet {
			hasRole = common.ContainsInt64Slice(ms.Roles, role)
		}
		c.GS.RUnlock()

		if !hasRole {
			return ""
		}
	}

	if delay > 0 {
		scheduledevents2.ScheduleRemoveRole(context.Background(), c.GS.ID, targetID, role, time.Now().Add(time.Second*time.Duration(delay)))
	} else {
		common.BotSession.GuildMemberRoleRemove(c.GS.ID, targetID, role)
	}

	return ""
}

func (c *Context) tmplTakeRoleName(target interface{}, name string, optionalArgs ...interface{}) string {
	if c.IncreaseCheckGenericAPICall() {
		return ""
	}

	delay := 0
	if len(optionalArgs) > 0 {
		delay = tmplToInt(optionalArgs[0])
	}

	targetID := targetUserID(target)
	if targetID == 0 {
		return ""
	}

	role := c.findRoleByName(name)
	if role == nil {
		return "no role by the name of " + name + " found"
	}

	// Maybe save an api request
	if delay <= 0 {
		ms := c.GS.Member(false, targetID)
		if ms != nil {
			if common.ContainsInt64Slice(ms.Roles, role.ID) {
				return ""
			}
		}
	}

	if delay > 0 {
		scheduledevents2.ScheduleRemoveRole(context.Background(), c.GS.ID, targetID, role.ID, time.Now().Add(time.Second*time.Duration(delay)))
	} else {
		common.BotSession.GuildMemberRoleRemove(c.GS.ID, targetID, role.ID)
	}

	return ""
}

func (c *Context) tmplAddRoleID(role interface{}, optionalArgs ...interface{}) (string, error) {
	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	delay := 0
	if len(optionalArgs) > 0 {
		delay = tmplToInt(optionalArgs[0])
	}

	if c.MS == nil {
		return "", nil
	}

	rid := ToInt64(role)
	if rid == 0 {
		return "", errors.New("No role id specified")
	}

	if delay > 0 {
		scheduledevents2.ScheduleAddRole(context.Background(), c.GS.ID, c.MS.ID, rid, time.Now().Add(time.Second*time.Duration(delay)))
	} else {
		if err := common.AddRoleDS(c.MS, rid); err != nil {
			return "", err
		}
	}

	return "", nil
}

func (c *Context) tmplAddRoleName(name string, optionalArgs ...interface{}) (string, error) {
	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}
	
	delay := 0
	if len(optionalArgs) > 0 {
		delay = tmplToInt(optionalArgs[0])
	}

	if c.MS == nil {
		return "", nil
	}

	role := c.findRoleByName(name)
	if role == nil {
		return "", errors.New("No Role with name " + name + " found")
	}

	if delay > 0 {
		scheduledevents2.ScheduleAddRole(context.Background(), c.GS.ID, c.MS.ID, role.ID, time.Now().Add(time.Second*time.Duration(delay)))
	} else {
		if err := common.AddRoleDS(c.MS, role.ID); err != nil {
			return "", err
		}
	}

	return "", nil
}

func (c *Context) tmplRemoveRoleID(role interface{}, optionalArgs ...interface{}) (string, error) {
	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	delay := 0
	if len(optionalArgs) > 0 {
		delay = tmplToInt(optionalArgs[0])
	}

	if c.MS == nil {
		return "", nil
	}

	rid := ToInt64(role)
	if rid == 0 {
		return "", errors.New("No role id specified")
	}

	if delay > 0 {
		scheduledevents2.ScheduleRemoveRole(context.Background(), c.GS.ID, c.MS.ID, rid, time.Now().Add(time.Second*time.Duration(delay)))
	} else {
		if err := common.RemoveRoleDS(c.MS, rid); err != nil {
			return "", err
		}
	}

	return "", nil
}

func (c *Context) tmplRemoveRoleName(name string, optionalArgs ...interface{}) (string, error) {
	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	delay := 0
	if len(optionalArgs) > 0 {
		delay = tmplToInt(optionalArgs[0])
	}

	if c.MS == nil {
		return "", nil
	}

	role := c.findRoleByName(name)
	if role == nil {
		return "", errors.New("No Role with name " + name + " found")
	}

	if delay > 0 {
		scheduledevents2.ScheduleRemoveRole(context.Background(), c.GS.ID, c.MS.ID, role.ID, time.Now().Add(time.Second*time.Duration(delay)))
	} else {
		if err := common.RemoveRoleDS(c.MS, role.ID); err != nil {
			return "", err
		}
	}

	return "", nil
}

func (c *Context) findRoleByID(id int64) *discordgo.Role {
 	c.GS.RLock()
 	defer c.GS.RUnlock()

 	for _, r := range c.GS.Guild.Roles {
 		if r.ID == id {
 			return r
 		}
 	}

 	return nil
 }
 
func (c *Context) findRoleByName(name string) *discordgo.Role {
	c.GS.RLock()
	defer c.GS.RUnlock()

	for _, r := range c.GS.Guild.Roles {
		if strings.EqualFold(r.Name, name) {
			return r
		}
	}

	return nil
}

func (c *Context) tmplGetRole(r interface{}) (*discordgo.Role, error) {
 	if c.IncreaseCheckGenericAPICall() {
 		return nil, ErrTooManyAPICalls
 	}

 	switch t := r.(type) {
 	case int, int64:
 		return c.findRoleByID(ToInt64(t)), nil
 	case string:
 		parsed, err := strconv.ParseInt(t, 10, 64)
 		if err == nil {
 			return c.findRoleByID(parsed), nil
 		}

 		if strings.HasPrefix(t, "<@&") && strings.HasSuffix(t, ">") {
 			re := regexp.MustCompile(`\d+`)
 			found := re.FindAllString(t, 1)
 			if len(found) > 0 {
 				parsedMention, err := strconv.ParseInt(found[0], 10, 64)
 				if err == nil {
 					return c.findRoleByID(parsedMention), nil
 				}
 			}
 		}

 		return c.findRoleByName(t), nil
 	default:
 		return nil, nil
 	}
 }


func (c *Context) tmplDelResponse(args ...interface{}) string {
	dur := 10
	if len(args) > 0 {
		dur = int(ToInt64(args[0]))
	}
	if dur > 86400 {
		dur = 86400
	}

	c.CurrentFrame.DelResponseDelay = dur
	c.CurrentFrame.DelResponse = true
	return ""
}

func (c *Context) tmplDelTrigger(args ...interface{}) string {
	if c.Msg != nil {
		return c.tmplDelMessage(c.Msg.ChannelID, c.Msg.ID, args...)
	}

	return ""
}

func (c *Context) tmplDelMessage(channel, msgID interface{}, args ...interface{}) string {
	cID := c.ChannelArgNoDM(channel)
	if cID == 0 {
		return ""
	}

	mID := ToInt64(msgID)

	dur := 10
	if len(args) > 0 {
		dur = int(ToInt64(args[0]))
	}

	if dur > 86400 {
		dur = 86400
	}

	MaybeScheduledDeleteMessage(c.GS.ID, cID, mID, dur)

	return ""
}

//Deletes reactions from a message either via reaction trigger or argument-set of emojis,
//needs channelID, messageID, userID, list of emojis - up to twenty
//can be run once per CC.
func (c *Context) tmplDelMessageReaction(values ...reflect.Value) (reflect.Value, error) {

	f := func(args []reflect.Value) (reflect.Value, error) {
		if len(args) < 4 {
			return reflect.Value{}, errors.New("Not enough arguments (need channelID, messageID, userID, emoji)")
		}

		var cArg interface{}
		if args[0].IsValid() {
			cArg = args[0].Interface()
		}

		cID := c.ChannelArg(cArg)
		if cID == 0 {
			return reflect.ValueOf("non-existing channel"), nil
		}

		mID := ToInt64(args[1].Interface())
		uID := targetUserID(args[2].Interface())

		for _, reaction := range args[3:] {

			if c.IncreaseCheckCallCounter("del_reaction_message", 10) {
				return reflect.Value{}, ErrTooManyCalls
			}

			if err := common.BotSession.MessageReactionRemove(cID, mID, reaction.String(), uID); err != nil {
				return reflect.Value{}, err
			}
		}
		return reflect.ValueOf(""), nil
	}

	return callVariadic(f, false, values...)
}

func (c *Context) tmplDelAllMessageReactions(values ...reflect.Value) (reflect.Value, error) {
	
	f := func(args []reflect.Value) (reflect.Value, error) {
		if len(args) < 2 {
			return reflect.Value{}, errors.New("Not enough arguments (need channelID, messageID, emojis[optional])")
		}

		var cArg interface{}
		if args[0].IsValid() {
			cArg = args[0].Interface()
		}

		cID := c.ChannelArg(cArg)
		if cID == 0 {
			return reflect.ValueOf("non-existing channel"), nil
		}

		mID := ToInt64(args[1].Interface())
		

		if len(args) > 2 {
			for _, emoji := range args[2:] {
				if c.IncreaseCheckCallCounter("del_reaction_message", 10) {
					return reflect.Value{}, ErrTooManyCalls
				}
			
				if err := common.BotSession.MessageReactionRemoveEmoji(cID, mID, emoji.String()); err != nil {
					return reflect.Value{}, err
				}
			}
			return reflect.ValueOf(""), nil
		}

		if c.IncreaseCheckGenericAPICall() {
			return reflect.Value{}, ErrTooManyAPICalls
		}
		common.BotSession.MessageReactionsRemoveAll(cID, mID)
		return reflect.ValueOf(""), nil
	}

	return callVariadic(f, false, values...)		
}

func (c *Context) tmplGetMessage(channel, msgID interface{}) (*discordgo.Message, error) {
	if c.IncreaseCheckGenericAPICall() {
		return nil, ErrTooManyAPICalls
	}

	cID := c.ChannelArgNoDM(channel)
	if cID == 0 {
		return nil, nil
	}

	mID := ToInt64(msgID)

	message, _ := common.BotSession.ChannelMessage(cID, mID)
	return message, nil
}

func (c *Context) tmplGetMember(target interface{}) (*CtxMember, error) {
	if c.IncreaseCheckGenericAPICall() {
		return nil, ErrTooManyAPICalls
	}

	mID := targetUserID(target)
	if mID == 0 {
		return nil, nil
	}

	member, _ := bot.GetMember(c.GS.ID, mID)
	if member == nil {
		return nil, nil
	}

	return CtxMemberFromMS(member), nil
}

func (c *Context) tmplGetChannel(channel interface{}) (*CtxChannel, error) {

	if c.IncreaseCheckGenericAPICall() {
		return nil, ErrTooManyAPICalls
	}

	cID := c.ChannelArg(channel)
	if cID == 0 {
		return nil, nil //dont send an error , a nil output would indicate invalid/unknown channel
	}

	cstate := c.GS.ChannelCopy(true, cID)

	if cstate == nil {
		return nil, errors.New("Channel not in state")
	}

	return CtxChannelFromCS(cstate), nil
}

func (c *Context) tmplAddReactions(values ...reflect.Value) (reflect.Value, error) {
	f := func(args []reflect.Value) (reflect.Value, error) {
		if c.Msg == nil {
			return reflect.Value{}, nil
		}

		for _, reaction := range args {
			if c.IncreaseCheckCallCounter("add_reaction_trigger", 20) {
				return reflect.Value{}, ErrTooManyCalls
			}

			if err := common.BotSession.MessageReactionAdd(c.Msg.ChannelID, c.Msg.ID, reaction.String()); err != nil {
				return reflect.Value{}, err
			}
		}
		return reflect.ValueOf(""), nil
	}

	return callVariadic(f, true, values...)
}

func (c *Context) tmplAddResponseReactions(values ...reflect.Value) (reflect.Value, error) {
	f := func(args []reflect.Value) (reflect.Value, error) {
		for _, reaction := range args {
			if c.IncreaseCheckCallCounter("add_reaction_response", 20) {
				return reflect.Value{}, ErrTooManyCalls
			}

			c.CurrentFrame.AddResponseReactionNames = append(c.CurrentFrame.AddResponseReactionNames, reaction.String())
		}
		return reflect.ValueOf(""), nil
	}

	return callVariadic(f, true, values...)
}

func (c *Context) tmplAddMessageReactions(values ...reflect.Value) (reflect.Value, error) {
	f := func(args []reflect.Value) (reflect.Value, error) {
		if len(args) < 2 {
			return reflect.Value{}, errors.New("Not enough arguments (need channel and message-id)")
		}

		// cArg := args[0].Interface()
		var cArg interface{}
		if args[0].IsValid() {
			cArg = args[0].Interface()
		}

		cID := c.ChannelArg(cArg)
		mID := ToInt64(args[1].Interface())

		if cID == 0 {
			return reflect.ValueOf(""), nil
		}

		for i, reaction := range args {
			if i < 2 {
				continue
			}

			if c.IncreaseCheckCallCounter("add_reaction_message", 20) {
				return reflect.Value{}, ErrTooManyCalls
			}

			if err := common.BotSession.MessageReactionAdd(cID, mID, reaction.String()); err != nil {
				return reflect.Value{}, err
			}
		}
		return reflect.ValueOf(""), nil
	}

	return callVariadic(f, false, values...)
}

func (c *Context) tmplCurrentUserAgeHuman() string {
	t := bot.SnowflakeToTime(c.MS.ID)

	humanized := common.HumanizeDuration(common.DurationPrecisionHours, time.Since(t))
	if humanized == "" {
		humanized = "Less than an hour"
	}

	return humanized
}

func (c *Context) tmplCurrentUserAgeMinutes() int {
	t := bot.SnowflakeToTime(c.MS.ID)
	d := time.Since(t)

	return int(d.Seconds() / 60)
}

func (c *Context) tmplCurrentUserCreated() time.Time {
	t := bot.SnowflakeToTime(c.MS.ID)
	return t
}

func (c *Context) tmplSleep(duration interface{}) (string, error) {
	seconds := tmplToInt(duration)
	if c.secondsSlept+seconds > 60 || seconds < 1 {
		return "", errors.New("can sleep for max 60 seconds combined")
	}

	c.secondsSlept += seconds
	time.Sleep(time.Duration(seconds) * time.Second)
	return "", nil
}

func (c *Context) compileRegex(r string) (*regexp.Regexp, error) {
	if c.RegexCache == nil {
		c.RegexCache = make(map[string]*regexp.Regexp)
	}

	cached, ok := c.RegexCache[r]
	if ok {
		return cached, nil
	}

	if len(c.RegexCache) >= 20 {
		return nil, ErrTooManyAPICalls
	}

	compiled, err := regexp.Compile(r)
	if err != nil {
		return nil, err
	}

	c.RegexCache[r] = compiled

	return compiled, nil
}

func (c *Context) reFind(r string, s string) (string, error) {
	compiled, err := c.compileRegex(r)
	if err != nil {
		return "", err
	}

	return compiled.FindString(s), nil
}

func (c *Context) reFindAll(r string, s string) ([]string, error) {
	compiled, err := c.compileRegex(r)
	if err != nil {
		return nil, err
	}

	return compiled.FindAllString(s, 1000), nil
}

func (c *Context) reFindAllSubmatches(r string, s string) ([][]string, error) {
	compiled, err := c.compileRegex(r)
	if err != nil {
		return nil, err
	}

	return compiled.FindAllStringSubmatch(s, 100), nil
}

func (c *Context) reReplace(r string, s string, repl string) (string, error) {
	compiled, err := c.compileRegex(r)
	if err != nil {
		return "", err
	}

	return compiled.ReplaceAllString(s, repl), nil
}

func (c *Context) reSplit(r, s string, i int) ([]string, error) {
 	compiled, err := c.compileRegex(r)
 	if err != nil {
 		return nil, err
 	}

 	return compiled.Split(s, i), nil
 }

func (c *Context) tmplEditChannelName(channel interface{}, newName string) (string, error) {
	if c.IncreaseCheckCallCounter("edit_channel", 10) {
		return "", ErrTooManyCalls
	}

	cID := c.ChannelArgNoDM(channel)
	if cID == 0 {
		return "", errors.New("Unknown channel")
	}

	if c.IncreaseCheckCallCounter("edit_channel_"+strconv.FormatInt(cID, 10), 2) {
		return "", ErrTooManyCalls
	}

	_, err := common.BotSession.ChannelEdit(cID, newName)
	return "", err
}

func (c *Context) tmplEditChannelTopic(channel interface{}, newTopic string) (string, error) {
	if c.IncreaseCheckCallCounter("edit_channel", 10) {
		return "", ErrTooManyCalls
	}

	cID := c.ChannelArgNoDM(channel)
	if cID == 0 {
		return "", errors.New("Unknown channel")
	}

	if c.IncreaseCheckCallCounter("edit_channel_"+strconv.FormatInt(cID, 10), 2) {
		return "", ErrTooManyCalls
	}

	edit := &discordgo.ChannelEdit{
		Topic: newTopic,
	}

	_, err := common.BotSession.ChannelEditComplex(cID, edit)
	return "", err
}

func (c *Context) tmplOnlineCount() (int, error) {
	if c.IncreaseCheckCallCounter("online_users", 1) {
		return 0, ErrTooManyCalls
	}

	online := 0
	c.GS.RLock()
	for _, v := range c.GS.Members {
		if v.PresenceSet && v.PresenceStatus != dstate.StatusOffline {
			online++
		}
	}
	c.GS.RUnlock()

	return online, nil
}

func (c *Context) tmplOnlineCountBots() (int, error) {
	if c.IncreaseCheckCallCounter("online_bots", 1) {
		return 0, ErrTooManyCalls
	}

	botCount := 0

	c.GS.RLock()
	defer c.GS.RUnlock()

	for _, v := range c.GS.Members {
		if v.Bot && v.PresenceSet && v.PresenceStatus != dstate.StatusOffline {
			botCount++
		}
	}

	return botCount, nil
}

func (c *Context) tmplEditNickname(user interface{}, nickname string) (string, error) {
  if c.IncreaseCheckCallCounter("edit_nick", 3) {
    return "", ErrTooManyCalls
  }

  id := targetUserID(user)
  if id == 0 {
    return "", nil
  }
  
  if c.MS.ID != id && !bot.IsSpecialGuild(c.GS.ID) {
 	 return "", errors.New("**Error:** server not whitelisted, can't change other users nicknames. ")
  }

  err := common.BotSession.GuildMemberNickname(c.GS.ID, id, nickname)
  return "", err
}

type StdDepth struct {
 	depth int
 }

 func newStdDepth() *StdDepth {
 	return &StdDepth{
 		depth: 0,
 	}
 }

 func (sd *StdDepth) Add() {
 	sd.depth += 1
 }

 func (c *Context) tmplStandardize(input interface{}) interface{} {
 	depth := newStdDepth()
 	return depth.StdInit(input)
 }

 func (sd *StdDepth) StdInit(input interface{}) interface{} {
 	val := reflect.ValueOf(input)
 	switch val.Kind() {
 	case reflect.Map:
 		switch t := input.(type) {
 		case Dict, SDict:
 			return t
 		case map[string]interface{}:
 			return sd.StdStringMap(t)
 		default:
 			return sd.StdMap(t)
 		}
 	default:
 		switch t := input.(type) {
 		case []interface{}:
 			return sd.StdSlice(t)
 		case *time.Time:
 			return *t
 		case float64:
 			i, err := strconv.ParseInt(fmt.Sprint(t), 10, 64)
 			if err != nil {
 				return t
 			}
 			return i
 		default:
 			return t
 		}
 	}
 }

 func (sd *StdDepth) StdMap(input interface{}) interface{} {
 	if sd.depth >= 10 {
 		return input
 	}
 	sd.Add()
 	out := make(Dict)
 	val := reflect.ValueOf(input)
 	switch val.Kind() {
 	case reflect.Map:
 		for _, k := range val.MapKeys() {
 			v := val.MapIndex(k)
 			switch t := v.Interface().(type) {
 			case map[interface{}]interface{}:
 				out.Set(k.Interface(), sd.StdMap(t))
 			case map[string]interface{}:
 				out.Set(k.Interface(), sd.StdStringMap(t))
 			case []interface{}:
 				out.Set(k.Interface(), sd.StdSlice(t))
 			default:
 				out.Set(k.Interface(), sd.StdInit(t))
 			}
 		}
 		return out
 	}
 	return nil
 }

 func (sd *StdDepth) StdStringMap(input interface{}) interface{} {
 	if sd.depth >= 10 {
 		return input
 	}
 	sd.Add()
 	out := make(SDict)
 	val := reflect.ValueOf(input)
 	switch val.Kind() {
 	case reflect.Map:
 		for _, k := range val.MapKeys() {
 			v := val.MapIndex(k)
 			switch t := v.Interface().(type) {
 			case map[interface{}]interface{}:
 				out.Set(k.Interface().(string), sd.StdMap(t))
 			case map[string]interface{}:
 				out.Set(k.Interface().(string), sd.StdStringMap(t))
 			case []interface{}:
 				out.Set(k.Interface().(string), sd.StdSlice(t))
 			default:
 				out.Set(k.Interface().(string), sd.StdInit(t))
 			}
 		}
 		return out
 	}
 	return nil
 }

 func (sd *StdDepth) StdSlice(input []interface{}) interface{} {
 	if sd.depth >= 10 {
 		return input
 	}
 	sd.Add()
 	var out Slice
 	for _, v := range input {
 		out = append(out, sd.StdInit(v))
 	}

 	return out
 }