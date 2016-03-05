package hosts

import (
	"testing"
)

const username = "__user__"

func expectHost(t *testing.T, desc, user, host string, port int) {
	h, err := toHost(desc, username)
	if err != nil {
		t.Fatal(err)
	}

	if h.Host != host {
		t.Fatalf("expected host of %s, got %s", host, h.Host)
	}
	if h.User != user {
		t.Fatalf("expected user of %s, got %s", user, h.User)
	}

	if h.Port != port {
		t.Fatalf("expected port of %d, got %d", port, h.Port)
	}
}

func TestToHost(t *testing.T) {
	expectHost(t, "kellegous.com", username, "kellegous.com", 22)
	expectHost(t, "kel@kellegous.com", "kel", "kellegous.com", 22)
	expectHost(t, "kel@kellegous.com:222", "kel", "kellegous.com", 222)
	expectHost(t, "kellegous.com:42", username, "kellegous.com", 42)
}
