package handlekey

import (
	"os"
	"regexp"
)

func UpdateKey(key string, userId string) string {
    if key == "" { return "Please submit a key." }

    // Shoutout to angeloped's amazing regex
    regex := `/^(-----BEGIN PGP PUBLIC KEY BLOCK-----).*([a-zA-Z0-9//\n\/\.\:\+\ \=]+).*(-----END PGP PUBLIC KEY BLOCK-----)$|^(-----BEGIN PGP PRIVATE KEY BLOCK-----).*([a-zA-Z0-9//\n\/\.\:\+\ \=]+).*(-----END PGP PRIVATE KEY BLOCK-----)$/`
    match, err := regexp.MatchString(regex, key)
    if err != nil { return "Server Regex Error" }
    if !match { return "Not a PGP key" }

    if _, err := os.Stat("/etc/pgpbot/" + userId); err == nil {
        err = os.Remove("/etc/pgpbot/" + userId)
        if err != nil { return "Server file error" }
    }
    file, err := os.Create("/etc/pgpbot/" + userId)
    if err != nil { return "Server file creation error. Your entry may have been deleted." }
    defer file.Close()

    _, err = file.Write([]byte(key))
    if err != nil { return "Server file write error. Your entry may hav been deleted." }

    return "Successful submission"
}
