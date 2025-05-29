package security

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// MTLSConfig holds the configuration for mutual TLS
type MTLSConfig struct {
	CertFile   string
	KeyFile    string
	CAFile     string
	ClientAuth tls.ClientAuthType
}

// SetupMTLS configures mutual TLS for the server
func SetupMTLS(config MTLSConfig) (*tls.Config, error) {
	// Load server certificate and key
	cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %w", err)
	}

	// Load CA certificate for client verification
	caCert, err := ioutil.ReadFile(config.CAFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA certificate")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   config.ClientAuth,
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		},
	}

	return tlsConfig, nil
}

// MTLSMiddleware validates client certificates and extracts identity
func MTLSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
			http.Error(w, "Client certificate required", http.StatusUnauthorized)
			return
		}

		clientCert := r.TLS.PeerCertificates[0]
		
		// Extract client identity from certificate
		clientID := extractClientIdentity(clientCert)
		if clientID == "" {
			http.Error(w, "Invalid client certificate", http.StatusUnauthorized)
			return
		}

		// Add client identity to request context
		ctx := r.Context()
		ctx = WithClientIdentity(ctx, clientID)
		r = r.WithContext(ctx)

		log.Printf("mTLS: Client authenticated: %s", clientID)
		next.ServeHTTP(w, r)
	})
}

// extractClientIdentity extracts the client identity from the certificate
func extractClientIdentity(cert *x509.Certificate) string {
	// Extract from Common Name or Subject Alternative Names
	if cert.Subject.CommonName != "" {
		return cert.Subject.CommonName
	}
	
	// Fallback to first DNS name in SAN
	if len(cert.DNSNames) > 0 {
		return cert.DNSNames[0]
	}
	
	return ""
}

// ValidateCertificateChain performs additional certificate validation
func ValidateCertificateChain(cert *x509.Certificate, intermediates *x509.CertPool, roots *x509.CertPool) error {
	opts := x509.VerifyOptions{
		Intermediates: intermediates,
		Roots:        roots,
		KeyUsages:    []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	_, err := cert.Verify(opts)
	return err
}
