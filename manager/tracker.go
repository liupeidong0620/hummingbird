package manager

import (
	"net"
	"time"

	"github.com/gofrs/uuid"
	"go.uber.org/atomic"

	"github.com/liupeidong0620/hummingbird/adapter"
)

type tracker interface {
	ID() string
	Close() error
}

type trackerInfo struct {
	Start time.Time `json:"start"`
	//End           time.Time         `json:"End"`
	UUID          uuid.UUID         `json:"id"`
	Metadata      *adapter.Metadata `json:"metadata"`
	UploadTotal   *atomic.Int64     `json:"upload"`
	DownloadTotal *atomic.Int64     `json:"download"`
}

type Tracker struct {
	net.Conn `json:"-"`

	*trackerInfo
	manager *Manager
}

func NewTracker(conn net.Conn, metadata *adapter.Metadata) *Tracker {
	id, _ := uuid.NewV4()

	t := &Tracker{
		Conn:    conn,
		manager: DefaultManager,
		trackerInfo: &trackerInfo{
			UUID:          id,
			Start:         time.Now(),
			Metadata:      metadata,
			UploadTotal:   atomic.NewInt64(0),
			DownloadTotal: atomic.NewInt64(0),
		},
	}

	DefaultManager.Join(t)
	return t
}

func (tt *Tracker) ID() string {
	return tt.UUID.String()
}

func (tt *Tracker) Read(b []byte) (int, error) {
	n, err := tt.Conn.Read(b)
	download := int64(n)
	//tt.manager.PushDownloaded(download)
	tt.DownloadTotal.Add(download)
	return n, err
}

func (tt *Tracker) Write(b []byte) (int, error) {
	n, err := tt.Conn.Write(b)
	upload := int64(n)
	//tt.manager.PushUploaded(upload)
	tt.UploadTotal.Add(upload)
	return n, err
}

func (tt *Tracker) Close() error {
	tt.manager.Leave(tt)
	return tt.Conn.Close()
}
