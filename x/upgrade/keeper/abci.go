package keeper

import (
	"context"
	"errors"
	"fmt"
	"time"

	consensusv1 "cosmossdk.io/api/cosmos/consensus/v1"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PreBlocker will check if there is a scheduled plan and if it is ready to be executed.
// If the current height is in the provided set of heights to skip, it will skip and clear the upgrade plan.
// If it is ready, it will execute it if the handler is installed, and panic/abort otherwise.
// If the plan is not ready, it will ensure the handler is not registered too early (and abort otherwise).
//
// The purpose is to ensure the binary is switched EXACTLY at the desired block, and to allow
// a migration to be executed if needed upon this switch (migration defined in the new binary)
// skipUpgradeHeightArray is a set of block heights for which the upgrade must be skipped
func (k Keeper) PreBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	blockHeight := k.environment.HeaderService.GetHeaderInfo(ctx).Height
	plan, err := k.GetUpgradePlan(ctx)
	if err != nil && !errors.Is(err, types.ErrNoUpgradePlanFound) {
		return err
	}
	found := err == nil

	sdkCtx := sdk.UnwrapSDKContext(ctx) // TODO remove with consensus messages
	if !k.DowngradeVerified() {
		k.SetDowngradeVerified(true)
		// This check will make sure that we are using a valid binary.
		// It'll panic in these cases if there is no upgrade handler registered for the last applied upgrade.
		// 1. If there is no scheduled upgrade.
		// 2. If the plan is not ready.
		// 3. If the plan is ready and skip upgrade height is set for current height.
		if !found || !plan.ShouldExecute(blockHeight) || (plan.ShouldExecute(blockHeight) && k.IsSkipHeight(blockHeight)) {
			lastAppliedPlan, _, err := k.GetLastCompletedUpgrade(ctx)
			if err != nil {
				return err
			}

			if lastAppliedPlan != "" && !k.HasHandler(lastAppliedPlan) {
				var appVersion uint64

				var res consensusv1.QueryParamsResponse
				if err := k.environment.RouterService.QueryRouterService().InvokeTyped(ctx, &consensusv1.QueryParamsRequest{}, &res); err != nil {
					return errors.New("failed to query consensus params")
				}
				if res.Params.Version != nil {
					appVersion = res.Params.Version.App
				}

				return fmt.Errorf("wrong app version %d, upgrade handler is missing for %s upgrade plan", appVersion, lastAppliedPlan)
			}
		}
	}

	if !found {
		return nil
	}

	logger := k.Logger(ctx)

	// To make sure clear upgrade is executed at the same block
	if plan.ShouldExecute(blockHeight) {
		// If skip upgrade has been set for current height, we clear the upgrade plan
		if k.IsSkipHeight(blockHeight) {
			skipUpgradeMsg := fmt.Sprintf("UPGRADE \"%s\" SKIPPED at %d: %s", plan.Name, plan.Height, plan.Info)
			logger.Info(skipUpgradeMsg)

			// Clear the upgrade plan at current height
			if err := k.ClearUpgradePlan(ctx); err != nil {
				return err
			}
			return nil
		}

		// Prepare shutdown if we don't have an upgrade handler for this upgrade name (meaning this software is out of date)
		if !k.HasHandler(plan.Name) {
			// Write the upgrade info to disk. The UpgradeStoreLoader uses this info to perform or skip
			// store migrations.
			err := k.DumpUpgradeInfoToDisk(blockHeight, plan)
			if err != nil {
				return fmt.Errorf("unable to write upgrade info to filesystem: %w", err)
			}

			upgradeMsg := BuildUpgradeNeededMsg(plan)
			logger.Error(upgradeMsg)

			// Returning an error will end up in a panic
			return errors.New(upgradeMsg)
		}

		// We have an upgrade handler for this upgrade name, so apply the upgrade
		logger.Info(fmt.Sprintf("applying upgrade \"%s\" at %s", plan.Name, plan.DueAt()))
		sdkCtx = sdkCtx.WithBlockGasMeter(storetypes.NewInfiniteGasMeter())
		if err := k.ApplyUpgrade(sdkCtx, plan); err != nil {
			return err
		}
		return nil
	}

	// if we have a pending upgrade, but it is not yet time, make sure we did not
	// set the handler already
	if k.HasHandler(plan.Name) {
		downgradeMsg := fmt.Sprintf("BINARY UPDATED BEFORE TRIGGER! UPGRADE \"%s\" - in binary but not executed on chain. Downgrade your binary", plan.Name)
		logger.Error(downgradeMsg)

		// Returning an error will end up in a panic
		return errors.New(downgradeMsg)
	}
	return nil
}

// BuildUpgradeNeededMsg prints the message that notifies that an upgrade is needed.
func BuildUpgradeNeededMsg(plan types.Plan) string {
	return fmt.Sprintf("UPGRADE \"%s\" NEEDED at %s: %s", plan.Name, plan.DueAt(), plan.Info)
}
