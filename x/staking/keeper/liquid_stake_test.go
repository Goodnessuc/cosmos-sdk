package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

// Tests Set/Get TotalLiquidStakedTokens
func (s *KeeperTestSuite) TestTotalLiquidStakedTokens() {
	ctx, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	// Update the total liquid staked
	total := sdk.NewInt(100)
	keeper.SetTotalLiquidStakedTokens(ctx, total)

	// Confirm it was updated
	require.Equal(total, keeper.GetTotalLiquidStakedTokens(ctx), "initial")
}

// Tests Increase/Decrease TotalValidatorLiquidShares
func (s *KeeperTestSuite) TestValidatorLiquidShares() {
	ctx, keeper := s.ctx, s.stakingKeeper

	// Create a validator address
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	valAddress := sdk.ValAddress(pubKey.Address())

	// Set an initial total
	initial := sdk.NewDec(100)
	validator := types.Validator{
		OperatorAddress: valAddress.String(),
		LiquidShares:    initial,
	}
	keeper.SetValidator(ctx, validator)
}

// Tests DecreaseTotalLiquidStakedTokens
func (s *KeeperTestSuite) TestDecreaseTotalLiquidStakedTokens() {
	ctx, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	intitialTotalLiquidStaked := sdk.NewInt(100)
	decreaseAmount := sdk.NewInt(10)

	// Set the total liquid staked to an arbitrary value
	keeper.SetTotalLiquidStakedTokens(ctx, intitialTotalLiquidStaked)

	// Decrease the total liquid stake and confirm the total was updated
	err := keeper.DecreaseTotalLiquidStakedTokens(ctx, decreaseAmount)
	require.NoError(err, "no error expected when decreasing total liquid staked tokens")
	require.Equal(intitialTotalLiquidStaked.Sub(decreaseAmount), keeper.GetTotalLiquidStakedTokens(ctx))

	// Attempt to decrease by an excessive amount, it should error
	err = keeper.DecreaseTotalLiquidStakedTokens(ctx, intitialTotalLiquidStaked)
	require.ErrorIs(err, types.ErrTotalLiquidStakedUnderflow)
}

