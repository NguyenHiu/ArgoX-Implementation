// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../openzeppelin-contracts/contracts/utils/cryptography/MessageHashUtils.sol";
import "../openzeppelin-contracts/contracts/utils/cryptography/ECDSA.sol";
import "../openzeppelin-contracts/contracts/utils/Strings.sol";

contract Onchain {
    struct Batch {
        bytes16 batchID;
        uint256 price;
        uint256 amount;
        bool side;
        address owner;
        bytes sign;
        uint256 time;
    }

    struct Order {
        bytes16 orderID;
        uint256 price;
        uint256 amount;
        bool side;
        bytes sign;
        address owner;
    }

    Batch[] _bidBatches;
    Batch[] _askBatches;
    uint256 _registerFee;
    uint256 _waitingTime;
    mapping(address => uint256) _matcherStakes; // matcher's address => stake amount
    mapping(bytes16 => Batch) _batchMapping; // batch's id ==> batch
    mapping(bytes16 => bytes16) _tradeMapping; // batch's id ==> id of the matched batch

    constructor() {
        _registerFee = 1 ether;
        _waitingTime = 0 seconds;
    }

    // events
    event PartialMatch(bytes16);
    event FullfilMatch(bytes16);
    event ReceivedBatchDetails(bytes16);
    event AcceptBatch(bytes16);
    event PunishMatcher(address);
    event RemoveBatchOutOfDate(bytes16);
    event WrongOrders(bytes16);
    event RevertBatch(bytes16);

    // modifiers
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

    // functions
    function getRegisterFee() public view returns (uint256) {
        return _registerFee;
    }

    function getWaitingTime() public view returns (uint256) {
        return _waitingTime;
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
        _tryRevertBatch(_tradeMapping[batchID]);
        emit RemoveBatchOutOfDate(batchID);
        emit PunishMatcher(_batchMapping[batchID].owner);
    }

    function isPending(bytes16 batchID) public view returns (bool) {
        return
            _batchMapping[batchID].time != 0 &&
            _batchMapping[batchID].time + _waitingTime < block.timestamp;
    }

    // function submitOrderDetails(
    //     bytes16 batchID,
    //     Order[] memory _ors
    // ) public isPendingBatch(batchID) isBatchOwner(batchID) {
    //     uint256 _temp = 0;
    //     bytes memory ordersHash;
    //     for (uint8 i = 0; i < _ors.length; i++) {
    //         // Hash each order
    //         bytes memory packedOrders = abi.encodePacked(
    //             _ors[i].orderID,
    //             _ors[i].price,
    //             _ors[i].amount,
    //             _ors[i].side
    //         );
    //         ordersHash = keccak256(abi.encodePacked(ordersHash, packedOrders));

    //         // Check order's signature
    //         if (
    //             ECDSA.recover(
    //                 MessageHashUtils.toEthSignedMessageHash(packedOrders),
    //                 _ors[i].sign
    //             ) != _ors[i].owner
    //         ) {
    //             emit WrongOrders(batchID);
    //             return;
    //         }
    //         _temp += _ors[i].amount;
    //     }

    //     if (_temp != _batchMapping[batchID].amount) {
    //         emit WrongOrders(batchID);
    //         return;
    //     }     

    //     // FIXME: There is a case where Matcher sends 'fake orders' that are still accepted,
    //     //          This action can be accepted if:
    //     //                  + These 'fake orders' are valid orders (having valid signatures)
    //     //                  + The cumulative amount of these fake orders is equal to the amount of real orders

    //     // SUGGEST: Verify batch's signature!
        


    //     _batchMapping[batchID].time = 0;
    //     emit ReceivedBatchDetails(batchID);
    // }

    function sendBatch(
        bytes16 batchID,
        uint256 price,
        uint256 amount,
        bool side,
        address owner,
        bytes memory sign
    ) public {
        emit AcceptBatch(batchID);
        Batch memory _nb = Batch(batchID, price, amount, side, owner, sign, 0);
        _sendBatch(_nb);
    }

    function deleteBatch(
        bytes16 batchID
    ) public isBatchOwner(batchID) isExistedBatch(batchID) {
        if (_batchMapping[batchID].side) _delete(_bidBatches, batchID);
        else _delete(_askBatches, batchID);
        delete _batchMapping[batchID];
    }

    function updateBatch(
        bytes16 batchID,
        uint256 amount
    ) public isBatchOwner(batchID) isExistedBatch(batchID) {
        if (_batchMapping[batchID].side) {
            _update(_bidBatches, batchID, amount);
        } else {
            _update(_askBatches, batchID, amount);
        }
        _batchMapping[batchID].amount = amount;
    }

    function _sendBatch(Batch memory batch) internal {
        if (batch.side) _insert(_bidBatches, batch);
        else _insert(_askBatches, batch);
        _match();
    }

    function _match() internal {
        while (
            _bidBatches.length != 0 &&
            _askBatches.length != 0 &&
            _bidBatches[0].price >= _askBatches[0].price
        ) {
            if (_bidBatches[0].amount > _askBatches[0].amount) {
                _bidBatches[0].amount -= _askBatches[0].amount;
                _delete(_askBatches, _askBatches[0].batchID);
            } else if (_bidBatches[0].amount < _askBatches[0].amount) {
                _askBatches[0].amount -= _bidBatches[0].amount;
                _delete(_bidBatches, _bidBatches[0].batchID);
            } else {
                _bidBatches[0].time = block.timestamp;
                _batchMapping[_bidBatches[0].batchID] = _bidBatches[0];
                _askBatches[0].time = block.timestamp;
                _batchMapping[_askBatches[0].batchID] = _askBatches[0];

                _tradeMapping[_bidBatches[0].batchID] = _askBatches[0].batchID;
                _tradeMapping[_askBatches[0].batchID] = _bidBatches[0].batchID;

                emit FullfilMatch(_bidBatches[0].batchID);
                emit FullfilMatch(_askBatches[0].batchID);

                _delete(_askBatches, _askBatches[0].batchID);
                _delete(_bidBatches, _bidBatches[0].batchID);
            }
        }
    }

    function _update(
        Batch[] storage arr,
        bytes16 batchID,
        uint256 amount
    ) internal {
        for (uint i = 0; i < arr.length; i++) {
            if (arr[i].batchID == batchID) {
                arr[i].amount = amount;
                break;
            }
        }
    }

    function _delete(Batch[] storage arr, bytes16 batchID) internal {
        bool mark = false;
        for (uint i = 0; i < arr.length - 1; i++) {
            if (!mark && arr[i].batchID == batchID) {
                mark = true;
            }
            if (mark) {
                arr[i] = arr[i + 1];
            }
        }
        arr.pop();
    }

    function _insert(
        Batch[] storage arr,
        Batch memory batch
    ) internal returns (bool) {
        int s = -1;
        if (batch.side) s = 1;

        // arr is empty OR batch is the worst
        if (
            arr.length == 0 ||
            (int(batch.price) * s <= int(arr[arr.length - 1].price) * s)
        ) {
            arr.push(batch);
        } else {
            for (uint i = 0; i < arr.length; i++) {
                if (int(batch.price) * s > int(arr[i].price) * s) {
                    arr.push(arr[arr.length - 1]);
                    for (uint j = arr.length - 2; j > i; j--) {
                        arr[j] = arr[j - 1];
                    }
                    arr[i] = batch;
                    break;
                }
            }
        }

        return true;
    }

    function _tryRevertBatch(bytes16 batchID) internal {
        if (_batchMapping[batchID].time == 0) {
            emit RevertBatch(batchID);
            Batch memory b = _batchMapping[batchID];
            _sendBatch(b);
        }
    }
}