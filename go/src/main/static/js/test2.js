var num=1;
var place=1;
function movea(id){
if(place==1){var op=document.getElementById("build");}
else if(place==2){var op=document.getElementById("setting");}
else if(place==3){var op=document.getElementById("log");}
op.className="";
if(id=="build"){place=1;}
else if(id=="setting"){place=2;}
else if(id=="log"){place=3;}
var np=document.getElementById(id);
np.className="active";
}
function addtext(id)
{
var opartdiv=document.getElementById(id); 
var partdiv=document.createElement("div");
partdiv.id=opartdiv.id.concat(num.toString());
partdiv.className="row";
partdiv.innerHTML=opartdiv.innerHTML; 
document.getElementById("div1").appendChild(partdiv);
var buttondiv=document.createElement("div");
buttondiv.id="deletediv"+num;
buttondiv.className="col-sm-2";
var partbutton=document.createElement("input");
partbutton.id="delete"+num;
partbutton.type="button";
partbutton.value="delete";
partbutton.className="btn btn-default";
partbutton.onclick=function(){deletetext(this)};
document.getElementById(partdiv.id).appendChild(buttondiv);
document.getElementById(buttondiv.id).appendChild(partbutton);
num=num+1;
}
function deletetext(obj){
    var strid=obj.id;  
    var o=document.getElementById(obj.id);  
    var z=o.parentElement;  
    var zz=document.getElementById(z.id).parentElement; 
    var stridone=zz.id;  
    var my = document.getElementById(stridone);  
    if (my != null){  
    my.parentNode.removeChild(my);}  
}
