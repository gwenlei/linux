var num=1;
function addtext(id)
{
var opartdiv=document.getElementById(id); 
var partdiv=document.createElement("div");
partdiv.id=opartdiv.id.concat(num.toString());
partdiv.innerHTML=opartdiv.innerHTML; 
document.getElementById("div1").appendChild(partdiv);
var partbutton=document.createElement("input");
partbutton.id="delete"+num;
partbutton.type="button";
partbutton.value="delete";
partbutton.onclick=function(){deletetext(this)};
document.getElementById(partdiv.id).appendChild(partbutton);
num=num+1;
}
function deletetext(obj){
    var strid=obj.id;  
    var o=document.getElementById(obj.id);  
    var z=o.parentElement;  
    var stridone=z.id;  
    var my = document.getElementById(stridone);  
    if (my != null){  
    my.parentNode.removeChild(my);}  
}
