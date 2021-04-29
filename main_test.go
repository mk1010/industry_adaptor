package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/mk1010/industry_adaptor/nclink"

	"github.com/apache/dubbo-go/common/logger"
	json "github.com/json-iterator/go"
	"golang.org/x/net/http2"
)

type mk struct {
	Name string `default:"123" json:"name"`
	Age  string `default:"123" json:"age"`
}

func (m *mk) Hi(s string) string {
	return "hello" + s
}

func TestChan(t *testing.T) {
	m := &mk{}
	ret := reflect.ValueOf(m).MethodByName("Hi").IsValid()
	err := json.Unmarshal([]byte(""), m)
	b, _ := json.Marshal(m)
	s := string(b)
	t.Log(err, m.Name, ret, s)
}

func TestNet(t *testing.T) {
	tcp, err := net.ListenTCP("tcp4", &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: 0,
	})
	t.Log(tcp.Addr().String(), err)
}

func Benchmark_test_reflect(b *testing.B) {
	m := &mk{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if t := reflect.ValueOf(*m).MethodByName("Hi"); t.IsValid() {
			t.Call([]reflect.Value{reflect.ValueOf("mk")})
		}
	}
}

func TestHttpsClient(t *testing.T) {
	req, err := http.NewRequest("GET", "https://www.haha.com", nil)
	if err == nil {
		t.Log(err)
	}
	req.Proto = "HTTP/2"
	pool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile("../idustry/ca.crt")
	if err != nil {
		t.Logf("Reading server certificate: %s", err)
	}
	pool.AppendCertsFromPEM(caCert)
	cliCrt, err := tls.LoadX509KeyPair("../idustry/test.pem", "../idustry/test.key")
	if err != nil {
		fmt.Println("LoadX509keypair err: ", err)
		return
	}
	tr := &http2.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
		},
	}
	client := &http.Client{Transport: tr}
	conn, err := client.Do(req)
	if err == nil {
		t.Log(err)
	}
	t.Logf("%+v,%v", conn, err)
}

func TestHttpClient(t *testing.T) {
	/*client := http.Client{
		// Skip TLS dial
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}*/

	resp, err := http.DefaultClient.Get("https://cn.bing.com")
	if err != nil {
		t.Log(fmt.Errorf("error making request: %v", err))
	}
	bytesFile, _ := ioutil.ReadAll(resp.Body)
	t.Log(string(bytesFile))
	defer resp.Body.Close()
}

func TestHttpSUnauthClient(t *testing.T) {
	tr := &http2.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", "https://127.0.0.1:8080/device/LogicID=123", nil)

	req.Proto = "HTTP/2"
	conn, err := client.Do(req)
	if err == nil {
		t.Log(err)
	}
	conn2, err := client.Do(req)
	if err == nil {
		t.Log(err)
	}
	conn3, err := client.Do(req)
	if err == nil {
		t.Log(err)
	}
	conn4, err := client.Do(req)
	if err == nil {
		t.Log(err)
	}
	t.Logf("%+v,%+v,%+v,%+v,%v", conn, conn2, conn3, conn4, err)
}

func TestSystemEdian(t *testing.T) {
	var i int = 0x1
	bs := (*[4]byte)(unsafe.Pointer(&i))
	if bs[0] == 0 {
		fmt.Println("system edian is little endian")
	} else {
		fmt.Println("system edian is big endian")
	}
	testBigEndian()
	testLittleEndian()
}

func testBigEndian() {
	// 0000 0000 0000 0000   0000 0001 1111 1111
	var testInt int32 = 256
	fmt.Printf("%d use big endian: \n", testInt)
	var testBytes []byte = make([]byte, 4)
	binary.BigEndian.PutUint32(testBytes, uint32(testInt))
	fmt.Println("int32 to bytes:", testBytes)

	convInt := binary.BigEndian.Uint32(testBytes)
	fmt.Printf("bytes to int32: %d\n\n", convInt)
}

func testLittleEndian() {
	// 0000 0000 0000 0000   0000 0001 1111 1111
	var testInt int32 = 256
	fmt.Printf("%d use little endian: \n", testInt)
	var testBytes []byte = make([]byte, 4)
	binary.LittleEndian.PutUint32(testBytes, uint32(testInt))
	fmt.Println("int32 to bytes:", testBytes)

	convInt := binary.LittleEndian.Uint32(testBytes)
	fmt.Printf("bytes to int32: %d\n\n", convInt)
}

func TestExe(t *testing.T) {
	/*
		val := func() {
			t.Log("execute")
		}
	*/
	val := 1
	rval := reflect.ValueOf(val)
	if rval.Kind() == reflect.Func {
		t.Log(rval.Call(nil))
	} else {
		t.Log("not exe")
	}
	jsonStrings := `{
		"mk":123,
		"reye":"12",
		"ee":[{
			"name":"1234",
			"age":"ss"
		},{
			"reyee":"123"
		}]
	}`
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStrings), &m)
	any := json.Get([]byte(jsonStrings), "ee", 0)
	t.Log(m, err)
	g := any.GetInterface()
	t.Log(g)
	now := time.Now()
	t.Log(now.Unix()*1e3 + int64(now.Nanosecond()/1e6))
}

func TestGrpc(t *testing.T) {
	l := new(nclink.NCLinkAdaptor)
	l.AdaptorId = "123"
	b, err := json.Marshal(l)
	t.Log(string(b), err)
	logger.Errorf("mk1", "234")
	//	nclink.ErrorWrap(fmt.Errorf("mk"), "123")
}

func TestTimeOut(t *testing.T) {
	ticker := time.NewTicker(100 * time.Millisecond)

	select {
	case <-ticker.C:
	}
}
