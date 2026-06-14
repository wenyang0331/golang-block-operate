// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title Counter - A simple counter contract
/// @author Generated for golang-block-operate project
/// @notice This contract implements a simple counter that can be incremented
contract Counter {
    // State variable to store the counter value
    uint256 private count;

    // Event emitted when the counter is incremented
    event Incremented(address indexed who, uint256 newValue);
    
    // Event emitted when the counter is reset
    event Reset(address indexed who);

    /// @notice Constructor to initialize the counter
    /// @param _initialCount Initial value for the counter
    constructor(uint256 _initialCount) {
        count = _initialCount;
    }

    /// @notice Increment the counter by 1
    function increment() public {
        count += 1;
        emit Incremented(msg.sender, count);
    }

    /// @notice Get the current counter value
    /// @return The current counter value
    function getCount() public view returns (uint256) {
        return count;
    }

    /// @notice Reset the counter to 0
    function reset() public {
        count = 0;
        emit Reset(msg.sender);
    }

    /// @notice Set the counter to a specific value
    /// @param _newCount The new counter value
    function setCount(uint256 _newCount) public {
        count = _newCount;
    }
}
