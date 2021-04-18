package cmd

const (
	// ExpendGasOperation is an operation that only affects the native currency.
	ExpendGasOperation = "Gas"

	// ContractCreationOperation is an operation that only affects the native currency.
	ContractCreationOperation = "ContractCreation"

	// NativeTransferOperation is an operation that only affects the native currency.
	NativeTransferOperation = "NativeTransfer"

	// NativeCrossShardTransferOperation is an operation that only affects the native currency.
	NativeCrossShardTransferOperation = "NativeCrossShardTransfer"

	// GenesisFundsOperation is a side effect operation for genesis block only.
	// Note that no transaction can be constructed with this operation.
	GenesisFundsOperation = "Genesis"

	// PreStakingBlockRewardOperation is a side effect operation for pre-staking era only.
	// Note that no transaction can be constructed with this operation.
	PreStakingBlockRewardOperation = "PreStakingBlockReward"

	// UndelegationPayoutOperation is a side effect operation for committee election block only.
	// Note that no transaction can be constructed with this operation.
	UndelegationPayoutOperation = "UndelegationPayout"
)
