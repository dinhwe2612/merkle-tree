import { SimpleMerkleTree, StandardMerkleTree } from "@openzeppelin/merkle-tree";
import { keccak256 } from "ethereum-cryptography/keccak";

const root = "0x388b02b93ee3b517ca794a0293ca294dcf222df1c4fb08e2cc498311e70745b7";
const value = "0xe25b1ca0956dcaefdeb1d3b1ac09beacd0c59a8da38d218beaafe304313ec5e3";
const proof = [
  "0x77bf017d3c7c57b13f4075b398f072fc239d6e184b9b72454b759d33070dac49", "0x3b4d86955e34e3c7a99d614089a837d7c2f7cd58bf7fcea6e3ef2b53f711a5af",
];
const types = ["bytes"];

// const hash = Buffer.from(keccak256(value[0]));
// console.log("Hash of value:", hash.toString("hex"));

const isValid = SimpleMerkleTree.verify(root, value, proof);
console.log("Merkle proof is valid?", isValid);


// const values = [
//   ["0x1a8cd71aeb2aa2af4b47bc876cbc93f6bbee71af2f83ffc9a3c4e2c860e6eff0"],
//   ["0x57580b9c14e7ed6fe1f1e6245fda0c26d7213a2904a1c32d55bf354eaac5ac39"]
// ];
// const types2 = ["bytes32"];
// const tree = StandardMerkleTree.of(values, types2);
// console.log("Merkle root (JS):", tree.root);
// console.log("Proof for leaf 0 (JS):", tree.getProof(0));
