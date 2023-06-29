package choosingstrategy

import (
	"context"

	"karma8_s3_hometask/src/utils"
)

type ServerCount int

type ServerChoosingStrategy interface {
	GetServerCount() ServerCount
	SetServerCount(ServerCount) // for dynamic changing the number of servers

	// Depending on its strategy, returns the array of servers to store files in
	Choose(ctx context.Context, sizes utils.Sizes) ([utils.FileParts]ServerCount, error)

	// Commits the choice. It is crucial to call it for the strategies
	// where the calls for Choose() depend on one another.
	// For example, if to the server[i] it has been uploaded something big,
	// may be in the next call, Choose() better return other servers to
	// distribute files equally enough.
	// However, if the choice is made randomly or the strategy can check
	// the server free spaces then it might be not important.
	CommitChoice(ctx context.Context, servers [utils.FileParts]ServerCount) error
}
