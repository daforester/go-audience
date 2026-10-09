package routers

import (
	"bitbucket.org/daforester/go-audience/apps/fantasy/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/fantasy/bindings"
	"bitbucket.org/daforester/go-audience/apps/fantasy/repositories"
	"bitbucket.org/daforester/go-audience/apps/fantasy/setup"
	"bytes"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/daforester/go-di-container/di"
	fws "github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
)

var testAuthKey = base32.StdEncoding.EncodeToString(bytes.Repeat([]byte{7}, 64))
var testEncKey = base32.StdEncoding.EncodeToString(bytes.Repeat([]byte{9}, 32))

type testServer struct {
	t    *testing.T
	app  *fiber.App
	repo *repositories.MemoryGameRepository
}

type joined struct {
	RoomKey     string
	UUID        string
	AccessToken string
}

func newTestServer(t *testing.T, settings map[string]interface{}) *testServer {
	t.Helper()
	viper.Reset()
	setup.Defaults()
	viper.Set("SESSION_AUTHKEY", testAuthKey)
	viper.Set("SESSION_ENCKEY", testEncKey)
	for k, v := range settings {
		viper.Set(k, v)
	}
	setup.Keys()

	// Config is global: before the next test resets it, wait for this test's sockets to close
	t.Cleanup(func() { waitForSocketsClosed(t) })

	container := di.New()
	bindings.Register(container, bindings.FLAG_NONE)
	repo := container.Make((*repositories.GameRepository)(nil)).(*repositories.MemoryGameRepository)

	return &testServer{t, NewServer(container), repo}
}

func waitForSocketsClosed(t *testing.T) {
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		sockets.LockCollections()
		open := len(sockets.WebSocketCollections)
		sockets.UnlockCollections()
		if open == 0 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Error("websockets still open after the test finished")
}

func (s *testServer) do(method string, path string, token string, body interface{}, headers map[string]string) (int, map[string]interface{}, http.Header) {
	s.t.Helper()
	var reader io.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		reader = bytes.NewReader(data)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if h, ok := headers["Host"]; ok {
		req.Host = h
	}

	resp, err := s.app.Test(req, fiber.TestConfig{Timeout: 5 * time.Second})
	if err != nil {
		s.t.Fatal(err)
	}
	defer resp.Body.Close()

	out := make(map[string]interface{})
	data, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(data, &out)
	return resp.StatusCode, out, resp.Header
}

func (s *testServer) createGame(name string) joined {
	s.t.Helper()
	code, out, _ := s.do("POST", "/api/game/create", "", map[string]string{"Name": name}, nil)
	if code != 200 {
		s.t.Fatalf("create game: %d %v", code, out)
	}
	return parseJoined(out)
}

func (s *testServer) joinGame(room string, name string) joined {
	s.t.Helper()
	code, out, _ := s.do("POST", "/api/game/join", "", map[string]string{"Name": name, "RoomCode": room}, nil)
	if code != 200 {
		s.t.Fatalf("join game: %d %v", code, out)
	}
	return parseJoined(out)
}

func parseJoined(out map[string]interface{}) joined {
	r := out["Response"].(map[string]interface{})
	p := r["Player"].(map[string]interface{})
	tk := r["Tokens"].(map[string]interface{})
	return joined{p["RoomKey"].(string), p["UUID"].(string), tk["AccessToken"].(string)}
}

