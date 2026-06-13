package lilianweng_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tamnd/lilianweng-cli/lilianweng"
)

func sampleFeed(baseURL string) string {
	return `<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0">
  <channel>
    <title>Lil'Log</title>
    <link>` + baseURL + `/</link>
    <item>
      <title>Why We Think</title>
      <link>` + baseURL + `/posts/2025-05-01-thinking/</link>
      <pubDate>Thu, 01 May 2025 00:00:00 +0000</pubDate>
      <description>&lt;p&gt;Test time compute and chain-of-thought reasoning.&lt;/p&gt;</description>
    </item>
    <item>
      <title>Attention Is All You Need</title>
      <link>` + baseURL + `/posts/2018-06-24-attention/</link>
      <pubDate>Sun, 24 Jun 2018 00:00:00 +0000</pubDate>
      <description>&lt;p&gt;Deep dive into the self-attention mechanism and transformers.&lt;/p&gt;</description>
    </item>
    <item>
      <title>Diffusion Models</title>
      <link>` + baseURL + `/posts/2021-07-11-diffusion-models/</link>
      <pubDate>Sun, 11 Jul 2021 00:00:00 +0000</pubDate>
      <description>&lt;p&gt;An introduction to diffusion probabilistic models for image generation.&lt;/p&gt;</description>
    </item>
  </channel>
</rss>`
}

func newTestServer(t *testing.T) (*httptest.Server, *lilianweng.Client) {
	t.Helper()
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") == "" {
			t.Error("request carried no User-Agent")
		}
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(sampleFeed(srv.URL)))
	}))
	cfg := lilianweng.DefaultConfig()
	cfg.BaseURL = srv.URL
	cfg.Rate = 0
	return srv, lilianweng.NewClient(cfg)
}

func TestLatest(t *testing.T) {
	srv, c := newTestServer(t)
	defer srv.Close()

	posts, err := c.Latest(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 3 {
		t.Fatalf("got %d posts, want 3", len(posts))
	}
	if posts[0].Title != "Why We Think" {
		t.Errorf("first title = %q, want %q", posts[0].Title, "Why We Think")
	}
	if posts[0].Rank != 1 {
		t.Errorf("rank = %d, want 1", posts[0].Rank)
	}
	if posts[0].Published != "2025-05-01" {
		t.Errorf("published = %q, want 2025-05-01", posts[0].Published)
	}
	if !strings.HasPrefix(posts[0].URL, srv.URL) {
		t.Errorf("unexpected URL %q", posts[0].URL)
	}
}

func TestLatestLimit(t *testing.T) {
	srv, c := newTestServer(t)
	defer srv.Close()

	posts, err := c.Latest(context.Background(), 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 2 {
		t.Fatalf("got %d posts, want 2", len(posts))
	}
}

func TestSearch(t *testing.T) {
	srv, c := newTestServer(t)
	defer srv.Close()

	posts, err := c.Search(context.Background(), "attention", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 1 {
		t.Fatalf("got %d results for 'attention', want 1", len(posts))
	}
	if posts[0].Title != "Attention Is All You Need" {
		t.Errorf("title = %q", posts[0].Title)
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	srv, c := newTestServer(t)
	defer srv.Close()

	posts, err := c.Search(context.Background(), "DIFFUSION", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 1 {
		t.Fatalf("got %d results for 'DIFFUSION', want 1", len(posts))
	}
}

func TestSearchNoMatch(t *testing.T) {
	srv, c := newTestServer(t)
	defer srv.Close()

	posts, err := c.Search(context.Background(), "kubernetes", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 0 {
		t.Errorf("got %d results, want 0", len(posts))
	}
}

func TestClientRetries(t *testing.T) {
	var hits int
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		if hits < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(sampleFeed(srv.URL)))
	}))
	defer srv.Close()

	cfg := lilianweng.DefaultConfig()
	cfg.BaseURL = srv.URL
	cfg.Rate = 0
	cfg.Retries = 5
	c := lilianweng.NewClient(cfg)

	posts, err := c.Latest(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) == 0 {
		t.Error("got no posts after retries")
	}
	if hits != 3 {
		t.Errorf("server saw %d hits, want 3", hits)
	}
}
