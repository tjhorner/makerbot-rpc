package makerbot

import (
	"strconv"
	"strings"
	"time"
)

// https://gist.github.com/alexmcroberts/219127816e7a16c7bd70
type epochTime time.Time

func (t epochTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

func (t *epochTime) UnmarshalJSON(s []byte) (err error) {
	r := strings.Replace(string(s), `"`, ``, -1)

	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q/1000, 0)
	return
}

func (t epochTime) String() string { return time.Time(t).String() }
