//go:build !darwin || cgo

package azure

import (
	"testing"

	"k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func Test_federatedTokenFile(t *testing.T) {
	const (
		fromCreds = "/creds/token"
		fromEnv   = "/env/token"
	)

	cases := map[string]struct {
		credentials map[string]string
		env         string
		want        string
	}{
		"CredentialsTakePrecedenceOverEnv": {
			credentials: map[string]string{CredentialsKeyFederatedTokenFile: fromCreds},
			env:         fromEnv,
			want:        fromCreds,
		},
		"CredentialsOnly": {
			credentials: map[string]string{CredentialsKeyFederatedTokenFile: fromCreds},
			want:        fromCreds,
		},
		"FallbackToEnvWhenKeyMissing": {
			credentials: map[string]string{},
			env:         fromEnv,
			want:        fromEnv,
		},
		"FallbackToEnvWhenKeyEmpty": {
			credentials: map[string]string{CredentialsKeyFederatedTokenFile: ""},
			env:         fromEnv,
			want:        fromEnv,
		},
		"NoCredentialsFallbackToEnv": {
			env:  fromEnv,
			want: fromEnv,
		},
		"NeitherSet": {
			credentials: map[string]string{},
			want:        "",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Setenv(envFederatedTokenFile, tc.env)
			if got := federatedTokenFile(tc.credentials); got != tc.want {
				t.Errorf("federatedTokenFile(...): want %q, got %q", tc.want, got)
			}
		})
	}
}

func Test_kubeloginTokenOptionsFromRESTConfig(t *testing.T) {
	t.Run("empty args", func(t *testing.T) {
		rc := &rest.Config{
			ExecProvider: &clientcmdapi.ExecConfig{
				Args: []string{},
			},
		}
		opts, err := kubeloginTokenOptionsFromRESTConfig(rc)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if opts == nil {
			t.Errorf("should not be nil")
		} else if opts.ServerID != "" {
			t.Errorf("unexpected server id: %q", opts.ServerID)
		}
	})

	t.Run("parsing server-id", func(t *testing.T) {
		serverID := "foo"
		rc := &rest.Config{
			ExecProvider: &clientcmdapi.ExecConfig{
				Args: []string{"--server-id", serverID},
			},
		}
		opts, err := kubeloginTokenOptionsFromRESTConfig(rc)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if opts == nil {
			t.Errorf("should not be nil")
		} else if opts.ServerID != serverID {
			t.Errorf("unexpected server id: %q, expected: %q", opts.ServerID, serverID)
		}
	})
}
