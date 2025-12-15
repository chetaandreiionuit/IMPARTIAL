package temporal

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/yourorg/truthweave/internal/domain/article"
	"go.temporal.io/sdk/testsuite"
)

// [RO] Suita de Teste Temporal
type WorkflowTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

func (s *WorkflowTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *WorkflowTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

// [RO] Test: Scenariul Fericit (Happy Path)
func (s *WorkflowTestSuite) TestOrchestrateNewsAnalysisWorkflow_Success() {
	// 1. Mock Activities
	activities := &NewsProcessingActivities{}

	// Înregistăm activitățile (chiar dacă sunt metode pe struct, le înregistrăm pointer-based)
	// Nota: Temporal TestEnv execută codul real al activităților dacă nu le facem Mock.
	// Aici VREM sa facem Mock la activități ca să testăm doar logica de workflow (orchestrarea).

	s.env.OnActivity(activities.ExtractWebPageContentActivity, mock.Anything, "http://test.com").Return("Raw Content", nil)
	s.env.OnActivity(activities.GenerateSemanticVectorActivity, mock.Anything, "Raw Content").Return([]float32{0.1, 0.2}, nil)

	// Simulăm că NU e duplicat (Score mic)
	simResult := &SimilarityCheckResult{ExistingArticle: nil, SimilarityScore: 0.1}
	s.env.OnActivity(activities.CheckForExistingDuplicatesActivity, mock.Anything, []float32{0.1, 0.2}).Return(simResult, nil)

	// Analiza AI
	aiResult := &article.AIAnalysisResult{
		RewrittenText: "Neutral Text",
		Score:         85.5,
		Location:      article.GaiaPoint{Latitude: 10, Longitude: 20},
	}
	s.env.OnActivity(activities.AnalyzeNewsContentActivity, mock.Anything, "Raw Content").Return(aiResult, nil)

	// Salvare DB și Graph
	s.env.OnActivity(activities.PersistAnalysisToDatabaseActivity, mock.Anything, mock.Anything).Return(nil)
	s.env.OnActivity(activities.ConnectKnowledgeGraphActivity, mock.Anything, mock.Anything).Return(nil)

	// 2. Execuție Workflow
	s.env.ExecuteWorkflow(OrchestrateNewsAnalysisWorkflow, "http://test.com")

	// 3. Verificare
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

// [RO] Test: Scenariul Duplicat (Deduplicare)
func (s *WorkflowTestSuite) TestOrchestrateNewsAnalysisWorkflow_DuplicateDetected() {
	activities := &NewsProcessingActivities{}

	s.env.OnActivity(activities.ExtractWebPageContentActivity, mock.Anything, "http://duplicate.com").Return("Duplicate Content", nil)
	s.env.OnActivity(activities.GenerateSemanticVectorActivity, mock.Anything, "Duplicate Content").Return([]float32{0.9, 0.9}, nil)

	// Simulăm că ESTE duplicat (Score mare > 0.90)
	existing := &article.NewsArticleEntity{ID: uuid.New()}
	simResult := &SimilarityCheckResult{ExistingArticle: existing, SimilarityScore: 0.98}
	s.env.OnActivity(activities.CheckForExistingDuplicatesActivity, mock.Anything, []float32{0.9, 0.9}).Return(simResult, nil)

	// Așteptare: Workflow-ul se oprește AICI. NU apelează Analyze, nici Save.
	// Dacă ar apela, testul ar crăpa cu "Unexpected call".

	s.env.ExecuteWorkflow(OrchestrateNewsAnalysisWorkflow, "http://duplicate.com")

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
	// Implicit: AssertExpectations verifică că nu s-au apelat alte activități.
}

func TestWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowTestSuite))
}