// Tests CheckExceedsValidatorBondCap
func (s *KeeperTestSuite) TestCheckExceedsValidatorBondCap() {
	ctx, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	testCases := []struct {
		name                string
		validatorShares     sdk.Dec
		validatorBondFactor sdk.Dec
		currentLiquidShares sdk.Dec
		newShares           sdk.Dec
		expectedExceeds     bool
	}{
		{
			// Validator Shares: 100, Factor: 1, Current Shares: 90 => 100 Max Shares, Capacity: 10
			// New Shares: 5 - below cap
			name:                "factor 1 - below cap",
			validatorShares:     sdk.NewDec(100),
			validatorBondFactor: sdk.NewDec(1),
			currentLiquidShares: sdk.NewDec(90),
			newShares:           sdk.NewDec(5),
			expectedExceeds:     false,
		},
		{
			// Validator Shares: 100, Factor: 1, Current Shares: 90 => 100 Max Shares, Capacity: 10
			// New Shares: 10 - at cap
			name:                "factor 1 - at cap",
			validatorShares:     sdk.NewDec(100),
			validatorBondFactor: sdk.NewDec(1),
			currentLiquidShares: sdk.NewDec(90),
			newShares:           sdk.NewDec(10),
			expectedExceeds:     false,
		},
		{
			// Validator Shares: 100, Factor: 1, Current Shares: 90 => 100 Max Shares, Capacity: 10
			// New Shares: 15 - above cap
			name:                "factor 1 - above cap",
			validatorShares:     sdk.NewDec(100),
			validatorBondFactor: sdk.NewDec(1),
			currentLiquidShares: sdk.NewDec(90),
			newShares:           sdk.NewDec(15),
			expectedExceeds:     true,
		},
		{
			// Validator Shares: 100, Factor: 2, Current Shares: 90 => 200 Max Shares, Capacity: 110
			// New Shares: 5 - below cap
			name:                "factor 2 - well below cap",
			validatorShares:     sdk.NewDec(100),
			validatorBondFactor: sdk.NewDec(2),
			currentLiquidShares: sdk.NewDec(90),
			newShares:           sdk.NewDec(5),
			expectedExceeds:     false,
		},
		{
			// Validator Shares: 100, Factor: 2, Current Shares: 90 => 200 Max Shares, Capacity: 110
			// New Shares: 100 - below cap
			name:                "factor 2 - below cap",
			validatorShares:     sdk.NewDec(100),
			validatorBondFactor: sdk.NewDec(2),
			currentLiquidShares: sdk.NewDec(90),
			newShares:           sdk.NewDec(100),
			expectedExceeds:     false,
		},
		{
			// Validator Shares: 100, Factor: 2, Current Shares: 90 => 200 Max Shares, Capacity: 110
			// New Shares: 110 - below cap
			name:                "factor 2 - at cap",
			validatorShares:     sdk.NewDec(100),
			validatorBondFactor: sdk.NewDec(2),
			currentLiquidShares: sdk.NewDec(90),
			newShares:           sdk.NewDec(110),
			expectedExceeds:     false,
		},
		{
			// Validator Shares: 100, Factor: 2, Current Shares: 90 => 200 Max Shares, Capacity: 110
			// New Shares: 111 - above cap
			name:                "factor 2 - above cap",
			validatorShares:     sdk.NewDec(100),
			validatorBondFactor: sdk.NewDec(2),
			currentLiquidShares: sdk.NewDec(90),
			newShares:           sdk.NewDec(111),
			expectedExceeds:     true,
		},
		{
			// Validator Shares: 100, Factor: 100, Current Shares: 90 => 10000 Max Shares, Capacity: 9910
			// New Shares: 100 - below cap
			name:                "factor 100 - below cap",
			validatorShares:     sdk.NewDec(100),
			validatorBondFactor: sdk.NewDec(100),
			currentLiquidShares: sdk.NewDec(90),
			newShares:           sdk.NewDec(100),
			expectedExceeds:     false,
		},
		{
			// Validator Shares: 100, Factor: 100, Current Shares: 90 => 10000 Max Shares, Capacity: 9910
			// New Shares: 9910 - at cap
			name:                "factor 100 - at cap",
			validatorShares:     sdk.NewDec(100),
			validatorBondFactor: sdk.NewDec(100),
			currentLiquidShares: sdk.NewDec(90),
			newShares:           sdk.NewDec(9910),
			expectedExceeds:     false,
		},
		{
			// Validator Shares: 100, Factor: 100, Current Shares: 90 => 10000 Max Shares, Capacity: 9910
			// New Shares: 9911 - above cap
			name:                "factor 100 - above cap",
			validatorShares:     sdk.NewDec(100),
			validatorBondFactor: sdk.NewDec(100),
			currentLiquidShares: sdk.NewDec(90),
			newShares:           sdk.NewDec(9911),
			expectedExceeds:     true,
		},
		{
			// Factor of -1 (disabled): Should always return false
			name:                "factor disabled",
			validatorShares:     sdk.NewDec(1),
			validatorBondFactor: sdk.NewDec(-1),
			currentLiquidShares: sdk.NewDec(1),
			newShares:           sdk.NewDec(1_000_000),
			expectedExceeds:     false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Update the validator bond factor
			params := keeper.GetParams(ctx)
			params.ValidatorBondFactor = tc.validatorBondFactor
			keeper.SetParams(ctx, params)

			// Create a validator with designated self-bond shares
			validator := types.Validator{
				LiquidShares:        tc.currentLiquidShares,
				ValidatorBondShares: tc.validatorShares,
			}

			// Check whether the cap is exceeded
			actualExceeds := keeper.CheckExceedsValidatorBondCap(ctx, validator, tc.newShares)
			require.Equal(tc.expectedExceeds, actualExceeds, tc.name)
		})
	}
}

