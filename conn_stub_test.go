package routerosclient

import (
	"fmt"
	"time"

	"github.com/go-routeros/routeros"
	"github.com/go-routeros/routeros/proto"
)

type ConnStub struct {
	q chan *routeros.Reply
}

func (c *ConnStub) RunArgs(s []string) (*routeros.Reply, error) {
	select {
	case r := <-c.q:
		return r, nil
	case <-time.After(time.Second):
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
		for _, v := range reply {
			r.Re = append(r.Re, &proto.Sentence{
				Word: "!re",
				Map:  v,
			})
		}
	}

	if done != nil {
		r.Done.Map = done
	}

	select {
	case c.q <- r:
	case <-time.After(time.Millisecond):
	}

	return nil, true
}
