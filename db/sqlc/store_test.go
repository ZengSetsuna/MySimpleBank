package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	// create two accounts
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">>> before transfer:", account1.Balance, account2.Balance)

	// check the account balance to ensure the transfer is correct
	// create a random account
	// make a transfer from account1 to account2
	// check the account balance to ensure the transfer is correct
	// check the transfer record to ensure it is correct
	// check the entry records to ensure they are correct

	// run n concurrent transfer transactions
	errs := make(chan error)
	results := make(chan TransferTxResult)

	n := 5
	amount := int64(10)
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, result.FromAccount.ID)
		require.Equal(t, account2.ID, result.ToAccount.ID)
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		toEntry := result.ToEntry

		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// TODO: check the account balance

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">>> after transfer:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, updatedAccount1.Balance, account1.Balance-int64(n)*amount)
	require.Equal(t, updatedAccount2.Balance, account2.Balance+int64(n)*amount)
}

func TestTransferDeadlock(t *testing.T) {

	// create two accounts
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">>> before transfer:", account1.Balance, account2.Balance)

	// check the account balance to ensure the transfer is correct
	// create a random account
	// make a transfer from account1 to account2
	// check the account balance to ensure the transfer is correct
	// check the transfer record to ensure it is correct
	// check the entry records to ensure they are correct

	// run n concurrent transfer transactions
	errs := make(chan error)

	n := 10
	amount := int64(10)
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		txName := fmt.Sprintf("tx %d", i)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">>> after transfer:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, updatedAccount1.Balance, account1.Balance)
	require.Equal(t, updatedAccount2.Balance, account2.Balance)
}