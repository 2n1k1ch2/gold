package Formatters

import (
	"bytes"
	"gold/internal/analyzer"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

var (
	alertTmplFull  *template.Template
	alertTmplShort *template.Template
	infoTmplShort  *template.Template
)

const (
	ShortAlertFormat = `----------------\n
	Alert in {{.ServiceName}}\n
		\t-Time: \t{{.Timestamp | formatTime}}\n
		\t-Type: \t{{.Type}}\n
		\t-Rate: \t{{.Rate | formatRate}}\n
	`
	ShortInfoFormat = `----------------\n
	ServiceName: {{.ServiceName}}\n
		\t-Time: \t{{.Timestamp | formatTime}}\n
		\t-Goroutine Count: \t{{.Goroutine}}\n
		\t-Alerts: \t{{.Alerts}}\n
	`

	//FullAlertFormat = `----------------\n
	//Alert in {{.ServiceName}}\n
	//	\t- Time:\t{{.Timestamp | formatTime}}\n
	//	\t- Hash:\t{{.Hash}}\n\n
	//	\t- Type:\t{{.Type}}\n
	//	\t- Description\t{{.Description}\n
	//	\t- Advice\t{{.Advice}}\n\n
	//Statistic:\n
	//	\t- Previous count: {{.PrevCount}}
	//	\t- New count: {{.NewCount}}
	//	\t- Rate: {{printf "%.2f" .Rate}}%
	//	\t- Cycles: {{.Cycles}}\n\n
	//Path:\n
	//	\t- Length:{{len .Frames}}
	//	{{Frames | formatFrames	}}
	//`
)

func init() {
	funcMap := template.FuncMap{
		"formatTime": func(t int64) string {
			return time.Unix(t, 0).Format("2006-01-02 15:04:05")
		},
		"formatRate": func(f float64) string {
			return strconv.FormatFloat(f, 'f', 2, 64)
		},
		"formatFrames": func(f []string) string {
			return strings.Join(f, "\n")
		},
	}
	alertTmplShort = template.Must(template.New("short").Funcs(funcMap).Parse(ShortAlertFormat))
	infoTmplShort = template.Must(template.New("info").Funcs(funcMap).Parse(ShortInfoFormat))
	//alertTmplFull = template.Must(template.New("full").Funcs(funcMap).Parse(FullAlertFormat))
}

var bufPool = sync.Pool{
	New: func() interface{} { return &bytes.Buffer{} },
}

type TextFormatter struct {
}

func (f *TextFormatter) FormatAlert(alert *analyzer.Alert) ([]byte, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)
	if err := alertTmplShort.Execute(buf, alert); err != nil {
		return nil, err
	}
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}
func (f *TextFormatter) FormatInfo(info *analyzer.Info) ([]byte, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)
	if err := infoTmplShort.Execute(buf, info); err != nil {
		return nil, err
	}
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}
