package internal

import (
	"fmt"

	s "go.uber.org/cadence/.gen/go/shared"
	"go.uber.org/cadence/internal/common"
	"go.uber.org/cadence/internal/common/backoff"
)

func convertRetryPolicy(retryPolicy *RetryPolicy) *s.RetryPolicy {
	if retryPolicy == nil {
		return nil
	}
	thriftRetryPolicy := s.RetryPolicy{
		InitialIntervalInSeconds:    common.Int32Ptr(common.Int32Ceil(retryPolicy.InitialInterval.Seconds())),
		MaximumIntervalInSeconds:    common.Int32Ptr(common.Int32Ceil(retryPolicy.MaximumInterval.Seconds())),
		BackoffCoefficient:          &retryPolicy.BackoffCoefficient,
		MaximumAttempts:             &retryPolicy.MaximumAttempts,
		NonRetriableErrorReasons:    retryPolicy.NonRetriableErrorReasons,
		ExpirationIntervalInSeconds: common.Int32Ptr(common.Int32Ceil(retryPolicy.ExpirationInterval.Seconds())),
	}
	if *thriftRetryPolicy.BackoffCoefficient == 0 {
		thriftRetryPolicy.BackoffCoefficient = common.Float64Ptr(backoff.DefaultBackoffCoefficient)
	}
	return &thriftRetryPolicy
}

func convertActiveClusterSelectionPolicy(policy *ActiveClusterSelectionPolicy) (*s.ActiveClusterSelectionPolicy, error) {
	if policy == nil {
		return nil, nil
	}

	switch policy.Strategy {
	case ActiveClusterSelectionStrategyRegionSticky:
		return &s.ActiveClusterSelectionPolicy{
			Strategy: s.ActiveClusterSelectionStrategyRegionSticky.Ptr(),
		}, nil
	case ActiveClusterSelectionStrategyExternalEntity:
		if policy.ExternalEntityType == "" {
			return nil, fmt.Errorf("external entity type is required for external entity strategy")
		}
		if policy.ExternalEntityKey == "" {
			return nil, fmt.Errorf("external entity key is required for external entity strategy")
		}
		return &s.ActiveClusterSelectionPolicy{
			Strategy:           s.ActiveClusterSelectionStrategyExternalEntity.Ptr(),
			ExternalEntityType: common.StringPtr(policy.ExternalEntityType),
			ExternalEntityKey:  common.StringPtr(policy.ExternalEntityKey),
		}, nil
	default:
		return nil, fmt.Errorf("invalid active cluster selection strategy: %d", policy.Strategy)
	}
}
