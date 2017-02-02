package meetup

import (
	"encoding/json"
	"fmt"
	"github.com/pmylund/go-cache"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	c      *cache.Cache
	timezone *time.Location
)

func init() {
	c = cache.New(5*time.Minute, 30*time.Second)
	timezone, _ = time.LoadLocation("America/Chicago")
	fmt.Println(timezone)
}

func FilterEventsByKeyword(keyword string) (results []Event) {
	keyword = strings.TrimSpace(strings.ToLower(keyword))
	r, found := c.Get(keyword)
	if !found {
		log.Println("Cache miss for [", keyword, "]")
		latest, err := getLatestEvents()
		if err != nil {
			log.Println(err)
			return results
		}

		for _, element := range latest.Results {
			if strings.Index(strings.ToLower(element.Name), keyword) > -1 {
				results = append(results, element)
			}
		}

		c.Set(keyword, results, cache.DefaultExpiration)
	} else {
		results = r.([]Event)
	}

	return results
}

func GetEventsForMonth(month int, year int) (results []Event) {
	t := time.Date(year, time.Month(month), 1, 0, 0,0,0, timezone)
	timeframe := fmt.Sprintf("%d,1m", t.Unix())
	r, found := c.Get(timeframe)
	if !found {
		log.Println("Cache miss for [", timeframe, "]")
		search, err := getEventsForTimeframe(timeframe)
		if err != nil {
			log.Println(err)
		}

		c.Set(timeframe, search, cache.DefaultExpiration)
		r = search
	}

	return r.([]Event)
}

func GetEventsForWeek(day int, month int, year int) ([]Event) {
	t := time.Date(year, time.Month(month), day, 0, 0,0,0, timezone)
	timeframe := fmt.Sprintf("%d,1w", t.Unix())
	r, found := c.Get(timeframe)
	if !found {
		log.Println("Cache miss for [", timeframe, "]")
		search, err := getEventsForTimeframe(timeframe)
		if err != nil {
			log.Println(err)
		}

		c.Set(timeframe, search, cache.DefaultExpiration)
		r = search
	}

	return r.([]Event)
}

func GetEventsForDay(day int, month int, year int) ([]Event) {
	t := time.Date(year, time.Month(month), day, 0, 0,0,0, timezone)
	timeframe := fmt.Sprintf("%d,1d", t.Unix())
	r, found := c.Get(timeframe)
	if !found {
		log.Println("Cache miss for [", timeframe, "]")
		search, err := getEventsForTimeframe(timeframe)
		if err != nil {
			log.Println(err)
		}

		c.Set(timeframe, search, cache.DefaultExpiration)
		r = search
	}

	return r.([]Event)
}

func getEventsForTimeframe(timeframe string) ([]Event, error) {
	group_id := os.Getenv("MEETUP_GROUP_ID")
	api_key := os.Getenv("MEETUP_API_KEY")

	url := fmt.Sprintf("https://api.meetup.com/2/events?group_id=%s&key=%s&time=%s", group_id, api_key, timeframe)
	var results []Event
	for url != "" {
		resp, _ := http.Get(url)

		events := new(Events)
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(events)
		results = append(results, (events.Results)...)
		url = events.Meta.Next

	}

	return results, nil
}

func getLatestEvents() (Events, error) {
	group_id := os.Getenv("MEETUP_GROUP_ID")
	api_key := os.Getenv("MEETUP_API_KEY")
	url := fmt.Sprintf("https://api.meetup.com/2/events?group_id=%s&key=%s", group_id, api_key)
	resp, err := http.Get(url)

	events := new(Events)
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(events)

	if err != nil {
		return *events, err
	}

	return *events, nil
}
