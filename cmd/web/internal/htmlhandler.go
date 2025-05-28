package internal

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	_ "time/tzdata"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var tplCache = make(map[string]*template.Template)

func HttpHtml(w http.ResponseWriter, r *http.Request, path string, v any, code int) {
	if _, hasCache := tplCache[path]; !hasCache {
		tpl, err := template.New(path).Funcs(template.FuncMap{
			"format":  func() *format { return &format{} },
			"convert": func() *convert { return &convert{} },
			"math":    func() *_math { return &_math{} },
			"route":   func() *route { return &route{req: r} },
		}).ParseFiles("web/templates/base.gohtml", path)

		if err != nil {
			HttpJsonError(w, err.Error(), http.StatusInternalServerError) // "Internal Server Error"
			return
		}

		tplCache[path] = tpl
	}

	var b bytes.Buffer
	err := tplCache[path].ExecuteTemplate(&b, "base", v)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // "Internal Server Error"
		return
	}

	w.Header().Del("Content-Length") // @see http.Error
	w.WriteHeader(code)
	b.WriteTo(w)
}

type format struct{}

func (f *format) Price(v uint64) string {
	return message.NewPrinter(language.Russian).Sprintf("%d", v)
}

func (f *format) P_1f(v float64) string {
	return fmt.Sprintf("%.1f", v)
}

func (f *format) Time(v time.Time) string {
	loc, _ := time.LoadLocation("Europe/Moscow")
	return v.In(loc).Format("2006-01-02 15:04:05 -07:00")
}

func (f *format) Bool(v *bool) string {
	if v == nil {
		return ""
	}

	if !*v {
		return "0"
	}

	return "1"
}

type convert struct{}

func (c *convert) Uint8PointerToValue(v *uint8) uint8 {
	return pointerToValue(v)
}

func pointerToValue[T comparable](v *T) T {
	if v == nil {
		return *new(T)
	}

	return *v
}

type _math struct{}

func (c *_math) SumUint64(i uint64, j uint64) uint64 {
	return i + j
}

// Возвращает номера страниц для отображения: 1 2 3 0 7 8 0 11 12, где 0 разделитель
func (c *_math) Pagination(page uint64, limit uint64, total uint64) []uint64 {
	totalPages := uint64(math.Ceil(float64(total) / float64(limit)))
	var pages []uint64

	for i := uint64(1); i <= 3; i++ {
		if totalPages >= i {
			pages = append(pages, i)
		}
	}

	if totalPages < 4 {
		return pages
	}

	middleStart := int64(page) - 2
	middleEnd := middleStart + 4

	if middleStart > 4 {
		pages = append(pages, 0)
	}

	for i := middleStart; i <= middleEnd; i++ {
		if int64(totalPages) >= i && i >= 4 {
			pages = append(pages, uint64(i))
		}
	}

	if totalPages < uint64(middleEnd) {
		return pages
	}

	lastStart := totalPages - 2
	lastNum := pages[len(pages)-1]

	if lastStart > (uint64(middleEnd) + 2) {
		pages = append(pages, 0)
	}

	for i := lastStart; i <= totalPages; i++ {
		if i > lastNum {
			pages = append(pages, i)
		}
	}

	return pages
}

type route struct {
	req *http.Request
}

// Возвращает ссылку на страницу относительно текущего url
func (r *route) Page(page uint64) string {
	cloneUrl := new(url.URL)
	*cloneUrl = *r.req.URL

	query := cloneUrl.Query()
	query.Set("page", strconv.FormatUint(page, 10))

	cloneUrl.RawQuery = query.Encode()
	return cloneUrl.String()
}

func (r *route) TargetObj(id uint64) string {
	return fmt.Sprintf("%s/objects/sale/flats/%d/", os.Getenv("OBJECT_DOMAIN"), id)
}
