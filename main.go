package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"math/rand"
	"strings"
	"time"
)

var rtm *slack.RTM
var quiteMessages = []string{
	"I hope you get cramp in your tongue",
	"What kind of noise is this?",
	"SILENCE!!!!! please?",
	"Do you talk in your sleep?",
	"I know you have a constitutional right to speak.",
	"I really miss Charlie Chaplin movies.",
	"Oh how I wish I were deaf!",
	"Your mom must be so proud of your social skills",
	"Somebody forget to mute that loudspeaker",
	"Are you punishing me?",
	"Volume needs turned down.",
	"hey, hey! hey! we get it, the knobs go up to eleven.",
	"You should try breathing.",
	"Do me a favour, can you move this noise to my enemies."}

func main() {
	api := slack.New("<a generated token from slack>")

	rtm = api.NewRTM()

	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		//fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore event
		case *slack.ConnectedEvent:
			// Ignore event
		case *slack.MessageEvent:
			//fmt.Printf("Message: %v\n", ev)
			processMessage(ev)
		case *slack.PresenceChangeEvent:
			//fmt.Printf("Presence Change: %v\n", ev)
		case *slack.LatencyReport:
			//fmt.Printf("Current latency: %v\n", ev.Value)
		case *slack.DesktopNotificationEvent:
			//fmt.Printf("Desktop Notification: %v\n", ev)
		case *slack.RTMError:
			//fmt.Printf("Error: %s\n", ev.Error())
		case *slack.InvalidAuthEvent:
			//fmt.Printf("Invalid credentials")
			return
		default:
			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

func processMessage(message *slack.MessageEvent) {
	if strings.HasPrefix(strings.ToLower(message.Text), "hush ") && message.User != "<this should be your slackbot user>" {
		fmt.Println("_______________________________")
		fmt.Printf("Recieved Message from User: %s\n", message.User)
		fmt.Printf("Message Channel : %s\n", message.Channel)
		fmt.Printf("Message Text: %s\n", message.Text)
		namesOrChannels := parseMessage(message.User, strings.TrimLeft(message.Text[4:], "hush "))
		messageProcessingSuccessfully := messagesSuccessfullyParsed(namesOrChannels)
		for _, nameOrChannel := range namesOrChannels {
			messageProcessingSuccessfully = messageProcessingSuccessfully && sendMessage(nameOrChannel)
		}
		if messageProcessingSuccessfully {
			notifySenderOfSuccess(message.User)
		}
	}
}

func parseMessage(senderId string, message string) []string {
	//fmt.Printf("message sent to be processed: %s\n", message)
	var processedUserAndChannelNames []string
	userOrChannelNames := strings.Fields(message)
	for _, userOrChannelName := range userOrChannelNames {
		userOrChannelName = strings.Trim(userOrChannelName, "< >")
		if userOrChannelName[0] == '@' {
			fmt.Printf("got user: %s\n", userOrChannelName)
			processedUserAndChannelNames = append(processedUserAndChannelNames, userOrChannelName)
		} else if userOrChannelName[0] == '#' {
			channelId := strings.Split(userOrChannelName, "|")
			processedChannelId := strings.ReplaceAll(channelId[0], "#", "")
			fmt.Printf("got channel: %s\n", processedChannelId)

			processedUserAndChannelNames = append(processedUserAndChannelNames, processedChannelId)
		} else {
			sendErrorMessage(senderId, userOrChannelName)
			return make([]string, 0)
		}
	}

	return processedUserAndChannelNames
}

func messagesSuccessfullyParsed(namesOrChannels []string) bool {
	return len(namesOrChannels) > 0
}

func sendMessage(channelOrNameID string) bool {
	fmt.Printf("sending message to %s\n", channelOrNameID)
	fmt.Println("_______________________________")

	_, _, err := rtm.PostMessage(channelOrNameID, slack.MsgOptionText(getRandomMessage(), false), slack.MsgOptionAsUser(true))
	if err != nil {
		fmt.Printf("%s\n", err)
		return false
	}
	return true
}

func getRandomMessage() string {
	rand.Seed(time.Now().UnixNano())
	return quiteMessages[rand.Intn(len(quiteMessages))]
}

func sendErrorMessage(senderId string, unprocessableString string) {
	fmt.Printf("sending error message %s\n", senderId)
	message := "there was an error processing the your request: " + unprocessableString + ". Messages to hushbot should say hush followed by single space references " +
		"to channels or users. eg: 'hush @user1 @user2 #channel1'"

	_, _, err := rtm.PostMessage(senderId, slack.MsgOptionText(message, false), slack.MsgOptionAsUser(true))
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

func notifySenderOfSuccess(userId string) {
	messageText := "Message(s) successfully sent."
	_, _, err := rtm.PostMessage(userId, slack.MsgOptionText(messageText, false), slack.MsgOptionAsUser(true))
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}
