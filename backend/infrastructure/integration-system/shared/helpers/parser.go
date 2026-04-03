package helpers

import (
	"encoding/json"
	"fmt"
	"time"
)

type IikoTime struct {
	time.Time
}

var iikoTimeLayouts = []string{
	"2006-01-02 15:04:05.000",
	time.RFC3339,
}

func (t *IikoTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" || len(b) == 0 {
		return nil
	}

	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}

	for _, layout := range iikoTimeLayouts {
		if parsed, err := time.Parse(layout, s); err == nil {
			t.Time = parsed
			return nil
		}
	}

	return fmt.Errorf("Не поддерживаемый формат datetime: %q", s)
}
