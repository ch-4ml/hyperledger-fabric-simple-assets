const express = require("express");
const path = require("path");

const app = express();
const port = 3000;

const assetsRouter = require("./routes/assets");

app.set('view engine', 'ejs');
app.set('views', path.join(__dirname, 'views'));

app.get("/", (req, res) =>{
    res.render("index");
});

app.use(express.urlencoded({extended: true}));
app.use(express.json()) // To parse the incoming requests with JSON payloads

app.use("/assets", assetsRouter);

app.listen(port, () => {
    console.log(`Application listening at http://localhost:${port}`);
});