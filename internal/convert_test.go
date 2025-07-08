package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	s "go.uber.org/cadence/.gen/go/shared"
	"go.uber.org/cadence/internal/common"
	"go.uber.org/cadence/internal/common/backoff"
)

func TestConvertRetryPolicy(t *testing.T) {
	tests := []struct {
		name         string
		retryPolicy  *RetryPolicy
		thriftPolicy *s.RetryPolicy
	}{
		{
			name:         "nil retry policy - return nil",
			retryPolicy:  nil,
			thriftPolicy: nil,
		},
		{
			name: "non-zero backoff coefficient - use provided",
			retryPolicy: &RetryPolicy{
				InitialInterval:          1 * time.Second,
				MaximumInterval:          10 * time.Second,
				BackoffCoefficient:       2.0,
				MaximumAttempts:          3,
				NonRetriableErrorReasons: []string{"error1", "error2"},
				ExpirationInterval:       10 * time.Second,
			},
			thriftPolicy: &s.RetryPolicy{
				InitialIntervalInSeconds:    common.Int32Ptr(1),
				MaximumIntervalInSeconds:    common.Int32Ptr(10),
				BackoffCoefficient:          common.Float64Ptr(2.0),
				MaximumAttempts:             common.Int32Ptr(3),
				NonRetriableErrorReasons:    []string{"error1", "error2"},
				ExpirationIntervalInSeconds: common.Int32Ptr(10),
			},
		},
		{
			name: "zero backoff coefficient - use default",
			retryPolicy: &RetryPolicy{
				InitialInterval:          1 * time.Second,
				MaximumInterval:          10 * time.Second,
				BackoffCoefficient:       0.0,
				MaximumAttempts:          3,
				NonRetriableErrorReasons: []string{"error1", "error2"},
				ExpirationInterval:       10 * time.Second,
			},
			thriftPolicy: &s.RetryPolicy{
				InitialIntervalInSeconds:    common.Int32Ptr(1),
				MaximumIntervalInSeconds:    common.Int32Ptr(10),
				BackoffCoefficient:          common.Float64Ptr(backoff.DefaultBackoffCoefficient),
				MaximumAttempts:             common.Int32Ptr(3),
				NonRetriableErrorReasons:    []string{"error1", "error2"},
				ExpirationIntervalInSeconds: common.Int32Ptr(10),
			},
		},
	}

	for _, test := range tests {
		thriftPolicy := convertRetryPolicy(test.retryPolicy)
		assert.Equal(t, test.thriftPolicy, thriftPolicy)
	}
}

func TestConvertActiveClusterSelectionPolicy(t *testing.T) {
	tests := []struct {
		name         string
		policy       *ActiveClusterSelectionPolicy
		thriftPolicy *s.ActiveClusterSelectionPolicy
		wantErr      bool
	}{
		{
			name:         "nil policy - return nil",
			policy:       nil,
			thriftPolicy: nil,
		},
		{
			name: "region sticky policy",
			policy: &ActiveClusterSelectionPolicy{
				Strategy: ActiveClusterSelectionStrategyRegionSticky,
			},
			thriftPolicy: &s.ActiveClusterSelectionPolicy{
				Strategy: s.ActiveClusterSelectionStrategyRegionSticky.Ptr(),
			},
		},
		{
			name: "external entity policy - success",
			policy: &ActiveClusterSelectionPolicy{
				Strategy:           ActiveClusterSelectionStrategyExternalEntity,
				ExternalEntityType: "test-type",
				ExternalEntityKey:  "test-key",
			},
			thriftPolicy: &s.ActiveClusterSelectionPolicy{
				Strategy:           s.ActiveClusterSelectionStrategyExternalEntity.Ptr(),
				ExternalEntityType: common.StringPtr("test-type"),
				ExternalEntityKey:  common.StringPtr("test-key"),
			},
		},
		{
			name: "external entity policy - missing type",
			policy: &ActiveClusterSelectionPolicy{
				Strategy:          ActiveClusterSelectionStrategyExternalEntity,
				ExternalEntityKey: "test-key",
			},
			wantErr: true,
		},
		{
			name: "external entity policy - missing key",
			policy: &ActiveClusterSelectionPolicy{
				Strategy:           ActiveClusterSelectionStrategyExternalEntity,
				ExternalEntityType: "test-type",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		thriftPolicy, err := convertActiveClusterSelectionPolicy(test.policy)
		if test.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.thriftPolicy, thriftPolicy)
		}
	}
}
