function SendEmail()
{
   let arr=[]
   Array.from(document.querySelectorAll("input[type=checkbox][name=type]:checked")).map(e => arr.push(e.value))


   if (arr.length ==0){
      alert("Select almount a order")
      return
   }
   let data={"mails": arr}
   console.log(JSON.stringify(data))
   var url = '/api/v1/ordersmails';
fetch(url, {
  method: 'POST',
  body: JSON.stringify(data), 
  headers:{
    'Content-Type': 'application/json'
  }
}).then(res => res.json())
.then(res => isSentMail(res))
.then(response => console.log('Success:', response))
.catch(error => console.error('Error:', error))

} 

function refreshOrders(){
  fetch("/updateOrder", {
    method: 'GET',
    headers:{
      'Content-Type': 'application/json'
    }
   
  }).then(res => res.json())
  .then(res => waitReaload(res))
  
  console.log("why are here")
}

function waitReaload(data)  {
  console.log(data)
  if(data.data.SI){     
    alert("Update success")   
    window.location.reload()
  }else{
    alert("There are not new orders!!!")
  }
}

function  isSentMail( data){
  if(data.data.No){
    alert("Those orders not have products from WERM")
  }else{    
    alert("Send success!!!")
    window.location.reload()
  }
}