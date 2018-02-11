package commands

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"yuzu"
)

// Cat f
type Cat struct{}

// IsOwnerOnly f
func (Cat) IsOwnerOnly() bool {
	return false
}

// Help f
func (Cat) Help() [2]string {
	return [2]string{"Makes a request to the random.cat api", ""}
}

// Process f
func (Cat) Process(ctx yuzu.Context) {
	res, err := http.Get("http://random.cat/meow")
	if err != nil {
		_, e := ctx.Error(err)
		if e != nil {
			return
		}
		return
	}
	defer res.Body.Close()

	catJSON := struct {
		File string `json:"file"`
	}{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		_, e := ctx.Error(err)
		if e != nil {
			return
		}
		return
	}
	if err := json.Unmarshal(body, &catJSON); err != nil {
		_, e := ctx.Error(err)
		if e != nil {
			return
		}
		return
	}
	ctx.Say(catJSON.File)
}
