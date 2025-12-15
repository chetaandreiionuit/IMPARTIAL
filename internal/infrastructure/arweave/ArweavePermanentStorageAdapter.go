package arweave

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
)

// [RO] Adaptor pentru Stocare Permanentă (Arweave)
//
// Arweave este un blockchain dedicat stocării datelor pe termen nelimitat ("Permaweb").
// Acest modul ne permite să salvăm dosarul articolului astfel încât să nu poată fi șters niciodată.
type ArweavePermanentStorageAdapter struct {
	client *goar.Client
	wallet *goar.Wallet
}

// [RO] Constructor Arweave
func NewArweavePermanentStorageAdapter(keyPath string, gateway string) (*ArweavePermanentStorageAdapter, error) {
	var wallet *goar.Wallet
	var err error

	if keyPath != "" {
		wallet, err = goar.NewWalletFromPath(keyPath, gateway)
		if err != nil {
			return nil, fmt.Errorf("[RO] Eroare: Nu am putut încărca portofelul Arweave: %w", err)
		}
	}

	client := goar.NewClient(gateway)

	return &ArweavePermanentStorageAdapter{
		client: client,
		wallet: wallet,
	}, nil
}

// [RO] Stochează Conținut Permanent
// Trimite datele către Arweave și returnează ID-ul tranzacției.
func (adapter *ArweavePermanentStorageAdapter) StorePermanentContent(executionContext context.Context, data []byte) (string, error) {
	if adapter.wallet == nil {
		return "", fmt.Errorf("[RO] Eroare: Portofelul nu este configurat. Nu putem semna tranzacția.")
	}

	// [RO] Etichetare (Metadata)
	// Adăugăm etichete pentru a putea găsi ușor datele mai târziu.
	tags := []types.Tag{
		{Name: "App-Name", Value: "TruthWeave"},
		{Name: "Content-Type", Value: "application/json"},
		{Name: "Timestamp", Value: strconv.FormatInt(time.Now().Unix(), 10)},
	}

	// [RO] Trimitere Date
	// Această funcție creează tranzacția, o semnează și o trimite la rețea.
	tx, err := adapter.wallet.SendData(data, tags)
	if err != nil {
		return "", fmt.Errorf("failed to send data to arweave: %w", err)
	}

	return tx.ID, nil
}

// [RO] Ancorează Hash (Stub)
// Metoda nu este suportată de Arweave (e treaba Solanei), dar trebuie să existe pentru interfață.
func (adapter *ArweavePermanentStorageAdapter) AnchorContentHash(executionContext context.Context, hash string) (string, error) {
	return "", fmt.Errorf("[RO] Eroare: Arweave nu suportă ancorarea de hash-uri scurte. Folosiți Solana.")
}
