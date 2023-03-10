// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package autoscaling

const (

	// ErrCodeActiveInstanceRefreshNotFoundFault for service response error code
	// "ActiveInstanceRefreshNotFound".
	//
	// The request failed because an active instance refresh or rollback for the
	// specified Auto Scaling group was not found.
	ErrCodeActiveInstanceRefreshNotFoundFault = "ActiveInstanceRefreshNotFound"

	// ErrCodeAlreadyExistsFault for service response error code
	// "AlreadyExists".
	//
	// You already have an Auto Scaling group or launch configuration with this
	// name.
	ErrCodeAlreadyExistsFault = "AlreadyExists"

	// ErrCodeInstanceRefreshInProgressFault for service response error code
	// "InstanceRefreshInProgress".
	//
	// The request failed because an active instance refresh already exists for
	// the specified Auto Scaling group.
	ErrCodeInstanceRefreshInProgressFault = "InstanceRefreshInProgress"

	// ErrCodeInvalidNextToken for service response error code
	// "InvalidNextToken".
	//
	// The NextToken value is not valid.
	ErrCodeInvalidNextToken = "InvalidNextToken"

	// ErrCodeIrreversibleInstanceRefreshFault for service response error code
	// "IrreversibleInstanceRefresh".
	//
	// The request failed because a desired configuration was not found or an incompatible
	// launch template (uses a Systems Manager parameter instead of an AMI ID) or
	// launch template version ($Latest or $Default) is present on the Auto Scaling
	// group.
	ErrCodeIrreversibleInstanceRefreshFault = "IrreversibleInstanceRefresh"

	// ErrCodeLimitExceededFault for service response error code
	// "LimitExceeded".
	//
	// You have already reached a limit for your Amazon EC2 Auto Scaling resources
	// (for example, Auto Scaling groups, launch configurations, or lifecycle hooks).
	// For more information, see DescribeAccountLimits (https://docs.aws.amazon.com/autoscaling/ec2/APIReference/API_DescribeAccountLimits.html)
	// in the Amazon EC2 Auto Scaling API Reference.
	ErrCodeLimitExceededFault = "LimitExceeded"

	// ErrCodeResourceContentionFault for service response error code
	// "ResourceContention".
	//
	// You already have a pending update to an Amazon EC2 Auto Scaling resource
	// (for example, an Auto Scaling group, instance, or load balancer).
	ErrCodeResourceContentionFault = "ResourceContention"

	// ErrCodeResourceInUseFault for service response error code
	// "ResourceInUse".
	//
	// The operation can't be performed because the resource is in use.
	ErrCodeResourceInUseFault = "ResourceInUse"

	// ErrCodeScalingActivityInProgressFault for service response error code
	// "ScalingActivityInProgress".
	//
	// The operation can't be performed because there are scaling activities in
	// progress.
	ErrCodeScalingActivityInProgressFault = "ScalingActivityInProgress"

	// ErrCodeServiceLinkedRoleFailure for service response error code
	// "ServiceLinkedRoleFailure".
	//
	// The service-linked role is not yet ready for use.
	ErrCodeServiceLinkedRoleFailure = "ServiceLinkedRoleFailure"
)
