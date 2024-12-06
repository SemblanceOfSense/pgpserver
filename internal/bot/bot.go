package bot

import (
	"log"
	"pgpserver/internal/handlekey"

	"github.com/bwmarrin/discordgo"
)

var (
    appId = ""
    guildId = ""
)

func Run(BotToken string) {
    discord, err := discordgo.New(("Bot " + BotToken))
    if err != nil { log.Fatal(err) }

    _, err = discord.ApplicationCommandBulkOverwrite(appId, guildId, []*discordgo.ApplicationCommand {
        {
            Name: "update-public-key",
            Description: "Submit a pgp key to the database tied to your discord account",
            Options: []*discordgo.ApplicationCommandOption {
                {
                    Type: discordgo.ApplicationCommandOptionString,
                    Name: "key",
                    Description: "PGP Public Key",
                    Required: true,
                },
            },
        },
    })
    if err != nil { log.Fatal(err) }

    discord.AddHandler(func (
        s *discordgo.Session,
        i *discordgo.InteractionCreate,
    ) {
        if i.Type == discordgo.InteractionApplicationCommand {
            data := i.ApplicationCommandData()
            switch data.Name {
            case "update-public-key":
                if i.Interaction.Member.User.ID == s.State.User.ID { return; }
                responseData := ""

                var key string
                for _, v := range i.Interaction.ApplicationCommandData().Options {
                    switch v.Name {
                    case "prompt":
                        key = v.StringValue()
                    }
                }
                responseData = handlekey.UpdateKey(key, i.Interaction.Member.User.ID);

                err = s.InteractionRespond(
                    i.Interaction,
                    &discordgo.InteractionResponse{
                        Type: discordgo.InteractionResponseChannelMessageWithSource,
                        Data: &discordgo.InteractionResponseData{
                            Content: responseData,
                        },
                    },
                )
            }
        }
    })
}
