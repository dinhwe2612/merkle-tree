import { StandardMerkleTree } from "@openzeppelin/merkle-tree";
import { keccak256 } from "ethereum-cryptography/keccak";
import { hexToBytes, bytesToHex } from "ethereum-cryptography/utils";
import { AbiCoder } from "@ethersproject/abi";

const root = "0x1c931c433cbe3c386bcb7254f224187ec475e4376a969de1532a8a1016892d02";
const value = [Uint8Array.from(Buffer.from("1a8cd71aeb2aa2af4b47bc876cbc93f6bbee71af2f83ffc9a3c4e2c860e6eff0", "hex"))];
const proof = [
  "0x22b6bb2b3c704a7fff4be0431c0b8da25612a2350d983977a339055903f55355",
];
const types = ["bytes32"];

const isValid = StandardMerkleTree.verify(root, types, value, proof);
console.log("Merkle proof is valid?", isValid);

const abi = new AbiCoder();
const encoded = abi.encode(types, value);
console.log("encoded:", encoded);
const leaf = keccak256(Buffer.from(encoded.slice(2), "hex"));

let current = leaf;
for (const sibling of proof) {
  const sib = hexToBytes(sibling);
  const pair = [current, sib].sort(Buffer.compare);
  current = keccak256(Buffer.concat(pair));
}
const isValidByHand = "0x" + bytesToHex(current) === root.toLowerCase();
console.log("Merkle proof is valid by hand?", isValidByHand);

const values = [
  [Uint8Array.from(Buffer.from("1a8cd71aeb2aa2af4b47bc876cbc93f6bbee71af2f83ffc9a3c4e2c860e6eff0", "hex"))],
  [Uint8Array.from(Buffer.from("57580b9c14e7ed6fe1f1e6245fda0c26d7213a2904a1c32d55bf354eaac5ac39", "hex"))]
];
const types2 = ["bytes32"];
const tree = StandardMerkleTree.of(values, types2);
console.log("Merkle root (JS):", tree.root);
console.log("Proof for leaf 0 (JS):", tree.getProof(0));
