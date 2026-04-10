package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func Logger(data []byte) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err == nil {
		fmt.Println("Pretty JSON Response:")
		fmt.Println(prettyJSON.String())
	} else {
		fmt.Println("Raw Response:", string(data))
	}
}
