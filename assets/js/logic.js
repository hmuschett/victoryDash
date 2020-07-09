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
   var url = 'http://localhost:3000/api/v1/ordersmails';
fetch(url, {
  method: 'POST',
  body: JSON.stringify(data), 
  headers:{
    'Content-Type': 'application/json'
  }
}).then(res => res.json())
.then(res => isSentMail(res))
.catch(error => console.error('Error:', error))
.then(response => console.log('Success:', response)); 

} 
function  isSentMail( data){
  console.log(data)
  if(data.data.No){
    alert("Those orders not have products from WERM")
  }
}
function refreshOrders(){
  fetch("/updateOrder", {
    method: 'GET',
    headers:{
      'Content-Type': 'application/json'
    }
   
  }).then(res => res.json()) 
  console.log("why are here")
}

function loadMore() {
  var xhr = new XMLHttpRequest();
  
  xhr.open("GET", "/updateOrder/", true);
  try { xhr.send(); } catch (err) { /* handle error */ }
}