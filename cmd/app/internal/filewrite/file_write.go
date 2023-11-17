package filewrite

import (
	"fmt"
	"os"
	"url-checker/cmd/app/internal/request"
)

var format = "\"%s\", \"%d\", \"%d\",\"%v\"\n"

func Print(file *os.File, urlsOrder []*string, results map[string]request.ResponseData) error {
	for _, url := range urlsOrder {
		info := results[*url]
		file.WriteString(fmt.Sprintf(format, *url, info.Size, info.Timing, info.Err))
	}

	return nil
}
