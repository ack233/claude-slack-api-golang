package requests

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"
	"slackapi/pkgs/initfunc"
	"slackapi/pkgs/zlog"

	"time"

	"github.com/go-resty/resty/v2"
)

const pioTimeout = 300
const pTimeOut = 100
const retry = 2

var Client *ClientStruct

type ClientStruct struct {
	*resty.Client
}

func init() {
	initfunc.RegisterInitFunc(
		func() {
			Client = CreateClinet()
		},
	)
}

func (c *ClientStruct) SetNotlsWithNewClinet() *ClientStruct {
	newclint := CreateClinet()
	newclint.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	return newclint
}

func (c *ClientStruct) SetLocalAddrWithNewClinet(ip string) *ClientStruct {
	newclint := CreateClinet()
	return newclint.SetLocalAddr(ip)

}

func (c *ClientStruct) SetLocalAddr(ip string) *ClientStruct {
	c.SetTransport(
		CustomTransport{
			CTimeout:  pTimeOut * time.Second,
			RWTimeout: pioTimeout * time.Second,
			LocalAddr: ip,
		}.Transport(),
	)
	return c
}

type RequestStruct struct {
	*resty.Request
}

func (r *RequestStruct) Notparse() *RequestStruct {
	r.SetDoNotParseResponse(true)
	return r
}

func (c *ClientStruct) R() *RequestStruct {
	return &RequestStruct{
		c.Client.R(),
	}
}

func CreateClinet() *ClientStruct {
	clientConfiguration := &ClientStruct{
		resty.New(),
	}
	clientConfiguration.
		SetTransport(CustomTransport{
			CTimeout:  pTimeOut * time.Second,
			RWTimeout: pioTimeout * time.Second,
			//LocalAddr: "your_ip_address",
		}.Transport(),
		).
		//SetTimeout(time.Duration(pTimeOut) * time.Second)// 总的请求时间,必须大于Transport.Dial: TimeoutDialer(30 * time.Second, 1 * time.Minute)
		SetRetryCount(retry).
		SetRetryWaitTime(100 * time.Nanosecond).
		AddRetryCondition(
			func(response *resty.Response, err error) bool {
				return !response.IsSuccess() || err != nil
			},
		).
		OnAfterResponse(
			func(c *resty.Client, resp *resty.Response) error {
				// Now you have access to Client and current Response object
				// manipulate it as per your need
				if !resp.IsSuccess() {
					return errors.New("request failed,http code is " + resp.Status())

				}
				return nil // if its success otherwise return error
			}).
		SetLogger(zlog.SugLog)
	return clientConfiguration
}

func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		//conn, err := net.DialTimeout(netw, addr, cTimeout)
		d := net.Dialer{
			Timeout:   cTimeout,
			DualStack: true}
		conn, err := d.Dial(netw, addr)

		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

func Parsebody_to_json(resp *resty.Response) map[string]interface{} {
	var v interface{}
	json.Unmarshal(resp.Body(), &v)
	return v.(map[string]interface{})
}

func Ecocde_json(v any) ([]byte, error) {
	e, err := json.Marshal(v)
	return e, err
}