func TestCreateJoinAndRoles(t *testing.T) {
	s := newTestServer(t, nil)
	admin := s.createGame("  The   Admin ")
	if len(admin.RoomKey) != 5 {
		t.Fatalf("unexpected room key %q", admin.RoomKey)
	}

	guest := s.joinGame(strings.ToLower(admin.RoomKey), "Guest")

	// Admin sees the players, audience does not
	_, out, _ := s.do("GET", "/api/game", admin.AccessToken, nil, nil)
	players := out["Response"].(map[string]interface{})["Game"].(map[string]interface{})["Players"].([]interface{})
	if len(players) != 2 || players[0].(map[string]interface{})["Name"] != "The Admin" {
		t.Fatalf("admin view wrong: %v", players)
	}
	_, out, _ = s.do("GET", "/api/game", guest.AccessToken, nil, nil)
	if _, ok := out["Response"].(map[string]interface{})["Game"].(map[string]interface{})["Players"]; ok {
		t.Fatal("audience must not see the player list")
	}

	// Audience cannot assign roles
	code, _, _ := s.do("POST", "/api/game/role/assign", guest.AccessToken, map[string]string{"UUID": guest.UUID, "PlayerClass": "Admin"}, nil)
	if code != fiber.StatusForbidden {
		t.Fatalf("audience assigning roles: expected 403, got %d", code)
	}

	// Admin can, but a director needs a team
	code, _, _ = s.do("POST", "/api/game/role/assign", admin.AccessToken, map[string]string{"UUID": guest.UUID, "PlayerClass": "Director"}, nil)
	if code != fiber.StatusBadRequest {
		t.Fatalf("director without team: expected 400, got %d", code)
	}
	code, _, _ = s.do("POST", "/api/game/role/assign", admin.AccessToken, map[string]string{"UUID": guest.UUID, "PlayerClass": "Staff"}, nil)
	if code != 200 {
		t.Fatalf("assign staff: %d", code)
	}

	// The role check uses the stored player, so the old token now carries Staff rights
	code, out, _ = s.do("GET", "/api/game/rejoin", guest.AccessToken, nil, nil)
	if code != 200 || out["Response"].(map[string]interface{})["Player"].(map[string]interface{})["Class"] != "Staff" {
		t.Fatalf("rejoin after role change: %d %v", code, out)
	}

	// The only admin cannot demote themselves
	code, _, _ = s.do("POST", "/api/game/role/assign", admin.AccessToken, map[string]string{"UUID": admin.UUID, "PlayerClass": "Audience"}, nil)
	if code != fiber.StatusBadRequest {
		t.Fatalf("last admin demotion: expected 400, got %d", code)
	}

	// Kick, then the kicked token no longer works
	code, _, _ = s.do("GET", "/api/game/kick/"+guest.UUID, admin.AccessToken, nil, nil)
	if code != 200 {
		t.Fatalf("kick: %d", code)
	}
	code, _, _ = s.do("GET", "/api/game/rejoin", guest.AccessToken, nil, nil)
	if code != fiber.StatusUnauthorized {
		t.Fatalf("kicked rejoin: expected 401, got %d", code)
	}

	// Last admin quitting ends the game
	code, _, _ = s.do("GET", "/api/game/quit", admin.AccessToken, nil, nil)
	if code != 200 || s.repo.GetByKey(admin.RoomKey) != nil {
		t.Fatalf("admin quit should delete the game: %d", code)
	}
}

func TestRefreshToken(t *testing.T) {
	s := newTestServer(t, nil)
	code, out, _ := s.do("POST", "/api/game/create", "", map[string]string{"Name": "Admin"}, nil)
	if code != 200 {
		t.Fatal(code)
	}
	refresh := out["Response"].(map[string]interface{})["Tokens"].(map[string]interface{})["RefreshToken"].(string)

	code, out, _ = s.do("GET", "/api/token/refresh", refresh, nil, nil)
	if code != 200 || out["Response"].(map[string]interface{})["AccessToken"] == "" {
		t.Fatalf("refresh: %d %v", code, out)
	}

	tampered := refresh[:len(refresh)-4] + "AAAA"
	code, _, _ = s.do("GET", "/api/token/refresh", tampered, nil, nil)
	if code != fiber.StatusUnauthorized {
		t.Fatalf("tampered refresh token: expected 401, got %d", code)
	}
}

// Request data must be copied, never pointed to: fasthttp reuses request buffers, so a
// stored reference would change under us as later requests arrive.
func TestConcurrentJoinsKeepTheirOwnData(t *testing.T) {
	s := newTestServer(t, nil)
	admin := s.createGame("Admin")

	const n = 50
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.joinGame(admin.RoomKey, fmt.Sprintf("Player-%03d-%s", i, strings.Repeat("x", i%15+1)))
		}(i)
	}
	wg.Wait()

	// Push more traffic through the same buffers after the joins
	for i := 0; i < n; i++ {
		s.do("POST", "/api/game/join", "", map[string]string{"Name": strings.Repeat("Z", 30), "RoomCode": "NOPE0"}, nil)
	}

	game := s.repo.GetByKey(admin.RoomKey)
	_ = s.repo.Lock(game)
	defer s.repo.Unlock(game)

	seen := make(map[string]bool)
	for _, p := range game.GetAudience() {
		var i int
		var rest string
		if _, err := fmt.Sscanf(p.Name, "Player-%03d-%s", &i, &rest); err != nil {
			t.Fatalf("stored name corrupted: %q", p.Name)
		}
		if p.Name != fmt.Sprintf("Player-%03d-%s", i, strings.Repeat("x", i%15+1)) {
			t.Fatalf("stored name corrupted: %q", p.Name)
		}
		if p.RoomKey != admin.RoomKey {
			t.Fatalf("stored room key corrupted: %q", p.RoomKey)
		}
		seen[p.Name] = true
	}
	if len(seen) != n {
		t.Fatalf("expected %d distinct players, got %d", n, len(seen))
	}
}

