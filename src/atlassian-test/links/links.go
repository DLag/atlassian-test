package links

import (
	fabric "atlassian-test/filter_fabric"
	log "github.com/Sirupsen/logrus"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
)

const FILTER_NAME = "links"
const TITLE_MAX_LEN = 50

// `(?m)` for multiline mode
var filter_re = regexp.MustCompile(`(?m)https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,4}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
var title_re = regexp.MustCompile(`(?m)<title>(.*)<\/title>`)

type linksResult struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

// Empty struct for filter object
type filter struct {
	get func(string) string
}

// Filter filters input message and returns founded objects
func (f *filter) Filter(input string) (output fabric.FilteredResult) {
	//Looking for matches regarding to filter expression
	m := filter_re.FindAllString(input, -1)
	log.WithField("filter", FILTER_NAME).Debugf("Found %d matches", len(m))
	for _, url := range m {
		//Adding founded objects
		var res linksResult
		res.Url = url
		// Trying to get web page and parse it to find <title>
		tm := title_re.FindAllStringSubmatch(f.get(url), 1)
		// If successful adding title to the resule
		if len(tm) > 0 {
			// Unescaping HTML symbols
			res.Title = html.UnescapeString(tm[0][1])
			// Cutting to max len
			if len(res.Title) > TITLE_MAX_LEN {
				res.Title = res.Title[:47] + "..."
			}
		}
		output = append(output, res)
	}
	return
}

func httpGet(url string) (result string) {
	r, err := http.Get(url)
	if err != nil {
		return
	}
	defer r.Body.Close()
	buf, _ := ioutil.ReadAll(r.Body)
	result = string(buf)
	return
}

func init() {
	//Registering in filters catalog
	fabric.Register(FILTER_NAME, &filter{get: httpGet})
}
