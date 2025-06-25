// This setup uses Hardhat Ignition to manage smart contract deployments.
// Learn more about it at https://hardhat.org/ignition

const hre= require('hardhat');

async function main() {
  // 替换为你的合约地址
  const contractAddress = '0xf9523a8efa656e264b0475ae3bfaf69b385e93da';

  // 获取合约工厂
  const MyContract = await hre.ethers.getContractFactory('Greeter');

  // 连接到已部署的合约
  const contract = await MyContract.attach(contractAddress);

  // 获取值
  const value = await contract.greet();
  console.log('Current value:', value);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });