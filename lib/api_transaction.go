package sebak

import (
	"net/http"

	"boscoin.io/sebak/lib/common"
	"boscoin.io/sebak/lib/error"
	"boscoin.io/sebak/lib/observer"
	"boscoin.io/sebak/lib/storage"
	"fmt"
	"github.com/gorilla/mux"
)

const GetTransactionsHandlerPattern = "/transactions"

func GetTransactionsHandler(storage *sebakstorage.LevelDBBackend) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		switch r.Header.Get("Accept") {
		case "text/event-stream":
			var readyChan = make(chan struct{})
			iterateId := sebakcommon.GetUniqueIDFromUUID()
			go func() {
				<-readyChan
				count := maxNumberOfExistingData
				iterFunc, closeFunc := GetBlockTransactions(storage, false)
				for {
					bt, hasNext := iterFunc()
					count--
					if !hasNext || count < 0 {
						break
					}
					observer.BlockTransactionObserver.Trigger(fmt.Sprintf("iterate-%s", iterateId), &bt)
				}
				closeFunc()
			}()

			callBackFunc := func(args ...interface{}) (account []byte, err error) {
				ba := args[1].(*BlockTransaction)
				if account, err = ba.Serialize(); err != nil {
					return []byte{}, sebakerror.ErrorBlockAccountDoesNotExists
				}
				return account, nil
			}
			event := "saved"
			event += " " + fmt.Sprintf("iterate-%s", iterateId)
			streaming(observer.BlockTransactionObserver, w, event, callBackFunc, readyChan)
		default:

			var s []byte
			var btl []BlockTransaction
			iterFunc, closeFunc := GetBlockTransactions(storage, false)
			for {
				bt, hasNext := iterFunc()
				if !hasNext {
					break
				}
				btl = append(btl, bt)
			}
			closeFunc()

			s, err = sebakcommon.EncodeJSONValue(btl)
			if _, err = w.Write(s); err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}
		}
	}
}

const GetTransactionByHashHandlerPattern = "/transactions/{txid}"

func GetTransactionByHashHandler(storage *sebakstorage.LevelDBBackend) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		key := vars["txid"]
		var err error

		switch r.Header.Get("Accept") {
		case "text/event-stream":

			var readyChan = make(chan struct{})
			iterateId := sebakcommon.GetUniqueIDFromUUID()
			go func() {
				<-readyChan
				var bt BlockTransaction
				if bt, err = GetBlockTransaction(storage, key); err != nil {
					http.Error(w, "Error reading request body", http.StatusInternalServerError)
					return
				}
				observer.BlockTransactionObserver.Trigger(fmt.Sprintf("iterate-%s", iterateId), &bt)
			}()

			callBackFunc := func(args ...interface{}) (account []byte, err error) {
				ba := args[1].(*BlockTransaction)
				if account, err = ba.Serialize(); err != nil {
					return []byte{}, sebakerror.ErrorBlockAccountDoesNotExists
				}
				return account, nil
			}

			event := fmt.Sprintf("iterate-%s", iterateId)
			event += " " + fmt.Sprintf("hash-%s", key)
			streaming(observer.BlockTransactionObserver, w, event, callBackFunc, readyChan)
		default:

			var s []byte
			if found, err := ExistBlockTransaction(storage, key); err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			} else if found {
				var bt BlockTransaction
				if bt, err = GetBlockTransaction(storage, key); err != nil {
					http.Error(w, "Error reading request body", http.StatusInternalServerError)
					return
				}
				if s, err = bt.Serialize(); err != nil {
					http.Error(w, "Error reading request body", http.StatusInternalServerError)
					return
				}
			} else {
				var bth BlockTransactionHistory
				if bth, err = GetBlockTransactionHistory(storage, key); err != nil {
					http.Error(w, "Error reading request body", http.StatusInternalServerError)
					return
				}
				if s, err = bth.Serialize(); err != nil {
					http.Error(w, "Error reading request body", http.StatusInternalServerError)
					return
				}
			}
			if _, err = w.Write(s); err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}
		}
	}
}
