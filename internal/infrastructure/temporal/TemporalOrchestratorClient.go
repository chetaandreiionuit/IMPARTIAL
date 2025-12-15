package temporal

import (
	"context"

	"go.temporal.io/sdk/client"
)

// [RO] Client Orchestrator Temporal
//
// Această componentă este "Telecomanda" prin care pornim procesele complexe în clusterul Temporal.
// Implementează interfața `WorkflowOrchestratorLauncher`.
type TemporalOrchestratorClient struct {
	client client.Client
}

// [RO] Constructor Client Temporal
func NewTemporalOrchestratorClient(c client.Client) *TemporalOrchestratorClient {
	return &TemporalOrchestratorClient{client: c}
}

// [RO] Execută Flux de Lucru
// Trimite semnalul de start către Temporal Server.
func (t *TemporalOrchestratorClient) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
	return t.client.ExecuteWorkflow(ctx, options, workflow, args...)
}
