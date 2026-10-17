// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "forge-std/Script.sol";
import "../src/MiniToken.sol";

contract DeployScript is Script {
    function run() external {
        uint256 deployerKey = vm.envUint("SIGNER_PRIVATE_KEY");
        vm.startBroadcast(deployerKey);

        MiniToken token = new MiniToken(1_000_000 ether);
        console.log("MiniToken deployed at:", address(token));

        vm.stopBroadcast();
    }
}