// Tests TestCheckExceedsValidatorLiquidStakingCap
func (s *KeeperTestSuite) TestCheckExceedsValidatorLiquidStakingCap() {
	ctx, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	testCases := []struct {
		name                  string
		validatorLiquidCap    sdk.Dec
		validatorLiquidShares sdk.Dec
		validatorTotalShares  sdk.Dec
		newLiquidShares       sdk.Dec
		tokenizingShares      bool
		expectedExceeds       bool
	}{
		{
			// Cap: 10% - Delegation Below Threshold
			// Liquid Shares: 5, Total Shares: 95, New Liquid Shares: 1
			// => Liquid Shares: 5+1=6, Total Shares: 95+1=96 => 6/96 = 6% < 10% cap
			name:                  "10 percent cap _ native delegation _ below cap",
			validatorLiquidCap:    sdk.MustNewDecFromStr("0.1"),
			validatorLiquidShares: sdk.NewDec(5),
			validatorTotalShares:  sdk.NewDec(95),
			newLiquidShares:       sdk.NewDec(1),
			tokenizingShares:      false,
			expectedExceeds:       false,
		},
		{
			// Cap: 10% - Delegation At Threshold
			// Liquid Shares: 5, Total Shares: 95, New Liquid Shares: 5
			// => Liquid Shares: 5+5=10, Total Shares: 95+5=100 => 10/100 = 10% == 10% cap
			name:                  "10 percent cap _ native delegation _ equals cap",
			validatorLiquidCap:    sdk.MustNewDecFromStr("0.1"),
			validatorLiquidShares: sdk.NewDec(5),
			validatorTotalShares:  sdk.NewDec(95),
			newLiquidShares:       sdk.NewDec(4),
			tokenizingShares:      false,
			expectedExceeds:       false,
		},
		{
			// Cap: 10% - Delegation Exceeds Threshold
			// Liquid Shares: 5, Total Shares: 95, New Liquid Shares: 6
			// => Liquid Shares: 5+6=11, Total Shares: 95+6=101 => 11/101 = 11% > 10% cap
			name:                  "10 percent cap _ native delegation _ exceeds cap",
			validatorLiquidCap:    sdk.MustNewDecFromStr("0.1"),
			validatorLiquidShares: sdk.NewDec(5),
			validatorTotalShares:  sdk.NewDec(95),
			newLiquidShares:       sdk.NewDec(6),
			tokenizingShares:      false,
			expectedExceeds:       true,
		},
		{
			// Cap: 20% - Delegation Below Threshold
			// Liquid Shares: 20, Total Shares: 220, New Liquid Shares: 29
			// => Liquid Shares: 20+29=49, Total Shares: 220+29=249 => 49/249 = 19% < 20% cap
			name:                  "20 percent cap _ native delegation _ below cap",
			validatorLiquidCap:    sdk.MustNewDecFromStr("0.2"),
			validatorLiquidShares: sdk.NewDec(20),
			validatorTotalShares:  sdk.NewDec(220),
			newLiquidShares:       sdk.NewDec(29),
			tokenizingShares:      false,
			expectedExceeds:       false,
		},
		{
			// Cap: 20% - Delegation At Threshold
			// Liquid Shares: 20, Total Shares: 220, New Liquid Shares: 30
			// => Liquid Shares: 20+30=50, Total Shares: 220+30=250 => 50/250 = 20% == 20% cap
			name:                  "20 percent cap _ native delegation _ equals cap",
			validatorLiquidCap:    sdk.MustNewDecFromStr("0.2"),
			validatorLiquidShares: sdk.NewDec(20),
			validatorTotalShares:  sdk.NewDec(220),
			newLiquidShares:       sdk.NewDec(30),
			tokenizingShares:      false,
			expectedExceeds:       false,
		},
		{
			// Cap: 20% - Delegation Exceeds Threshold
			// Liquid Shares: 20, Total Shares: 220, New Liquid Shares: 31
			// => Liquid Shares: 20+31=51, Total Shares: 220+31=251 => 51/251 = 21% > 20% cap
			name:                  "20 percent cap _ native delegation _ exceeds cap",
			validatorLiquidCap:    sdk.MustNewDecFromStr("0.2"),
			validatorLiquidShares: sdk.NewDec(20),
			validatorTotalShares:  sdk.NewDec(220),
			newLiquidShares:       sdk.NewDec(31),
			tokenizingShares:      false,
			expectedExceeds:       true,
		},
		{
			// Cap: 50% - Native Delegation - Delegation At Threshold
			// Liquid shares: 0, Total Shares: 100, New Liquid Shares: 50
			// Total Liquid Shares: 0+50=50, Total Shares: 100+50=150
			// => 50/150 = 33% < 50% cap
			name:                  "50 percent cap _ native delegation _ delegation equals cap",
			validatorLiquidCap:    sdk.MustNewDecFromStr("0.5"),
			validatorLiquidShares: sdk.NewDec(0),
			validatorTotalShares:  sdk.NewDec(100),
			newLiquidShares:       sdk.NewDec(50),
			tokenizingShares:      false,
			expectedExceeds:       false,
		},
		{
			// Cap: 50% - Tokenized Delegation - Delegation At Threshold
			// Liquid shares: 0, Total Shares: 100, New Liquid Shares: 50
			// Total Liquid Shares => 0+50=50, Total Shares: 100,  New Liquid Shares: 50
			// => 50 / 100 = 50% == 50% cap
			name:                  "50 percent cap _ tokenized delegation _ delegation equals cap",
			validatorLiquidCap:    sdk.MustNewDecFromStr("0.5"),
			validatorLiquidShares: sdk.NewDec(0),
			validatorTotalShares:  sdk.NewDec(100),
			newLiquidShares:       sdk.NewDec(50),
			tokenizingShares:      true,
			expectedExceeds:       false,
		},
		{
			// Cap: 50% - Native Delegation - Delegation At Threshold
			// Liquid shares: 0, Total Shares: 100, New Liquid Shares: 51
			// Total Liquid Shares: 0+51=51, Total Shares: 100+51=151
			// => 51/150 = 33% < 50% cap
			name:                  "50 percent cap _ native delegation _ delegation equals cap",
			validatorLiquidCap:    sdk.MustNewDecFromStr("0.5"),
			validatorLiquidShares: sdk.NewDec(0),
			validatorTotalShares:  sdk.NewDec(100),
			newLiquidShares:       sdk.NewDec(51),
			tokenizingShares:      false,
			expectedExceeds:       false,
		},
		{
			// Cap: 50% - Tokenized Delegation - Delegation At Threshold
			// Liquid shares: 0, Total Shares: 100, New Liquid Shares: 50
			// Total Liquid Shares => 0+51=51, Total Shares: 100,  New Liquid Shares: 51
			// => 51 / 100 = 51% > 50% cap
			name:                  "50 percent cap _ tokenized delegation _ delegation equals cap",
			validatorLiquidCap:    sdk.MustNewDecFromStr("0.5"),
			validatorLiquidShares: sdk.NewDec(0),
			validatorTotalShares:  sdk.NewDec(100),
			newLiquidShares:       sdk.NewDec(51),
			tokenizingShares:      true,
			expectedExceeds:       true,
		},
		{
			// Cap of 0% - everything should exceed
			name:                  "0 percent cap",
			validatorLiquidCap:    sdk.ZeroDec(),
			validatorLiquidShares: sdk.NewDec(0),
			validatorTotalShares:  sdk.NewDec(1_000_000),
			newLiquidShares:       sdk.NewDec(1),
			tokenizingShares:      false,
			expectedExceeds:       true,
		},
		{
			// Cap of 100% - nothing should exceed
			name:                  "100 percent cap",
			validatorLiquidCap:    sdk.OneDec(),
			validatorLiquidShares: sdk.NewDec(1),
			validatorTotalShares:  sdk.NewDec(1_000_000),
			newLiquidShares:       sdk.NewDec(1),
			tokenizingShares:      false,
			expectedExceeds:       false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Update the validator liquid staking cap
			params := keeper.GetParams(ctx)
			params.ValidatorLiquidStakingCap = tc.validatorLiquidCap
			keeper.SetParams(ctx, params)

			// Create a validator with designated self-bond shares
			validator := types.Validator{
				LiquidShares:    tc.validatorLiquidShares,
				DelegatorShares: tc.validatorTotalShares,
			}

			// Check whether the cap is exceeded
			actualExceeds := keeper.CheckExceedsValidatorLiquidStakingCap(ctx, validator, tc.newLiquidShares, tc.tokenizingShares)
			require.Equal(tc.expectedExceeds, actualExceeds, tc.name)
		})
	}
}

