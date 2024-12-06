package handlekey

import (
	"os"
	"regexp"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

func UpdateKey(key string, userId string) string {
    if key == "" { return "Please submit a key." }

    _, err := crypto.NewKeyFromArmored(key)
    if err != nil { return "Not a PGP key" }

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
