package main

import (
	"fmt"
	"strings"
)

type Member struct {
	Name string
	// ID field can be 0, which is unknown ID.
	// We know the ID of a member when we get the member information from etcd,
	// but not from Kubernetes pod list.
	ID uint64
}

type MemberSet map[string]Member

// the set of all members of s1 that are not members of s2
func (ms MemberSet) Diff(other MemberSet) MemberSet {
	diff := MemberSet{}
	for n, m := range ms {
		if _, ok := other[n]; !ok {
			diff[n] = m
		}
	}
	return diff
}

func (ms MemberSet) Size() int {
	return len(ms)
}

func (ms MemberSet) String() string {
	var mstring []string

	for m := range ms {
		mstring = append(mstring, m)
	}
	return strings.Join(mstring, ",")
}

func (ms MemberSet) PickOne() Member {
	for _, m := range ms {
		return m
	}
	panic("empty")
}

func (ms MemberSet) PeerURLPairs() []string {
	ps := make([]string, 0)
	for _, m := range ms {
		ps = append(ps, fmt.Sprintf("%s=%s", m.Name, makeEtcdPeerAddr(m.Name)))
	}
	return ps
}

func (ms MemberSet) Add(m Member) {
	ms[m.Name] = m
}

func (ms MemberSet) Remove(name string) {
	delete(ms, name)
}

func (ms MemberSet) ClientURLs() []string {
	endpoints := make([]string, len(ms))
	i := 0
	for _, m := range ms {
		endpoints[i] = makeClientAddr(m.Name)
		i++
	}
	return endpoints
}
