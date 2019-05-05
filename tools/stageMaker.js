const divideWidth = 16
const divideHeight = 16

const app = new Vue({
    el : '#app',
    data : {
        picPath : ""
    },
    methods : {
        handleClick : function() {
            this.divideImage(this.picPath)
        },
        divideImage : function(imgSrc) {
            //画像の取得からロード後処理まで
            const img = new Image();
            img.origin = 'anonymous';
            tmp = this
            img.onload = function() {
                alert(img)
                tmp.segmentationImage(img, divideWidth, divideHeight);
            }
            img.src = imgSrc
        },
        segmentationImage : function(img,divideWidth,divideHeight) {
            // 分割用のキャンバスを作成する
            // 画面には表示されない
            let canvas = $("<canvas width=" + divideWidth + " height=" + divideHeight + ">").get(0);
            let ctx = canvas.getContext("2d");
            // 分割後のデータを保存する
            let imgList = [];
            // 縦横の個数を取得する
            let wLength = img.width / divideWidth;
            let hLength = img.height / divideHeight;
            // 分割数だけリストに入れる
            for(let num = 0; num < wLength * hLength; num++) {
                ctx.clearRect(0,0,canvas.width,canvas.divideHeight);
                ctx.drawImage(img, divideWidth * (num % wLength), divideHeight* Math.floor(num / wLength), divideWidth, divideHeight, 0, 0, divideWidth, divideHeight);
                imgList.push(canvas.toDataURL());
            }
            for(let num = 0; num < imgList.length; num++) {
                let list = $("<li></li>"); 
                list.css("background-image", "url(" + imgList[num] + ")"); 
                $("ul").append(list);
            };
        }    
    }
})