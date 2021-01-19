INSERT INTO transactions
    (
    txHash,
    txType,
    totalValue,
    fee,
    senderAddr,
    outCount,
    outs,
    signature,
    timestamp
    )
VALUES
    (
        :txHash,
        :txType,
        :totalValue,
        :fee,
        :senderAddr,
        :outCount,
        :outs,
        :signature,
        :timestamp
    );
