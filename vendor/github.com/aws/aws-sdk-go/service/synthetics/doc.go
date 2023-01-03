// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

// Package synthetics provides the client and types for making API
// requests to Synthetics.
//
// You can use Amazon CloudWatch Synthetics to continually monitor your services.
// You can create and manage canaries, which are modular, lightweight scripts
// that monitor your endpoints and APIs from the outside-in. You can set up
// your canaries to run 24 hours a day, once per minute. The canaries help you
// check the availability and latency of your web services and troubleshoot
// anomalies by investigating load time data, screenshots of the UI, logs, and
// metrics. The canaries seamlessly integrate with CloudWatch ServiceLens to
// help you trace the causes of impacted nodes in your applications. For more
// information, see Using ServiceLens to Monitor the Health of Your Applications
// (https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/ServiceLens.html)
// in the Amazon CloudWatch User Guide.
//
// Before you create and manage canaries, be aware of the security considerations.
// For more information, see Security Considerations for Synthetics Canaries
// (https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/servicelens_canaries_security.html).
//
// See https://docs.aws.amazon.com/goto/WebAPI/synthetics-2017-10-11 for more information on this service.
//
// See synthetics package documentation for more information.
// https://docs.aws.amazon.com/sdk-for-go/api/service/synthetics/
//
// # Using the Client
//
// To contact Synthetics with the SDK use the New function to create
// a new service client. With that client you can make API requests to the service.
// These clients are safe to use concurrently.
//
// See the SDK's documentation for more information on how to use the SDK.
// https://docs.aws.amazon.com/sdk-for-go/api/
//
// See aws.Config documentation for more information on configuring SDK clients.
// https://docs.aws.amazon.com/sdk-for-go/api/aws/#Config
//
// See the Synthetics client Synthetics for more
// information on creating client for this service.
// https://docs.aws.amazon.com/sdk-for-go/api/service/synthetics/#New
package synthetics
