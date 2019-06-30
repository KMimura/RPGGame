const fs = require("fs");
require('date-utils');

//まずは実行する前に、nodejsをいれてください。npmでfs、date-utilsも入れてください。
//変数motoにはTiled Map Editorが出力するjsonを張ること
//CUIでnode change.jsとすると同一フォルダにmain.jsonが出来上がる

const moto = { "height":9,
 "infinite":false,
 "layers":[
        {
         "data":[90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 152, 153, 90, 90, 90, 90, 90, 49, 50, 50, 51, 90, 90, 90, 90, 90, 152, 153, 90, 90, 90, 90, 164, 165, 90, 22, 23, 24, 90, 61, 85, 86, 63, 90, 18, 19, 20, 90, 164, 165, 90, 90, 90, 90, 176, 177, 90, 34, 35, 36, 90, 61, 97, 98, 63, 90, 30, 31, 32, 90, 176, 177, 90, 90, 90, 90, 188, 189, 90, 46, 47, 48, 90, 73, 74, 74, 75, 90, 42, 43, 44, 90, 188, 189, 90, 90, 90, 90, 199, 200, 201, 202, 199, 200, 201, 202, 199, 200, 201, 202, 199, 200, 201, 202, 199, 200, 201, 202, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90],
         "height":9,
         "id":4,
         "name":"\u30bf\u30a4\u30eb\u30fb\u30ec\u30a4\u30e4\u30fc 2",
         "opacity":1,
         "type":"tilelayer",
         "visible":true,
         "width":22,
         "x":0,
         "y":0
        }],
 "nextlayerid":5,
 "nextobjectid":1,
 "orientation":"orthogonal",
 "renderorder":"right-down",
 "tiledversion":"1.2.4",
 "tileheight":16,
 "tilesets":[
        {
         "columns":12,
         "firstgid":1,
         "image":"..\/RPGGame\/assets\/pics\/overworld_tileset_grass.png",
         "imageheight":336,
         "imagewidth":192,
         "margin":0,
         "name":"overworld_tileset_grass",
         "spacing":0,
         "tilecount":252,
         "tileheight":16,
         "tilewidth":16
        }],
 "tilewidth":16,
 "type":"map",
 "version":1.2,
 "width":22
}
function splitArray(array, part) {
    var tmp = [];
    for(var i = 0; i < array.length; i += part) {
        tmp.push(array.slice(i, i + part));
    }
    return tmp;
}
const rows = splitArray(moto.layers[0].data, moto.width)
var celldata = [];
for(const row of rows){
    celldata.push(row.map(e => ({"cell":e,"portal":false,"obstacle":false,"enemy":false})))
}
var saki = {
    "meta-data":{
        "id":0,
        "player-initial-positions":{
            "A":{"X":5,"Y":5},
            "B":{"X":26,"Y":32}
        },
        "spritesheet":"pics/overworld_tileset_grass.png"
    }
}
saki["cell-data"] = celldata;
console.log(saki)
var dt = new Date();
var formatted = dt.toFormat("YYYYMMDDHH24MISS");
try{
    fs.rename('./main.json', './backup'+formatted+'.json', function (err) {
    });
}catch(err){
}
fs.writeFileSync(formatted+".json", JSON.stringify(saki),'utf8');