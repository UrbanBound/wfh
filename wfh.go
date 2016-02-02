package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "net/url"
  "os"
  "os/user"
  "path/filepath"
  "time"
  "strings"

  "golang.org/x/net/context"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
  "google.golang.org/api/calendar/v3"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
  cacheFile, err := tokenCacheFile()
  if err != nil {
    log.Fatalf("Unable to get path to cached credential file. %v", err)
  }
  tok, err := tokenFromFile(cacheFile)
  if err != nil {
    tok = getTokenFromWeb(config)
    saveToken(cacheFile, tok)
  }
  return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
  authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
  fmt.Printf("Go to the following link in your browser then type the "+
    "authorization code: \n%v\n", authURL)

  var code string
  if _, err := fmt.Scan(&code); err != nil {
    log.Fatalf("Unable to read authorization code %v", err)
  }

  tok, err := config.Exchange(oauth2.NoContext, code)
  if err != nil {
    log.Fatalf("Unable to retrieve token from web %v", err)
  }
  return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
  usr, err := user.Current()
  if err != nil {
    return "", err
  }
  tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
  os.MkdirAll(tokenCacheDir, 0700)
  return filepath.Join(tokenCacheDir,
    url.QueryEscape("calendar-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
  f, err := os.Open(file)
  if err != nil {
    return nil, err
  }
  t := &oauth2.Token{}
  err = json.NewDecoder(f).Decode(t)
  defer f.Close()
  return t, err
}

func makeDayMap() *map[string]time.Time {
  today := time.Now()
  // tomorrow = time.AddDate(0,0,1)
  // days := ["sun", "mon", "tue", "wed", "thu", "fri", "sat"]

  dayMap := make(map[string]time.Time)
  for i := 0; i < 7; i++ {
    mappedTime := today.AddDate(0,0,i)
    key := strings.ToLower(mappedTime.Weekday().String())[:3]
    dayMap[key] = mappedTime
  }
  return &dayMap
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
  fmt.Printf("Saving credential file to: %s\n", file)
  f, err := os.Create(file)
  if err != nil {
    log.Fatalf("Unable to cache oauth token: %v", err)
  }
  defer f.Close()
  json.NewEncoder(f).Encode(token)
}

func getDateString() *string {
  dayMap := *makeDayMap()
  yourTime := dayMap[os.Args[1]]
  if yourTime.IsZero() {
    log.Fatalf("Your day must be any of 'sun', 'mon', 'tue', 'wed', 'thu', 'fri', or 'sat'")
  }
  timeIdx := strings.Index(yourTime.Format(time.RFC3339), "T")
  dateString := yourTime.Format(time.RFC3339)[:timeIdx]
  return &dateString
}

func calendarIdFromFile() string {
  return os.Getenv("PRODUCT_OOO_CALENDAR_ID")
}

func main() {
  ctx := context.Background()

  b, err := ioutil.ReadFile("client_secret.json")
  if err != nil {
    log.Fatalf("Unable to read client secret file: %v", err)
  }

  config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
  if err != nil {
    log.Fatalf("Unable to parse client secret file to config: %v", err)
  }
  client := getClient(ctx, config)

  srv, err := calendar.New(client)
  if err != nil {
    log.Fatalf("Unable to retrieve calendar Client %v", err)
  }



  args := os.Args[1:]
  if len(args) != 1 {
    log.Fatalf("Must provide a day you are working from home")
  }

  dateString := getDateString()

  event := &calendar.Event{
    Summary: "Max K WFH",
    Location: "Home",
    Description: "Working from home",
    Start: &calendar.EventDateTime{
      Date: *dateString,
    },
    End: &calendar.EventDateTime{
      Date: *dateString,
    },
  }

  calendarId := calendarIdFromFile()
  event, err = srv.Events.Insert(calendarId, event).Do()
  if err != nil {
    log.Fatalf("Unable to create event. %v\n", err)
  }
  fmt.Printf("Event created: %s\n", event.HtmlLink)
}
