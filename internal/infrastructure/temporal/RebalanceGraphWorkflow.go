package temporal

import (
	"context"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/yourorg/truthweave/internal/infrastructure/gemini"
)

type RebalanceInput struct {
	TargetEventID   string
	CandidateEvents []gemini.PotentialCause
}

// [RO] Activitate: Analiză Cauzală Retroactivă
func (activities *NewsProcessingActivities) CalculateCausalityActivity(ctx context.Context, input RebalanceInput) (*gemini.CausalityAnalysisResult, error) {
	// Obținem summary-ul target event-ului (Simulat, ar trebui un DB fetch)
	targetSummary := "Event " + input.TargetEventID + " summary placeholder."

	return activities.ArtificialIntelligence.DetermineCausality(ctx, targetSummary, input.CandidateEvents)
}

// [RO] Activitate: Actualizare Graf
type GraphMutationParams struct {
	ChildID         string
	CausalityResult gemini.CausalityAnalysisResult
}

// [RO] Activitate: Actualizare Graf
func (activities *NewsProcessingActivities) ApplyGraphMutationsActivity(ctx context.Context, params GraphMutationParams) error {
	if params.CausalityResult.IsConsequence {
		return activities.KnowledgeGraph.CreateCausalEdge(ctx, params.CausalityResult.ParentEventID, params.ChildID, params.CausalityResult.RelationshipType)
	}
	return nil
}

// [RO] Workflow: Rebalansare Graf (Retroactive Causality)
// Acest workflow este declanșat atunci când o știre nouă are potențialul de a explica evenimente trecute,
// sau invers, când vrem să legăm o știre nouă de un context istoric (7 zile).
func RebalanceGraphWorkflow(ctx workflow.Context, input RebalanceInput) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 2,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval: time.Second,
			MaximumAttempts: 3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	logger := workflow.GetLogger(ctx)

	var tools *NewsProcessingActivities

	// 1. Determină Cauzalitatea folosind AI
	var causalityResult gemini.CausalityAnalysisResult
	if err := workflow.ExecuteActivity(ctx, tools.CalculateCausalityActivity, input).Get(ctx, &causalityResult); err != nil {
		logger.Error("Eșec analiză cauzalitate", "Error", err)
		return err
	}

	// 2. Dacă există o legătură, actualizează graful
	if causalityResult.IsConsequence {
		logger.Info("Legătură cauzală găsită", "Parent", causalityResult.ParentEventID, "Type", causalityResult.RelationshipType)

		mutationParams := GraphMutationParams{
			ChildID:         input.TargetEventID,
			CausalityResult: causalityResult,
		}

		if err := workflow.ExecuteActivity(ctx, tools.ApplyGraphMutationsActivity, mutationParams).Get(ctx, nil); err != nil {
			return err
		}
	} else {
		logger.Info("Nicio legătură cauzală detectată.")
	}

	return nil
}
