package handlekey

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

func UpdateKey(url string, userId string) string {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil { return "Failed to make attachment request" }

    resp, err := http.DefaultClient.Do(req)
    if err != nil { return "Failed to get attachment" }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil { return "Bad data" }
    key := string(body)

    fmt.Println(key)
    if key == "" { return "Please submit a key." }

    _, err = crypto.NewKeyFromArmored(key)
    if err != nil { fmt.Println(err); return "Not a PGP key" }

    if _, err := os.Stat("/etc/pgpbot/" + userId); err == nil {
        err = os.Remove("/etc/pgpbot/" + userId)
        if err != nil { return "Server file error" }
    }
    file, err := os.Create("/etc/pgpbot/" + userId)
    if err != nil { fmt.Println(err); return "Server file creation error. Your entry may have been deleted." }
    defer file.Close()

    _, err = file.Write([]byte(key))
    if err != nil { return "Server file write error. Your entry may hav been deleted." }

    return "Successful submission"
}
