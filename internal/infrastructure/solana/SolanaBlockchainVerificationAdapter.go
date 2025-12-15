package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	compute_budget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/memo"
	"github.com/gagliardetto/solana-go/rpc"
)

// [RO] Adaptor de Verificare Blockchain (Solana)
//
// Acest modul se comportă ca un "Notar Public Digital".
// Rolul său este să ia amprenta unică a unui articol (Hash-ul) și să o scrie
// într-o tranzacție pe Blockchain-ul Solana.
// Aceasta oferă dovada imutabilă că articolul exista la acea dată în acea formă.
type SolanaBlockchainVerificationAdapter struct {
	privateKey solana.PrivateKey
	rpcClient  *rpc.Client
}

// [RO] Constructor pentru Adaptorul Solana
func NewSolanaBlockchainVerificationAdapter(privateKeyBase58 string, rpcURL string) (*SolanaBlockchainVerificationAdapter, error) {
	privKey, err := solana.PrivateKeyFromBase58(privateKeyBase58)
	if err != nil {
		return nil, fmt.Errorf("[RO] Cheie privată invalidă: %w", err)
	}

	client := rpc.New(rpcURL)

	return &SolanaBlockchainVerificationAdapter{
		privateKey: privKey,
		rpcClient:  client,
	}, nil
}

// [RO] Ancorează Hash-ul (Semnătura Digitală)
// Trimite o tranzacție pe Solana care conține un mesaj text (Memo) cu hash-ul știrii.
// Returnează ID-ul tranzacției (Signature) ca dovadă.
func (adapter *SolanaBlockchainVerificationAdapter) AnchorContentHash(executionContext context.Context, contentHash string) (string, error) {
	// [RO] 1. Setăm Bugetul de Procesare (Priority Fee)
	// Plătim o mică taxă extra (în MicroLamports) pentru ca tranzacția să fie procesată rapid
	// de către validatori.
	limit := uint32(200_000)      // Limita de resurse de calcul
	microLamports := uint64(1000) // Taxa de prioritate (bacșișul pentru rețea)

	computeBudgetLimit := compute_budget.NewSetComputeUnitLimitInstruction(limit).Build()
	computeBudgetPrice := compute_budget.NewSetComputeUnitPriceInstruction(microLamports).Build()

	// [RO] 2. Instrucțiunea Memo
	// Aici scriem efectiv informația pe blockchain: "TruthWeave Verif: <HASH>"
	memoContent := fmt.Sprintf("TruthWeave Verif: %s", contentHash)

	memoInstruction := memo.NewMemoInstruction(
		[]byte(memoContent),
		adapter.privateKey.PublicKey(),
	).Build()

	// [RO] 3. Obținem cel mai recent Blockhash
	// Necesar pentru a preveni reluarea tranzacțiilor vechi.
	recentBlock, err := adapter.rpcClient.GetLatestBlockhash(executionContext, rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("failed to get latest blockhash: %w", err)
	}

	// [RO] 4. Creăm Tranzacția
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			computeBudgetLimit,
			computeBudgetPrice,
			memoInstruction,
		},
		recentBlock.Value.Blockhash,
		solana.TransactionPayer(adapter.privateKey.PublicKey()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %w", err)
	}

	// [RO] 5. Semnăm Tranzacția
	// Folosim cheia noastră privată pentru a autoriza plata taxelor si mesajul.
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if adapter.privateKey.PublicKey().Equals(key) {
				return &adapter.privateKey
			}
			return nil
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// [RO] 6. Trimitem către Rețea
	signature, err := adapter.rpcClient.SendTransactionWithOpts(
		executionContext,
		tx,
		rpc.TransactionOpts{
			SkipPreflight: false,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	return signature.String(), nil
}
