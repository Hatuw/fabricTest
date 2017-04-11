$(document).ready(function(){
	var user = getUrlParam('name')
	$("#accountName").text(user)
	query(user)
	
	$("#topup").click(function(){
		$("#txType").val("topup");
		$("#amount").attr("placeholder","100 USD");
		$("#fee").val("0 USD");
	});

	$("#invest").click(function(){
		$("#txType").val("invest");
		$("#amount").attr("placeholder","100 GP Coins");
		$("#fee").val("5% USD");
	});

	$("#cashout").click(function(){
		$("#txType").val("cashout");
		$("#amount").attr("placeholder","100 GP Coins");
		$("#fee").val("5% USD");

	});

	$("#transfer").click(function(){
		$(".transfer").show()
		$("#txType").val("transfer");
		$("#amount").attr("placeholder","100 GP Coins");
		$("#fee").val("0 USD");	

	});
	
	$("#middleTx").click(function(){
		  user = getUrlParam("name")
			func   = $("#txType").val()
			amount = $("#amount").val()
			args = {
				"User" : user,
				"Amount": amount,
				}
			middleTx(func, args)
			
			 $('.dialog').hide();
      // $('.dialog2').hide();
	});
	
	$("#transferTx").click(function(){
			user = getUrlParam("name")
			amount   = $("#amount2").val()
			receiver = $("#party").val()
			args = {
			"From"   :   user,
			"Amount" :   amount,
			"To" : receiver,	
			}
			transferTx(args)
			$(".dialog").hide();
			$(".dialog2").hide();
	});

	$("#blcokInfo").click(function(){
			window.open("./page2.html");		
	});
	
})

function query(user){
		$.post("http://54.179.182.63:8888/query", {"User": user},
			function(result){
					$("#GPCoin").text(result.gpcoin + " GpCoin");
					$("#USD").text(result.usd + " USD");
			}
	);
}

function middleTx(func, args){
	$.post("http://54.179.182.63:8888/"+func, args)
	//alert(args.User+func+args.Amount)
	query(args.User)
}

function transferTx(args){
	$.post("http://54.179.182.63:8888/transfer", args)
	//alert(args.From+args.Amount+args.Receiver)
	query(args.From)
}

function getUrlParam(name) {
            var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)"); 
            var r = window.location.search.substr(1).match(reg);  
            if (r != null) return unescape(r[2]); return null; 
        }
