// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "forge-std/Test.sol";
import "../src/MiniToken.sol";

contract MiniTokenTest is Test {
    MiniToken public token;
    address public owner;
    address public alice;

    function setUp() public {
        owner = address(this);
        alice = makeAddr("alice");
        token = new MiniToken(1_000_000 ether);
    }

    function test_InitialSupply() public view {
        assertEq(token.totalSupply(), 1_000_000 ether);
        assertEq(token.balanceOf(owner), 1_000_000 ether);
    }

    function test_Name() public view {
        assertEq(token.name(), "MiniToken");
        assertEq(token.symbol(), "MINI");
    }

    function test_Transfer() public {
        token.transfer(alice, 100 ether);
        assertEq(token.balanceOf(alice), 100 ether);
        assertEq(token.balanceOf(owner), 999_900 ether);
    }

    function test_Mint() public {
        token.mint(alice, 500 ether);
        assertEq(token.balanceOf(alice), 500 ether);
        assertEq(token.totalSupply(), 1_000_500 ether);
    }

    function test_MintOnlyOwner() public {
        vm.prank(alice);
        vm.expectRevert();
        token.mint(alice, 100 ether);
    }

    function test_TransferEvent() public {
        vm.expectEmit(true, true, false, true);
        emit IERC20.Transfer(owner, alice, 100 ether);
        token.transfer(alice, 100 ether);
    }
}