// Tests SafelyIncreaseValidatorLiquidShares
func (s *KeeperTestSuite) TestSafelyIncreaseValidatorLiquidShares() {
	ctx, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	// Generate a test validator address
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	valAddress := sdk.ValAddress(pubKey.Address())

	// Helper function to check the validator's liquid shares
	checkValidatorLiquidShares := func(expected sdk.Dec, description string) {
		actualValidator, found := keeper.GetValidator(ctx, valAddress)
		require.True(found)
		require.Equal(expected.TruncateInt64(), actualValidator.LiquidShares.TruncateInt64(), description)
	}

	// Start with the following:
	//   Initial Liquid Shares: 0
	//   Validator Bond Shares: 10
	//   Validator TotalShares: 75
	//
	// Initial Caps:
	//   ValidatorBondFactor: 1 (Cap applied at 10 shares)
	//   ValidatorLiquidStakingCap: 25% (Cap applied at 25 shares)
	//
	// Cap Increases:
	//   ValidatorBondFactor: 10 (Cap applied at 100 shares)
	//   ValidatorLiquidStakingCap: 40% (Cap applied at 50 shares)
	initialLiquidShares := sdk.NewDec(0)
	validatorBondShares := sdk.NewDec(10)
	validatorTotalShares := sdk.NewDec(75)

	firstIncreaseAmount := sdk.NewDec(20)
	secondIncreaseAmount := sdk.NewDec(10) // total increase of 30

	initialBondFactor := sdk.NewDec(1)
	finalBondFactor := sdk.NewDec(10)
	initialLiquidStakingCap := sdk.MustNewDecFromStr("0.25")
	finalLiquidStakingCap := sdk.MustNewDecFromStr("0.4")

	// Create a validator with designated self-bond shares
	initialValidator := types.Validator{
		OperatorAddress:     valAddress.String(),
		LiquidShares:        initialLiquidShares,
		ValidatorBondShares: validatorBondShares,
		DelegatorShares:     validatorTotalShares,
	}
	keeper.SetValidator(ctx, initialValidator)

	// Set validator bond factor to a small number such that any delegation would fail,
	// and set the liquid staking cap such that the first stake would succeed, but the second
	// would fail
	params := keeper.GetParams(ctx)
	params.ValidatorBondFactor = initialBondFactor
	params.ValidatorLiquidStakingCap = initialLiquidStakingCap
	keeper.SetParams(ctx, params)

	// Attempt to increase the validator liquid shares, it should throw an
	// error that the validator bond cap was exceeded
	_, err := keeper.SafelyIncreaseValidatorLiquidShares(ctx, valAddress, firstIncreaseAmount, false)
	require.ErrorIs(err, types.ErrInsufficientValidatorBondShares)
	checkValidatorLiquidShares(initialLiquidShares, "shares after low bond factor")

	// Change validator bond factor to a more conservative number, so that the increase succeeds
	params.ValidatorBondFactor = finalBondFactor
	keeper.SetParams(ctx, params)

	// Try the increase again and check that it succeeded
	expectedLiquidSharesAfterFirstStake := initialLiquidShares.Add(firstIncreaseAmount)
	_, err = keeper.SafelyIncreaseValidatorLiquidShares(ctx, valAddress, firstIncreaseAmount, false)
	require.NoError(err)
	checkValidatorLiquidShares(expectedLiquidSharesAfterFirstStake, "shares with cap loose bond cap")

	// Attempt another increase, it should fail from the liquid staking cap
	_, err = keeper.SafelyIncreaseValidatorLiquidShares(ctx, valAddress, secondIncreaseAmount, false)
	require.ErrorIs(err, types.ErrValidatorLiquidStakingCapExceeded)
	checkValidatorLiquidShares(expectedLiquidSharesAfterFirstStake, "shares after liquid staking cap hit")

	// Raise the liquid staking cap so the new increment succeeds
	params.ValidatorLiquidStakingCap = finalLiquidStakingCap
	keeper.SetParams(ctx, params)

	// Finally confirm that the increase succeeded this time
	expectedLiquidSharesAfterSecondStake := expectedLiquidSharesAfterFirstStake.Add(secondIncreaseAmount)
	_, err = keeper.SafelyIncreaseValidatorLiquidShares(ctx, valAddress, secondIncreaseAmount, false)
	require.NoError(err, "no error expected after increasing liquid staking cap")
	checkValidatorLiquidShares(expectedLiquidSharesAfterSecondStake, "shares after loose liquid stake cap")
}

