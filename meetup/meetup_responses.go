package meetup

type Events struct {
  Results []Event
  Meta ResponseMetadata
}

type Event struct {
  UtcOffset int `json:"utc_offset"`
  Venue Venue `json:"venue"`
  Headcount int `json:"headcount"`
  Visibility string `json:"visibility"`
  WaitlistCount int `json:"waitlist_count"`
  Created int64 `json:"created"`
  MaybeRsvpCount int `json:"maybe_rsvp_count"`
  Description string `json:"description"`
  EventUrl string `json:"event_url"`
  YesRsvpCount int `json:"yes_rsvp_count"`
  Announced bool `json:"announced"`
  Name string `json:"name"`
  Id string `json:"id"`
  Time int64 `json:"time"`
  Updated int64 `json:"updated"`
  Group Group `json:"group"`
  Status string `json:"status"`
}

type Venue struct {
  Country string `json:"country"`
  City string `json:"city"`
  Address1 string `json:"address_1"`
  Address2 string `json:"address_2"`
  Name string `json:"name"`
  Lon float64 `json:"lon"`
  Lat float64 `json:"lat"`
  Id int `json:"id"`
  Repinned bool `json:"repinned"`
}

type Group struct {
  JoinMode string `json:"join_mode"`
  Created int `json:"created"`
  Name string `json:"name"`
  GroupLon float64 `json:"group_lon"`
  GroupLat float64 `json:"group_lat"`
  Id int `json:"id"`
  Urlname string `json:"url_name"`
  Who string `json:"who"`
}

type ResponseMetadata struct {
  Next string `json:"next"`
  Method string `json:"method"`
  TotalCount int `json:"total_count"`
  Link string `json:"link"`
  Count int `json:"count"`
  Description string `json:"description"`
  Lon float64 `json:"lon"`
  Title string `json:"title"`
  Url string `json:"url"`
  Id string `json:"id"`
  Updated int `json:"updated"`
  Lat float64 `json:"lat"`
}
