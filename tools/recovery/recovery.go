package recovery

import (
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"

	bsmt "github.com/bnb-chain/zkbas-smt"
	"github.com/bnb-chain/zkbas/tools/recovery/internal/config"
	"github.com/bnb-chain/zkbas/tools/recovery/internal/svc"
	"github.com/bnb-chain/zkbas/tree"
)

func RecoveryTreeDB(
	configFile string,
	blockHeight int64,
	serviceName string,
	batchSize int,
) {
	var c config.Config
	conf.MustLoad(configFile, &c)
	ctx := svc.NewServiceContext(c)
	logx.MustSetup(c.LogConf)
	logx.DisableStat()
	proc.AddShutdownListener(func() {
		logx.Close()
	})

	// dbinitializer tree database
	treeCtx := &tree.Context{
		Name:          serviceName,
		Driver:        c.TreeDB.Driver,
		LevelDBOption: &c.TreeDB.LevelDBOption,
		RedisDBOption: &c.TreeDB.RedisDBOption,
		Reload:        true,
	}
	treeCtx.SetOptions(bsmt.InitializeVersion(bsmt.Version(blockHeight) - 1))
	treeCtx.SetBatchReloadSize(batchSize)
	err := tree.SetupTreeDB(treeCtx)
	if err != nil {
		logx.Errorf("Init tree database failed: %s", err)
		return
	}

	// dbinitializer accountTree and accountStateTrees
	_, _, err = tree.InitAccountTree(
		ctx.AccountModel,
		ctx.AccountHistoryModel,
		blockHeight,
		treeCtx,
	)
	if err != nil {
		logx.Error("InitMerkleTree error:", err)
		return
	}
	// dbinitializer liquidityTree
	_, err = tree.InitLiquidityTree(
		ctx.LiquidityHistoryModel,
		blockHeight,
		treeCtx)
	if err != nil {
		logx.Errorf("InitLiquidityTree error: %s", err.Error())
		return
	}
	// dbinitializer nftTree
	_, err = tree.InitNftTree(
		ctx.NftHistoryModel,
		blockHeight,
		treeCtx)
	if err != nil {
		logx.Errorf("InitNftTree error: %s", err.Error())
		return
	}
}
