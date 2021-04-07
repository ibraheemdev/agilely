// Package test defines implemented interfaces for testing modules
package test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/ibraheemdev/agilely/pkg/mailer"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ engine.Database = &Database{}

// Database ...
type Database struct {
	Collections map[string]*Collection
}

// Aggregate ...
func (d *Database) Aggregate(ctx context.Context, pipeline map[string]interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	panic("not implemented")
}

// Collection ...
func (d *Database) Collection(name string) engine.Collection {
	if val, ok := d.Collections[name]; ok {
		return val
	}
	c := &Collection{
		Documents: make([]map[string]interface{}, 1),
		name:      name,
		db:        d,
	}

	d.Collections[name] = c
	return c
}

// Collection ...
type Collection struct {
	db        *Database
	name      string
	Documents []map[string]interface{}
}

// Name ...
func (c *Collection) Name() string {
	return c.name
}

// Database ...
func (c *Collection) Database() engine.Database {
	return c.db
}

// Drop ...
func (c *Collection) Drop(_ context.Context) error {
	c = nil
	return nil
}

// Aggregate ...
func (c *Collection) Aggregate(ctx context.Context, pipeline map[string]interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	panic("not implemented") // TODO: Implement
}

// Find ...
func (c *Collection) Find(ctx context.Context, filter map[string]interface{}, opts ...*options.FindOptions) (engine.Cursor, error) {
	var results []map[string]interface{}
	for _, d := range c.Documents {
		if contains(d, filter) {
			results = append(results, d)
			break
		}
	}
	return newCursor(results), nil
}

type cursor struct {
	pos  int
	docs []map[string]interface{}
}

func newCursor(docs []map[string]interface{}) *cursor {
	return &cursor{
		pos:  0,
		docs: docs,
	}
}

func (c *cursor) Next(_ context.Context) bool {
	if c.docs[c.pos] == nil {
		return false
	}
	c.pos++
	return true
}

func (c *cursor) Err() error {
	return nil
}

func (c *cursor) Decode(v interface{}) error {
	bytes, err := json.Marshal(c.docs[c.pos])
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		return err
	}

	return nil
}

func (c *cursor) Close(_ context.Context) error {
	return nil
}

func (c *cursor) All(_ context.Context, results interface{}) error {
	slice := reflect.ValueOf(results).Elem()

	slice.Set(reflect.MakeSlice(slice.Type(), len(c.docs), len(c.docs)))

	for i, doc := range c.docs {
		bytes, err := json.Marshal(doc)
		if err != nil {
			return err
		}

		err = json.Unmarshal(bytes, slice.Index(i).Addr().Interface())
		if err != nil {
			return err
		}
	}

	return nil
}

// FindOne ...
func (c *Collection) FindOne(ctx context.Context, filter map[string]interface{}, opts ...*options.FindOneOptions) engine.SingleResult {
	var result map[string]interface{}
	for _, d := range c.Documents {
		if contains(d, filter) {
			result = d
			break
		}
	}
	if result == nil {
		return &singleResult{result, mongo.ErrNoDocuments}
	}
	return &singleResult{result, nil}
}

type singleResult struct {
	doc map[string]interface{}
	err error
}

func (s *singleResult) Decode(v interface{}) error {
	bytes, err := json.Marshal(s.doc)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, v); err != nil {
		return err
	}
	return nil
}

func (s *singleResult) Err() error {
	return s.err
}

func contains(doc, filter map[string]interface{}) bool {
	for k, v := range filter {
		if doc[k] != v {
			return false
		}
	}
	return true
}

// FindOneAndDelete ...
func (c *Collection) FindOneAndDelete(ctx context.Context, filter map[string]interface{}, opts ...*options.FindOneAndDeleteOptions) engine.SingleResult {
	var result map[string]interface{}
	for _, d := range c.Documents {
		if contains(d, filter) {
			result = d
			d = nil
			break
		}
	}
	return &singleResult{result, nil}
}

// FindOneAndUpdate ...
func (c *Collection) FindOneAndUpdate(ctx context.Context, filter map[string]interface{}, update map[string]interface{}, opts ...*options.FindOneAndUpdateOptions) engine.SingleResult {
	panic("not implemented") // TODO: Implement
}

// InsertMany ...
func (c *Collection) InsertMany(ctx context.Context, documents []map[string]interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	panic("not implemented") // TODO: Implement
}

// InsertOne ...
func (c *Collection) InsertOne(ctx context.Context, document map[string]interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	panic("not implemented") // TODO: Implement
}

// UpdateMany ...
func (c *Collection) UpdateMany(ctx context.Context, filter map[string]interface{}, update map[string]interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	panic("not implemented") // TODO: Implement
}

// UpdateOne ...
func (c *Collection) UpdateOne(ctx context.Context, filter map[string]interface{}, update map[string]interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	panic("not implemented") // TODO: Implement
}

