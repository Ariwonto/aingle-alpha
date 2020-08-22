package coordinator

import (
	"github.com/Ariwonto/aingle-alpha/pkg/model/aingle"
	"github.com/Ariwonto/aingle-alpha/pkg/model/milestone"
)

// CheckpointCaller is used to signal issued checkpoints.
func CheckpointCaller(handler interface{}, params ...interface{}) {
	handler.(func(checkpointIndex int, tipIndex int, tipsTotal int, txHash aingle.Hash))(params[0].(int), params[1].(int), params[2].(int), params[3].(aingle.Hash))
}

// MilestoneCaller is used to signal issued milestones.
func MilestoneCaller(handler interface{}, params ...interface{}) {
	handler.(func(index milestone.Index, tailTxHash aingle.Hash))(params[0].(milestone.Index), params[1].(aingle.Hash))
}
