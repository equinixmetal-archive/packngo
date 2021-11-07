package metadata

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ServerMock() (baseURL string, mux *http.ServeMux, teardownFn func()) {
	mux = http.NewServeMux()
	srv := httptest.NewServer(mux)
	return srv.URL, mux, srv.Close
}

func Test_Deserialization(t *testing.T) {
	baseURL, mux, teardown := ServerMock()
	defer teardown()

	mux.HandleFunc("/metadata", func(w http.ResponseWriter, req *http.Request) {
		b, err := ioutil.ReadFile("testdata/" + "device.json")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(b); err != nil {
			t.Fatal(err)
		}
	})

	device, err := GetMetadataFromURL(baseURL)
	assert.Nil(t, err)
	assert.NotNil(t, device)

	assert.Equal(t, "9307dc37-7f39-400b-9cd2-009087434a95", device.ID)
	assert.Equal(t, "spcqvzylz6-worker-2409003", device.Hostname)

	volumes := device.Volumes
	assert.Equal(t, 1, len(volumes))
	assert.Equal(t, "volume-b7f8e13c", volumes[0].Name)
	assert.Equal(t, "iqn.2013-05.com.daterainc:tc:01:sn:60448a8a63a20a82", volumes[0].IQN)
	assert.Equal(t, 2, len(volumes[0].IPs))
	assert.Equal(t, "10.144.35.132", volumes[0].IPs[0].String())
	assert.Equal(t, "10.144.51.11", volumes[0].IPs[1].String())
	assert.Equal(t, 10, volumes[0].Capacity.Size)
	assert.Equal(t, "gb", volumes[0].Capacity.Unit)

}

func TestIPNet_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		i       *IPNet
		args    args
		want    *IPNet
		wantErr bool
	}{
		{
			name:    "ipv4",
			args:    args{data: []byte("192.168.0.0/24")},
			want:    &IPNet{Mask: net.CIDRMask(24, 32), IP: net.IPv4(192, 168, 0, 0)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IPNet{}
			err := i.UnmarshalJSON(tt.args.data)

			if i != tt.want {
				t.Errorf("IPNet.UnmarshalJSON() want = %v, got %v", tt.want, i)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("IPNet.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
