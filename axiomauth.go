// Package axiomauth provides a Go client for the AxiomAuth REST API.
// Documentation: https://axiomauth.com/docs
package axiomauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	DefaultBaseURL = "https://api.axiomauth.com/v1"
	sdkVersion     = "1.1.0"
)

// Client is the AxiomAuth API client.
type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client

	Users    *UsersService
	Sessions *SessionsService
	Config   *ConfigService
	Audit    *AuditService
}

// New creates a new AxiomAuth client.
func New(apiKey string) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
	c.Users = &UsersService{c}
	c.Sessions = &SessionsService{c}
	c.Config = &ConfigService{c}
	c.Audit = &AuditService{c}
	return c
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	var r io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, r)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "axiomauth-go/"+sdkVersion)

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return fmt.Errorf("axiomauth: API error %d", res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(out)
}

// ListParams are common pagination parameters.
type ListParams struct {
	Page    int
	PerPage int
}

// --- Users ---

type UsersService struct{ c *Client }
type User struct {
	ID         string    \`json:"id"\`
	Email      string    \`json:"email"\`
	Name       string    \`json:"name"\`
	Role       string    \`json:"role"\`
	LastLogin  time.Time \`json:"last_login"\`
	MFAEnabled bool      \`json:"mfa_enabled"\`
}
type UserList struct {
	Data    []User \`json:"data"\`
	Total   int    \`json:"total"\`
	Page    int    \`json:"page"\`
	PerPage int    \`json:"per_page"\`
}

func (s *UsersService) List(p *ListParams) (*UserList, error) {
	out := &UserList{}
	return out, s.c.do(context.Background(), "GET", "/users", nil, out)
}
func (s *UsersService) Get(id string) (*User, error) {
	out := &User{}
	return out, s.c.do(context.Background(), "GET", "/users/"+id, nil, out)
}
func (s *UsersService) Deprovision(id string) error {
	return s.c.do(context.Background(), "DELETE", "/users/"+id, nil, nil)
}

// --- Sessions ---

type SessionsService struct{ c *Client }
type Session struct {
	ID        string    \`json:"id"\`
	UserID    string    \`json:"user_id"\`
	CreatedAt time.Time \`json:"created_at"\`
	ExpiresAt time.Time \`json:"expires_at"\`
}
type SessionList struct {
	Data  []Session \`json:"data"\`
	Total int       \`json:"total"\`
}

func (s *SessionsService) List(p *ListParams) (*SessionList, error) {
	out := &SessionList{}
	return out, s.c.do(context.Background(), "GET", "/sessions", nil, out)
}
func (s *SessionsService) Revoke(id string) error {
	return s.c.do(context.Background(), "DELETE", "/sessions/"+id, nil, nil)
}

// --- Config ---

type ConfigService struct{ c *Client }
type OrgConfig struct {
	OrgID               string   \`json:"org_id"\`
	OrgName             string   \`json:"org_name"\`
	SSOEnabled          bool     \`json:"sso_enabled"\`
	SAMLMetadataURL     string   \`json:"saml_metadata_url"\`
	SCIMEndpoint        string   \`json:"scim_endpoint"\`
	AllowedDomains      []string \`json:"allowed_domains"\`
	MFARequired         bool     \`json:"mfa_required"\`
	SessionDurationHours int     \`json:"session_duration_hours"\`
}
type ConfigUpdate struct {
	SessionDurationHours int  \`json:"session_duration_hours,omitempty"\`
	MFARequired          bool \`json:"mfa_required,omitempty"\`
}

func (s *ConfigService) Get() (*OrgConfig, error) {
	out := &OrgConfig{}
	return out, s.c.do(context.Background(), "GET", "/config", nil, out)
}
func (s *ConfigService) Update(u *ConfigUpdate) (*OrgConfig, error) {
	out := &OrgConfig{}
	return out, s.c.do(context.Background(), "PATCH", "/config", u, out)
}

// --- Audit ---

type AuditService struct{ c *Client }
type AuditEvent struct {
	ID        string    \`json:"id"\`
	Action    string    \`json:"action"\`
	ActorID   string    \`json:"actor_id"\`
	ActorEmail string   \`json:"actor_email"\`
	CreatedAt time.Time \`json:"created_at"\`
	Metadata  any       \`json:"metadata"\`
}
type AuditList struct {
	Data  []AuditEvent \`json:"data"\`
	Total int          \`json:"total"\`
}
type AuditParams struct {
	Action string
	From   string
	To     string
	Limit  int
}

func (s *AuditService) List(p *AuditParams) (*AuditList, error) {
	out := &AuditList{}
	return out, s.c.do(context.Background(), "GET", "/audit", nil, out)
}
