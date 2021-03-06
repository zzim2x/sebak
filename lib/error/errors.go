package sebakerror

var (
	ErrorBlockAlreadyExists               = NewError(100, "already exists in block")
	ErrorHashDoesNotMatch                 = NewError(101, "`Hash` does not match")
	ErrorSignatureVerificationFailed      = NewError(102, "signature verification failed")
	ErrorBadPublicAddress                 = NewError(103, "failed to parse public address")
	ErrorInvalidFee                       = NewError(104, "invalid fee")
	ErrorInvalidOperation                 = NewError(105, "invalid operation")
	ErrorNewButKnownMessage               = NewError(106, "received new, but known message")
	ErrorInvalidState                     = NewError(107, "found invalid state")
	ErrorInvalidVotingThresholdPolicy     = NewError(108, "invalid `VotingThresholdPolicy`")
	ErrorBallotEmptyMessage               = NewError(109, "init state ballot does not have `Message`")
	ErrorInvalidHash                      = NewError(110, "invalid `Hash`")
	ErrorInvalidMessage                   = NewError(111, "invalid `Message`")
	ErrorBallotHasMessage                 = NewError(112, "none-init state ballot must not have `Message`")
	ErrorVotingResultAlreadyExists        = NewError(113, "`VotingResult` already exists")
	ErrorVotingResultNotFound             = NewError(114, "`VotingResult` not found")
	ErrorVotingResultFailedToSetState     = NewError(115, "failed to set the new state to `VotingResult`")
	ErrorVotingResultNotInBox             = NewError(116, "ballot is not in here")
	ErrorBallotNoVoting                   = NewError(118, "ballot has no `Voting`")
	ErrorBallotNoNodeKey                  = NewError(119, "ballot has no `NodeKey`")
	ErrorVotingThresholdInvalidValidators = NewError(120, "invalid validators")
	ErrorBallotHasInvalidState            = NewError(121, "ballot has invalid state")
	ErrorVotingResultFailedToClose        = NewError(122, "failed to close `VotingResult`")
	ErrorTransactionEmptyOperations       = NewError(123, "operations needs in transaction")
	ErrorAlreadySaved                     = NewError(124, "already saved")
	ErrorDuplicatedOperation              = NewError(125, "duplicated operations in transaction")
	ErrorUnknownOperationType             = NewError(126, "unknown operation type")
	ErrorTypeOperationBodyNotMatched      = NewError(127, "operation type and it's type does not match")
	ErrorBlockAccountDoesNotExists        = NewError(128, "account does not exists in block")
	ErrorBlockAccountAlreadyExists        = NewError(129, "account already exists in block")
	ErrorAccountBalanceUnderZero          = NewError(130, "account balance will be under zero")
	ErrorMaximumBalanceReached            = NewError(131, "monetary amount would be greater than the total supply of coins")
	ErrorStorageRecordDoesNotExist        = NewError(132, "record does not exist in storage")
	ErrorTransactionInvalidCheckpoint     = NewError(133, "invalid checkpoint found")
	ErrorBlockTransactionDoesNotExists    = NewError(134, "transaction does not exists in block")
	ErrorBlockOperationDoesNotExists      = NewError(135, "operation does not exists in block")
)
