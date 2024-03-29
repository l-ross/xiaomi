// Package miio providea a Golang implementation of the Xiaomi MIIO protocol
package miio

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"net"
	"sync"
)

const (
	DefaultIP   = "192.168.1.1"
	DefaultPort = 54321
)

// A Client for talking to a device using the Xiaomi MIIO protocol
type Client struct {
	token []byte

	h hash.Hash

	deviceID uint32
	stamp    uint32

	conn    net.Conn
	rwMutex sync.RWMutex

	blockSize int
	createEnc func() cipher.BlockMode
	createDec func() cipher.BlockMode

	options *Options
}

type Options struct {
	IP   string
	Port int
}

type Option func(*Options) error

// defaultOptions defines the default Options for a Client. To override these provide the appropriate
// Option when calling New.
func defaultOptions() *Options {
	return &Options{
		IP:   DefaultIP,
		Port: DefaultPort,
	}
}

// SetIP sets the destination IP of this Client
func SetIP(ip string) Option {
	return func(o *Options) error {
		o.IP = ip
		return nil
	}
}

// SetPort sets the destionation port of this Client
func SetPort(port int) Option {
	return func(o *Options) error {
		o.Port = port
		return nil
	}
}

func (c *Client) Connected() bool {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()
	return c.conn != nil
}

// New constructs a new Client
func New(token string, opts ...Option) (*Client, error) {
	// Validate token
	if len(token) != 32 {
		return nil, fmt.Errorf("token must be 32 characters")
	}

	t, err := hex.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex token: %w", err)
	}

	c := &Client{
		token:   t,
		h:       md5.New(),
		options: defaultOptions(),
	}

	for _, opt := range opts {
		err := opt(c.options)
		if err != nil {
			return nil, fmt.Errorf("error setting option: %w", err)
		}
	}

	// Set up encryption / decryption
	key := c.md5(t)
	iv := c.md5(append(key, t...))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	c.blockSize = block.BlockSize()

	c.createEnc = func() cipher.BlockMode {
		return cipher.NewCBCEncrypter(block, iv)
	}

	c.createDec = func() cipher.BlockMode {
		return cipher.NewCBCDecrypter(block, iv)
	}

	return c, nil
}

// Connect to the device.
//
// If Connect has already been called and Close has not then this is a no-op and will
// re-use the existing connection.
func (c *Client) Connect() error {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	if c.conn != nil {
		return nil
	}

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", c.options.IP, c.options.Port))
	if err != nil {
		return err
	}

	// Send hello packet to receive device ID and stamp ID
	deviceID, stamp, err := c.hello(conn)
	if err != nil {
		return err
	}

	c.deviceID = deviceID
	c.stamp = stamp
	c.conn = conn

	return nil
}

// Close the connection to the device
func (c *Client) Close() error {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	if c.conn == nil {
		return nil
	}

	err := c.conn.Close()
	c.conn = nil
	return err
}

// Send will perform the necessary handshake and then send the provided payload, response data
// is returned.
func (c *Client) Send(payload []byte) ([]byte, error) {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()

	if c.conn == nil {
		return nil, errors.New("client is not connected")
	}

	// Create and send request
	req, err := c.createRequest(payload)
	if err != nil {
		return nil, err
	}

	sent, err := c.conn.Write(req)
	if err != nil {
		return nil, err
	}
	if sent != len(req) {
		return nil, fmt.Errorf("expected to write %d but wrote %d", len(req), sent)
	}

	// Receive response
	var (
		// Response to return
		rsp []byte
		// Response is read in 4K chunks, for the first chunk we want to read the expected
		// length of the entire response.
		once = sync.Once{}
		// Default rspLen to max value until we have read the first chunk and know
		// the actual length we expect
		rspLen = ^uint16(0)
	)

	// Keep reading until we have the entire response
	for x := uint16(0); x < rspLen; {
		chunk := make([]byte, 4096)
		read, err := c.conn.Read(chunk)
		if err != nil {
			return nil, err
		}

		once.Do(func() {
			// First chunk, read the expected length of the entire response.
			rspLen = binary.BigEndian.Uint16(chunk[2:4])
		})

		rsp = append(rsp, chunk[0:read]...)
		x += uint16(read)
	}

	return c.decodeResponse(rsp)
}

