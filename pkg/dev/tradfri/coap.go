package tradfri

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	coap "github.com/go-ocf/go-coap"
	codes "github.com/go-ocf/go-coap/codes"
	dtls "github.com/pion/dtls/v2"
)

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// coapConnectWith creates a secure COAP connection
func coapConnectWith(addr, key, value string, timeout time.Duration) (*coap.ClientConn, error) {
	if key == "" || value == "" {
		return nil, gopi.ErrBadParameter.WithPrefix("CoapConnect")
	} else if conn, err := coap.DialDTLSWithTimeout("udp", addr, &dtls.Config{
		PSK: func(hint []byte) ([]byte, error) {
			return []byte(value), nil
		},
		PSKIdentityHint: []byte(key),
		CipherSuites:    []dtls.CipherSuiteID{dtls.TLS_PSK_WITH_AES_128_CCM_8},
	}, timeout); err != nil {
		return nil, fmt.Errorf("%w (addr: %s)", err, addr)
	} else {
		return conn, nil
	}
}

// coapAuthenticate performs the gateway authentication and returns JSON response
func coapAuthenticate(conn *coap.ClientConn, id string, timeout time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Set identity
	path := filepath.Join("/", ROOT_GATEWAY, ATTR_AUTH)
	body := strings.NewReader(fmt.Sprintf("{%q:%q}", ATTR_IDENTITY, id))
	if response, err := conn.PostWithContext(ctx, path, coap.AppJSON, body); err != nil {
		return nil, err
	} else if response.Code() != codes.Created {
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(fmt.Sprint(response.Code()))
	} else {
		// Return success
		payload := response.Payload()
		return payload[0 : len(payload)-3], nil
	}
}

// coapRequestIdsForPath requests id's for a path and returns array response
func coapRequestIdsForPath(ctx context.Context, conn *coap.ClientConn, path ...string) ([]uint, error) {
	var ids []uint

	path_ := filepath.Join(append([]string{"/"}, path...)...)
	if response, err := conn.GetWithContext(ctx, path_); err != nil {
		return nil, err
	} else if response.Code() != codes.Content {
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(response.Code(), strconv.Quote(path_))
	} else if err := json.Unmarshal(response.Payload(), &ids); err != nil {
		return nil, err
	}

	// Success
	return ids, nil
}

// coapRequestObjForPath requests object for a path
func coapRequestObjForPath(ctx context.Context, conn *coap.ClientConn, obj interface{}, path ...string) error {
	path_ := filepath.Join(append([]string{"/"}, path...)...)
	if response, err := conn.GetWithContext(ctx, path_); err != nil {
		return err
	} else if response.Code() != codes.Content {
		return gopi.ErrUnexpectedResponse.WithPrefix(response.Code(), strconv.Quote(path_))
	} else if err := json.Unmarshal(response.Payload(), obj); err != nil {
		return err
	}

	// Success
	return nil
}
