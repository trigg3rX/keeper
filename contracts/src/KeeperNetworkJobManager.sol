// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;


import "@eigenlayer-middleware/src/interfaces/IServiceManager.sol";
import {BLSApkRegistry} from "@eigenlayer-middleware/src/BLSApkRegistry.sol";
import {RegistryCoordinator} from "@eigenlayer-middleware/src/RegistryCoordinator.sol";
import {BLSSignatureChecker, IRegistryCoordinator} from "@eigenlayer-middleware/src/BLSSignatureChecker.sol";
import {OperatorStateRetriever} from "@eigenlayer-middleware/src/OperatorStateRetriever.sol";
import "@eigenlayer-middleware/src/libraries/BN254.sol";
import "./IKeeperNetworkJobManager.sol";


contract JobCreator is 
    IKeeperNetworkJobManager
{
    address public owner;
    mapping(uint32 => Job) public jobs;
    uint32 public jobCount;

    modifier onlyOwner() {
        require(msg.sender == owner, "Not the owner");
        _;
    }

    constructor() {
        owner = msg.sender;
    }

    // EVM STACK OVERFLOW ERROR
    // This function single-handedly causes the error. Due to this, wherever this is used, the error follows.
    function createJob(JobParams calldata params) external override {
        jobCount++;
        jobs[jobCount] = Job({
            jobId: jobCount,
            jobType: params.jobType,
            status: params.status,
            quorumNumbers: params.quorumNumbers,
            quorumThresholdPercentage: params.quorumThresholdPercentage,
            timeframe: params.timeframe,
            blockNumber: block.number,
            contract_add: params.contract_add,
            chain_id: params.chain_id,
            target_fnc: params.target_fnc
        });

        emit JobCreated(jobCount, params.jobType, params.contract_add, params.chain_id);
    }

    function deleteJob(uint32 jobId) external override onlyOwner {
        require(jobs[jobId].jobId != 0, "Job does not exist");
        delete jobs[jobId];
        emit JobDeleted(jobId);
    }

    function updateJobStatus(uint32 jobId, string calldata status) external override onlyOwner {
        require(jobs[jobId].jobId != 0, "Job does not exist");
        jobs[jobId].status = status;
        emit JobStatusUpdated(jobId, status);
    }

    function stake() external payable override {
        emit Staked(msg.sender, msg.value);
    }

    function withdraw(uint256 amount) external override {
        require(amount <= address(this).balance, "Insufficient balance");
        payable(msg.sender).transfer(amount);
        emit Withdrawn(msg.sender, amount);
    }

    function joobNumber() external view override returns (uint32) {
        return jobCount;
    }

    function respondToJob(
        uint32 jobId,
        JobResponse calldata jobResponse,
        JobResponseMetadata calldata jobResponseMetadata,
        BN254.G1Point[] memory pubkeysOfNonSigningOperators
    ) external {
        require(jobs[jobId].jobId != 0, "Job does not exist");
        // Logic to handle job response
        emit JobResponded(jobResponse, jobResponseMetadata);
    }
}