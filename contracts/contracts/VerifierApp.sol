// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.0;

import "../perun-eth-contracts/contracts/App.sol";

contract VerifierApp is App {
    function validTransition(
        Channel.Params calldata params,
        Channel.State calldata from,
        Channel.State calldata to,
        uint256 actorIdx
    ) external pure override {
        
    }
}