package osexec

import (
	"bytes"
	"encoding/xml"
	"io"
	"os/exec"
	"strings"
	"testing"
)

type Payload struct {
	Message string `xml:"message"`
}

func GetData(data io.Reader) string {
	var payload Payload
	xml.NewDecoder(data).Decode(&payload)

	return strings.ToUpper(payload.Message)
}

func GetXmlFromCommand() io.Reader {
	cmd := exec.Command("cat", "msg.xml")
	out, _ := cmd.StdoutPipe()

	cmd.Start()
	data, _ := io.ReadAll(out)
	cmd.Wait()

	return bytes.NewReader(data)
}

func TestGetData(t *testing.T) {
	got := GetData(GetXmlFromCommand())
	want := "HAPPY NEW YEAR!"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestGetDataUnit(t *testing.T) {
	input := strings.NewReader(`
<payload>
    <message>Cats are the best animal</message>
</payload>
    `)

	got := GetData(input)
	want := "CATS ARE THE BEST ANIMAL"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
