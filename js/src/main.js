const { ethers } = require("hardhat");

(async () => {

    const signers = await ethers.getSigners();

    const TOKEN = await ethers.getContractFactory("Token", signers[0]);
    const token = await TOKEN.deploy();
    await token.deployed();

    await (
        await token.mint(signers[1].address, 100)
    ).wait()

    await (
        await token.mint(signers[2].address, 100)
    ).wait()

    await (
        await token.mint(signers[3].address, 100)
    ).wait()
})()