// hello sends the hello packet and returns the device ID and stamp
func (c *Client) hello(conn net.Conn) (uint32, uint32, error) {
	helloPayload := []byte{
		// Magic number
		0x21, 0x31,
		// Length
		0x00, 0x20,
		// All the Fs
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
	}

	sent, err := conn.Write(helloPayload)
	if err != nil {
		return 0, 0, err
	}
	if sent != 32 {
		return 0, 0, fmt.Errorf("expected to write 32 but wrote %d", sent)
	}

	rsp := make([]byte, 32)

	read, err := conn.Read(rsp)
	if err != nil {
		return 0, 0, err
	}
	if read != 32 {
		return 0, 0, fmt.Errorf("error sending hello packet, expected to receive 32 bytes but got %d", read)
	}

	deviceID := binary.BigEndian.Uint32(rsp[8:12])
	stamp := binary.BigEndian.Uint32(rsp[12:16])

	return deviceID,
		stamp,
		nil
}

func (c *Client) createRequest(payload []byte) ([]byte, error) {
	// Append null byte if it's missing from the payload.
	if payload[len(payload)-1] != 0x00 {
		payload = append(payload, 0x00)
	}

	// Pad and encrypt the payload
	payload, err := c.pkcs7Pad(payload)
	if err != nil {
		return nil, err
	}

	enc := c.createEnc()
	enc.CryptBlocks(payload, payload)

	// Construct request payload
	d := newData()
	d.writeMagic()
	// Write length
	d.writeUint16(uint16(len(payload) + 32))
	// Always 0
	d.writeUint32(0)
	d.writeUint32(c.deviceID)
	d.writeUint32(c.stamp)

	// Write token in place of hash
	d.write(c.token)
	d.write(payload)

	// Calculate MD5 and overwrite token with it.
	d.seek(16)
	d.write(c.md5(d.bytes()))

	// Increment stamp
	c.stamp++

	return d.bytes(), nil
}

func (c *Client) decodeResponse(rsp []byte) ([]byte, error) {
	var err error

	if len(rsp) < 32 {
		return nil, fmt.Errorf("invalid response, expected at least 32 bytes but got %d", len(rsp))
	}

	d := newDataFromByte(rsp)

	// Read the MD5 from the response
	d.seek(16)
	responseMD5, err := d.readBytes(16)
	if err != nil {
		return nil, err
	}

	// Write the token over the MD5 and compute the packets MD5
	d.seek(16)
	d.write(c.token)
	expectedMD5 := c.md5(d.bytes())

	// Verify the MD5s match
	if !bytes.Equal(responseMD5, expectedMD5) {
		responseMD5String := hex.EncodeToString(responseMD5)
		expectedMD5String := hex.EncodeToString(expectedMD5)
		return nil, fmt.Errorf("token mismatch, expected %s got %s", responseMD5String, expectedMD5String)
	}

	// Decrypt the body
	body := d.bytes()[32:]
	if len(body) == 0 {
		return nil, errors.New("empty response body")
	}

	dec := c.createDec()
	dec.CryptBlocks(body, body)

	// Remove padding
	body, err = c.pkcs7Unpad(body)
	if err != nil {
		return nil, err
	}

	// Remove null byte
	body = body[0 : len(body)-1]

	return body, nil
}

func (c *Client) md5(b []byte) []byte {
	c.h.Reset()
	c.h.Write(b)
	return c.h.Sum(nil)
}

// Source: https://github.com/go-web/tokenizer/blob/master/pkcs7.go
func (c *Client) pkcs7Pad(b []byte) ([]byte, error) {
	if c.blockSize <= 0 {
		return nil, fmt.Errorf("invalid blocksize")
	}
	if b == nil || len(b) == 0 {
		return nil, fmt.Errorf("input cannot be empty")
	}
	n := c.blockSize - (len(b) % c.blockSize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

// Source: https://github.com/go-web/tokenizer/blob/master/pkcs7.go
func (c *Client) pkcs7Unpad(b []byte) ([]byte, error) {
	if c.blockSize <= 0 {
		return nil, fmt.Errorf("invalid blocksize")
	}
	if b == nil || len(b) == 0 {
		return nil, fmt.Errorf("input cannot be empty")
	}
	if len(b)%c.blockSize != 0 {
		return nil, fmt.Errorf("invalid padding on input")
	}
	paddingLen := b[len(b)-1]
	n := int(paddingLen)
	if n == 0 || n > len(b) {
		return nil, fmt.Errorf("invalid padding on input")
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != paddingLen {
			return nil, fmt.Errorf("invalid padding on input")
		}
	}
	return b[:len(b)-n], nil
}
