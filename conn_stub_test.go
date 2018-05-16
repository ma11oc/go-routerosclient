package routerosclient

import (
	"fmt"
	"time"

	"github.com/go-routeros/routeros"
	"github.com/go-routeros/routeros/proto"
)

type ConnStub struct {
	repliesQueue chan *routeros.Reply
}

func (c *ConnStub) RunArgs(s []string) (*routeros.Reply, error) {
	select {
	case reply := <-c.repliesQueue:
		return reply, nil
	case <-time.After(3 * time.Second):
		return nil, fmt.Errorf("timeout exceeded getting reply from queue")
	}
}

func (c *ConnStub) Close() {}

func (c *ConnStub) buildReply(reply []map[string]string, done map[string]string) (error, bool) {
	r := &routeros.Reply{
		Re: []*proto.Sentence{},
		Done: &proto.Sentence{
			Word: "!done",
			Map:  make(map[string]string),
		},
	}

	if len(reply) > 0 {
		for i := 0; i < len(reply); i++ {
			r.Re = append(r.Re, &proto.Sentence{
				Word: "!re",
				Map:  reply[i],
			})
		}
	}

	if done != nil {
		r.Done.Map = done
	}

	go func() {
		c.repliesQueue <- r
	}()

	return nil, true
}
