const express = require("express");
const router = express.Router();

// 조회
router.get('/', (req, res) =>{
    res.render('get');
}); 

// 생성
router.get('/new', (req, res) =>{
    res.render('set');
});

module.exports = router;