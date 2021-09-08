const express = require("express");
const router = express.Router();

const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');

const ccpPath = path.resolve(__dirname, '..', '..', 'basic-network', 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

// 조회 페이지로 이동
router.get('/', (req, res) =>{
  res.render('get');
}); 

// 생성 페이지로 이동
router.get('/new', (req, res) =>{
  res.render('set');
});

// 조회 요청
router.get('/:key', async (req, res) => {
  const key = req.params.key;

  try {
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }

    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });

    const network = await gateway.getNetwork('mychannel');

    const contract = network.getContract('sacc');
    
    const result = await contract.evaluateTransaction('get', key);

    res.status(200).send({ data: result.toString() });

  } catch(err) {
    console.log(err);
    res.status(500).send();
  }
});

// 생성 요청
router.post('/new', async (req, res) => {

  const { key, value } = req.body;

  try {
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }

    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });

    const network = await gateway.getNetwork('mychannel');

    const contract = network.getContract('sacc');
    
    const result = await contract.submitTransaction('set', key, value);

    console.log(`Key ${key}에 해당하는 데이터가 성공적으로 등록되었습니다: ${result}`);

    res.status(200).send();

  } catch(err) {
    console.log(err);
    res.status(500).send();
  }
});

module.exports = router;