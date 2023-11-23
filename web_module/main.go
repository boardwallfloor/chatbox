package main

import (
	"chatbox/web_module/pb"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type MessageBuilder struct {
	msg pb.ChatMessage
}

func (bld *MessageBuilder) AddMessage(msg string) *MessageBuilder {
	bld.msg.Content = msg
	return bld
}

type HtmlTemplate struct{}

func (h *HtmlTemplate) Render(w io.Writer, name string, data interface{}, ec echo.Context) error {
	tmpl, err := template.ParseFiles(name)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return tmpl.Execute(w, data)
}

func sendDummyChat(ws *websocket.Conn, errChan chan<- error) {
	type ChatMessage struct {
		Sender string
		Msg    string
	}
	chatMessages := []ChatMessage{
		{Sender: "User1", Msg: "Do you know how to make a great cup of coffee?"},
		{Sender: "User2", Msg: "Absolutely! There are various methods. The French press is one of my favorites."},
		{Sender: "User3", Msg: "I prefer pour-over coffee. It allows for more control over the brewing process."},
		{Sender: "User4", Msg: "Espresso is my go-to. It's all about that strong, concentrated flavor."},
		{Sender: "User1", Msg: "I've heard about the French press. How do you use it?"},
		{
			Sender: "User2",
			Msg:    "It's simple! Coarsely grind your coffee beans, add hot water, and let it steep for a few minutes before pressing the plunger.",
		},
		{
			Sender: "User3",
			Msg:    "With pour-over, it's crucial to use the right water temperature and pour slowly to ensure even extraction.",
		},
		{Sender: "User4", Msg: "Espresso requires finely ground coffee and precise timing. The pressure creates a rich crema."},
		{Sender: "User1", Msg: "Sounds good! What about Turkish coffee?"},
		{
			Sender: "User2",
			Msg:    "Turkish coffee is unique! You need finely ground coffee, sugar, and cardamom. It's brewed in a special pot called a cezve.",
		},
		{Sender: "User3", Msg: "Turkish coffee is all about tradition. The key is to brew it slowly over low heat."},
		{Sender: "User4", Msg: "I love Turkish coffee's strong flavor and the ritual of brewing it."},
		{Sender: "User1", Msg: "Thanks for the tips! I'll give those methods a try."},
		{Sender: "User2", Msg: "You're welcome! Enjoy experimenting with different brewing methods."},
		{Sender: "User3", Msg: "Let us know how your coffee adventures go!"},
		{Sender: "User4", Msg: "Coffee brings people together. Share your coffee stories with us!"},
	}
	tmpl, err := template.ParseFiles("./template/msg_chat.html")
	if err != nil {
		log.Fatal(tmpl)
	}
	for _, msg := range chatMessages {
		chat := struct {
			Message string
		}{
			Message: msg.Sender + " : " + msg.Msg,
		}
		var resp strings.Builder
		err := tmpl.Execute(&resp, chat)
		respString := resp.String()
		err = ws.WriteMessage(websocket.TextMessage, []byte(respString))
		if err != nil {
			errChan <- err
		}
		time.Sleep(5 * time.Second)
	}
}

// func getUnsentMessages()

func main() {
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	ec := echo.New()
	ec.Static("assets", "./template/static/")
	ec.Renderer = &HtmlTemplate{}
	ec.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "./template/home.html", "")
	})
	ec.POST("/enter-room", func(c echo.Context) error {
		roomName := c.FormValue("roomName")
		roomData := struct {
			RoomTitle string
		}{
			RoomTitle: roomName,
		}
		return c.Render(http.StatusOK, "./template/client.html", roomData)
	})
	ec.GET("/ws", func(c echo.Context) error {
		// check for unseen / not received message
		// - call grpc function to get latest unseen message by sending last seen message date(?)
		// -

		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return nil
		}
		defer ws.Close()
		for {
			errChan := make(chan error)
			go sendDummyChat(ws, errChan)
			select {
			case err := <-errChan:
				if err != nil {
					return err
				}
			default:
				mstype, msg, err := ws.ReadMessage()
				if err != nil {
					return nil
				}
				log.Println(string(msg))

				resp := `<div hx-swap-oob="beforeend:#chat-container"><p>` + string(msg) + `</p></div>`
				err = ws.WriteMessage(mstype, []byte(resp))
				if err != nil {
					return nil
				}
			}
		}
	})
	ec.Logger.Fatal(ec.Start(":8080"))
}