// Tests DecreaseValidatorLiquidShares
func (s *KeeperTestSuite) TestDecreaseValidatorLiquidShares() {
	ctx, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	initialLiquidShares := sdk.NewDec(100)
	decreaseAmount := sdk.NewDec(10)

	// Create a validator with designated self-bond shares
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	valAddress := sdk.ValAddress(pubKey.Address())

	initialValidator := types.Validator{
		OperatorAddress: valAddress.String(),
		LiquidShares:    initialLiquidShares,
	}
	keeper.SetValidator(ctx, initialValidator)

	// Decrease the validator liquid shares, and confirm the new share amount has been updated
	_, err := keeper.DecreaseValidatorLiquidShares(ctx, valAddress, decreaseAmount)
	require.NoError(err, "no error expected when decreasing validator liquid shares")

	actualValidator, found := keeper.GetValidator(ctx, valAddress)
	require.True(found)
	require.Equal(initialLiquidShares.Sub(decreaseAmount), actualValidator.LiquidShares, "liquid shares")

	// Attempt to decrease by a larger amount than it has, it should fail
	_, err = keeper.DecreaseValidatorLiquidShares(ctx, valAddress, initialLiquidShares)
	require.ErrorIs(err, types.ErrValidatorLiquidSharesUnderflow)
}

// Tests SafelyDecreaseValidatorBond
func (s *KeeperTestSuite) TestSafelyDecreaseValidatorBond() {
	ctx, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	// Initial Bond Factor: 100, Initial Validator Bond: 10
	// => Max Liquid Shares 1000 (Initial Liquid Shares: 200)
	initialBondFactor := sdk.NewDec(100)
	initialValidatorBondShares := sdk.NewDec(10)
	initialLiquidShares := sdk.NewDec(200)

	// Create a validator with designated self-bond shares
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	valAddress := sdk.ValAddress(pubKey.Address())

	initialValidator := types.Validator{
		OperatorAddress:     valAddress.String(),
		ValidatorBondShares: initialValidatorBondShares,
		LiquidShares:        initialLiquidShares,
	}
	keeper.SetValidator(ctx, initialValidator)

	// Set the bond factor
	params := keeper.GetParams(ctx)
	params.ValidatorBondFactor = initialBondFactor
	keeper.SetParams(ctx, params)

	// Decrease the validator bond from 10 to 5 (minus 5)
	// This will adjust the cap (factor * shares)
	// from (100 * 10 = 1000) to (100 * 5 = 500)
	// Since this is still above the initial liquid shares of 200, this will succeed
	decreaseAmount, expectedBondShares := sdk.NewDec(5), sdk.NewDec(5)
	err := keeper.SafelyDecreaseValidatorBond(ctx, valAddress, decreaseAmount)
	require.NoError(err)

	actualValidator, found := keeper.GetValidator(ctx, valAddress)
	require.True(found)
	require.Equal(expectedBondShares, actualValidator.ValidatorBondShares, "validator bond shares shares")

	// Now attempt to decrease the validator bond again from 5 to 1 (minus 4)
	// This time, the cap will be reduced to (factor * shares) = (100 * 1) = 100
	// However, the liquid shares are currently 200, so this should fail
	decreaseAmount, expectedBondShares = sdk.NewDec(4), sdk.NewDec(1)
	err = keeper.SafelyDecreaseValidatorBond(ctx, valAddress, decreaseAmount)
	require.ErrorIs(err, types.ErrInsufficientValidatorBondShares)

	// Finally, disable the cap and attempt to decrease again
	// This time it should succeed
	params.ValidatorBondFactor = types.ValidatorBondCapDisabled
	keeper.SetParams(ctx, params)

	err = keeper.SafelyDecreaseValidatorBond(ctx, valAddress, decreaseAmount)
	require.NoError(err)

	actualValidator, found = keeper.GetValidator(ctx, valAddress)
	require.True(found)
	require.Equal(expectedBondShares, actualValidator.ValidatorBondShares, "validator bond shares shares")
}

