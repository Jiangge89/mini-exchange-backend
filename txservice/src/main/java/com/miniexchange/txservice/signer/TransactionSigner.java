package com.miniexchange.txservice.signer;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.web3j.crypto.Credentials;
import org.web3j.crypto.RawTransaction;
import org.web3j.crypto.TransactionEncoder;
import org.web3j.utils.Numeric;

import java.math.BigInteger;

@Component
public class TransactionSigner {

    private final Credentials credentials;
    private final long chainId;

    public TransactionSigner(
            @Value("${signer.private-key}") String privateKey,
            @Value("${chain.id}") long chainId) {
        this.credentials = Credentials.create(privateKey);
        this.chainId = chainId;
    }

    public String getAddress() {
        return credentials.getAddress();
    }

    public String signTransaction(BigInteger nonce, BigInteger gasPrice, BigInteger gasLimit,
                                   String to, BigInteger value) {
        RawTransaction rawTx = RawTransaction.createEtherTransaction(
                nonce, gasPrice, gasLimit, to, value);
        byte[] signed = TransactionEncoder.signMessage(rawTx, chainId, credentials);
        return Numeric.toHexString(signed);
    }
}
