package choosingstrategy

import (
	"context"
	"hash/maphash"
	"math/rand"

	utils "karma8_s3_hometask/src/utils"
)

type randomChoosingStrategy struct {
	serverCount ServerCount
	random      *rand.Rand
}

func NewEqualUploadStrategy(serverCount ServerCount) ServerChoosingStrategy {
	res := &randomChoosingStrategy{}
	res.SetServerCount(serverCount)
	res.random = rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
	return res
}

func (r *randomChoosingStrategy) GetServerCount() ServerCount {
	return r.serverCount
}

func (r *randomChoosingStrategy) SetServerCount(serverCount ServerCount) {
	if serverCount < utils.FileParts {
		panic("number of servers must be >= dividable file parts")
	}
	r.serverCount = serverCount
}

func (r *randomChoosingStrategy) Choose(ctx context.Context, size utils.Sizes) ([utils.FileParts]ServerCount, error) {
	// TODO: Optimize! Cache in to some synchronized pools;
	// e.g. create 32 pools with the same size and return each at a time,
	// after random.Shuffle (below), return it back to the pools.
	// Don't fotget to recreate cached slices on serverCount change.
	servers := make([]ServerCount, r.serverCount)
	for i := range servers {
		servers[i] = ServerCount(i)
	}

	r.random.Shuffle(int(r.serverCount), func(i, j int) {
		t := servers[i]
		servers[i] = servers[j]
		servers[j] = t
	})

	shuffledServers := [utils.FileParts]ServerCount{}
	for i := range shuffledServers {
		shuffledServers[i] = servers[i]
	}
	return shuffledServers, nil
}

func (r *randomChoosingStrategy) CommitChoice(ctx context.Context, servers [utils.FileParts]ServerCount) error {
	// since the choices are random, they don't depend on one another
	return nil
}
