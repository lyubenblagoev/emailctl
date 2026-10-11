package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/lyubenblagoev/emailctl"
	"github.com/lyubenblagoev/goprsc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestPaginatedCommandsMapFlagsToAPIRequests(t *testing.T) {
	tests := []struct {
		name string
		path string
		args []string
		body string
		cmd  func() *cobra.Command
	}{
		{
			name: "domains",
			path: "/api/v1/domains",
			body: `{"items":[{"id":1,"name":"example.com","enabled":true}],"page":2,"pageSize":50,"total":73}`,
			cmd: func() *cobra.Command {
				return findListCommand(t, CreateDomainCommand())
			},
		},
		{
			name: "accounts",
			path: "/api/v1/domains/example.com/accounts",
			args: []string{"example.com"},
			body: `{"items":[{"id":1,"username":"user","domain":"example.com","enabled":true}],"page":2,"pageSize":50,"total":73}`,
			cmd: func() *cobra.Command {
				return findListCommand(t, CreateAccountCommand())
			},
		},
		{
			name: "aliases for domain",
			path: "/api/v1/domains/example.com/aliases",
			args: []string{"example.com"},
			body: `{"items":[{"id":1,"alias":"contact","email":"user@example.com","enabled":true}],"page":2,"pageSize":50,"total":73}`,
			cmd: func() *cobra.Command {
				return findListCommand(t, CreateAliasCommand())
			},
		},
		{
			name: "recipients for alias",
			path: "/api/v1/domains/example.com/aliases/contact",
			args: []string{"example.com", "contact"},
			body: `{"items":[{"id":1,"alias":"contact","email":"user@example.com","enabled":true}],"page":2,"pageSize":50,"total":73}`,
			cmd: func() *cobra.Command {
				return findListCommand(t, CreateAliasCommand())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != test.path {
					t.Errorf("expected request path %q, got %q", test.path, r.URL.Path)
				}
				if r.URL.Query().Get("page") != "2" {
					t.Errorf("expected page=2, got %q", r.URL.Query().Get("page"))
				}
				if r.URL.Query().Get("pageSize") != "50" {
					t.Errorf("expected pageSize=50, got %q", r.URL.Query().Get("pageSize"))
				}
				fmt.Fprint(w, test.body)
			}))
			defer server.Close()

			previousClient := client
			client = newCommandTestClient(t, server)
			t.Cleanup(func() { client = previousClient })

			cmd := test.cmd()
			if err := cmd.Flags().Set("page", "2"); err != nil {
				t.Fatal(err)
			}
			if err := cmd.Flags().Set("page-size", "50"); err != nil {
				t.Fatal(err)
			}
			if err := cmd.PreRunE(cmd, test.args); err != nil {
				t.Fatal(err)
			}
			output := captureStdout(t, func() { cmd.Run(cmd, test.args) })
			if !strings.Contains(output, "page 2, page size 50; 73 total") {
				t.Fatalf("expected returned page metadata in output, got:\n%s", output)
			}
		})
	}
}

func TestPaginationPageSizeCanBeSetWithoutPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("page"); got != "" {
			t.Errorf("expected page to be omitted, got %q", got)
		}
		if got := r.URL.Query().Get("pageSize"); got != "25" {
			t.Errorf("expected pageSize=25, got %q", got)
		}
		fmt.Fprint(w, `{"items":[],"page":1,"pageSize":25,"total":0}`)
	}))
	defer server.Close()

	previousClient := client
	client = newCommandTestClient(t, server)
	t.Cleanup(func() { client = previousClient })

	cmd := findListCommand(t, CreateDomainCommand())
	if err := cmd.Flags().Set("page-size", "25"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.PreRunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	captureStdout(t, func() { cmd.Run(cmd, nil) })
}

