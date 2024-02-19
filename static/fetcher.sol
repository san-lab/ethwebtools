// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

contract Fetcher {
    address payable public owner;
    constructor() {
        owner = payable(msg.sender);
    }
//Sweep all the balance to the owner
      function sweep() public {
        // Note that "to" is declared as payable
        (bool success, ) = owner.call{value: address(this).balance}("");
        require(success, "Failed to send Ether");
    }
    event Received(address, uint);
    receive() external payable {
        emit Received(msg.sender, msg.value);
    }
}