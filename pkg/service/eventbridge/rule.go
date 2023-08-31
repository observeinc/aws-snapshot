package eventbridge

import (
	"context"
	"fmt"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/eventbridge"
)

// ListRuleOutput combines targets as part of the rules, since we want to avoid
// having to scrape rules twice in order to get targets
type ListRuleOutput struct {
	Rule    *eventbridge.Rule     `json:"rule"`
	Targets []*eventbridge.Target `json:"targets"`
}

func (o *ListRuleOutput) Records() (records []*api.Record) {
	records = append(records, &api.Record{
		ID:   o.Rule.Arn,
		Data: o,
	})
	return
}

type ListRules struct {
	API
}

var _ api.RequestBuilder = &ListRules{}

// New implements api.RequestBuilder
// NOTE: events:ListRules
func (fn *ListRules) New(name string, config interface{}) ([]api.Request, error) {
	var listRulesInput eventbridge.ListRulesInput
	if err := api.DecodeConfig(config, &listRulesInput); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var countRules int
		r, _ := ctx.Value("runner_config").(api.Runner)
		for {
			listRulesOutput, err := fn.ListRulesWithContext(ctx, &listRulesInput)
			if err != nil {
				return fmt.Errorf("failed to list rules: %w", err)
			}

			if r.Stats {
				countRules += len(listRulesOutput.Rules)
			} else {
				for _, rule := range listRulesOutput.Rules {
					var targets []*eventbridge.Target

					listTargetsByRuleInput := eventbridge.ListTargetsByRuleInput{Rule: rule.Name}
					for {
						listTargetsByRuleOutput, err := fn.ListTargetsByRuleWithContext(ctx, &listTargetsByRuleInput)
						if err != nil {
							return fmt.Errorf("failed to list targets for rule %s: %w", *rule.Name, err)
						}

						targets = append(targets, listTargetsByRuleOutput.Targets...)
						if listTargetsByRuleOutput.NextToken == nil {
							break
						}
						listTargetsByRuleInput.NextToken = listTargetsByRuleOutput.NextToken
					}

					if err := api.SendRecords(ctx, ch, name, &ListRuleOutput{
						Rule:    rule,
						Targets: targets,
					}); err != nil {
						return err
					}
				}
			}
			if listRulesOutput.NextToken == nil {
				break
			}
			listRulesInput.NextToken = listRulesOutput.NextToken
		}
		if r.Stats {
			err := api.SendRecords(ctx, ch, name, &api.CountRecords{Count: countRules})
			if err != nil {
				return err
			}
		}
		return nil
	}

	return []api.Request{call}, nil
}
