const divideWidth = 16
const divideHeight = 16

const app = new Vue({
    el : '#app',
    data : {
        picPath : "",
        imgList : [],
        contents:[],
        picSelected:"",
    },
    methods : {
        clicked : function(){
            $("#previewArea").empty();
            for (let i = 0; i < this.contents.length; i++) {
                let tr = $("<tr></tr>");
                for (let j= 0; j < this.contents[i].length; j++) {
                    let td = $("<td></td>");
                    if (this.contents[i][j] != undefined ){
                        td.css("background-image", "url(" + this.imgList[this.contents[i][j]] + ")")
                    }else{
                        td.css("display","none")
                    }
                    tr.append(td);
                }
                $("#previewArea").append(tr);
            }
        },
        print : function(){
            let printObj = [];
            for (let i = 0; i < this.contents.length; i++){
                let tmpRow = [];
                for (let j = 0; j < this.contents[i].length; j++){
                    if(this.contents[i][j] == undefined){
                        continue
                    } else {
                        tmp = {}
                        tmp.cell = Number(this.contents[i][j]);
                        tmp.portal = false;
                        tmp.obstacle = false;
                        tmp.enemy = false;
                        tmpRow[tmpRow.length] = tmp
                    };
                }
                if (tmpRow.length != 0) {
                    printObj[printObj.length] = tmpRow;
                }
            }
            printObj = {"cell-data":printObj}
            $("#resultArea").removeClass("display-none");
            $("#resultArea").text(JSON.stringify(printObj));
        },
        divideImage : function(imgSrc) {
            //画像の取得からロード後処理まで
            const img = new Image();
            img.origin = 'anonymous';
            tmp = this
            img.onload = function() {
                tmp.segmentationImage(img, divideWidth, divideHeight);
            }
            img.src = imgSrc
        },
        segmentationImage : function(img,divideWidth,divideHeight) {
            // 分割用のキャンバスを作成する
            // 画面には表示されない
            let canvas = $("<canvas width=" + divideWidth + " height=" + divideHeight + ">").get(0);
            let ctx = canvas.getContext("2d");
            // 縦横の個数を取得する
            let wLength = img.width / divideWidth;
            let hLength = img.height / divideHeight;
            // 分割数だけリストに入れる
            for(let num = 0; num < wLength * hLength; num++) {
                ctx.clearRect(0,0,canvas.width,canvas.divideHeight);
                ctx.drawImage(img, divideWidth * (num % wLength), divideHeight* Math.floor(num / wLength), divideWidth, divideHeight, 0, 0, divideWidth, divideHeight);
                this.imgList.push(canvas.toDataURL());
            }
            for (let y = 0; y < this.imgList.length / 15; y++){
                // 番号を表示
                numTr = $("<tr></tr>");
                for (let x = y*15; x < y*15+15; x++){
                    numTr.append($("<td class='indexTd'>" + x + "</td>"))
                }
                $("#sampleArea").append(numTr);                
                let tr = $("<tr></tr>");
                for (let x = y*15; x < y*15+15; x++){
                    if (x >= this.imgList.length){
                        break;
                    }
                    let td = $("<td></td>")
                    td.css("background-image", "url(" + this.imgList[x] + ")")
                    tr.append(td);
                }
                $("#sampleArea").append(tr);
            }
        }    
    },
    created () {
        for (let i = 0; i < 50; i++){
            tmp = []
            for (let j = 0; j < 50; j++){
                tmp.push(undefined)
            }
            this.contents.push(tmp)
        }
        this.divideImage("overworld_tileset_grass.png");
    },
    watch: {
        picSelected: function (val) {
            this.imgList = [];
            this.contents = [];
            for (let i = 0; i < 50; i++){
                tmp = [];
                for (let j = 0; j < 50; j++){
                    tmp.push(undefined);
                }
                this.contents.push(tmp);
            }
            $("#sampleArea").empty();  
            this.divideImage(val);
        }
    }
})