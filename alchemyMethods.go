package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/schollz/progressbar/v3"
)

type token struct {
	TokenID string `json:"tokenId"`
}

type responseOwned struct {
	Tokens  []token `json:"ownedNfts"`
	Name    string  `json:"name"`
	PageKey string  `json:pageKey`
}

type response1 struct {
	Tokens  []token `json:"nfts"`
	Name    string  `json:"name"`
	PageKey string  `json:pageKey`
}

func getNftsForOwnerInCollection(owner string, collection string, pageKey string, metadata bool) []byte {

	url := "https://eth-mainnet.g.alchemy.com/nft/v3/" +
		API_KEY +
		"/getNFTsForOwner?owner=" +
		owner +
		"&contractAddresses[]=" +
		LEGACY_NT +
		"&withMetadata=" +
		strconv.FormatBool(metadata) +
		"&pageKey=" +
		pageKey +
		"&pageSize=100"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	// Handle Error
	if err != nil {
		fmt.Println(err.Error())
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return body
}
func getNftsForCollection(collectionAddress string, pageKey string, metadata bool) []byte {
	url := "https://eth-mainnet.g.alchemy.com/nft/v3/" +
		API_KEY +
		"/getNFTsForContract?contractAddress=" +
		collectionAddress +
		"&withMetadata=" +
		strconv.FormatBool(metadata) +
		"&startToken=" +
		pageKey +
		"&limit=100"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)

	// Handle Error
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return body

}

func getBurnedNfts(metadata bool) []token {
	nftsForOwnerOfCollection := []token{}
	pageKey, prevKey, i := "", "1", 1
	bar := progressbar.Default(1820)
	for pageKey != prevKey {
		if len(nftsForOwnerOfCollection) >= 1802 {
			break
		}
		res := getNftsForOwnerInCollection(BURN, LEGACY_NT, pageKey, metadata)
		resultOwned := responseOwned{}
		if err := json.Unmarshal(res, &resultOwned); err != nil {
			panic(err)
		}
		nftsForOwnerOfCollection = append(nftsForOwnerOfCollection, resultOwned.Tokens...)
		bar.Add(len(resultOwned.Tokens))

		prevKey = pageKey
		pageKey = resultOwned.PageKey
		i++
	}
	writeToJsonFile(nftsForOwnerOfCollection, "burnedNT.json")
	return nftsForOwnerOfCollection
}

func getLegacyNftsInCollection(fileName string, tokensLimit int, metadata bool) []token {
	nftsForCollection := []token{}
	pageKey, prevKey, i := "", "1", 1
	bar := progressbar.Default(int64(tokensLimit))
	for pageKey != prevKey {
		if len(nftsForCollection) >= tokensLimit {
			break
		}
		res := getNftsForCollection(LEGACY_NT, pageKey, metadata)
		result := response1{}
		if err := json.Unmarshal(res, &result); err != nil {
			panic(err)
		}

		nftsForCollection = append(nftsForCollection, result.Tokens...)
		bar.Add(len(result.Tokens))

		prevKey = pageKey
		pageKey = result.PageKey
		i++
	}
	writeToJsonFile(nftsForCollection, fileName)
	return nftsForCollection
}
