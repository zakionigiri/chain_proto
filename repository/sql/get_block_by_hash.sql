SELECT
    height,
    hash,
    prevBlockHash,
    merkleRoot,
    extraNonce,
    timestamp,
    bits,
    nonce,
    transactions
FROM blocks
WHERE prevBlockHash=:hash;