func TestCors(t *testing.T) {
	strict := newTestServer(t, map[string]interface{}{"CORS_ALLOW_ORIGINS": []string{"https://allowed.example"}})

	_, _, h := strict.do("OPTIONS", "/api/game/create", "", nil, map[string]string{
		"Origin": "https://allowed.example", "Access-Control-Request-Method": "POST",
	})
	if h.Get("Access-Control-Allow-Origin") != "https://allowed.example" {
		t.Fatalf("allow-listed origin not allowed: %v", h)
	}

	_, _, h = strict.do("OPTIONS", "/api/game/create", "", nil, map[string]string{
		"Origin": "https://evil.example", "Access-Control-Request-Method": "POST",
	})
	if h.Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("unlisted origin allowed: %v", h)
	}

	// Naughty mode is off: matching host is not enough
	_, _, h = strict.do("OPTIONS", "/api/game/create", "", nil, map[string]string{
		"Origin": "http://game.local:8080", "Host": "game.local:8080", "Access-Control-Request-Method": "POST",
	})
	if h.Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("host match allowed without naughty mode: %v", h)
	}

	naughty := newTestServer(t, map[string]interface{}{"CORS_ALLOW_ORIGINS": []string{}, "CORS_BE_NAUGHTY": true})

	code, _, h := naughty.do("OPTIONS", "/api/game/create", "", nil, map[string]string{
		"Origin": "http://game.local:8080", "Host": "game.local:8080", "Access-Control-Request-Method": "POST",
	})
	if code != fiber.StatusNoContent || h.Get("Access-Control-Allow-Origin") != "http://game.local:8080" || h.Get("Access-Control-Allow-Credentials") != "true" {
		t.Fatalf("naughty host match not allowed: %d %v", code, h)
	}

	_, _, h = naughty.do("OPTIONS", "/api/game/create", "", nil, map[string]string{
		"Origin": "http://evil.example", "Host": "game.local:8080", "Access-Control-Request-Method": "POST",
	})
	if h.Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("naughty mismatch allowed: %v", h)
	}

	// Without TRUST_PROXY a spoofed X-Forwarded-Host is ignored
	_, _, h = naughty.do("OPTIONS", "/api/game/create", "", nil, map[string]string{
		"Origin": "http://evil.example", "Host": "game.local:8080", "X-Forwarded-Host": "evil.example", "Access-Control-Request-Method": "POST",
	})
	if h.Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("spoofed X-Forwarded-Host trusted: %v", h)
	}

	_, _, h = naughty.do("POST", "/api/game/create", "", map[string]string{"Name": "A"}, map[string]string{
		"Origin": "http://game.local:8080", "Host": "game.local:8080",
	})
	if h.Get("Access-Control-Allow-Origin") != "http://game.local:8080" {
		t.Fatalf("naughty simple request missing CORS headers: %v", h)
	}
	if h.Get("Access-Control-Allow-Origin") == "*" {
		t.Fatal("never answer with a wildcard origin")
	}
}

func listen(t *testing.T, s *testServer) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go func() { _ = s.app.Listener(ln, fiber.ListenConfig{DisableStartupMessage: true}) }()
	t.Cleanup(func() { _ = s.app.Shutdown() })
	return ln.Addr().String()
}

func dial(addr string, path string, origin string) (*fws.Conn, int, error) {
	header := http.Header{}
	if origin != "" {
		header.Set("Origin", origin)
	}
	conn, resp, err := fws.DefaultDialer.Dial("ws://"+addr+path, header)
	status := 0
	if resp != nil {
		status = resp.StatusCode
	}
	return conn, status, err
}

