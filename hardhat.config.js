require('dotenv').config();
// Initialize hardhat-tenderly plugin for automatic contract verification
const tdly=require('@tenderly/hardhat-tenderly')
// Your private key and tenderly devnet URL (which contains our secret id)
// We read both from the .env file so we can push this config to git and/or share publicly
const privateKey = process.env.PRIVATE_KEY;
const tenderlyUrl = process.env.TENDERLY_URL;

module.exports = {
  defaultNetwork: "mainvr",
  networks: {
    hardhat: {
    },
    mainvr: {
      url: tenderlyUrl,
      accounts: [`0x${privateKey}`]
    }
  },
  solidity: {
    version: "0.8.0",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200
      }
    }
  },
  paths: {
    sources: "./contracts",
    tests: "./test",
    cache: "./cache",
    artifacts: "./artifacts"
  },
  tenderly:{
	 project: 'suc',
	// Replace with your Tenderly username
	username: 'Lxx',
	// This will automatically verify contracts on Tenderly after deployment
	privateVerification: false,
  }
}


task("accounts", "Prints the list of accounts", async (taskArgs, hre) => {
  const accounts = await hre.ethers.getSigners();

  for (const account of accounts) {
    console.log(account.address);
  }
});


task("balance", "Prints an account's balance")
  .addParam("account", "The account's address")
  .setAction(async (taskArgs,hre) => {
    const balance = await hre.ethers.provider.getBalance(taskArgs.account);

    console.log(ethers.formatEther(balance), "ETH");
  });


