package main

import (
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/cbednarski/hostess/hostess"
)

// CopyHostsFile creates a temporary hosts file in the system temp directory,
// sets the HOSTESS_PATH environment variable, and returns the file path and a
// cleanup function
func CopyHostsFile(t *testing.T) (string, func()) {
	t.Helper()

	fixture, err := os.Open("testdata/ubuntu.hosts")
	if err != nil {
		t.Fatal(err)
	}

	temp, err := ioutil.TempFile("", "hostess-test-*")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := io.Copy(temp, fixture); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv(hostess.EnvHostessPath, temp.Name()); err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		os.Remove(temp.Name())
	}

	return temp.Name(), cleanup
}

func TestFormat(t *testing.T) {
	temp, cleanup := CopyHostsFile(t)
	defer cleanup()

	if err := wrappedMain(strings.Split("hostess fmt", " ")); err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadFile(temp)
	if err != nil {
		t.Fatal(err)
	}
	output := string(data)

	expected := `127.0.0.1 localhost myapp.local
127.0.1.1 ubuntu
192.168.0.30 raspberrypi
::1 ip6-localhost ip6-loopback
fe00:: ip6-localnet
ff00:: ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
`

	if runtime.GOOS == "windows" {
		expected = `127.0.0.1 localhost
127.0.0.1 myapp.local
127.0.1.1 ubuntu
192.168.0.30 raspberrypi
::1 ip6-localhost
::1 ip6-loopback
fe00:: ip6-localnet
ff00:: ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
`
	}

	if output != expected {
		t.Errorf("--- Expected ---\n%s\n--- Found ---\n%s\n", expected, output)
	}
}

func TestAddHostname(t *testing.T) {
	temp, cleanup := CopyHostsFile(t)
	defer cleanup()

	if err := wrappedMain(strings.Split("hostess add my.new.website 127.0.0.1", " ")); err != nil {
		t.Fatal(err)
	}
	if err := wrappedMain(strings.Split("hostess add mediaserver 192.168.0.82", " ")); err != nil {
		t.Fatal(err)
	}
	if err := wrappedMain(strings.Split("hostess add myapp.local 10.20.0.23", " ")); err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadFile(temp)
	if err != nil {
		t.Fatal(err)
	}
	output := string(data)

	expected := `127.0.0.1 localhost my.new.website
127.0.1.1 ubuntu
10.20.0.23 myapp.local
192.168.0.30 raspberrypi
192.168.0.82 mediaserver
::1 ip6-localhost ip6-loopback
fe00:: ip6-localnet
ff00:: ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
`

	if runtime.GOOS == "windows" {
		expected = `127.0.0.1 localhost
127.0.0.1 my.new.website
127.0.1.1 ubuntu
10.20.0.23 myapp.local
192.168.0.30 raspberrypi
192.168.0.82 mediaserver
::1 ip6-localhost
::1 ip6-loopback
fe00:: ip6-localnet
ff00:: ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
`
	}

	if output != expected {
		t.Errorf("--- Expected ---\n%s\n--- Found ---\n%s\n", expected, output)
	}
}

func TestRemoveHostname(t *testing.T) {
	temp, cleanup := CopyHostsFile(t)
	defer cleanup()

	if err := wrappedMain(strings.Split("hostess rm myapp.local", " ")); err != nil {
		t.Fatal(err)
	}
	if err := wrappedMain(strings.Split("hostess rm raspberrypi", " ")); err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadFile(temp)
	if err != nil {
		t.Fatal(err)
	}
	output := string(data)

	expected := `127.0.0.1 localhost
127.0.1.1 ubuntu
::1 ip6-localhost ip6-loopback
fe00:: ip6-localnet
ff00:: ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
`

	if runtime.GOOS == "windows" {
		expected = `127.0.0.1 localhost
127.0.1.1 ubuntu
::1 ip6-localhost
::1 ip6-loopback
fe00:: ip6-localnet
ff00:: ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
`
	}

	if output != expected {
		t.Errorf("--- Expected ---\n%s\n--- Found ---\n%s\n", expected, output)
	}
}

func TestHostnameOff(t *testing.T) {
	temp, cleanup := CopyHostsFile(t)
	defer cleanup()

	if err := wrappedMain(strings.Split("hostess off myapp.local", " ")); err != nil {
		t.Fatal(err)
	}
	if err := wrappedMain(strings.Split("hostess off raspberrypi", " ")); err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadFile(temp)
	if err != nil {
		t.Fatal(err)
	}
	output := string(data)

	expected := `127.0.0.1 localhost
# 127.0.0.1 myapp.local
127.0.1.1 ubuntu
# 192.168.0.30 raspberrypi
::1 ip6-localhost ip6-loopback
fe00:: ip6-localnet
ff00:: ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
`

	if runtime.GOOS == "windows" {
		expected = `127.0.0.1 localhost
# 127.0.0.1 myapp.local
127.0.1.1 ubuntu
# 192.168.0.30 raspberrypi
::1 ip6-localhost
::1 ip6-loopback
fe00:: ip6-localnet
ff00:: ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
`
	}

	if output != expected {
		t.Errorf("--- Expected ---\n%s\n--- Found ---\n%s\n", expected, output)
	}
}
