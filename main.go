package main

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    "io"
    "log"
    "os"
    "os/exec"
    "strings"
    "time"
)

var(
    Token         = "Bot "+os.Getenv("DISCORD_TOKEN")
    WebHookURL    = os.Getenv("DISCORD_WEB_HOOK_URL")
    stopBot       = make(chan bool)
    NoticeCommand = "!notice "
)

func main()  {
    discord, err := discordgo.New()
    if err != nil {
        fmt.Println("Error logging in")
        fmt.Println(err)
    }
    discord.Token = Token

    discord.AddHandler(onMessageCreate)
    err = discord.Open()
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println("Listening...")
    <- stopBot
    return
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate){
    c, err := s.State.Channel(m.ChannelID)
    if err != nil {
        log.Println("Error getting channel: ", err)
        return
    }
    fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

    switch {
    case strings.HasPrefix(m.Content, NoticeCommand):
        content := strings.Replace(m.Content, NoticeCommand, "", -1)
        command := strings.SplitN(content, " ", 2)
        if len(command) != 2 {
            sendMessage(s, c, "Invalid arguments. ")
            return
        }
        reserveNotice(command[0], command[1])
    }
}

func sendMessage(session *discordgo.Session, channel *discordgo.Channel, message string)  {
    _, err := session.ChannelMessageSend(channel.ID, message)

    log.Println(message)
    if err != nil {
        log.Fatal("Error sending message: ", err)
    }
}

func reserveNotice(time string, message string){
    curl := exec.Command("at", time)
    stdin, err := curl.StdinPipe()
    if err != nil {
        log.Fatal(err)
    }
    go func() {
        defer stdin.Close()
        _, err := io.WriteString(stdin, "curl -H 'Accept: application/json' -H 'Content-type: application/json' -X POST -d '{\"content\": \""+message+"\"}' "+WebHookURL+"\n")
        if err != nil {
            log.Fatal(err)
        }
    }()
    out, err := curl.CombinedOutput()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(out))
}