// DeleteMany ...
func (c *Collection) DeleteMany(ctx context.Context, filter map[string]interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	panic("not implemented") // TODO: Implement
}

// DeleteOne ...
func (c *Collection) DeleteOne(ctx context.Context, filter map[string]interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	panic("not implemented") // TODO: Implement
}

// NewDatabase ...
func NewDatabase() *Database {
	return &Database{}
}

// ClientState is used for testing the client stores on context
type ClientState struct {
	Values        map[string]string
	GetShouldFail bool
}

// NewClientState constructs a ClientStorer
func NewClientState(data ...string) *ClientState {
	if len(data) != 0 && len(data)%2 != 0 {
		panic("It should be a key value list of arguments.")
	}

	values := make(map[string]string)

	for i := 0; i < len(data)-1; i += 2 {
		values[data[i]] = data[i+1]
	}

	return &ClientState{Values: values}
}

// Get a key's value
func (m *ClientState) Get(key string) (string, bool) {
	if m.GetShouldFail {
		return "", false
	}

	v, ok := m.Values[key]
	return v, ok
}

// Put a value
func (m *ClientState) Put(key, val string) { m.Values[key] = val }

// Del a key/value pair
func (m *ClientState) Del(key string) { delete(m.Values, key) }

// ClientStateRW stores things that would originally
// go in a session, or a map, in memory!
type ClientStateRW struct {
	ClientValues map[string]string
}

// NewClientRW takes the data from a client state
// and returns.
func NewClientRW() *ClientStateRW {
	return &ClientStateRW{
		ClientValues: make(map[string]string),
	}
}

// ReadState from memory
func (c *ClientStateRW) ReadState(*http.Request) (engine.ClientState, error) {
	return &ClientState{Values: c.ClientValues}, nil
}

// WriteState to memory
func (c *ClientStateRW) WriteState(w http.ResponseWriter, cstate engine.ClientState, cse []engine.ClientStateEvent) error {
	for _, e := range cse {
		switch e.Kind {
		case engine.ClientStateEventPut:
			c.ClientValues[e.Key] = e.Value
		case engine.ClientStateEventDel:
			delete(c.ClientValues, e.Key)
		case engine.ClientStateEventDelAll:
			c.ClientValues = make(map[string]string)
		}
	}

	return nil
}

