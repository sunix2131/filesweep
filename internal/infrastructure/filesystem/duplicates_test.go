package filesystem

import (
	"context"
	"testing"

	scan "filesweep/internal/domain/scan"
	"github.com/stretchr/testify/require"
)

func TestCancellationStopsHashQueueWithZeroWorkerSetting(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	calls := 0
	files := make([]scan.File, 100)
	result := ApplyHashesContext(ctx, files, 0, func(string) (HashResult, error) {
		calls++
		cancel()
		return HashResult{}, context.Canceled
	}, nil)
	require.Equal(t, 1, calls)
	require.Len(t, result, 100)
	require.Equal(t, context.Canceled.Error(), result[0].ScanError)
	require.Empty(t, result[1].SHA256)
}
