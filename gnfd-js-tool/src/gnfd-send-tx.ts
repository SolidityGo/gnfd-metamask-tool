import {StargateClient} from "@cosmjs/stargate";
import {ethers} from "ethers";

let client: StargateClient

const log = console.log

const args = process.argv.slice(2)
const txRawBytes = args[0]

const mnemonic = 'solar fetch fancy fish female hundred vocal honey desert essay install predict neck tube biology sugar fork regular area rely mystery inherit human question'
let wallet1 = ethers.Wallet.fromMnemonic(mnemonic)
log(wallet1.address, wallet1.privateKey)

const main = async () => {
    let url = 'http://localhost:26750'
    client = await StargateClient.connect(url);
    const chainId = await client.getChainId()
    log('chainId', chainId)

    let balance1 = await client.getBalance(wallet1.address, 'BNB');
    log(wallet1.address, 'balance: ', balance1)
    let balance2 = await client.getBalance('0x0000000000000000000000000000000000000001', 'BNB');
    log('0x0000000000000000000000000000000000000001', 'balance: ', balance2)

    await sendRawTx(txRawBytes)

    balance1 = await client.getBalance(wallet1.address, 'BNB');
    log(wallet1.address, 'balance: ', balance1)
    balance2 = await client.getBalance('0x0000000000000000000000000000000000000001', 'BNB');
    log('0x0000000000000000000000000000000000000001', 'balance: ', balance2)
}

const sendRawTx = async (txBytesHex: string) => {
    if (txBytesHex.startsWith('0x')) {
        txBytesHex = txBytesHex.slice(2)
    }
    const rawTxBytes: Uint8Array = Uint8Array.from(Buffer.from(txBytesHex, 'hex'));
    const txResponse = await client.broadcastTx(rawTxBytes)
    log('tx response', txResponse)
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
