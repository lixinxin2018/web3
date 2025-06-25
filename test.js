const { ethers } = require('ethers');
require('dotenv').config();
async function executeMethod() {

  // Initialize an ethers instance
  const provider = new ethers.providers.JsonRpcProvider(process.env.TENDERLY_URL);

  // Execute method
  const blockNumber = await provider.getBlockNumber();

  // Print the output to console
  console.log(blockNumber);
}

executeMethod();