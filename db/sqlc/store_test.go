package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xiao-yangg/simplebank/util"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// run n concurrent transfer transactions (using Go routines)
	n := util.RandomInt(1,5)
	amount := util.RandomInt(1, 10)

	// channel collects result
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < int(n); i++ {
		// txName := fmt.Sprintf("tx %d - ", i)
		
		// concurrent Go routine
		go func() {
			// ctx := context.WithValue(context.Background(), txKey, txName) // modified context
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})

			// send to channel
			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool) // keeps loop idx

	for i := 0; i < int(n); i++ {
		err := <- errs
		require.NoError(t, err)

		result := <- results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		fmt.Printf(">> tx %d: %v %v\n", i, account1.Balance, account2.Balance)

		// check accounts' balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := account2.Balance - toAccount.Balance
		require.Equal(t, diff1, -diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= int(n))
		require.NotContains(t, existed, k)
		existed[k] = true		
	}

	// check final updated balances
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updateAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance - n * amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance + n * amount, updatedAccount2.Balance)
}

func TestTransferTx2(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// run n concurrent transfer transactions (using Go routines)
	n := 10
	amount := util.RandomInt(1, 10)

	// channel collects result
	errs := make(chan error)

	for i := 0; i < n; i++ {
		// txName := fmt.Sprintf("tx %d - ", i)

		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i % 2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		// concurrent Go routine
		go func() {
			// ctx := context.WithValue(context.Background(), txKey, txName) // modified context
			ctx := context.Background()
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID: toAccountID,
				Amount: amount,
			})
			
			// send to channel
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)
	}

	// check final updated balances
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updateAccount1.Balance, updatedAccount2.Balance)

	fmt.Println(account1.ID, updateAccount1.ID, account2.ID, updatedAccount2.ID)

	// net=0 cash-flow
	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}