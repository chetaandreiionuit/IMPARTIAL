package ports

import (
	"context"

	"github.com/yourorg/truthweave/internal/domain/article"
	"go.temporal.io/sdk/client"
)

// [RO] Interfațe pentru Sisteme Externe
// Aceste contracte definesc cum interacționăm cu lumea exterioară.

// [RO] Poarta către Inteligența Artificială (Oracolul)
type ArtificialIntelligenceGateway interface {
	AnalyzeAndNeutralizeNewsContent(ctx context.Context, rawContent string) (*article.AIAnalysisResult, error)
	GenerateSemanticVector(ctx context.Context, text string) ([]float32, error)
	ChatWithContext(ctx context.Context, query string, context string) (string, error)
}

// [RO] Poarta către Blockchain (Notarul Digital)
type BlockchainGateway interface {
	StorePermanentContent(ctx context.Context, data []byte) (string, error) // Arweave
	AnchorContentHash(ctx context.Context, hash string) (string, error)     // Solana/Ethereum
}

// [RO] Poarta către Orchestrare (Temporal)
type WorkflowOrchestratorLauncher interface {
	ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error)
}

// [RO] Poarta către Știri Globale (NewsAPI)
type GlobalNewsAggregator interface {
	FetchGlobalHeadlines(ctx context.Context, query string) ([]string, error)
}
