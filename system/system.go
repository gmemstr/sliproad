package system

import (
	"encoding/json"
	"github.com/gmemstr/nas/common"
	"io/ioutil"
	"net/http"
	"syscall"
)

type Config struct {
	ColdStorage string
	HotStorage  string
}

type UsageStat struct {
	Available int64
	Free      int64
	Total     int64
}

type UsageStats struct {
	ColdStorage UsageStat
	HotStorage  UsageStat
}

func DiskUsages() common.Handler {

	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		var statHot syscall.Statfs_t
		var statCold syscall.Statfs_t

		d, err := ioutil.ReadFile("assets/config/config.json")
		if err != nil {
			panic(err)
		}

		var config Config
		err = json.Unmarshal(d, &config)
		if err != nil {
			panic(err)
		}

		// Default to hot storage
		storage := config.HotStorage
		err = syscall.Statfs(storage, &statHot)
		if err != nil {
			panic(err)
		}
		hotStats := UsageStat{
			Free: statHot.Bsize * int64(statHot.Bfree),
			Total: statHot.Bsize * int64(statHot.Blocks),
		}
		storage = config.ColdStorage
		err = syscall.Statfs(storage, &statCold)
		if err != nil {
			panic(err)
		}
		coldStats := UsageStat{
			Free: statCold.Bsize * int64(statCold.Bfree),
			Total: statCold.Bsize * int64(statCold.Blocks),
		}
		usages := UsageStats{
			HotStorage: hotStats,
			ColdStorage: coldStats,
		}
		// Available blocks * size per block = available space in bytes
		resultJson, err := json.Marshal(usages)
		w.Write(resultJson)
		return nil
	}
}