// Tests Add/Remove/Get/SetTokenizeSharesLock
func (s *KeeperTestSuite) TestTokenizeSharesLock() {
	ctx, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	addresses := simtestutil.CreateIncrementalAccounts(2)
	addressA, addressB := addresses[0], addresses[1]

	unlocked := types.TOKENIZE_SHARE_LOCK_STATUS_UNLOCKED.String()
	locked := types.TOKENIZE_SHARE_LOCK_STATUS_LOCKED.String()
	lockExpiring := types.TOKENIZE_SHARE_LOCK_STATUS_LOCK_EXPIRING.String()

	// Confirm both accounts start unlocked
	status, _ := keeper.GetTokenizeSharesLock(ctx, addressA)
	require.Equal(unlocked, status.String(), "addressA unlocked at start")

	status, _ = keeper.GetTokenizeSharesLock(ctx, addressB)
	require.Equal(unlocked, status.String(), "addressB unlocked at start")

	// Lock the first account
	keeper.AddTokenizeSharesLock(ctx, addressA)

	// The first account should now have tokenize shares disabled
	// and the unlock time should be the zero time
	status, _ = keeper.GetTokenizeSharesLock(ctx, addressA)
	require.Equal(locked, status.String(), "addressA locked")

	status, _ = keeper.GetTokenizeSharesLock(ctx, addressB)
	require.Equal(unlocked, status.String(), "addressB still unlocked")

	// Update the lock time and confirm it was set
	expectedUnlockTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	keeper.SetTokenizeSharesUnlockTime(ctx, addressA, expectedUnlockTime)

	status, actualUnlockTime := keeper.GetTokenizeSharesLock(ctx, addressA)
	require.Equal(lockExpiring, status.String(), "addressA lock expiring")
	require.Equal(expectedUnlockTime, actualUnlockTime, "addressA unlock time")

	// Confirm B is still unlocked
	status, _ = keeper.GetTokenizeSharesLock(ctx, addressB)
	require.Equal(unlocked, status.String(), "addressB still unlocked")

	// Remove the lock
	keeper.RemoveTokenizeSharesLock(ctx, addressA)
	status, _ = keeper.GetTokenizeSharesLock(ctx, addressA)
	require.Equal(unlocked, status.String(), "addressA unlocked at end")

	status, _ = keeper.GetTokenizeSharesLock(ctx, addressB)
	require.Equal(unlocked, status.String(), "addressB unlocked at end")
}

// Tests GetAllTokenizeSharesLocks
func (s *KeeperTestSuite) TestGetAllTokenizeSharesLocks() {
	ctx, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	addresses := simtestutil.CreateIncrementalAccounts(4)

	// Set 2 locked accounts, and two accounts with a lock expiring
	keeper.AddTokenizeSharesLock(ctx, addresses[0])
	keeper.AddTokenizeSharesLock(ctx, addresses[1])

	unlockTime1 := time.Date(2023, 1, 1, 1, 0, 0, 0, time.UTC)
	unlockTime2 := time.Date(2023, 1, 2, 1, 0, 0, 0, time.UTC)
	keeper.SetTokenizeSharesUnlockTime(ctx, addresses[2], unlockTime1)
	keeper.SetTokenizeSharesUnlockTime(ctx, addresses[3], unlockTime2)

	// Defined expected locks after GetAll
	expectedLocks := map[string]types.TokenizeShareLock{
		addresses[0].String(): {
			Status: types.TOKENIZE_SHARE_LOCK_STATUS_LOCKED.String(),
		},
		addresses[1].String(): {
			Status: types.TOKENIZE_SHARE_LOCK_STATUS_LOCKED.String(),
		},
		addresses[2].String(): {
			Status:         types.TOKENIZE_SHARE_LOCK_STATUS_LOCK_EXPIRING.String(),
			CompletionTime: unlockTime1,
		},
		addresses[3].String(): {
			Status:         types.TOKENIZE_SHARE_LOCK_STATUS_LOCK_EXPIRING.String(),
			CompletionTime: unlockTime2,
		},
	}

	// Check output from GetAll
	actualLocks := keeper.GetAllTokenizeSharesLocks(ctx)
	require.Len(actualLocks, len(expectedLocks), "number of locks")

	for i, actual := range actualLocks {
		expected, ok := expectedLocks[actual.Address]
		require.True(ok, "address %s not expected", actual.Address)
		require.Equal(expected.Status, actual.Status, "tokenize share lock #%d status", i)
		require.Equal(expected.CompletionTime, actual.CompletionTime, "tokenize share lock #%d completion time", i)
	}
}

