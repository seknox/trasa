package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/seknox/trasa/server/initdb"
	"github.com/vulcand/oxy/forward"
	"golang.org/x/crypto/acme/autocert"

	"github.com/seknox/trasa/server/api/my"

	"github.com/seknox/trasa/server/api/auth/serviceauth"

	"github.com/go-chi/chi"
	"github.com/go-chi/hostrouter"
	"github.com/rs/cors"
	webproxy "github.com/seknox/trasa/server/accessproxy/http"
	"github.com/seknox/trasa/server/accessproxy/rdpproxy"
	"github.com/seknox/trasa/server/accessproxy/sshproxy"
	"github.com/seknox/trasa/server/api/accesscontrol"
	"github.com/seknox/trasa/server/api/accessmap"
	"github.com/seknox/trasa/server/api/auth"
	"github.com/seknox/trasa/server/api/devices"
	"github.com/seknox/trasa/server/api/groups"
	"github.com/seknox/trasa/server/api/logs"
	"github.com/seknox/trasa/server/api/misc"
	"github.com/seknox/trasa/server/api/notif"
	"github.com/seknox/trasa/server/api/orgs"
	"github.com/seknox/trasa/server/api/policies"
	"github.com/seknox/trasa/server/api/providers/ca"
	"github.com/seknox/trasa/server/api/providers/sidp"
	"github.com/seknox/trasa/server/api/providers/uidp"
	"github.com/seknox/trasa/server/api/providers/vault"
	"github.com/seknox/trasa/server/api/providers/vault/tsxvault"
	"github.com/seknox/trasa/server/api/redis"
	"github.com/seknox/trasa/server/api/services"
	"github.com/seknox/trasa/server/api/stats"
	"github.com/seknox/trasa/server/api/system"
	"github.com/seknox/trasa/server/api/users"
	"github.com/seknox/trasa/server/global"
	"github.com/seknox/trasa/server/utils"
	"github.com/sirupsen/logrus"
)

// StartServer starts trasa core server
func StartServer() {

	state := global.InitDBSTORE()

	rdpproxy.InitStore(state, accesscontrol.TrasaUAC)
	sshproxy.InitStore(state, accesscontrol.TrasaUAC)
	serviceauth.InitStore(state, accesscontrol.TrasaUAC)

	accesscontrol.InitStore(state, accesscontrol.TrasaUAC)

	accessmap.InitStore(state)

	auth.InitStore(state)
	vault.InitStore(state)
	tsxvault.InitStore(state)
	devices.InitStore(state)
	groups.InitStore(state)
	uidp.InitStore(state)
	logs.InitStore(state)
	misc.InitStore(state)
	my.InitStore(state)
	notif.InitStore(state)
	orgs.InitStore(state)
	policies.InitStore(state)
	redis.InitStore(state)
	services.InitStore(state)
	system.InitStore(state)
	stats.InitStore(state)
	users.InitStore(state)

	uidp.InitStore(state)

	sidp.InitStore(state)
	ca.InitStore(state)

	initdb.InitDB()

	closeChan := make(chan bool, 1)
	go func() {
		err := sshproxy.ListenSSH(closeChan)
		if err != nil {
			logrus.Error(err)
		}
		closeChan <- true
	}()

	webproxy.PrepareProxyConfig()

	// Init chi router
	r := chi.NewRouter()
	hr := hostrouter.New()
	trasaListenAddr := global.GetConfig().Trasa.ListenAddr

	// domain below should be read and passed from config file. If no value is supplied, should listen on localhost by default
	logrus.Trace("trasa listen addr: ", trasaListenAddr)

	coreRouter := chi.NewRouter()
	hr.Map(trasaListenAddr, CoreAPIRouter(coreRouter))

	hr.Map("*", ProxyRouter())
	r.Mount("/", hr)

	go http.ListenAndServe(":80", http.HandlerFunc(redirect))

	go StartRadiusServer(closeChan)

	s := http.Server{
		Addr:    ":443",
		Handler: r,
	}

	done := make(chan struct{})
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, os.Kill)
		sig := <-quit
		logrus.Infof("trasa-server: shutting down ... %v", sig)
		logs.Store.RemoveAllActiveSessions()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := s.Shutdown(ctx); err != nil {
			logrus.Infof("trasa-server: server shutdown with error: %v", err)
		}
		cancel()
		done <- struct{}{}
	}()

	var err error
	if global.GetConfig().Trasa.AutoCert {
		logrus.Infof("HTTPs server started.")
		fmt.Printf("Open TRASA dashboard at https://%s \n", trasaListenAddr)
		err = s.Serve(autocert.NewListener(trasaListenAddr))
	} else {
		certPath := filepath.Join(utils.GetETCDir(), "trasa", "certs", "trasa-server.crt")
		keyPath := filepath.Join(utils.GetETCDir(), "trasa", "certs", "trasa-server.key")

		err = checkIfCertExists(certPath, keyPath)
		// If they are not available, generate new ones.
		if err != nil {
			err = generateCerts(certPath, keyPath, trasaListenAddr)
			if err != nil {
				logrus.Fatal("Error: Couldn't create https certs.")
			}
		}

		logrus.Infof("HTTPs server started. ")
		fmt.Printf("Open TRASA dashboard at https://%s \n", trasaListenAddr)
		err = s.ListenAndServeTLS(certPath, keyPath)
	}

	if err != nil {
		fmt.Println(err)
		logrus.Error(err)

	}

}

