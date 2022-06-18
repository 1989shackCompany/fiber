package fiber

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3/utils"
	"github.com/valyala/fasthttp/fasthttputil"
)

// - [ ] Test_Start_Graceful_Shutdown

// go test -run Test_Start
func Test_Start(t *testing.T) {
	app := New()

	utils.AssertEqual(t, false, app.Start(":99999") == nil)

	go func() {
		time.Sleep(1000 * time.Millisecond)
		utils.AssertEqual(t, nil, app.Shutdown())
	}()

	utils.AssertEqual(t, nil, app.Start(":4003", StartConfig{DisableStartupMessage: true}))
}

// go test -run Test_Start_Prefork
func Test_Start_Prefork(t *testing.T) {
	testPreforkMaster = true

	app := New()

	utils.AssertEqual(t, nil, app.Start(":99999", StartConfig{DisableStartupMessage: true, EnablePrefork: true}))
}

// go test -run Test_Start_TLS
func Test_Start_TLS(t *testing.T) {
	app := New()

	// invalid port
	utils.AssertEqual(t, false, app.Start(":99999", StartConfig{
		CertFile:    "./.github/testdata/ssl.pem",
		CertKeyFile: "./.github/testdata/ssl.key",
	}) == nil)

	go func() {
		time.Sleep(1000 * time.Millisecond)
		utils.AssertEqual(t, nil, app.Shutdown())
	}()

	utils.AssertEqual(t, nil, app.Start(":0", StartConfig{
		CertFile:    "./.github/testdata/ssl.pem",
		CertKeyFile: "./.github/testdata/ssl.key",
	}))
}

// go test -run Test_Start_TLS_Prefork
func Test_Start_TLS_Prefork(t *testing.T) {
	testPreforkMaster = true

	app := New()

	// invalid key file content
	utils.AssertEqual(t, false, app.Start(":0", StartConfig{
		DisableStartupMessage: true,
		EnablePrefork:         true,
		CertFile:              "./.github/testdata/ssl.pem",
		CertKeyFile:           "./.github/testdata/template.tmpl",
	}) == nil)

	go func() {
		time.Sleep(1000 * time.Millisecond)
		utils.AssertEqual(t, nil, app.Shutdown())
	}()

	utils.AssertEqual(t, nil, app.Start(":99999", StartConfig{
		DisableStartupMessage: true,
		EnablePrefork:         true,
		CertFile:              "./.github/testdata/ssl.pem",
		CertKeyFile:           "./.github/testdata/ssl.key",
	}))
}

// go test -run Test_Start_MutualTLS
func Test_Start_MutualTLS(t *testing.T) {
	app := New()

	// invalid port
	utils.AssertEqual(t, false, app.Start(":99999", StartConfig{
		CertFile:       "./.github/testdata/ssl.pem",
		CertKeyFile:    "./.github/testdata/ssl.key",
		CertClientFile: "./.github/testdata/ca-chain.cert.pem",
	}) == nil)

	go func() {
		time.Sleep(1000 * time.Millisecond)
		utils.AssertEqual(t, nil, app.Shutdown())
	}()

	utils.AssertEqual(t, nil, app.Start(":0", StartConfig{
		CertFile:       "./.github/testdata/ssl.pem",
		CertKeyFile:    "./.github/testdata/ssl.key",
		CertClientFile: "./.github/testdata/ca-chain.cert.pem",
	}))
}

// go test -run Test_Start_MutualTLS_Prefork
func Test_Start_MutualTLS_Prefork(t *testing.T) {
	testPreforkMaster = true

	app := New()

	// invalid key file content
	utils.AssertEqual(t, false, app.Start(":0", StartConfig{
		DisableStartupMessage: true,
		EnablePrefork:         true,
		CertFile:              "./.github/testdata/ssl.pem",
		CertKeyFile:           "./.github/testdata/template.html",
		CertClientFile:        "./.github/testdata/ca-chain.cert.pem",
	}) == nil)

	go func() {
		time.Sleep(1000 * time.Millisecond)
		utils.AssertEqual(t, nil, app.Shutdown())
	}()

	utils.AssertEqual(t, nil, app.Start(":99999", StartConfig{
		DisableStartupMessage: true,
		EnablePrefork:         true,
		CertFile:              "./.github/testdata/ssl.pem",
		CertKeyFile:           "./.github/testdata/ssl.key",
		CertClientFile:        "./.github/testdata/ca-chain.cert.pem",
	}))
}

