// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/utils/math/SafeMath.sol";

contract Token is ERC20 {
    using SafeMath for uint256;

    constructor() ERC20("Gavin", "GVN") { 
        mint(msg.sender, 1000 ether);
    }

    // functions
    function mint(address account, uint256 amount) public {
        _mint(account, amount);
    }
}