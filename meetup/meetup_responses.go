package meetup

type Events struct {
  Results []Event
  Meta ResponseMetadata
}

type Event struct {
  UtcOffset int `json:utc_offset`
  Venue Venue
  Headcount int
  Visibility string
  WaitlistCount int `json:waitlist_count`
  Created int
  MaybeRsvpCount int `json:maybe_rsvp_count`
  Description string
  EventUrl string `json:event_url`
  YesRsvpCount int `json:yes_rsvp_count`
  Announced bool
  Name string
  Id string
  Time int
  Updated int
  Group Group
  Status string
}

type Venue struct {
  Country string
  City string
  Address1 string `json:address_1`
  Address2 string `json:address_2`
  Name string
  Lon float64
  Lat float64
  Id int
  Repinned bool
}

type Group struct {
  JoinMode string `json:join_mode`
  Created int
  Name string
  GroupLon float64
  GroupLat float64
  Id int
  Urlname string
  Who string
}

type ResponseMetadata struct {
  Next string
  Method string
  TotalCount int `json:total_count`
  Link string
  Count int
  Description string
  Lon float64
  Title string
  Url string
  Id string
  Updated int
  Lat float64
}
