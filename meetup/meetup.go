package meetup

import (
	"encoding/json"
	"fmt"
	"github.com/pmylund/go-cache"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	c        *cache.Cache
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
	t := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, timezone)
	timeframe := fmt.Sprintf("%d,%d", t.Unix()*1000, t.AddDate(0, 1, 0).Unix()*1000)
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

func GetEventsForWeek(day int, month int, year int) []Event {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, timezone)
	timeframe := fmt.Sprintf("%d,%d", t.Unix()*1000, t.AddDate(0, 0, 7).Unix()*1000)
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

func GetEventsForDay(day int, month int, year int) []Event {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, timezone)
	timeframe := fmt.Sprintf("%d,%d", t.Unix()*1000, t.AddDate(0, 0, 1).Unix()*1000)
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

func GetEventsForTimerange(timeframe string, groups ...string) []Event {
	if groups == nil {
		groups = []string{"memphis-technology-user-groups"}
	}
	cacheKey := timeframe + "-" + strings.Join(groups, ",")
	r, found := c.Get(cacheKey)
	if !found {
		log.Println("Cache miss for [", cacheKey, "]")
		search, err := getEventsForTimeframe(timeframe, groups...)
		if err != nil {
			log.Println(err)
		}

		c.Set(timeframe, search, cache.DefaultExpiration)
		r = search
	}

	return r.([]Event)
}

func getEventsForTimeframe(timeframe string, groups ...string) ([]Event, error) {
	api_key := os.Getenv("MEETUP_API_KEY")
	var results []Event
	if groups == nil {
		groups = []string{"memphis-technology-user-groups"}
	}

	for _, group := range groups {
		url := fmt.Sprintf("https://api.meetup.com/2/events?group_urlname=%s&key=%s&time=%s&status=upcoming,past", group, api_key, timeframe)
		for url != "" {
			resp, _ := http.Get(url)

			htmlData, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(os.Stdout, string(htmlData))

			events := new(Events)
			json.Unmarshal(htmlData, &events)
			results = append(results, events.Results...)
			url = events.Meta.Next

		}
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
