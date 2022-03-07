package plugin

import (
	"context"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("escort")
var regIpDash = regexp.MustCompile(`^(\d{1,3}-\d{1,3}-\d{1,3}-\d{1,3})(-\d+)?\.`)

// Dump implement the plugin interface.
type Escort struct {
	Next  plugin.Handler
	state *State
}

// taken from https://github.com/wenerme/coredns-ipin
func (e Escort) Resolve(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (*dns.Msg, int, error) {
	state := request.Request{W: w, Req: r}

	a := new(dns.Msg)
	a.SetReply(r)
	a.Compress = true
	a.Authoritative = true

	matches := regIpDash.FindStringSubmatch(state.QName())
	if len(matches) > 1 {
		ip := matches[1]
		ip = strings.Replace(ip, "-", ".", -1)

		var rr dns.RR
		rr = new(dns.A)
		rr.(*dns.A).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()}
		rr.(*dns.A).A = net.ParseIP(ip).To4()

		a.Answer = []dns.RR{rr}

		if len(matches[2]) > 0 {
			srv := new(dns.SRV)
			srv.Hdr = dns.RR_Header{Name: "_port." + state.QName(), Rrtype: dns.TypeSRV, Class: state.QClass()}
			if state.QName() == "." {
				srv.Hdr.Name = "_port." + state.QName()
			}
			port, _ := strconv.Atoi(matches[2][1:])
			srv.Port = uint16(port)
			srv.Target = "."

			a.Extra = []dns.RR{srv}
		}
	} else {
		// return empty
	}

	return a, 0, nil
}

// ServeDNS implements the plugin.Handler interface.
func (e Escort) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := &request.Request{W: w, Req: r}
	_ = state
	log.Info("resolving dns request")
	a, i, err := e.Resolve(ctx, w, r)
	if err != nil {
		return i, err
	}
	return 0, w.WriteMsg(a)
}

// Name implements the Handler interface.
func (e Escort) Name() string { return "blowback" }

func New() *Escort { return &Escort{} }