// go test -run Test_Start_CustomListener
func Test_Start_CustomListener(t *testing.T) {
	app := New()

	go func() {
		time.Sleep(500 * time.Millisecond)
		utils.AssertEqual(t, nil, app.Shutdown())
	}()

	ln := fasthttputil.NewInmemoryListener()
	utils.AssertEqual(t, nil, app.Start(ln))
}

// go test -run Test_Start_CustomListener_Prefork
func Test_Start_CustomListener_Prefork(t *testing.T) {
	testPreforkMaster = true

	app := New()

	ln := fasthttputil.NewInmemoryListener()
	utils.AssertEqual(t, nil, app.Start(ln, StartConfig{DisableStartupMessage: true, EnablePrefork: true}))
}

// go test -run Test_Start_CustomTLSListener
func Test_Start_CustomTLSListener(t *testing.T) {
	// Create tls certificate
	cer, err := tls.LoadX509KeyPair("./.github/testdata/ssl.pem", "./.github/testdata/ssl.key")
	if err != nil {
		utils.AssertEqual(t, nil, err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}

	ln, err := tls.Listen(NetworkTCP4, ":0", config)
	utils.AssertEqual(t, nil, err)

	app := New()

	go func() {
		time.Sleep(time.Millisecond * 500)
		utils.AssertEqual(t, nil, app.Shutdown())
	}()

	utils.AssertEqual(t, nil, app.Start(ln))
}

// go test -run Test_Start_TLSConfigFunc
func Test_Start_TLSConfigFunc(t *testing.T) {
	var callTLSConfig bool
	app := New()

	go func() {
		time.Sleep(1000 * time.Millisecond)
		utils.AssertEqual(t, nil, app.Shutdown())
	}()

	utils.AssertEqual(t, nil, app.Start(":0", StartConfig{
		DisableStartupMessage: true,
		TLSConfigFunc: func(tlsConfig *tls.Config) {
			callTLSConfig = true
		},
		CertFile:    "./.github/testdata/ssl.pem",
		CertKeyFile: "./.github/testdata/ssl.key",
	}))

	utils.AssertEqual(t, true, callTLSConfig)
}

// go test -run Test_Start_ListenerAddrFunc
func Test_Start_ListenerAddrFunc(t *testing.T) {
	var network string
	app := New()

	go func() {
		time.Sleep(1000 * time.Millisecond)
		utils.AssertEqual(t, nil, app.Shutdown())
	}()

	utils.AssertEqual(t, nil, app.Start(":0", StartConfig{
		DisableStartupMessage: true,
		ListenerAddrFunc: func(addr net.Addr) {
			network = addr.Network()
		},
		CertFile:    "./.github/testdata/ssl.pem",
		CertKeyFile: "./.github/testdata/ssl.key",
	}))

	utils.AssertEqual(t, "tcp", network)
}

// go test -run Test_Start_BeforeServeFunc
func Test_Start_BeforeServeFunc(t *testing.T) {
	var handlers uint32
	app := New()

	go func() {
		time.Sleep(1000 * time.Millisecond)
		utils.AssertEqual(t, nil, app.Shutdown())
	}()

	utils.AssertEqual(t, errors.New("test"), app.Start(":0", StartConfig{
		DisableStartupMessage: true,
		BeforeServeFunc: func(fiber *App) error {
			handlers = fiber.HandlersCount()

			return errors.New("test")
		},
	}))

	utils.AssertEqual(t, uint32(0), handlers)
}

// go test -run Test_Start_ListenerNetwork
func Test_Start_ListenerNetwork(t *testing.T) {
	var network string
	app := New()

	go func() {
		time.Sleep(1000 * time.Millisecond)
		utils.AssertEqual(t, nil, app.Shutdown())
	}()

	utils.AssertEqual(t, nil, app.Start(":0", StartConfig{
		DisableStartupMessage: true,
		ListenerNetwork:       NetworkTCP6,
		ListenerAddrFunc: func(addr net.Addr) {
			network = addr.String()
		},
	}))

	utils.AssertEqual(t, true, strings.Contains(network, "[::]:"))

	go func() {
		time.Sleep(1000 * time.Millisecond)
		utils.AssertEqual(t, nil, app.Shutdown())
	}()

	utils.AssertEqual(t, nil, app.Start(":0", StartConfig{
		DisableStartupMessage: true,
		ListenerNetwork:       NetworkTCP4,
		ListenerAddrFunc: func(addr net.Addr) {
			network = addr.String()
		},
	}))

	utils.AssertEqual(t, true, strings.Contains(network, "0.0.0.0:"))
}

// go test -run Test_Start_Master_Process_Show_Startup_Message
func Test_Start_Master_Process_Show_Startup_Message(t *testing.T) {
	cfg := StartConfig{
		EnablePrefork: true,
	}

	startupMessage := captureOutput(func() {
		New().
			startupMessage(":3000", true, strings.Repeat(",11111,22222,33333,44444,55555,60000", 10), cfg)
	})
	fmt.Println(startupMessage)
	utils.AssertEqual(t, true, strings.Contains(startupMessage, "https://127.0.0.1:3000"))
	utils.AssertEqual(t, true, strings.Contains(startupMessage, "(bound on host 0.0.0.0 and port 3000)"))
	utils.AssertEqual(t, true, strings.Contains(startupMessage, "Child PIDs"))
	utils.AssertEqual(t, true, strings.Contains(startupMessage, "11111, 22222, 33333, 44444, 55555, 60000"))
	utils.AssertEqual(t, true, strings.Contains(startupMessage, "Prefork ........ Enabled"))
}

// go test -run Test_Start_Master_Process_Show_Startup_MessageWithAppName
func Test_Start_Master_Process_Show_Startup_MessageWithAppName(t *testing.T) {
	cfg := StartConfig{
		EnablePrefork: true,
	}

	app := New(Config{AppName: "Test App v3.0.0"})
	startupMessage := captureOutput(func() {
		app.startupMessage(":3000", true, strings.Repeat(",11111,22222,33333,44444,55555,60000", 10), cfg)
	})
	fmt.Println(startupMessage)
	utils.AssertEqual(t, "Test App v3.0.0", app.Config().AppName)
	utils.AssertEqual(t, true, strings.Contains(startupMessage, app.Config().AppName))
}

// go test -run Test_Start_Print_Route
func Test_Start_Print_Route(t *testing.T) {
	app := New()
	app.Get("/", emptyHandler).Name("routeName")

	printRoutesMessage := captureOutput(func() {
		app.printRoutesMessage()
	})

	fmt.Println(printRoutesMessage)

	utils.AssertEqual(t, true, strings.Contains(printRoutesMessage, "GET"))
	utils.AssertEqual(t, true, strings.Contains(printRoutesMessage, "/"))
	utils.AssertEqual(t, true, strings.Contains(printRoutesMessage, "emptyHandler"))
	utils.AssertEqual(t, true, strings.Contains(printRoutesMessage, "routeName"))
}

// go test -run Test_Start_Print_Route_With_Group
func Test_Start_Print_Route_With_Group(t *testing.T) {
	app := New()
	app.Get("/", emptyHandler)

	v1 := app.Group("v1")
	v1.Get("/test", emptyHandler).Name("v1")
	v1.Post("/test/fiber", emptyHandler)
	v1.Put("/test/fiber/*", emptyHandler)

	printRoutesMessage := captureOutput(func() {
		app.printRoutesMessage()
	})

	fmt.Println(printRoutesMessage)

	utils.AssertEqual(t, true, strings.Contains(printRoutesMessage, "GET"))
	utils.AssertEqual(t, true, strings.Contains(printRoutesMessage, "/"))
	utils.AssertEqual(t, true, strings.Contains(printRoutesMessage, "emptyHandler"))
	utils.AssertEqual(t, true, strings.Contains(printRoutesMessage, "/v1/test"))
	utils.AssertEqual(t, true, strings.Contains(printRoutesMessage, "POST"))
	utils.AssertEqual(t, true, strings.Contains(printRoutesMessage, "/v1/test/fiber"))
	utils.AssertEqual(t, true, strings.Contains(printRoutesMessage, "PUT"))
	utils.AssertEqual(t, true, strings.Contains(printRoutesMessage, "/v1/test/fiber/*"))
}

func emptyHandler(c *Ctx) error {
	return nil
}