// Request returns a new request with optional key-value body (form-post)
func Request(method string, postKeyValues ...string) *http.Request {
	var body io.Reader
	location := "http://localhost"

	if len(postKeyValues) > 0 {
		urlValues := make(url.Values)
		for i := 0; i < len(postKeyValues); i += 2 {
			urlValues.Set(postKeyValues[i], postKeyValues[i+1])
		}

		if method == "POST" || method == "PUT" {
			body = strings.NewReader(urlValues.Encode())
		} else {
			location += "?" + urlValues.Encode()
		}
	}

	req, err := http.NewRequest(method, location, body)
	if err != nil {
		panic(err.Error())
	}

	if len(postKeyValues) > 0 {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return req
}

// Mailer helps simplify mailer testing by storing the last sent email
type Mailer struct {
	Last    mailer.Email
	SendErr string
}

// NewMailer constructs a  mailer
func NewMailer() *Mailer {
	return &Mailer{}
}

// Send an e-mail
func (m *Mailer) Send(ctx context.Context, email mailer.Email) error {
	if len(m.SendErr) > 0 {
		return errors.New(m.SendErr)
	}

	m.Last = email
	return nil
}

// AfterCallback is a callback that knows if it was called
type AfterCallback struct {
	HasBeenCalled bool
	Fn            engine.AuthEventHandler
}

// NewAfterCallback constructs a new aftercallback.
func NewAfterCallback() *AfterCallback {
	m := AfterCallback{}

	m.Fn = func(http.ResponseWriter, *http.Request, bool) (bool, error) {
		m.HasBeenCalled = true
		return false, nil
	}

	return &m
}

// Renderer mock
type Renderer struct {
	Pages []string

	// Render call variables
	Context context.Context
	Page    string
	Data    engine.HTMLData
}

// HasLoadedViews ensures the views were loaded
func (r *Renderer) HasLoadedViews(pages ...string) error {
	if len(r.Pages) != len(pages) {
		return fmt.Errorf("want: %d loaded views, got: %d", len(pages), len(r.Pages))
	}

	for i, want := range pages {
		got := r.Pages[i]
		if want != got {
			return fmt.Errorf("want: %s [%d], got: %s", want, i, got)
		}
	}

	return nil
}

// Load nothing but store the pages that were loaded
func (r *Renderer) Load(pages ...string) error {
	r.Pages = append(r.Pages, pages...)
	return nil
}

// Render nothing, but record the arguments into the renderer
func (r *Renderer) Render(ctx context.Context, page string, data engine.HTMLData) ([]byte, string, error) {
	r.Context = ctx
	r.Page = page
	r.Data = data
	return nil, "text/html", nil
}

// Responder records how a request was responded to
type Responder struct {
	Status int
	Page   string
	Data   engine.HTMLData
}

// Respond stores the arguments in the struct
func (r *Responder) Respond(w http.ResponseWriter, req *http.Request, code int, page string, data engine.HTMLData) error {
	r.Status = code
	r.Page = page
	r.Data = data

	return nil
}

// Redirector stores the redirect options passed to it and writes the Code
// to the ResponseWriter.
type Redirector struct {
	Options engine.RedirectOptions
}

// Redirect a request
func (r *Redirector) Redirect(w http.ResponseWriter, req *http.Request, ro engine.RedirectOptions) error {
	r.Options = ro
	if len(ro.RedirectPath) == 0 {
		panic("no redirect path on redirect call")
	}
	http.Redirect(w, req, ro.RedirectPath, ro.Code)
	return nil
}

// Emailer that holds the options it was given
type Emailer struct {
	Email mailer.Email
}

// Send an e-mail
func (e *Emailer) Send(ctx context.Context, email mailer.Email) error {
	e.Email = email
	return nil
}

// BodyReader reads the body of a request and returns some values
type BodyReader struct {
	Return engine.Validator
}

// Read the return values
func (b BodyReader) Read(page string, r *http.Request) (engine.Validator, error) {
	return b.Return, nil
}

// Values is returned from the BodyReader
type Values struct {
	PID         string
	Password    string
	Token       string
	Code        string
	Recovery    string
	PhoneNumber string
	Remember    bool

	Errors []error
}

// GetPID from values
func (v Values) GetPID() string {
	return v.PID
}

// GetPassword from values
func (v Values) GetPassword() string {
	return v.Password
}

// GetToken from values
func (v Values) GetToken() string {
	return v.Token
}

// GetCode from values
func (v Values) GetCode() string {
	return v.Code
}

// GetPhoneNumber from values
func (v Values) GetPhoneNumber() string {
	return v.PhoneNumber
}

// GetRecoveryCode from values
func (v Values) GetRecoveryCode() string {
	return v.Recovery
}

// GetShouldRemember gets the value that tells
// the remember module if it should remember the user
func (v Values) GetShouldRemember() bool {
	return v.Remember
}

// Validate the values
func (v Values) Validate() []error {
	return v.Errors
}

// ArbValues is arbitrary value storage
type ArbValues struct {
	Values map[string]string
	Errors []error
}

// GetPID gets the pid
func (a ArbValues) GetPID() string {
	return a.Values["email"]
}

// GetPassword gets the password
func (a ArbValues) GetPassword() string {
	return a.Values["password"]
}

// GetValues returns all values
func (a ArbValues) GetValues() map[string]string {
	return a.Values
}

// Validate nothing
func (a ArbValues) Validate() []error {
	return a.Errors
}

// Logger logs to the void
type Logger struct {
}

// Info logging
func (l Logger) Info(string) {}

// Error logging
func (l Logger) Error(string) {}

// Router records the routes that were registered
type Router struct {
	Gets    []string
	Posts   []string
	Deletes []string
	Puts    []string
}

// ServeHTTP does nothing
func (Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

// GET records the path in the router
func (r *Router) GET(path string, _ http.Handler) {
	r.Gets = append(r.Gets, path)
}

// POST records the path in the router
func (r *Router) POST(path string, _ http.Handler) {
	r.Posts = append(r.Posts, path)
}

// DELETE records the path in the router
func (r *Router) DELETE(path string, _ http.Handler) {
	r.Deletes = append(r.Deletes, path)
}

// PUT records the path in the router
func (r *Router) PUT(path string, _ http.Handler) {
	r.Puts = append(r.Puts, path)
}

// HasGets ensures all gets routes are present
func (r *Router) HasGets(gets ...string) error {
	return r.hasRoutes(gets, r.Gets)
}

// HasPosts ensures all gets routes are present
func (r *Router) HasPosts(posts ...string) error {
	return r.hasRoutes(posts, r.Posts)
}

// HasDeletes ensures all gets routes are present
func (r *Router) HasDeletes(deletes ...string) error {
	return r.hasRoutes(deletes, r.Deletes)
}

// HasPuts ensures all gets routes are present
func (r *Router) HasPuts(puts ...string) error {
	return r.hasRoutes(puts, r.Puts)
}

func (r *Router) hasRoutes(want []string, got []string) error {
	if len(got) != len(want) {
		return fmt.Errorf("want: %d get routes, got: %d", len(want), len(got))
	}

	for i, w := range want {
		g := got[i]
		if w != g {
			return fmt.Errorf("wanted route: %s [%d], but got: %s", w, i, g)
		}
	}

	return nil
}

// ErrorHandler just holds the last error
type ErrorHandler struct {
	Error error
}

// Wrap an http method
func (e *ErrorHandler) Wrap(handler func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			e.Error = err
		}
	})
}
