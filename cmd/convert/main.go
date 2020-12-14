// Copyright 2019 The NeuralChain Authors
// This file is part of NeuralChain.
//
// NeuralChain is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// NeuralChain is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with NeuralChain. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/lvbin2012/NeuralChain/cmd/utils"
	"github.com/lvbin2012/NeuralChain/common"
	"github.com/lvbin2012/NeuralChain/log"
)

const (
	defaultKeyfileName = "./keyfile.json"
)

func main() {
	var (
		oldAddress     = flag.String("oldaddress", "", "old address format")
		inkeyfilepath  = flag.String("inkeyfilepath", "", "old keyfile path")
		outkeyfilepath = flag.String("outkeyfilepath", "", "new keyfile path")
		verbosity      = flag.Int("verbosity", int(log.LvlInfo), "log verbosity (0-9)")
		vmodule        = flag.String("vmodule", "", "log verbosity pattern")
	)
	flag.Parse()

	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	glogger.Verbosity(log.Lvl(*verbosity))
	glogger.Vmodule(*vmodule)
	log.Root().SetHandler(glogger)

	// Convert old address to neuralChain address
	if *oldAddress != "" {
		oldAddressArray := strings.Split(*oldAddress, ",")
		for index, addrStr := range oldAddressArray {
			fmt.Printf("[%d] old address string:\"%s\"\n", index, addrStr)
			neutAddressStr := addressToNeutAddress(addrStr)
			if neutAddressStr == common.NeutEmptyAddress {
				fmt.Println("\tInput address string convert into empty address")
			}
			fmt.Printf("\tConvert to NeutAddress.:\"%s\"\n", neutAddressStr)
		}
	}

	// Convert old keyfile to neuralChain keyfile
	if *inkeyfilepath != "" {
		if *outkeyfilepath == "" {
			*outkeyfilepath = defaultKeyfileName
		}
		if _, err := os.Stat(*outkeyfilepath); err == nil {
			utils.Fatalf("Keyfile already exists at %s.", *outkeyfilepath)
		} else if !os.IsNotExist(err) {
			utils.Fatalf("Error checking if keyfile exists: %v", err)
		}
		// Read key from file.
		content, err := ioutil.ReadFile(*inkeyfilepath)
		if err != nil {
			utils.Fatalf("Failed to read input keyfile %s: %v", inkeyfilepath, err)
		}
		var jsonContent map[string]interface{}
		err = json.Unmarshal(content, &jsonContent)
		if err != nil {
			utils.Fatalf("Failed to unmarshal json bytes: %v", err)
		}
		preAddrStr, ok := jsonContent["address"].(string)
		if !ok {
			utils.Fatalf("Failed convert interface{} to string")
		}
		neutAddressStr := addressToNeutAddress(preAddrStr)
		jsonContent["address"] = neutAddressStr
		content, err = json.Marshal(jsonContent)
		if err != nil {
			utils.Fatalf("Failed to marshal object: %v", err)
		}
		// Store the file to disk.
		if err := os.MkdirAll(filepath.Dir(*outkeyfilepath), 0700); err != nil {
			utils.Fatalf("Could not create directory %s", filepath.Dir(*outkeyfilepath))
		}
		if err := ioutil.WriteFile(*outkeyfilepath, content, 0600); err != nil {
			utils.Fatalf("Failed to write keyfile to %s: %v", outkeyfilepath, err)
		}
		fmt.Printf("old keyfile path:\"%s\", old address:\"%s\"\n", *inkeyfilepath, preAddrStr)
		if neutAddressStr == common.NeutEmptyAddress {
			fmt.Println("\told address convert into empty address")
		}
		fmt.Printf("\tConvert to NeutAddress.:\"%s\"\n", neutAddressStr)
		fmt.Printf("\tNeutAddress. keypath:\"%s\"\n", *outkeyfilepath)
	}
}

// addressToNeutAddress: convert old address to neut addresss
func addressToNeutAddress(addr string) string {
	address := common.HexToAddress(addr)
	return common.AddressToNeutAddressString(address)
}