func readText(t *testing.T, conn *fws.Conn) string {
	t.Helper()
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, data, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	return string(data)
}

func TestWebSocket(t *testing.T) {
	s := newTestServer(t, map[string]interface{}{"CORS_ALLOW_ORIGINS": []string{"http://allowed.example"}})
	addr := listen(t, s)

	admin := s.createGame("Admin")
	guest := s.joinGame(admin.RoomKey, "Guest")
	path := func(p joined, token string) string {
		return "/ws/" + p.RoomKey + "/" + p.UUID + "?token=" + token
	}

	if _, status, err := dial(addr, path(guest, guest.AccessToken), "http://evil.example"); err == nil || status != fiber.StatusForbidden {
		t.Fatalf("bad origin: expected 403, got %d %v", status, err)
	}
	if _, status, err := dial(addr, path(guest, guest.AccessToken), ""); err == nil || status != fiber.StatusForbidden {
		t.Fatalf("missing origin: expected 403, got %d %v", status, err)
	}
	if _, status, err := dial(addr, path(guest, ""), "http://allowed.example"); err == nil || status != fiber.StatusUnauthorized {
		t.Fatalf("missing token: expected 401, got %d %v", status, err)
	}
	if _, status, err := dial(addr, path(guest, admin.AccessToken), "http://allowed.example"); err == nil || status != fiber.StatusUnauthorized {
		t.Fatalf("someone else's token: expected 401, got %d %v", status, err)
	}

	conn, _, err := dial(addr, path(guest, guest.AccessToken), "http://allowed.example")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	if msg := readText(t, conn); msg != "CODE 20" {
		t.Fatalf("expected game code on connect, got %q", msg)
	}
	if msg := readText(t, conn); !strings.HasPrefix(msg, `{"Command":"GAME_UPDATE"`) {
		t.Fatalf("expected the player's game view on connect, got %q", msg)
	}

	_ = conn.WriteMessage(fws.TextMessage, []byte("PING 123"))
	if msg := readText(t, conn); msg != "PONG 123" {
		t.Fatalf("expected PONG, got %q", msg)
	}

	// Oversized client messages close the connection
	big, _, err := dial(addr, path(guest, guest.AccessToken), "http://allowed.example")
	if err != nil {
		t.Fatal(err)
	}
	readText(t, big)
	readText(t, big)
	_ = big.WriteMessage(fws.TextMessage, bytes.Repeat([]byte("x"), 10000))
	_ = big.SetReadDeadline(time.Now().Add(3 * time.Second))
	if _, _, err := big.ReadMessage(); err == nil {
		t.Fatal("oversized message should close the socket")
	}

	// Kicking closes the socket server side, after a QUIT
	code, _, _ := s.do("GET", "/api/game/kick/"+guest.UUID, admin.AccessToken, nil, nil)
	if code != 200 {
		t.Fatal(code)
	}
	if msg := readText(t, conn); msg != "QUIT" {
		t.Fatalf("expected QUIT, got %q", msg)
	}
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	if _, _, err := conn.ReadMessage(); err == nil {
		t.Fatal("kicked socket should be closed by the server")
	}
}

func TestNaughtyWebSocket(t *testing.T) {
	s := newTestServer(t, map[string]interface{}{"CORS_ALLOW_ORIGINS": []string{}, "CORS_BE_NAUGHTY": true})
	addr := listen(t, s)
	admin := s.createGame("Admin")
	path := "/ws/" + admin.RoomKey + "/" + admin.UUID + "?token=" + admin.AccessToken

	if _, status, err := dial(addr, path, "http://evil.example"); err == nil || status != fiber.StatusForbidden {
		t.Fatalf("naughty mismatch: expected 403, got %d %v", status, err)
	}

	conn, _, err := dial(addr, path, "http://"+addr)
	if err != nil {
		t.Fatalf("naughty host match should connect: %v", err)
	}
	defer conn.Close()
	if msg := readText(t, conn); msg != "CODE 20" {
		t.Fatalf("expected game code, got %q", msg)
	}
}

func TestVersion(t *testing.T) {
	s := newTestServer(t, nil)
	_, out, _ := s.do("GET", "/api/version", "", nil, nil)
	if out["Response"].(map[string]interface{})["Version"] != setup.Version {
		t.Fatalf("version endpoint: %v", out)
	}
}
