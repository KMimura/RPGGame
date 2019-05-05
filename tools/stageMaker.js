const divideWidth = 16
const divideHeight = 16

const app = new Vue({
    el : '#app',
    data : {
        picPath : "",
        imgList : [],
        contents:[],
    },
    methods : {
        clicked : function(){
            $("#previewArea").empty();
            console.log(this.contents)
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
            let printStr = "[\n";
            for (let i = 0; i < this.contents.length; i++){
                printStr += "{"
                for (let j = 0; j < this.contents[i].length; j++){
                    if(this.contents[i][j] == undefined){
                        printStr += ""
                    } else {
                        printStr += this.contents[i][j]
                    }
                    if (j < this.contents[i].length - 1){
                        printStr += ","
                    }
                }
                printStr += "}"
                if (i < this.contents.length - 1) {
                    printStr += ",\n"
                }
            }
            printStr += "\n]";
            $("#resultArea").removeClass("display-none");
            $("#resultArea").text(printStr);
            console.log(printStr)
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
            console.log(this.imgList)
            for (let y = 0; y < this.imgList.length / 15; y++){
                let tr = $("<tr></tr>");
                for (let x = y*15; x < y*15+14; x++){
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
        this.divideImage("overworld_tileset_grass.png")
    }
})