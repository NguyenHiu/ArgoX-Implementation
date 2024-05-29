require("@nomicfoundation/hardhat-toolbox");


task("balances", "")
  .addParam("ta", "")
  .setAction(async(args) => {
    const signers = await ethers.getSigners();

    const token = await ethers.getContractAt("Token", args.ta)

    // aliceBal = await token.balanceOf(signers[1].address)
    aliceBal = await token.balanceOf(signers[2].address)
    bobBal = await token.balanceOf(signers[3].address)

    console.log(`Alice: ${aliceBal}`);
    console.log(`Bob: ${bobBal}`);
  })

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
	solidity: "0.8.24",
	defaultNetwork: "ganache",
	networks: {
		ganache: {
			url: "http://localhost:8545",
			chainId: 1337,
		}
	}
};