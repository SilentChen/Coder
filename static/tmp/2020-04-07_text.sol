pragma solidity ^0.5.0;
contract SimpleStorage {
    uint storedData;
        function set(uint x) public {
bytes32 funName='set';
        storedData = x;
    }
        function get() public view returns (uint) {
bytes32 funName='get';
        return storedData;
    }
        function gettwice() public view returns (uint) {
bytes32 funName='gettwice';
        return storedData*2;
    }
}
