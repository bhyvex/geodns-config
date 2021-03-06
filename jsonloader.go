package dnsconfig

import (
	"encoding/json"
	"fmt"
	"github.com/abh/errorutil"
	"log"
	"os"
	"strconv"
)

type objMap map[string]interface{}

func jsonLoader(fileName string, objmap objMap, fn func() error) error {
	fh, err := os.Open(fileName)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(fh)
	if err = decoder.Decode(&objmap); err != nil {
		extra := ""
		if serr, ok := err.(*json.SyntaxError); ok {
			if _, serr := fh.Seek(0, os.SEEK_SET); serr != nil {
				log.Fatalf("seek error: %v", serr)
			}
			line, col, highlight := errorutil.HighlightBytePosition(fh, serr.Offset)
			extra = fmt.Sprintf(":\nError at line %d, column %d (file offset %d):\n%s",
				line, col, serr.Offset, highlight)
		}
		return fmt.Errorf("error parsing JSON object in config file %s%s\n%v",
			fh.Name(), extra, err)
	}

	err = fn()
	return err

}

func toInt(i interface{}) (int, error) {
	switch i.(type) {
	case string:
		return strconv.Atoi(i.(string))
	case float64:
		return int(i.(float64)), nil
	case nil:
		return 0, nil
	}
	return 0, fmt.Errorf("Unknown type %T", i)
}

func toBool(i interface{}) (bool, error) {
	n, err := toInt(i)
	if err != nil {
		return false, err
	}
	if n > 0 {
		return true, nil
	}
	return false, nil
}
