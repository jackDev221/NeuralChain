// Copyright 2015 The NeuralChain Authors
// This file is part of the NeuralChain library .
//
// The NeuralChain library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The NeuralChain library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the NeuralChain library . If not, see <http://www.gnu.org/licenses/>.

// Contains the metrics collected by the downloader.

package downloader

import (
	"github.com/lvbin2012/NeuralChain/metrics"
)

var (
	headerInMeter      = metrics.NewRegisteredMeter("neut/downloader/headers/in", nil)
	headerReqTimer     = metrics.NewRegisteredTimer("neut/downloader/headers/req", nil)
	headerDropMeter    = metrics.NewRegisteredMeter("neut/downloader/headers/drop", nil)
	headerTimeoutMeter = metrics.NewRegisteredMeter("neut/downloader/headers/timeout", nil)

	bodyInMeter      = metrics.NewRegisteredMeter("neut/downloader/bodies/in", nil)
	bodyReqTimer     = metrics.NewRegisteredTimer("neut/downloader/bodies/req", nil)
	bodyDropMeter    = metrics.NewRegisteredMeter("neut/downloader/bodies/drop", nil)
	bodyTimeoutMeter = metrics.NewRegisteredMeter("neut/downloader/bodies/timeout", nil)

	receiptInMeter      = metrics.NewRegisteredMeter("neut/downloader/receipts/in", nil)
	receiptReqTimer     = metrics.NewRegisteredTimer("neut/downloader/receipts/req", nil)
	receiptDropMeter    = metrics.NewRegisteredMeter("neut/downloader/receipts/drop", nil)
	receiptTimeoutMeter = metrics.NewRegisteredMeter("neut/downloader/receipts/timeout", nil)

	stateInMeter   = metrics.NewRegisteredMeter("neut/downloader/states/in", nil)
	stateDropMeter = metrics.NewRegisteredMeter("neut/downloader/states/drop", nil)
)
