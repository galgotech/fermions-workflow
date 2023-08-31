package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/galgotech/fermions-workflow/pkg/bus"
	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/galgotech/fermions-workflow/pkg/setting"
)

func New(s setting.Setting, busEvent bus.Bus) (*Server, error) {
	return &Server{
		log:      log.New("server"),
		setting:  s,
		busEvent: busEvent,
	}, nil
}

type Server struct {
	log      log.Logger
	setting  setting.Setting
	busEvent bus.Bus
}

func (s *Server) Execute() error {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := gin.New()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		s.log.Debug("endpoint", "httpMethod", httpMethod, "absolutePath", absolutePath, "handlerName", handlerName, "nuHandlers", nuHandlers)
	}

	// Middleware
	router.Use(Logger(s.log, []string{}))
	router.Use(gin.Recovery())

	websocketHandler, err := NewCentrifuge(s.busEvent)
	if err != nil {
		return nil
	}

	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})

	router.GET("/ws", gin.WrapH(authMiddleware(websocketHandler)))

	srv := &http.Server{
		Addr:           ":4443",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	serverCrt, serverKey, err := GenX509KeyPair()
	if err != nil {
		return err
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		err := srv.ListenAndServeTLS(serverCrt, serverKey)
		// err := srv.ListenAndServe();
		if err != nil && err != http.ErrServerClosed {
			s.log.Fatal("listen", "err", err.Error())
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	s.log.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return errors.New("Server forced to shutdown: " + err.Error())
	}

	s.log.Info("Server exiting")
	return nil
}

func GenX509KeyPair() (string, string, error) {
	crtPath := "wrserver-crt.pem"
	keyPath := "wrserver-key.pem"

	_, errCrt := os.Stat(crtPath)
	_, errKey := os.Stat(keyPath)
	if errCrt == nil && errKey == nil {
		return crtPath, keyPath, nil
	}

	now := time.Now()
	template := &x509.Certificate{
		SerialNumber: big.NewInt(now.Unix()),
		Subject: pkix.Name{
			CommonName:         "workflow-runtime.galgo.tech",
			Country:            []string{"BR"},
			Organization:       []string{"galgo.tech"},
			OrganizationalUnit: []string{"workflow-runtime"},
		},
		NotBefore:             now,
		NotAfter:              now.AddDate(1, 0, 0), // Valid for one year
		SubjectKeyId:          []byte("galgotech"),
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, template, priv.Public(), priv)
	if err != nil {
		return "", "", err
	}

	certPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert,
		},
	)
	err = ioutil.WriteFile(crtPath, certPem, 0640)
	if err != nil {
		return "", "", err
	}

	privPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)
	err = ioutil.WriteFile(keyPath, privPem, 0640)
	if err != nil {
		return "", "", err
	}

	return crtPath, keyPath, nil
}
