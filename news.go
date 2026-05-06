package news

import (
    "encoding/xml"
    "errors"
    "net/http"
    "time"
)

const feedURL = "https://www.cleveland.com/arc/outboundfeeds/rss/?outputType=xml"

type rss struct {
    Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
    Title string    `xml:"title"`
    Items []rssItem `xml:"item"`
}

type rssItem struct {
    Title       string `xml:"title"`
    Link        string `xml:"link"`
    Description string `xml:"description"`
    PubDate     string `xml:"pubDate"`
}

type Item struct {
    Title       string
    Link        string
    Description string
    PubDate     time.Time
}

func GetLatest() ([]Item, error) {
    resp, err := http.Get(feedURL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, errors.New("non-200 from cleveland.com RSS")
    }

    var r rss
    if err := xml.NewDecoder(resp.Body).Decode(&r); err != nil {
        return nil, err
    }

    items := make([]Item, 0, len(r.Channel.Items))
    for _, it := range r.Channel.Items {
        t, _ := parsePubDate(it.PubDate)
        items = append(items, Item{
            Title:       it.Title,
            Link:        it.Link,
            Description: it.Description,
            PubDate:     t,
        })
    }
    return items, nil
}

func parsePubDate(s string) (time.Time, error) {
    // Cleveland.com uses standard RSS-style dates; adjust if needed.
    layouts := []string{
        time.RFC1123Z,
        time.RFC1123,
        time.RFC822Z,
        time.RFC822,
    }
    var lastErr error
    for _, layout := range layouts {
        if t, err := time.Parse(layout, s); err == nil {
            return t, nil
        } else {
            lastErr = err
        }
    }
    return time.Time{}, lastErr
}
