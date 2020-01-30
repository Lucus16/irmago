// Package test contains functionality that should be available to
// all unit tests (which live in separate packages).
package test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/privacybydesign/irmago/internal/fs"
	"github.com/privacybydesign/irmago/internal/keysharecore"
	"github.com/privacybydesign/irmago/server/keyshareserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkError(t *testing.T, err error) {
	if err == nil {
		return
	}
	if t != nil {
		require.NoError(t, err)
	} else {
		panic(err)
	}
}

var keyshareServ *http.Server

func StartKeyshareServer(t *testing.T) {
	db := keyshareserver.NewMemoryDatabase()
	db.NewUser(keyshareserver.KeyshareUserData{
		Username: "",
		Coredata: keysharecore.EncryptedKeysharePacket{},
	})
	var ep keysharecore.EncryptedKeysharePacket
	p, err := base64.StdEncoding.DecodeString("YWJjZB7irkDzwMWtBC6PTItWmO2AgAGm1/gFOyrd+nyt3/0GaHLY5Z1S1TM6N5nzb1Jh+Nqx0z0c3f9R2UyoYuy+pnrerTpYL1mpoZZfz8MPqcrAMsmVdb2kHH0BuAGSC0V28tp1BCVzhYnfMJyrUlNWonsTWSn68Av1BwpIBOGxqBXYfW0JzaffuSmZIyubImmTN7p32ASbseJSNwu0Rg==")
	require.NoError(t, err)
	copy(ep[:], p)
	db.NewUser(keyshareserver.KeyshareUserData{
		Username: "testusername",
		Coredata: ep,
	})

	testdataPath := FindTestdataFolder(t)
	s, err := keyshareserver.New(&keyshareserver.Configuration{
		SchemesPath:           filepath.Join(testdataPath, "irma_configuration"),
		URL:                   "http://localhost:8080/irma_keyshare_server/",
		DB:                    db,
		JwtKeyId:              0,
		JwtPrivateKeyFile:     filepath.Join(testdataPath, "jwtkeys", "kss-sk.pem"),
		StoragePrimaryKeyFile: filepath.Join(testdataPath, "keyshareStorageTestkey"),
		KeyshareCredential:    "test.test.mijnirma",
		KeyshareAttribute:     "email",
	})
	require.NoError(t, err)

	r := chi.NewRouter()
	r.Mount("/irma_keyshare_server/", s.Handler())

	keyshareServ = &http.Server{
		Addr:    "localhost:8080",
		Handler: r,
	}

	go func() {
		err := keyshareServ.ListenAndServe()
		if err == http.ErrServerClosed {
			err = nil
		}
		assert.NoError(t, err)
	}()
}

func StopKeyshareServer(t *testing.T) {
	err := keyshareServ.Shutdown(context.Background())
	assert.NoError(t, err)
}

var schemeServer *http.Server
var badServer *http.Server
var badServerCount int

func StartSchemeManagerHttpServer() {
	path := FindTestdataFolder(nil)
	schemeServer = &http.Server{Addr: ":48681", Handler: http.FileServer(http.Dir(path))}
	go func() {
		schemeServer.ListenAndServe()
	}()
	time.Sleep(100 * time.Millisecond) // Give server time to start
}

func StopSchemeManagerHttpServer() {
	schemeServer.Close()
}

// StartBadHttpServer starts an HTTP server that times out and returns 500 on the first few times.
func StartBadHttpServer(count int, timeout time.Duration, success string) {
	badServer = &http.Server{Addr: ":48682", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if badServerCount >= count {
			_, _ = fmt.Fprintln(w, success)
			return
		} else {
			badServerCount++
			time.Sleep(timeout)
		}
	})}

	go func() {
		_ = badServer.ListenAndServe()
	}()
	time.Sleep(100 * time.Millisecond) // Give server time to start
}

func StopBadHttpServer() {
	_ = badServer.Close()
}

// FindTestdataFolder finds the "testdata" folder which is in . or ..
// depending on which package is calling us.
func FindTestdataFolder(t *testing.T) string {
	path := "testdata"

	for i := 0; i < 3; i++ {
		exists, err := fs.PathExists(path)
		checkError(t, err)
		if exists {
			return path
		}
		path = filepath.Join("..", path)
	}

	checkError(t, errors.New("testdata folder not found"))
	return ""
}

// ClearTestStorage removes any output from previously run tests to ensure a clean state;
// some of the tests don't like it when there is existing state in storage.
func ClearTestStorage(t *testing.T) {
	path := filepath.Join(FindTestdataFolder(t), "storage", "test")
	err := os.RemoveAll(path)
	checkError(t, err)
}

func CreateTestStorage(t *testing.T) {
	ClearTestStorage(t)
	path := filepath.Join(FindTestdataFolder(t), "storage")

	// EnsureDirectoryExists eventually uses mkdir from the OS which is not recursive
	// so we have to create the temporary test storage by two function calls.
	// We ignore any error possibly returned by creating the first one, because if it errors,
	// then the second one certainly will as well.
	_ = fs.EnsureDirectoryExists(path)
	err := fs.EnsureDirectoryExists(filepath.Join(path, "test"))
	checkError(t, err)
}

func SetupTestStorage(t *testing.T) {
	CreateTestStorage(t)
	path := FindTestdataFolder(t)
	err := fs.CopyDirectory(filepath.Join(path, "teststorage"), filepath.Join(path, "storage", "test"))
	checkError(t, err)
}

func PrettyPrint(t *testing.T, ob interface{}) string {
	b, err := json.MarshalIndent(ob, "", "  ")
	require.NoError(t, err)
	return string(b)
}