// Test Get/SetPendingTokenizeShareAuthorizations
func (s *KeeperTestSuite) TestPendingTokenizeShareAuthorizations() {
	ctx, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	// Create dummy accounts and completion times

	addresses := simtestutil.CreateIncrementalAccounts(4)
	addressStrings := []string{}
	for _, address := range addresses {
		addressStrings = append(addressStrings, address.String())
	}

	timeA := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	timeB := timeA.Add(time.Hour)

	// There should be no addresses returned originally
	authorizationsA := keeper.GetPendingTokenizeShareAuthorizations(ctx, timeA)
	require.Empty(authorizationsA.Addresses, "no addresses at timeA expected")

	authorizationsB := keeper.GetPendingTokenizeShareAuthorizations(ctx, timeB)
	require.Empty(authorizationsB.Addresses, "no addresses at timeB expected")

	// Store addresses for timeB
	keeper.SetPendingTokenizeShareAuthorizations(ctx, timeB, types.PendingTokenizeShareAuthorizations{
		Addresses: addressStrings,
	})

	// Check addresses
	authorizationsA = keeper.GetPendingTokenizeShareAuthorizations(ctx, timeA)
	require.Empty(authorizationsA.Addresses, "no addresses at timeA expected at end")

	authorizationsB = keeper.GetPendingTokenizeShareAuthorizations(ctx, timeB)
	require.Equal(addressStrings, authorizationsB.Addresses, "address length")
}

// Test QueueTokenizeSharesAuthorization and RemoveExpiredTokenizeShareLocks
func (s *KeeperTestSuite) TestTokenizeShareAuthorizationQueue() {
	ctx, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	// Create dummy accounts and completion times

	// We'll start by adding the following addresses to the queue
	//   Time 0: [address0]
	//   Time 1: []
	//   Time 2: [address1, address2, address3]
	//   Time 3: [address4, address5]
	//   Time 4: [address6]
	addresses := simtestutil.CreateIncrementalAccounts(7)
	addressesByTime := map[int][]sdk.AccAddress{
		0: {addresses[0]},
		1: {},
		2: {addresses[1], addresses[2], addresses[3]},
		3: {addresses[4], addresses[5]},
		4: {addresses[6]},
	}

	// Set the unbonding time to 1 day
	unbondingPeriod := time.Hour * 24
	params := keeper.GetParams(ctx)
	params.UnbondingTime = unbondingPeriod
	keeper.SetParams(ctx, params)

	// Add each address to the queue and then increment the block time
	// such that the times line up as follows
	//   Time 0: 2023-01-01 00:00:00
	//   Time 1: 2023-01-01 00:01:00
	//   Time 2: 2023-01-01 00:02:00
	//   Time 3: 2023-01-01 00:03:00
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	ctx = ctx.WithBlockTime(startTime)
	blockTimeIncrement := time.Hour

	for timeIndex := 0; timeIndex <= 4; timeIndex++ {
		for _, address := range addressesByTime[timeIndex] {
			keeper.QueueTokenizeSharesAuthorization(ctx, address)
		}
		ctx = ctx.WithBlockTime(ctx.BlockTime().Add(blockTimeIncrement))
	}

	// We'll unlock the tokens using the following progression
	// The "alias'"/keys for these times assume a starting point of the Time 0
	// from above, plus the Unbonding Time
	//   Time -1  (2023-01-01 23:59:99): []
	//   Time  0  (2023-01-02 00:00:00): [address0]
	//   Time  1  (2023-01-02 00:01:00): []
	//   Time 2.5 (2023-01-02 00:02:30): [address1, address2, address3]
	//   Time 10  (2023-01-02 00:10:00): [address4, address5, address6]
	unlockBlockTimes := map[string]time.Time{
		"-1":  startTime.Add(unbondingPeriod).Add(-time.Second),
		"0":   startTime.Add(unbondingPeriod),
		"1":   startTime.Add(unbondingPeriod).Add(blockTimeIncrement),
		"2.5": startTime.Add(unbondingPeriod).Add(2 * blockTimeIncrement).Add(blockTimeIncrement / 2),
		"10":  startTime.Add(unbondingPeriod).Add(10 * blockTimeIncrement),
	}
	expectedUnlockedAddresses := map[string][]string{
		"-1":  {},
		"0":   {addresses[0].String()},
		"1":   {},
		"2.5": {addresses[1].String(), addresses[2].String(), addresses[3].String()},
		"10":  {addresses[4].String(), addresses[5].String(), addresses[6].String()},
	}

	// Now we'll remove items from the queue sequentially
	// First check with a block time before the first expiration - it should remove no addresses
	actualAddresses := keeper.RemoveExpiredTokenizeShareLocks(ctx, unlockBlockTimes["-1"])
	require.Equal(expectedUnlockedAddresses["-1"], actualAddresses, "no addresses unlocked from time -1")

	// Then pass in (time 0 + unbonding time) - it should remove the first address
	actualAddresses = keeper.RemoveExpiredTokenizeShareLocks(ctx, unlockBlockTimes["0"])
	require.Equal(expectedUnlockedAddresses["0"], actualAddresses, "one address unlocked from time 0")

	// Now pass in (time 1 + unbonding time) - it should remove no addresses since
	// the address at time 0 was already removed
	actualAddresses = keeper.RemoveExpiredTokenizeShareLocks(ctx, unlockBlockTimes["1"])
	require.Equal(expectedUnlockedAddresses["1"], actualAddresses, "no addresses unlocked from time 1")

	// Now pass in (time 2.5 + unbonding time) - it should remove the three addresses from time 2
	actualAddresses = keeper.RemoveExpiredTokenizeShareLocks(ctx, unlockBlockTimes["2.5"])
	require.Equal(expectedUnlockedAddresses["2.5"], actualAddresses, "addresses unlocked from time 2.5")

	// Finally pass in a block time far in the future, which should remove all the remaining locks
	actualAddresses = keeper.RemoveExpiredTokenizeShareLocks(ctx, unlockBlockTimes["10"])
	require.Equal(expectedUnlockedAddresses["10"], actualAddresses, "addresses unlocked from time 10")
}

