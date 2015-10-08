package config

import (
	"encoding/json"
	"regexp"
)

//Regex is a wrapper for regexp.Regexp but implements json marshaler interface
type Regex struct {
	*regexp.Regexp
}

//UnmarshalJSON turns a string into proper regex
func (r *Regex) UnmarshalJSON(b []byte) error {
	str := new(string)
	json.Unmarshal(b, str)
	reg, err := regexp.Compile(*str)
	if err != nil {
		return err
	}
	r.Regexp = reg
	return nil
}

//MarshalJSON turns regex into string/json
func (r *Regex) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Regexp.String())
}
