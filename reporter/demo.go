package reporter

import (
	"time"

	"github.com/google/uuid"
)

func (r *Reporter) ReportABatch(batchID uuid.UUID) {
	if time.Now().Unix()-r.PendingBatches[batchID] >= r.WaitingTime {
		_logger.Info("Report batch::%v\n", batchID.String())
		r.prepareNonceAndGasPrice(0, 2000000)
		batchID, _ := batchID.MarshalBinary()
		var _batchID [16]byte
		copy(_batchID[:], batchID[:16])
		_, err := r.OnchainInstance.ReportMissingDeadline(r.Auth, _batchID)
		if err != nil {
			_logger.Error("Reporting error, err: %v\n", err)
		}
	}
}