// Tests DelegatorIsLiquidStaker
func (s *KeeperTestSuite) TestDelegatorIsLiquidStaker() {
	_, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	// Create base and ICA accounts
	baseAccountAddress := sdk.AccAddress("base-account")
	icaAccountAddress := sdk.AccAddress(
		address.Derive(authtypes.NewModuleAddress("icahost"), []byte("connection-0"+"icahost")),
	)

	// Only the ICA module account should be considered a liquid staking provider
	require.False(keeper.DelegatorIsLiquidStaker(baseAccountAddress), "base account")
	require.True(keeper.DelegatorIsLiquidStaker(icaAccountAddress), "ICA module account")
}

func TestCheckVestedDelegationInVestingAccount(t *testing.T) {

	var (
		vestingAcct     *vestingtypes.ContinuousVestingAccount
		startTime       = time.Now()
		endTime         = startTime.Add(24 * time.Hour)
		originalVesting = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100_000)))
	)

	testCases := []struct {
		name         string
		setupAcct    func()
		blockTime    time.Time
		coinRequired sdk.Coin
		expRes       bool
	}{
		{
			name:         "vesting account has zero delegations",
			setupAcct:    func() {},
			blockTime:    endTime,
			coinRequired: sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt()),
			expRes:       false,
		},
		{
			name: "vested delegations exist but for a different coin",
			setupAcct: func() {
				vestingAcct.DelegatedFree = sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(100_000)))
			},
			blockTime:    endTime,
			coinRequired: sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt()),
			expRes:       false,
		},
		{
			name: "all delegations are vesting",
			setupAcct: func() {
				vestingAcct.DelegatedVesting = vestingAcct.OriginalVesting
			},
			blockTime:    startTime,
			coinRequired: sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt()),
			expRes:       false,
		},
		{
			name: "not enough vested coin",
			setupAcct: func() {
				vestingAcct.DelegatedFree = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(80_000)))
			},
			blockTime:    endTime,
			coinRequired: sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100_000)),
			expRes:       false,
		},
		{
			name: "account is vested and have vested delegations",
			setupAcct: func() {
				vestingAcct.DelegatedFree = vestingAcct.OriginalVesting
			},
			blockTime:    endTime,
			coinRequired: sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100_000)),
			expRes:       true,
		},
		{
			name: "vesting account partially vested and have vesting and vested delegations",
			setupAcct: func() {
				vestingAcct.DelegatedFree = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(50_000)))
				vestingAcct.DelegatedVesting = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(50_000)))
			},
			blockTime:    startTime.Add(18 * time.Hour), // vest 3/4 vesting period
			coinRequired: sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(75_000)),

			expRes: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			baseAcc := authtypes.NewBaseAccount(sdk.AccAddress([]byte("addr")), secp256k1.GenPrivKey().PubKey(), 0, 0)
			vestingAcct = vestingtypes.NewContinuousVestingAccount(
				baseAcc,
				originalVesting,
				startTime.Unix(),
				endTime.Unix(),
			)

			tc.setupAcct()

			require.Equal(
				t,
				tc.expRes, keeper.CheckVestedDelegationInVestingAccount(
					vestingAcct,
					tc.blockTime,
					tc.coinRequired,
				),
			)
		})
	}
}