// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../openzeppelin-contracts/contracts/utils/cryptography/MessageHashUtils.sol";
import "../openzeppelin-contracts/contracts/utils/cryptography/ECDSA.sol";
import "../openzeppelin-contracts/contracts/utils/Strings.sol";
import "../openzeppelin-contracts/contracts/token/ERC20/IERC20.sol";

contract Onchain {
    struct Batch {
        bytes16 batchID;
        uint256 price;
        uint256 amount;
        bool side;
        address owner;
        bytes signature;
        uint256 time;
    }

    struct Order {
        uint256 price;
        uint256 amount;
        bool side;
        bytes16 from;
        bytes32 tradeHash;
        bytes32 originalOrderHash;
        address owner;
    }

    struct ShadowOrder {
        uint256 price;
        uint256 amount;
        bool side;
        address owner;
    }

    uint256 _registerFee;
    uint256 _waitingTime;
    mapping(address => uint256) _matcherStakes; // matcher's address => stake amount
    mapping(bytes16 => Batch) _batchMapping; // batch's id ==> batch
    mapping(bytes16 => bytes16) _tradeMapping; // batch's id ==> id of the matched batch
    mapping(address => uint256) _depositAmount;
    mapping(bytes16 => ShadowOrder[]) _validBatch;
    address _GVNToken;
    address _owner;

    constructor(address _token) {
        _GVNToken = _token;
        _registerFee = 1 ether;
        _waitingTime = 10 seconds;
        _owner = msg.sender;
    }

    /**
     * Events
     */
    event FullfilMatch(bytes16);
    event ReceivedBatchDetails(bytes16, uint256);
    event AcceptBatch(bytes16, uint256, uint256, bool);
    event PunishMatcher(address);
    event RemoveBatchOutOfDate(bytes16);
    event InvalidOrder(bytes16);
    event InvalidBatch(bytes16);
    event RevertBatch(bytes16);
    // Statistical
    event BatchTimestamp(bytes16, uint256);
    event BatchMatchAmountAndProfit(uint256, uint256);
    event MatchedPrice(uint256);
    // --- Log Events ---
    event LogString(string);
    event LogBytes32(bytes32);
    event LogBytes16(bytes16);
    event LogBytes(bytes);
    event LogAddress(address);
    event LogUint256(uint256);
    event LogRecoverError(ECDSA.RecoverError);

    /**
     * Test Funcs
     */
    function addressToString(
        address _addr
    ) public pure returns (string memory) {
        return Strings.toHexString(uint256(uint160(_addr)), 20);
    }

    function uintToString(uint256 num) public pure returns (string memory) {
        return Strings.toString(num);
    }

    /**
     * Modifiers
     */
    modifier isPendingBatch(bytes16 batchID) {
        require(_batchMapping[batchID].time != 0, "the batch is not pending");
        _;
    }

    modifier isBatchOwner(bytes16 batchID) {
        require(
            _batchMapping[batchID].owner == msg.sender,
            "require batch's o"
        );
        _;
    }

    modifier isExistedBatch(bytes16 batchID) {
        require(
            _batchMapping[batchID].owner != address(0),
            "the batch does not exist"
        );
        _;
    }

    /**
     * Functions
     */
    function getRegisterFee() public view returns (uint256) {
        return _registerFee;
    }

    function getWaitingTime() public view returns (uint256) {
        return _waitingTime;
    }

    function myDeposit() public payable {
        _depositAmount[msg.sender] = msg.value;
    }

    function register(address _m) public payable {
        require(msg.value >= _registerFee, "register fee is not enough");
        _matcherStakes[_m] = _registerFee;
    }

    function isMatcher(address _m) public view returns (bool) {
        return _matcherStakes[_m] != 0;
    }

    function reportMissingDeadline(
        bytes16 batchID
    ) public isPendingBatch(batchID) {
        require(
            _batchMapping[batchID].time + _waitingTime < block.timestamp,
            "the batch is not out-of-date"
        );
        _matcherStakes[_batchMapping[batchID].owner] = 0; // punish: take all the stake token of the matcher
        _batchMapping[batchID].time = 0;

        // revert
        bytes16 _tradedBatch = _tradeMapping[batchID];
        delete _tradeMapping[batchID];
        _tryRevertBatch(_tradedBatch);
        emit RemoveBatchOutOfDate(batchID);
        emit PunishMatcher(_batchMapping[batchID].owner);
    }

    function isPending(bytes16 batchID) public view returns (bool) {
        return
            _batchMapping[batchID].time != 0 &&
            _batchMapping[batchID].time + _waitingTime < block.timestamp;
    }

    function submitOrderDetails(
        bytes16 batchID,
        Order[] memory _ords
    ) public isPendingBatch(batchID) isBatchOwner(batchID) {
        // Statistical
        emit BatchTimestamp(batchID, block.timestamp);

        uint256 _temp = 0;
        bytes memory ordersHash;
        for (uint8 i = 0; i < _ords.length; i++) {
            _temp += _ords[i].amount;
            // Prepare for batch verification
            ordersHash = abi.encodePacked(
                ordersHash,
                keccak256(
                    abi.encodePacked(
                        _ords[i].price,
                        _ords[i].amount,
                        _ords[i].side,
                        _ords[i].from
                    )
                ),
                _ords[i].tradeHash,
                _ords[i].originalOrderHash
            );
        }

        if (_temp != _batchMapping[batchID].amount) {
            emit InvalidBatch(batchID);
            return;
        }

        address decryptedAddr = ECDSA.recover(
            keccak256(
                abi.encodePacked(
                    batchID,
                    _batchMapping[batchID].price,
                    _batchMapping[batchID].amount,
                    _batchMapping[batchID].side,
                    uint8(_ords.length),
                    ordersHash,
                    _batchMapping[batchID].owner
                )
            ),
            _batchMapping[batchID].signature
        );
        if (decryptedAddr != _batchMapping[batchID].owner) {
            emit InvalidBatch(batchID);
            return;
        }

        // Is first
        if (_validBatch[_tradeMapping[batchID]].length == 0) {
            for (uint8 i = 0; i < _ords.length; i++) {
                ShadowOrder memory so = ShadowOrder(
                    _batchMapping[batchID].price,
                    _ords[i].amount,
                    _ords[i].side,
                    _ords[i].owner
                );
                _validBatch[batchID].push(so);
            }
        } else {
            // Is Second
            ShadowOrder[] memory fsOrders = _validBatch[_tradeMapping[batchID]];
            if (fsOrders.length > 0 && fsOrders[0].side == true) {
                for (uint8 i = 0; i < fsOrders.length; i++) {
                    // buy
                    _depositAmount[fsOrders[i].owner] -=
                        fsOrders[i].price *
                        fsOrders[i].amount;
                    IERC20(_GVNToken).transferFrom(
                        _owner,
                        fsOrders[i].owner,
                        fsOrders[i].amount
                    );
                }
            } else {
                for (uint8 i = 0; i < fsOrders.length; i++) {
                    // sell
                    IERC20(_GVNToken).transferFrom(
                        fsOrders[i].owner,
                        _owner,
                        fsOrders[i].amount
                    );
                    _depositAmount[fsOrders[i].owner] +=
                        fsOrders[i].price *
                        fsOrders[i].amount;
                }
            }

            if (_ords.length > 0 && _ords[0].side == true) {
                for (uint8 i = 0; i < _ords.length; i++) {
                    // buy
                    _depositAmount[_ords[i].owner] -=
                        _ords[i].price *
                        _ords[i].amount;
                    IERC20(_GVNToken).transferFrom(
                        _owner,
                        _ords[i].owner,
                        _ords[i].amount
                    );
                }
            } else {
                for (uint8 i = 0; i < _ords.length; i++) {
                    // sell
                    IERC20(_GVNToken).transferFrom(
                        _ords[i].owner,
                        _owner,
                        _ords[i].amount
                    );
                    _depositAmount[_ords[i].owner] +=
                        _ords[i].price *
                        _ords[i].amount;
                }
            }
        }

        _batchMapping[batchID].time = 0;
        emit ReceivedBatchDetails(batchID, _ords.length);
    }

    function sendBatch(
        bytes16 batchID,
        uint256 price,
        uint256 amount,
        bool side,
        address owner,
        bytes memory sign
    ) public {
        // Statistical
        emit BatchTimestamp(batchID, block.timestamp);

        emit AcceptBatch(batchID, price, amount, side);
        Batch memory _nb = Batch(batchID, price, amount, side, owner, sign, 0);
        _sendBatch(_nb);
    }

    function _sendBatch(Batch memory newBatch) internal {
        _batchMapping[newBatch.batchID] = newBatch;
    }

    function matching(
        bytes16 bidBatchID,
        bytes16 askBatchID
    ) public isExistedBatch(bidBatchID) isExistedBatch(askBatchID) {
        Batch memory bidBatch = _batchMapping[bidBatchID];
        Batch memory askBatch = _batchMapping[askBatchID];
        require(bidBatch.batchID == bidBatchID, "bid batch doesn't exsit");
        require(bidBatch.time == 0, "bid batch is pending");
        require(askBatch.batchID == askBatchID, "ask batch odesn't exsit");
        require(askBatch.time == 0, "ask batch is pending");

        // Fulfill Match
        if (
            bidBatch.price >= askBatch.price &&
            bidBatch.amount == askBatch.amount
        ) {
            // Statistical
            uint256 matchPrice = (bidBatch.price + askBatch.price) / 2;
            emit BatchMatchAmountAndProfit(bidBatch.amount, bidBatch.price - matchPrice); 
            emit MatchedPrice(matchPrice);

            _batchMapping[bidBatchID].time = block.timestamp;
            _batchMapping[askBatchID].time = block.timestamp;
            _tradeMapping[bidBatch.batchID] = askBatch.batchID;
            _tradeMapping[askBatch.batchID] = bidBatch.batchID;
            emit FullfilMatch(bidBatch.batchID);
            emit FullfilMatch(askBatch.batchID);
        }
    }

    function _tryRevertBatch(bytes16 batchID) internal {
        if (
            _batchMapping[batchID].time == 0 &&
            _tradeMapping[batchID] != 0x00000000000000000000000000000000
        ) {
            emit RevertBatch(batchID);
            Batch memory b = _batchMapping[batchID];
            _sendBatch(b);
        }
    }
}
