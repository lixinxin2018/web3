const { ethers } = require('hardhat');
//const { ethers }=require("@nomicfoundation/hardhat-ethers");
require('dotenv').config();
require("@tenderly/hardhat-tenderly");



async function main() {
  let greeter = await ethers.deployContract("Greeter", ["Hello, Hardhat!"]);

  greeter = await greeter.waitForDeployment();//自动验证合约
}

// Do the thing!
main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
