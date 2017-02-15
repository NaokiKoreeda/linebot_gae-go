package main

import (
    "net/http"

    "github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/appengine"
    "google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func init() {

    g := e.Group("/bot")
    g.Use(middleware.CORS())

    g.GET("", helloMsg)
    g.POST("", postBot)
}

func helloMsg(c echo.Context) error {
    return c.JSON(http.StatusOK, "hello,bot!")
}

func postBot(c echo.Context) error {
    channelSecret := "{your_ChannelSecret}"
    channelAccessToken := "{your_ChannelAccessToken}"

    cx := appengine.NewContext(c.Request().(*standard.Request).Request)

    client := urlfetch.Client(cx)
    bot, err := linebot.New(channelSecret, channelAccessToken, linebot.WithHTTPClient(client))
    if err != nil {
        log.Infof(cx, err.Error())
        return c.JSON(http.StatusInternalServerError, err)
    }

    received, err := bot.ParseRequest(c.Request().(*standard.Request).Request)
    if err != nil {
        log.Infof(cx, err.Error())
        return c.JSON(http.StatusInternalServerError, err)
    }

    for _, event := range received {
        if event.Type == linebot.EventTypeMessage {
            switch message := event.Message.(type) {
            case *linebot.TextMessage:
                log.Infof(cx, "TextMessage %#v", message)
                postMessage := linebot.NewTextMessage(message.Text)
                if _, err = bot.ReplyMessage(event.ReplyToken, postMessage).Do(); err != nil {
                    log.Errorf(cx, "send error: %v", err)
                }
            }
        }
    }

    return c.JSON(http.StatusOK, "success")
}