func TestPaginationPageCanBeSetWithoutPageSize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("expected page=2, got %q", got)
		}
		if got := r.URL.Query().Get("pageSize"); got != "" {
			t.Errorf("expected pageSize to be omitted, got %q", got)
		}
		fmt.Fprint(w, `{"items":[],"page":2,"pageSize":50,"total":0}`)
	}))
	defer server.Close()

	previousClient := client
	client = newCommandTestClient(t, server)
	t.Cleanup(func() { client = previousClient })

	cmd := findListCommand(t, CreateDomainCommand())
	if err := cmd.Flags().Set("page", "2"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.PreRunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	captureStdout(t, func() { cmd.Run(cmd, nil) })
}

func TestPaginationFlagsValidateValues(t *testing.T) {
	tests := []struct {
		name    string
		flag    string
		value   string
		wantErr bool
	}{
		{name: "zero page", flag: "page", value: "0", wantErr: true},
		{name: "negative page", flag: "page", value: "-1", wantErr: true},
		{name: "zero page size", flag: "page-size", value: "0", wantErr: true},
		{name: "negative page size", flag: "page-size", value: "-1", wantErr: true},
		{name: "page size above default maximum", flag: "page-size", value: "201"},
		{name: "page size above default maximum by a wider margin", flag: "page-size", value: "500"},
		{name: "minimum page size", flag: "page-size", value: "1"},
		{name: "default maximum page size", flag: "page-size", value: "200"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := &Command{Command: &cobra.Command{}}
			pagination := &paginationFlags{}
			pagination.addFlags(cmd)
			if err := cmd.Flags().Set(test.flag, test.value); err != nil {
				t.Fatal(err)
			}
			err := cmd.PreRunE(cmd.Command, nil)
			if (err != nil) != test.wantErr {
				t.Fatalf("expected error %t, got %v", test.wantErr, err)
			}
		})
	}
}

func TestPageWritersDisplayItemsAndMetadata(t *testing.T) {
	metadata := goprsc.PaginationMetadata{Page: 3, PageSize: 25, Total: 88}

	var domains bytes.Buffer
	writeDomainPage(&domains, &emailctl.DomainsPage{
		Items:              []emailctl.Domain{{Domain: &goprsc.Domain{ID: 17, Name: "example.test", Enabled: true}}},
		PaginationMetadata: metadata,
	})
	assertPageOutput(t, domains.String(), "Domains (page 3, page size 25; 88 total):", "17   example.test")

	var accounts bytes.Buffer
	writeAccountPage(&accounts, "example.test", &emailctl.AccountsPage{
		Items:              []emailctl.Account{{Account: &goprsc.Account{ID: 23, Username: "user", Enabled: true}}},
		PaginationMetadata: metadata,
	})
	assertPageOutput(t, accounts.String(), "Accounts for 'example.test' (page 3, page size 25; 88 total):", "23   user@example.test")

	var aliases bytes.Buffer
	writeAliasPage(&aliases, "Aliases for 'example.test'", &emailctl.AliasesPage{
		Items:              []emailctl.Alias{{Alias: &goprsc.Alias{ID: 31, Name: "support", Email: "help@example.test", Enabled: true}}},
		PaginationMetadata: metadata,
	})
	assertPageOutput(t, aliases.String(), "Aliases for 'example.test' (page 3, page size 25; 88 total):", "31   support")
	if !strings.Contains(aliases.String(), "help@example.test") {
		t.Fatalf("expected recipient from Items in output, got:\n%s", aliases.String())
	}

	var empty bytes.Buffer
	writeDomainPage(&empty, &emailctl.DomainsPage{Items: []emailctl.Domain{}, PaginationMetadata: goprsc.PaginationMetadata{Page: 1, PageSize: 50, Total: 0}})
	assertPageOutput(t, empty.String(), "Domains (page 1, page size 50; 0 total):", "ID   Name")
	if strings.Contains(empty.String(), "example.test") {
		t.Fatalf("empty page unexpectedly displayed a row:\n%s", empty.String())
	}
}

func TestListCommandReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message":"page size exceeds the configured maximum"}`)
	}))
	defer server.Close()

	testClient := newCommandTestClient(t, server)
	err := listDomains(testClient, nil, []goprsc.PaginationOptions{{PageSize: 500}})
	if err == nil {
		t.Fatal("expected API error")
	}
	var apiError *goprsc.ErrorResponse
	if !errors.As(err, &apiError) {
		t.Fatalf("expected goprsc.ErrorResponse, got %T: %v", err, err)
	}
	if apiError.Response.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected HTTP 400, got %d", apiError.Response.StatusCode)
	}
}

func TestAliasBulkMutationsProcessEveryRecipientPage(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		apply         func(*emailctl.Client) error
		wantAliasName string
	}{
		{
			name:   "delete all",
			method: http.MethodDelete,
			apply: func(client *emailctl.Client) error {
				return client.Aliases.DeleteAll("example.com", "contact")
			},
		},
		{
			name:          "rename all",
			method:        http.MethodPut,
			wantAliasName: "new-contact",
			apply: func(client *emailctl.Client) error {
				return client.Aliases.RenameAll("example.com", "contact", "new-contact")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var changedRecipients []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					switch r.URL.Query().Get("page") {
					case "":
						fmt.Fprint(w, `{"items":[{"id":1,"alias":"contact","email":"first@example.com","enabled":true}],"page":1,"pageSize":1,"total":2}`)
					case "2":
						if got := r.URL.Query().Get("pageSize"); got != "1" {
							t.Errorf("expected next page size 1, got %q", got)
						}
						fmt.Fprint(w, `{"items":[{"id":2,"alias":"contact","email":"second@example.com","enabled":true}],"page":2,"pageSize":1,"total":2}`)
					default:
						t.Errorf("unexpected page request: %q", r.URL.Query().Get("page"))
						fmt.Fprint(w, `{"items":[],"page":1,"pageSize":1,"total":0}`)
					}
					return
				}

				if r.Method != test.method {
					t.Errorf("expected %s, got %s", test.method, r.Method)
				}
				changedRecipients = append(changedRecipients, r.URL.Path)
				if test.wantAliasName != "" {
					var request goprsc.AliasUpdateRequest
					if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
						t.Errorf("decode rename request: %v", err)
					}
					if request.Name != test.wantAliasName {
						t.Errorf("expected alias name %q, got %q", test.wantAliasName, request.Name)
					}
				}
			}))
			defer server.Close()

			client := newCommandTestClient(t, server)
			if err := test.apply(client); err != nil {
				t.Fatal(err)
			}
			if len(changedRecipients) != 2 {
				t.Fatalf("expected both recipients to be processed, got %v", changedRecipients)
			}
			if !strings.HasSuffix(changedRecipients[0], "first@example.com") || !strings.HasSuffix(changedRecipients[1], "second@example.com") {
				t.Fatalf("unexpected recipients processed: %v", changedRecipients)
			}
		})
	}
}

func newCommandTestClient(t *testing.T, server *httptest.Server) *emailctl.Client {
	t.Helper()
	parsed, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	for key, value := range map[string]interface{}{
		"host":         parsed.Hostname(),
		"port":         parsed.Port(),
		"https":        false,
		"login":        "",
		"authToken":    "",
		"refreshToken": "",
	} {
		viper.Set(key, value)
	}

	client, err := emailctl.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func findListCommand(t *testing.T, parent *Command) *cobra.Command {
	t.Helper()
	cmd, _, err := parent.Find([]string{"list"})
	if err != nil {
		t.Fatal(err)
	}
	return cmd
}

func captureStdout(t *testing.T, run func()) string {
	t.Helper()
	previousStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writer
	defer func() { os.Stdout = previousStdout }()
	run()
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	output, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	return string(output)
}

func assertPageOutput(t *testing.T, output, header, row string) {
	t.Helper()
	if !strings.Contains(output, header) {
		t.Fatalf("expected %q in output, got:\n%s", header, output)
	}
	if !strings.Contains(output, row) {
		t.Fatalf("expected %q in output, got:\n%s", row, output)
	}
}