func redirect(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req,
		"https://"+req.Host+req.URL.String(),
		http.StatusMovedPermanently)
}

func ProxyRouter() chi.Router {

	r := chi.NewRouter()

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		pxy := webproxy.Proxy()

		pxy.ServeHTTP(w, r)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		logrus.Debugf("Reached not found in auth: %s", r.URL)
	})

	return r
}

func CoreAPIRouter(r *chi.Mux) chi.Router {

	// Cors Handler
	cors := cors.New(cors.Options{

		AllowedOrigins: []string{"*"},

		AllowedMethods: []string{"GET", "POST", "DELETE", "Put", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-SESSION", "X-CSRF",
			"Sec-Websocket-Key", "Sec-Websocket-Extensions", "Sec-Websocket-Protocol",
			"Sec-Websocket-Version", "Trasa-Extoken", "Trasa-SessionID", "Upgrade"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Sskey"},
		MaxAge:           300,
	})
	r.Use(cors.Handler)

	r = CoreAPIRoutes(r)

	logrus.Trace("Proxying dashboard: ", global.GetConfig().Trasa.ProxyDashboard)
	if global.GetConfig().Trasa.ProxyDashboard == true {

		r.Get("/*", func(w http.ResponseWriter, req *http.Request) {
			logrus.Trace("Forwarding non api request to dashboard: ", global.GetConfig().Trasa.DashboardAddr)
			transport := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}

			requestURI, err := url.ParseRequestURI(global.GetConfig().Trasa.DashboardAddr)
			if err != nil {
				logrus.Error(err)
				// TODO respond with error notification?
				return
			}
			req.URL = requestURI
			fwd, err := forward.New(forward.RoundTripper(transport), forward.Logger(global.OxyLog))

			if err != nil {
				logrus.Error(err)
			}
			fwd.ServeHTTP(w, req)
		})

		return r

	}

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		logrus.Trace("Not Found ROOT URL: serving ROOT: ", req.URL.Path)
		w.Header().Set("Cache-Control", "public, max-age=8176000")
		http.FileServer(http.Dir(filepath.Join(utils.GetVarDir(), "trasa", "dashboard"))).ServeHTTP(w, req)
	})

	r.Get("/static*", func(w http.ResponseWriter, req *http.Request) {
		logrus.Trace("Found static URL: serving STATIC : ", req.URL.Path)
		w.Header().Set("Cache-Control", "public, max-age=8176000")
		http.FileServer(http.Dir(filepath.Join(utils.GetVarDir(), "trasa", "dashboard"))).ServeHTTP(w, req)
	})

	r.Get("/assets*", func(w http.ResponseWriter, req *http.Request) {
		logrus.Trace("Found static URL: serving ASSETS : ", req.URL.Path)
		w.Header().Set("Cache-Control", "public, max-age=8176000")
		http.FileServer(http.Dir(filepath.Join(utils.GetVarDir(), "trasa", "dashboard"))).ServeHTTP(w, req)
	})

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		logrus.Trace("Not Found URL: serving Index File : ", req.URL.Path)
		w.Header().Set("Cache-Control", "no-store")
		http.ServeFile(w, req, filepath.Join(utils.GetVarDir(), "trasa", "dashboard", "index.html"))

	})

	return r
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string) {

	// if strings.ContainsAny(path, "{}*") {
	// 	panic("FileServer does not permit any URL parameters.")
	// }

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"
	r.Get(path, serveFile)

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		logrus.Debugf("Reached not found in auth: %s", req.URL)
		http.ServeFile(w, req, filepath.Join(utils.GetVarDir(), "trasa", "dashboard", "index.html"))
	})
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	rctx := chi.RouteContext(r.Context())

	//workDir, _ := os.Getwd()

	// filesDir := http.Dir(filepath.Join(workDir, filepath.Join(utils.GetVarDir(),"trasa","dashboard")))

	//fmt.Println("context: ", rctx.RoutePattern())
	pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
	// fmt.Println("serving: ", pathPrefix)
	fs := http.StripPrefix(pathPrefix, http.FileServer(http.Dir(filepath.Join(utils.GetVarDir(), "trasa", "dashboard"))))

	fs.ServeHTTP(w, r)
}

func fileServeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next(w, r)
	})
}
