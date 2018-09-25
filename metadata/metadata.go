package metadata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/golang/glog"
)

type SponsorLevel string

// all different SponsorLevels
const (
	Kilo SponsorLevel = "kilo"
	Mega SponsorLevel = "mega"
	Giga SponsorLevel = "giga"
	Tera SponsorLevel = "tera"
	// Aux refers to an auxilary type of sponsor who doesn't
	// belong to the categories above
	Aux SponsorLevel = "aux"
)

type AuburnHacks struct {
	metaFileURL string
	// mu gaurds all the variables below.
	mu sync.RWMutex

	AboutUs     string            `json:"about_us"`
	InfoCards   []*InfoCard       `json:"info_cards"`
	Sponsors    []*Sponsor        `json:"sponsors"`
	SocialMedia map[string]string `json:"social_media"`
}

type InfoCard struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Sponsor struct {
	Name  string       `json:"name"`
	Logo  string       `json:"logo"`
	Level SponsorLevel `json:"level"`
}

func New(metaFileURL string) *AuburnHacks {
	auHack := &AuburnHacks{
		metaFileURL: metaFileURL,
	}

	// not updating metadata in this function as the watch
	// go routine must be called immediately after object
	// allocation

	return auHack
}

// Watch is a function that runs as a goroutine and polls a
// metadata file on google drive
func (c *AuburnHacks) Watch(d time.Duration) {
	for {
		glog.V(2).Infof("polling metadata...")

		metadata, err := c.getMetadata()
		if err != nil {
			glog.Errorf("error getting metadata: %v", err)
			continue
		}

		if err = c.updateMetadata(metadata); err != nil {
			glog.Errorf("error updating metadata: %v", err)
			continue
		}

		glog.V(2).Infof("raw metadata: %s", metadata)
		glog.V(2).Infof("struct after: %+v", c)

		time.Sleep(d)
	}
}

func (c *AuburnHacks) RLock() {
	c.mu.RLock()
}

func (c *AuburnHacks) RUnlock() {
	c.mu.RUnlock()
}

func (c *AuburnHacks) updateMetadata(newData []byte) error {
	// grab the lock
	c.mu.Lock()
	defer c.mu.Unlock()

	newMeta := New(c.metaFileURL)
	if err := json.Unmarshal(newData, newMeta); err != nil {
		return err
	}

	// update the stuff
	c.AboutUs = newMeta.AboutUs
	c.InfoCards = newMeta.InfoCards
	c.Sponsors = newMeta.Sponsors
	c.SocialMedia = newMeta.SocialMedia

	return nil
}

func (c *AuburnHacks) getMetadata() ([]byte, error) {
	hc := http.Client{}
	req, err := http.NewRequest("GET", c.metaFileURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bb, nil
}
