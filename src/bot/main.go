package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"net"
	"strings"

	"log"
	"net/http"

	slackbot "github.com/BeepBoopHQ/go-slackbot"
	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
)

const (
	WithTyping    = slackbot.WithTyping
	WithoutTyping = slackbot.WithoutTyping

	HelpText = "I will respond to the following messages: \n" +
		"`bot hi` for a simple message.\n" +
		"`bot attachment` to see a Slack attachment message.\n" +
		"`hey @<your bot's name>` to demonstrate detecting a mention.\n" +
		"`:smile:` some art.\n" +
		"`:wink:` some art.\n" +
		"`bot help` to see this again."
)

var greetingPrefixes = []string{"Hi", "Hello", "Howdy", "Wazzzup", "Hey"}
var bot *slackbot.Bot

func main() {


	bot = slackbot.New(os.Getenv("SLACK_TOKEN"))
	toMe := bot.Messages(slackbot.DirectMessage, slackbot.DirectMention).Subrouter()

	hi := "hi|hello|bot hi|bot hello"
	toMe.Hear(hi).MessageHandler(HelloHandler)
	bot.Hear(hi).MessageHandler(HelloHandler)
	bot.Hear("help|bot help").MessageHandler(HelpHandler)
	bot.Hear("attachment|bot attachment").MessageHandler(AttachmentsHandler)
	bot.Hear(`<@([a-zA-z0-9]+)?>`).MessageHandler(MentionHandler)
	bot.Hear("(bot ).*").MessageHandler(CatchAllHandler)
	bot.Hear(":wink:").MessageHandler(WinkHandler)
	bot.Hear(":smile:").MessageHandler(SmileHandler)
	bot.Hear("getAddress").MessageHandler(AddressHandler)
	go bot.Run()

	router := NewRouter(bot)
	http.Handle("/", router)
	errHTTP :=  http.ListenAndServe(":8080", nil)
	if errHTTP != nil {
		panic("Error: " + errHTTP.Error())
	}

}


type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(bot *slackbot.Bot) *mux.Router {
    router := mux.NewRouter().StrictSlash(false)
    for _, route := range routes {

        var handler http.Handler
        handler = route.HandlerFunc
				handler = Logger(handler, route.Name)
				handler = Botter(handler, bot)

        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)
    }
    return router
}

var routes = Routes{
  Route{
      "get",
      "GET",
      "/get",
      emptyHandler_func,
  },
	Route{
      "get2",
      "GET",
      "/ge",
      emptyHandler_func,
  },
}

func emptyHandler_func(rw http.ResponseWriter, req *http.Request){
	log.Printf("emptyHandler_func")
}

func Botter(inner http.Handler, bot *slackbot.Bot) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        inner.ServeHTTP(w, r)
				if(r.RequestURI == "/get"){
					bot.RTM.NewOutgoingMessage("Hello","@ascii-art-bot")
				}else{
					log.Printf("noooo",)
				}

    })
}

func Logger(inner http.Handler, name string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        inner.ServeHTTP(w, r)

        log.Printf(
            "%s\t%s\t%s\t%s",
            r.Method,
            r.RequestURI,
            name,
            time.Since(start),
        )
    })
}

func GetOutboundIP() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
				fmt.Println(err)
		}
    defer conn.Close()

    localAddr := conn.LocalAddr().String()
    idx := strings.LastIndex(localAddr, ":")

    return localAddr[0:idx]
}

func AddressHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	log.Printf("%s",evt,)
	bot.Reply(evt, GetOutboundIP(), WithTyping)
	bot.Reply(evt, "https://beepboophq.com/proxy/7024263b0799480bb2ebba99a059beac", WithTyping)
}

func HelloHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	rand.Seed(time.Now().UnixNano())
	msg := greetingPrefixes[rand.Intn(len(greetingPrefixes))] + " <@" + evt.User + ">!"
	bot.Reply(evt, msg, WithTyping)

	if slackbot.IsDirectMessage(evt) {
		dmMsg := "It's nice to talk to you directly."
		bot.Reply(evt, dmMsg, WithoutTyping)
	}

	bot.Reply(evt, "If you'd like to talk some more, "+HelpText, WithTyping)
}

func WinkHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	msg := fmt.Sprintf("░░░░██░░████████░░██░░░░░░░░░░░░░░░░░░░░░░░░░░\n"+
										 "░░██░░██░░░░░░░░██░░██░░░░░░░░░░░░░░░░░░░░░░░░\n"+
										 "░░██░░░░░░░░░░░░░░░░██░░░░░░░░░░░░░░░░░░░░░░░░\n"+
										 "░░██░░░░░░░░░░░░░░░░░░██░░░░░░░░░░░░░░░░░░░░░░\n"+
										 "██░░░░██░░░░██░░░░░░░░░░██░░░░░░░░░░░░░░░░░░░░\n"+
										 "██░░░░░░░░░░░░░░░░░░░░░░░░████░░░░░░░░░░░░░░░░\n"+
										 "██░░░░░░████░░░░░░░░░░░░░░░░░░██████████████░░\n"+
										 "██░░██░░██░░░░██░░░░░░░░░░░░░░░░░░░░░░░░░░░░██\n"+
										 "██░░░░████████░░░░░░░░░░░░░░░░░░░░░░░░██████░░\n"+
										 "██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██░░░░░░\n"+
										 "██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██░░░░░░\n"+
										 "██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██░░░░░░\n"+
										 "██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██░░░░░░\n"+
										 "██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██░░░░░░░░\n"+
										 "░░██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██░░░░░░░░\n"+
										 "░░██░░░░████░░░░████████░░░░████░░░░██░░░░░░░░\n"+
										 "░░██░░░░████░░██░░░░░░██░░██░░██░░░░██░░░░░░░░\n"+
										 "░░██░░██░░░░██░░░░░░░░░░██░░░░██░░██░░░░░░░░░░\n"+
										 "░░░░██░░░░░░░░░░░░░░░░░░░░░░░░░░██░░░░░░░░░░░░")
	bot.Reply(evt, msg, WithTyping)
}

func SmileHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	msg := fmt.Sprintf("░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░MMM88&&&,░░░░░░░░░\n"+
										" ░░░░░░,MMM8&&&.░░░░░░░░░░░░░░░░░░░░░░░░░░`'MMM88&&&,░░░░░\n"+
										" ░░░░░MMMMM88&&&&░░░░░░░░░░░░░░░░░░░░░░░░░░'MMM88&&&,░░░░░\n"+
										" ░░░░MMMMM88&&&&&&░░░░░░░░░░░░░░░░░░░░░░░░░░'MMM88&&&,░░░░\n"+
										" ░░░░MMMMM88&&&&&&░░░░░░░░░░░░░░░░░░░░░░░░░░░░'MMM88&&&░░░\n"+
										" ░░░░MMMMM88&&&&&&░░░░░░░░░░░░░░░░░░░░░░░░░░░░'MMM88&&&░░░\n"+
										" ░░░░░MMMMM88&&&&░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ MMM88&&&░░\n"+
										" ░░░░░░'MMM8&&&'░░░░░░░░░░░░MMMM888&&&&░░░░░░░░░░'MM88&&&░\n"+
										" ░░░░░░░░░░░░░░░░░░░░░░░░░░MMMM88&&&&&░░░░░░░░░░MM88&&&░░░\n"+
										" ░░░░░░░░░░░░░░░░░░░░░░░░░░MMMM88&&&&&░░░░░░░░░░MM88&&&░░░\n"+
										" ░░░░░░,MMM8&&&.░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░MM88&&&░░\n"+
										" ░░░░░MMMMM88&&&&░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░,MM88&&&░░░\n"+
										" ░░░░MMMMM88&&&&&&░░░░░░░░░░░░░░░░░░░░░░░░░░░░░MMM88&&&'░░\n"+
										" ░░░░MMMMM88&&&&&&░░░░░░░░░░░░░░░░░░░░░░░░░░░░░MMM88&&&'░░\n"+
										" ░░░░MMMMM88&&&&&&░░░░░░░░░░░░░░░░░░░░░░░░░░░░MMM88&&&'░░░\n"+
										" ░░░░MMMMM88&&&&░░░░░░░░░░░░░░░░░░░░░░░░░░░░░MMM88&&&'░░░░\n"+
										" ░░░░░'MMM8&&&'░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░MMM88&&&'░░░░\n"+
										" ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░MMM88&&&░░░░░░░'\n")
	bot.Reply(evt, msg, WithTyping)
}

func CatchAllHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	msg := fmt.Sprintf("I'm sorry, I don't know how to: `%s`.\n%s", evt.Text, HelpText)
	bot.Reply(evt, msg, WithTyping)
}

func MentionHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	if slackbot.IsMentioned(evt, bot.BotUserID()) {
		bot.Reply(evt, "You really do care about me. :heart:", WithTyping)
	}
}

func HelpHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	bot.Reply(evt, HelpText, WithTyping)
}

func AttachmentsHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	txt := "Beep Beep Boop is a ridiculously simple hosting platform for your Slackbots."
	attachment := slack.Attachment{
		Pretext:   "We bring bots to life. :sunglasses: :thumbsup:",
		Title:     "Host, deploy and share your bot in seconds.",
		TitleLink: "https://beepboophq.com/",
		Text:      txt,
		Fallback:  txt,
		ImageURL:  "https://storage.googleapis.com/beepboophq/_assets/bot-1.22f6fb.png",
		Color:     "#7CD197",
	}

	// supports multiple attachments
	attachments := []slack.Attachment{attachment}
	bot.ReplyWithAttachments(evt, attachments, WithTyping)
}
