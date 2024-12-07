package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"pgpserver/internal/handlekey"

	"github.com/bwmarrin/discordgo"
)

var (
    appId = "1314714143382831144"
    guildId = ""
)

func Run(BotToken string) {
    discord, err := discordgo.New(("Bot " + BotToken))
    if err != nil { fmt.Println("Bot 1"); log.Fatal(err) }

    _, err = discord.ApplicationCommandBulkOverwrite(appId, guildId, []*discordgo.ApplicationCommand {
        {
            Name: "update-key",
            Description: "Submit a pgp key to the database tied to your discord account",
            Options: []*discordgo.ApplicationCommandOption {
                {
                    Type: discordgo.ApplicationCommandOptionAttachment,
                    Name: "key",
                    Description: "pgp public key",
                    Required: true,
                },
            },
        },
        {
            Name: "get-key",
            Description: "Get pgp key of user by username",
            Options: []*discordgo.ApplicationCommandOption {
                {
                    Type: discordgo.ApplicationCommandOptionString,
                    Name: "username",
                    Description: "username of requested user",
                    Required: true,
                },
            },
        },
    })
    if err != nil { fmt.Println("Bot 2"); log.Fatal(err) }

    discord.AddHandler(func (
        s *discordgo.Session,
        i *discordgo.InteractionCreate,
    ) {
        if i.Type == discordgo.InteractionApplicationCommand {
            data := i.ApplicationCommandData()
            switch data.Name {
            case "update-key":
                if i.Interaction.Member.User.ID == s.State.User.ID { return; }
                responseData := ""

                var attachmentUrl string
                for _, v := range i.Interaction.ApplicationCommandData().Options {
                    switch v.Name {
                    case "key":
                        attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
                        attachmentUrl = i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL
                    }
                }
                responseData = handlekey.UpdateKey(attachmentUrl, i.Interaction.Member.User.Username);

                err = s.InteractionRespond(
                    i.Interaction,
                    &discordgo.InteractionResponse{
                        Type: discordgo.InteractionResponseChannelMessageWithSource,
                        Data: &discordgo.InteractionResponseData{
                            Flags: 1 << 6,
                            Content: responseData,
                        },
                    },
                )
                if err != nil { fmt.Println("Bot 3"); log.Fatal(err) }
            case "get-key":
                if i.Interaction.Member.User.ID == s.State.User.ID { return; }
                responseData := ""

                var username string
                for _, v := range i.Interaction.ApplicationCommandData().Options {
                    switch v.Name {
                    case "username":
                        username = v.StringValue()
                    }
                }
                responseData = handlekey.GetKey(username);

                err = s.InteractionRespond(
                    i.Interaction,
                    &discordgo.InteractionResponse{
                        Type: discordgo.InteractionResponseChannelMessageWithSource,
                        Data: &discordgo.InteractionResponseData{
                            Flags: 1 << 6,
                            Content: responseData,
                        },
                    },
                )
            }
        }
    })


    err = discord.Open()
    if err != nil { log.Fatal(err) }

    stop := make (chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt)
    log.Println("Press Ctrl+C to Exit")
    <-stop

    err = discord.Close()
    if err != nil { log.Fatal(err) }
}
