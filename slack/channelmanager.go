package slack

import (
	"sync"
	"time"
)

type ChannelManager struct {
	channels sync.Map
}

type ChannelData struct {
	sync.RWMutex
	Ch     chan string
	Closed bool
	Timer  *time.Timer
}

func (cd *ChannelData) CloseChannel() {
	cd.Lock()
	defer cd.Unlock()
	close(cd.Ch)
	cd.Closed = true
}

func (cd *ChannelData) SafeSend(message string) bool {
	cd.RLock()
	defer cd.RUnlock()
	if cd.Closed {
		return false
	}
	cd.Ch <- message
	return true
}

func (m *ChannelManager) CreateChannel(id string) (ch chan string) {

	if _, exists := m.channels.Load(id); !exists {
		ch = make(chan string, 1)

		m.channels.Store(id, &ChannelData{
			Ch: ch,
		})
	}
	return
}

func (m *ChannelManager) DeleteChannel(id string) {
	chData, exists := m.channels.Load(id)
	if exists {
		chData.(*ChannelData).CloseChannel()
		m.channels.Delete(id)
	}
}

func (m *ChannelManager) GetChannel(id string) (chan string, bool) {
	chData, exists := m.channels.Load(id)
	if exists {
		return chData.(*ChannelData).Ch, true
	}
	return nil, false
}

func (m *ChannelManager) SendMessage(id string, message string) bool {
	chData, ok := m.channels.Load(id)
	if ok {
		return chData.(*ChannelData).SafeSend(message)
	}
	return false
}
