
// 获取时间戳
var timestamp = new Date().getTime();

$(document).ready(function(){
	getlikesUpHtml();
	getlikesDownHtml();
	showClickNum();
});
var up = parseInt((timestamp/10000)*0.000075);
var down = parseInt((timestamp/10000)*0.000025);
function getlikesUpHtml(){
	var likesup_html = [];
	likesup_html.push('<img class="col-xs-5" onClick="getLikeUpNum()" src="../static/img/up.png">');
	likesup_html.push('<h4 class="text-center" id="likeUp">'+ up +'</h4>');	
	$("#up").html(likesup_html.join(""));
}
function getlikesDownHtml(){
	var likesdown_html = [];
	likesdown_html.push('<img class="col-xs-5" onClick="getLikeDownNum()" src="../static/img/down.png">');
	likesdown_html.push('<h4 class="text-center"  id="likeDown">'+ down +'</h4>');
	$("#down").html(likesdown_html.join(""));
}
function getLikeUpNum() {
	var upNum = parseInt($("#likeUp").text()) + 1;
	up = upNum;
	getlikesUpHtml();
}
function getLikeDownNum() {
	var dwonNum = parseInt($("#likeDown").text()) - 1;
	down = dwonNum;
	getlikesDownHtml();
}

function showClickNum(){
	var showClickHtml = [];
	var arr = [];

	for(var i=0; i<15; i++){
		var a = Math.floor(Math.random()*10);
		a = parseInt((a*(timestamp/10000))*0.000025);
		arr.push(a);
	}
	console.log(arr);
	$.each(arr, function(i,c) {
		showClickHtml.push('<div class="col-xs-4 show_content"> <div class="index_img"><img src="../static/img/home.jpg">');
		showClickHtml.push('<div class="img_text"><small>Opinions</small>');
		showClickHtml.push('<small class="num"><i class="fa fa-eye"></i>'+ c +'</small></div></div></div>');
	})
	$("#index_bottom").html(showClickHtml.join(""))

	
}            


         
            
            
              
              
      