package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/memgo/api/meetup"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Host string `default:"0.0.0.0" envconfig:"HOST"`
	Port string `default:"8080" envconfig:"PORT"`
}

var (
	config Config
)

func main() {
	err := envconfig.Process("API", &config)
	if err != nil {
		log.Printf("Error processing config: %v\n", err.Error())
	}
	log.Println("Listening on port ", config.Port)

	r := mux.NewRouter()
	r.Handle("/", http.RedirectHandler("/calendar.json?keyword=memphis+ruby", 302))
	r.Handle("/favicon.ico", http.NotFoundHandler())
	r.HandleFunc("/slack/meetup", slackMeetup)
	r.HandleFunc("/calendar.json", calendarJson)
	r.HandleFunc("/calendar/day.json", calendarJsonDay)
	r.HandleFunc("/calendar/week.json", calendarJsonWeek)
	r.HandleFunc("/calendar/month.json", calendarJsonMonth)
	r.HandleFunc("/calendar/range.json", calendarJsonTimerange)
	r.HandleFunc("/{meetup}", meetupRedir)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", config.Host, config.Port), nil))
}

// Return a JSON blob of matching upcoming meetups
func calendarJson(w http.ResponseWriter, r *http.Request) {
	keyword := string(r.FormValue("keyword"))

	events := meetup.FilterEventsByKeyword(keyword)

	type CalendarJsonResponse struct {
		PR      string         `json:"_pull_requests_appreciated"`
		Meetups []meetup.Event `json:"meetups"`
	}

	data := CalendarJsonResponse{"https://github.com/memgo/api", events}
	marsh, _ := json.Marshal(data)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(marsh)
}

func calendarJsonDay(w http.ResponseWriter, r *http.Request) {
	year, _ := strconv.Atoi(r.FormValue("year"))
	month, _ := strconv.Atoi(r.FormValue("month"))
	day, _ := strconv.Atoi(r.FormValue("day"))

	events := meetup.GetEventsForDay(day, month, year)

	marsh, err := json.Marshal(events)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(marsh)
	}
}

func calendarJsonWeek(w http.ResponseWriter, r *http.Request) {
	year, _ := strconv.Atoi(r.FormValue("year"))
	month, _ := strconv.Atoi(r.FormValue("month"))
	day, _ := strconv.Atoi(r.FormValue("day"))

	events := meetup.GetEventsForWeek(day, month, year)

	marsh, err := json.Marshal(events)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(marsh)
	}

}

func calendarJsonMonth(w http.ResponseWriter, r *http.Request) {
	year, _ := strconv.Atoi(r.FormValue("year"))
	month, _ := strconv.Atoi(r.FormValue("month"))

	events := meetup.GetEventsForMonth(month, year)

	marsh, err := json.Marshal(events)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(marsh)
	}
}

func calendarJsonTimerange(w http.ResponseWriter, r *http.Request) {
	timerange := r.FormValue("timerange")
	groupsVal := r.FormValue("groups")
	if groupsVal == "" {
		groupsVal = "memphis-technology-user-groups"
	}
	groups := strings.Split(groupsVal, ",")

	events := meetup.GetEventsForTimerange(timerange, groups...)

	marsh, err := json.Marshal(events)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(marsh)
	}
}

// Handler for Slack's outgoing webhook for meetups: !meetup
func slackMeetup(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	trigger := r.FormValue("trigger_word")

	text = strings.Replace(text, trigger, "", -1)

	events := meetup.FilterEventsByKeyword(text)
	var response string

	if len(events) > 0 {
		e := events[0]
		t := time.Unix(0, e.Time*int64(time.Millisecond)).In(time.Local)
		layout := "Jan 2, 2006 at 3:04pm (MST)"
		response = fmt.Sprint(e.Name, " | ", t.Format(layout), " @ ", e.Venue.Name, " | ", e.EventUrl)
	} else {
		response = "No matching meetup found."
	}

	data := struct {
		Text string `json:"text"`
	}{response}

	output, _ := json.Marshal(data)
	log.Println(string(output))

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(output)
}

// Handle redirecting to the event page for a meetup
func meetupRedir(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	keyword := string(vars["meetup"])

	events := meetup.FilterEventsByKeyword(keyword)

	url := "http://www.meetup.com/memphis-technology-user-groups/"
	if len(events) > 0 {
		url = events[0].EventUrl
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	http.Redirect(w, r, url, 302)
}
