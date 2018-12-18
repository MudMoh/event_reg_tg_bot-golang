package post

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"telegram"
)

const postURL = "https://" //my server url

func SendPost(data *telegram.EventSign) error { //POST request to server
	params := fmt.Sprintf(
		"name=%s&email=%s&phone=%s&event=%s",
		data.Name,
		data.Email,
		data.Phone,
		data.Event)
	buffer := bytes.NewBufferString(params)
	res, err := http.Post(
		postURL,
		"application/x-www-form-urlencoded",
		buffer)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	ba, err := ioutil.ReadAll(res.Body)
	fmt.Printf("response: %s\n", ba)

	if res.StatusCode != 200 {
		return errors.New("it's not 200 response")
	}
	return nil
}
