import fetch from "node-fetch";
import fs from "fs";
import cliProgress from "cli-progress";
import nfts from "./nft-Ids.json" assert { type: "json" };
import owners from "./nft-owners.json" assert { type: "json" };

var requestOptions = {
  method: "get",
  reredirect: "follow",
};



// "0x0000000000000000000000000000000000000000000000000000000000000001"
async function getAllNfts(
  requestUrl,
  contractAddress,
  startToken,
  totalLimit,
  limitPerRequest,
  withMetadata,
  requestOptions
) {
  let startId = startToken;
  let nftCount = 0;
  let allNfts = [];
  console.log("======Fetching Nfts...=========");
  // Create Progress bar
  const progressBar = new cliProgress.SingleBar({}, cliProgress.Presets.rect);
  // start the progress bar with a total value of 200 and start value of 0
  progressBar.start(totalLimit, 0);

  while (nftCount < totalLimit) {
    let response = await getNftsForCollection(
      requestUrl,
      contractAddress,
      startId,
      limitPerRequest,
      withMetadata,
      requestOptions
    );
    let nfts = await response.nfts;
    allNfts.push(...nfts.map((x) => x.id.tokenId));
    nftCount += nfts.length;
    startId = response.nextToken;
    // For Logging
    // console.log("New start ID: ", startId);
    // console.log("NFT count: ", nftCount);
    progressBar.update(nftCount);
  }
  progressBar.stop();
  return allNfts;
}


async function getNftsForCollection(
  requestUrl,
  contractAddress,
  startToken,
  limit,
  withMetadata,
  requestOptions
) {
  const requestNFTsUrl = `${requestUrl}/getNFTsForCollection/?contractAddress=${contractAddress}&withMetadata=${withMetadata}&startToken=${startToken}&limit=${limit}`;
  const response = await fetch(requestNFTsUrl, requestOptions);
  return response.json();
}

async function getOwnerForToken(
  requestUrl,
  contractAddress,
  tokenId,
  requestOptions
) {
  const getOwnerRequestUrl = `${requestUrl}/getOwnersForToken/?contractAddress=${contractAddress}&tokenId=${tokenId}`;
  let response = await fetch(getOwnerRequestUrl, requestOptions);
  return response.json();
}

async function getOwnersForTokens(requestUrl, contractAddress, nfts) {
  let owners = [];
  console.log("========Fetching Owners....========");
  // Create Progress bar
  const progressBar = new cliProgress.SingleBar({}, cliProgress.Presets.rect);
  // start the progress bar with a total value of 1000 and start value of 0
  progressBar.start(nfts.length, 0);
  for (const id of nfts) {
    const ownerAddress = await getOwnerForToken(
      requestUrl,
      contractAddress,
      id
    );
    owners.push(ownerAddress.owners.pop());
    progressBar.update(owners.length);
  }
  progressBar.stop();
  return owners;
}

const writeAddressesToJson = () => {
  fs.writeFile("nft-Ids.json", JSON.stringify(nfts), "utf8", callback);
  fs.writeFile("nft-owners.json", JSON.stringify(owners), "utf8", callback);

  fs.writeFile("nft-owners.json", owners, "utf8", callback);
  fs.close();
};

function callback(err) {
  if (err) console.log(err);
  else {
    console.log("File written successfully\n");
  }
}

function getNftCount(owners) {
  const counts = {};

  for (const num of owners) {
    counts[num] = (counts[num] || 0) + 1;
  }
  console.log(counts);
  console.log("Number of Unique holders: ", Object.keys(counts).length);
}

// TODO: remove contract address
const CONTRACT_ADDRESS = "0x07ce82f414a42d9a73b0bd9ec23c249d446a0109";
const REQUEST_URL =
  "https://eth-mainnet.g.alchemy.com/v2/-aA3zbORKR83VEUCgMAhdRFFwpqxKSt7";
const startToken =
  "0x0000000000000000000000000000000000000000000000000000000000000001";
const totalLimit = 1000;
const limitPerRequest = 100;
// Metadata inclusion flag
const withMetadata = "false";

const nfts = await getAllNfts(
  REQUEST_URL,
  CONTRACT_ADDRESS,
  startToken,
  totalLimit,
  limitPerRequest,
  withMetadata,
  requestOptions
);

const owners = await getOwnersForTokens(REQUEST_URL, CONTRACT_ADDRESS, nfts);

console.log(nfts.length);
console.log(owners.length);
// Unique wallets
let uniqueWallets = new Set(owners);
// Console
console.log(uniqueWallets);
console.log(uniqueWallets.size);
getNftCount(owners);
