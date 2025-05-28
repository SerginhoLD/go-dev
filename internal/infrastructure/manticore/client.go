package manticore

import (
	"context"
	"exampleapp/internal/infrastructure/logger"
	"fmt"
	"log/slog"
	"os"
	"time"

	manticoresearch "github.com/manticoresoftware/manticoresearch-go"
)

type Client struct {
	client *manticoresearch.APIClient
}

func NewClient() *Client {
	conf := manticoresearch.NewConfiguration()
	conf.Servers[0].URL = os.Getenv("MANTICORE_HTTP")
	client := manticoresearch.NewAPIClient(conf)

	return &Client{client}
}

func (c *Client) Replace(ctx context.Context, table string, id uint64, doc map[string]interface{}) error {
	slog.DebugContext(ctx, fmt.Sprintf("manticore: replace into %s", table), "id", id, "doc", doc)

	request := *manticoresearch.NewInsertDocumentRequest(table, doc)
	request.SetId(int64(id))

	_, _, err := c.client.IndexAPI.Replace(ctx).InsertDocumentRequest(request).Execute()

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("manticore: %s", err))
	}

	return err
}

// Возвращает построитель запроса
func (c *Client) NewSearch(table string, page uint64, limit uint64) *search {
	return &search{
		client: c,
		table:  table,
		page:   page,
		limit:  limit,
		match:  make(map[string]string),
		equals: make(map[string]any),
		ranges: make(map[string]map[string]uint64),
		in:     make(map[string]any),
		notIn:  make(map[string]any),
	}
}

type search struct {
	client *Client
	table  string
	page   uint64
	limit  uint64
	match  map[string]string            // "_all": "комн"
	equals map[string]any               // "price": 500
	ranges map[string]map[string]uint64 // "price": {"gte": 500,"lte": 1000}
	in     map[string]any               // "metro": ["a", "b"]
	notIn  map[string]any               // "metro": ["a", "b"]
}

// Match("_all", "комн")
// https://manual.manticoresearch.com/Searching/Full_text_matching/Basic_usage#match
func (s *search) Match(f string, v string) *search {
	s.match[f] = v
	return s
}

// Equals("price", 500)
// https://manual.manticoresearch.com/Searching/Filters#Equality-filters
func (s *search) Equals(f string, v any) *search {
	s.equals[f] = v
	return s
}

// Range("price", "gte", 123)
// https://manual.manticoresearch.com/Searching/Filters#Range-filters
func (s *search) Range(f string, t string, v uint64) *search {
	if s.ranges[f] == nil {
		s.ranges[f] = make(map[string]uint64)
	}

	s.ranges[f][t] = v
	return s
}

// In("metro", "a", "b")
// https://manual.manticoresearch.com/Searching/Filters?#Set-filters
func (s *search) In(f string, v ...any) *search {
	s.in[f] = v
	return s
}

func (s *search) NotIn(f string, v ...any) *search {
	s.notIn[f] = v
	return s
}

// client.NewSearch(...).Execute(ctx)
// client.NewSearch(...).Equals(...).Execute(ctx)
func (s *search) Execute(ctx context.Context) (*manticoresearch.SearchResponse, error) {
	request := *manticoresearch.NewSearchRequest(s.table)

	from := int32((s.page - 1) * s.limit)
	request.Offset = &from
	size := int32(s.limit)
	request.Limit = &size
	maxMatches := int32(20000)
	request.MaxMatches = &maxMatches
	request.Sort = map[string]string{"id": "desc"}

	if len(s.match) > 0 {
		s.addBoolMust(&request, manticoresearch.QueryFilter{Match: s.match})
	}

	for f, v := range s.equals {
		s.addBoolMust(&request, manticoresearch.QueryFilter{Equals: map[string]any{f: v}})
	}

	for f, v := range s.in {
		s.addBoolMust(&request, manticoresearch.QueryFilter{In: map[string]any{f: v}})
	}

	for f, v := range s.ranges {
		s.addBoolMust(&request, manticoresearch.QueryFilter{Range: map[string]any{f: v}})
	}

	for f, v := range s.notIn {
		s.addBoolMustNot(&request, &manticoresearch.QueryFilter{In: map[string]any{f: v}})
	}

	json, _ := request.MarshalJSON()
	slog.DebugContext(ctx, fmt.Sprintf("manticore: search %s", s.table), "request", string(json))
	start := time.Now()

	resp, _, err := s.client.client.SearchAPI.Search(ctx).SearchRequest(request).Execute()

	logger.AddMetric("app_search_request_duration_ms", float64(time.Since(start).Milliseconds()))

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("manticore: %s", err))
	}

	return resp, err
}

func (s *search) addBoolMust(r *manticoresearch.SearchRequest, f manticoresearch.QueryFilter) {
	if !r.HasQuery() {
		r.SetQuery(r.GetQuery())
	}

	if !r.Query.HasBool() {
		r.Query.SetBool(r.Query.GetBool())
	}

	r.Query.Bool.Must = append(r.Query.Bool.Must, f)
}

func (s *search) addBoolMustNot(r *manticoresearch.SearchRequest, f *manticoresearch.QueryFilter) {
	if !r.HasQuery() {
		r.SetQuery(r.GetQuery())
	}

	if !r.Query.HasBool() {
		r.Query.SetBool(r.Query.GetBool())
	}

	r.Query.Bool.MustNot = append(r.Query.Bool.MustNot, f)
}

// https://manual.manticoresearch.com/Searching/Faceted_search#HTTP-JSON
func (c *Client) AggsTerm(ctx context.Context, table string, f string) ([]*aggsBucket, error) {
	request := *manticoresearch.NewSearchRequest(table)

	size := int32(0)
	request.Limit = &size
	termSize := int32(10000)

	request.Aggs = map[string]manticoresearch.Aggregation{
		f: {
			Terms: &manticoresearch.AggTerms{
				Field: f,
				Size:  &termSize,
			},
		},
	}

	json, _ := request.MarshalJSON()
	slog.DebugContext(ctx, fmt.Sprintf("manticore: aggs %s", table), "request", string(json))
	resp, _, err := c.client.SearchAPI.Search(ctx).SearchRequest(request).Execute()

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("manticore: %s", err))
		return []*aggsBucket{}, err
	}

	//slog.DebugContext(ctx, fmt.Sprintf("manticore: %#v", resp.Aggregations[f].(map[string]any)["buckets"].([]any)))
	var buckets []*aggsBucket

	for _, b := range resp.Aggregations[f].(map[string]any)["buckets"].([]any) {
		buckets = append(buckets, &aggsBucket{
			Key:   b.(map[string]any)["key"],
			Count: uint64(b.(map[string]any)["doc_count"].(float64)),
		})
	}
	return buckets, nil
}

type aggsBucket struct {
	Key   any // string | float64
	Count uint64
}
