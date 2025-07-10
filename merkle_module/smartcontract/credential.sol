// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

/**
 * @title NDACredential
 * @dev Gas-optimized version with minimal storage and efficient operations
 */
contract NDACredential is AccessControl {
    bytes32 public constant WRITER_ROLE = keccak256("WRITER_ROLE");

    // Core storage: issuer address => tree index => root hash
    mapping(address => mapping(uint256 => bytes32)) public issuerTrees;
    
    // Events
    event TreeUpdated(address indexed issuer, uint256 indexed treeIndex, bytes32 newRoot);
    event BatchTreesUpdated(address[] issuers, uint256[] treeIndices, bytes32[] newRoots);
    event RevokeStatusListUpdated(address indexed issuer, string indexed revokeStatusListId, bytes revokeStatusList);
    event BatchRevokeStatusListUpdated(address indexed issuer, string[] revokeStatusListIds, bytes[] revokeStatusLists);
    
    // Custom errors
    error TreeNotExists();
    error InvalidProof();
    error EmptyRoot();
    error ArrayLengthMismatch();
    error EmptyStatusListId();
    error EmptyRevokeStatusList();

    // Allows issuer to modify its tree root
    modifier onlyWriterOrIssuer(address issuer) {
        if (!hasRole(WRITER_ROLE, msg.sender) && msg.sender != issuer) {
            revert AccessControlUnauthorizedAccount(msg.sender, WRITER_ROLE);
        }
        _;
    }

    constructor() {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(WRITER_ROLE, msg.sender);
    }    

    /**
     * @dev Update tree root (create new or update existing)
     * @param issuer Address of the issuer (tree owner)
     * @param treeIndex Index of the tree
     * @param newRoot New root hash
     */
    function updateTreeRoot(address issuer, uint256 treeIndex, bytes32 newRoot) external onlyWriterOrIssuer(issuer) {
        if (newRoot == bytes32(0)) revert EmptyRoot();            
        issuerTrees[issuer][treeIndex] = newRoot;
        emit TreeUpdated(issuer, treeIndex, newRoot);
    }
    
    /**
     * @dev Batch update multiple tree roots for different issuers
     * @param issuers Array of issuer addresses
     * @param treeIndices Array of tree indices
     * @param newRoots Array of new root hashes
     */
    function batchUpdateTreeRoots(
        address[] calldata issuers,
        uint256[] calldata treeIndices, 
        bytes32[] calldata newRoots
    ) external onlyRole(WRITER_ROLE) {
        uint256 length = issuers.length;
        if (length != treeIndices.length || length != newRoots.length) {
            revert ArrayLengthMismatch();
        }
        
        // Cache mapping reference
        mapping(address => mapping(uint256 => bytes32)) storage trees = issuerTrees;
        
        for (uint256 i; i < length;) {
            address issuer = issuers[i];
            bytes32 newRoot = newRoots[i];
            
            if (newRoot == bytes32(0)) revert EmptyRoot();
            
            trees[issuer][treeIndices[i]] = newRoot;
            
            unchecked { ++i; }
        }
        
        emit BatchTreesUpdated(issuers, treeIndices, newRoots);
    }
    
    /**
     * @dev Update revoke status list for an issuer (event only, no storage)
     * @param issuer Address of the issuer
     * @param revokeStatusListId Unique identifier for the bitstring status list
     * @param revokeStatusList The bitstring status list data
     */
    function updateRevokeStatusList(
        address issuer,
        string calldata revokeStatusListId,
        bytes calldata revokeStatusList
    ) external onlyWriterOrIssuer(issuer) {
        if (bytes(revokeStatusListId).length == 0) revert EmptyStatusListId();
        if (revokeStatusList.length == 0) revert EmptyRevokeStatusList();
        
        emit RevokeStatusListUpdated(issuer, revokeStatusListId, revokeStatusList);
    }
    
    /**
     * @dev Batch update multiple revoke status lists for an issuer (events only, no storage)
     * @param issuer Address of the issuer
     * @param revokeStatusListIds Array of unique identifiers for bitstring status lists
     * @param revokeStatusLists Array of bitstring status list data
     */
    function batchUpdateRevokeStatusList(
        address issuer,
        string[] calldata revokeStatusListIds,
        bytes[] calldata revokeStatusLists
    ) external onlyWriterOrIssuer(issuer) {
        uint256 length = revokeStatusListIds.length;
        if (length != revokeStatusLists.length) {
            revert ArrayLengthMismatch();
        }
        
        for (uint256 i; i < length;) {
            if (bytes(revokeStatusListIds[i]).length == 0) revert EmptyStatusListId();
            if (revokeStatusLists[i].length == 0) revert EmptyRevokeStatusList();
            
            unchecked { ++i; }
        }
        
        emit BatchRevokeStatusListUpdated(issuer, revokeStatusListIds, revokeStatusLists);
    }
    
    /**
     * @dev Verify a Verifiable Credential using Merkle proof
     * @param issuer Address of the issuer
     * @param treeIndex Index of the Merkle tree
     * @param leaf The leaf node to verify (hash of the VC)
     * @param proof Merkle proof for the leaf
     * @return True if the proof is valid
     */
    function verifyVC(
        address issuer,
        uint256 treeIndex,
        bytes32 leaf,
        bytes32[] calldata proof
    ) external view returns (bool) {
        bytes32 root = issuerTrees[issuer][treeIndex];
        if (root == bytes32(0)) revert TreeNotExists();
        
        return MerkleProof.verify(proof, root, leaf);
    }
    
    /**
     * @dev Get tree root
     * @param issuer Address of the issuer
     * @param treeIndex Index of the tree
     * @return root Root hash of the tree
     */
    function getTreeRoot(address issuer, uint256 treeIndex) external view returns (bytes32 root) {
        return issuerTrees[issuer][treeIndex];
    }
    
    /**
     * @dev Check if tree exists (has non-zero root)
     * @param issuer Address of the issuer
     * @param treeIndex Index of the tree
     * @return True if tree exists
     */
    function treeExists(address issuer, uint256 treeIndex) external view returns (bool) {
        return issuerTrees[issuer][treeIndex] != bytes32(0);
    }

}