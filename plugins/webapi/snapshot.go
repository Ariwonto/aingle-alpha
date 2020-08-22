package webapi

import (
	"fmt"
	milestone2 "github.com/Ariwonto/aingle-alpha/pkg/model/milestone"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"

	"github.com/Ariwonto/aingle-alpha/pkg/config"
	"github.com/Ariwonto/aingle-alpha/pkg/model/milestone"
	"github.com/Ariwonto/aingle-alpha/plugins/snapshot"
)

func init() {
	addEndpoint("createSnapshotFile", createSnapshotFile, implementedAPIcalls)
}

func createSnapshotFile(i interface{}, c *gin.Context, abortSignal <-chan struct{}) {
	e := ErrorReturn{}
	query := &CreateSnapshotFile{}

	if err := mapstructure.Decode(i, query); err != nil {
		e.Error = fmt.Sprintf("%v: %v", ErrInternalError, err)
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	snapshotFilePath := filepath.Join(filepath.Dir(config.NodeConfig.GetString(config.CfgLocalSnapshotsPath)), fmt.Sprintf("export_%d.bin", query.TargetIndex))

	if err := snapshot.CreateLocalSnapshot(milestone2.Index(milestone.Index(query.TargetIndex)), snapshotFilePath, false, abortSignal); err != nil {
		e.Error = err.Error()
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	c.JSON(http.StatusOK, CreateSnapshotFileReturn{})
}
