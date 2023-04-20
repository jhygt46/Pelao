var cartvisible = 0;
var BuscarPropConf = {};
$(document).ready(function(){

    $('.conthtml').css('min-height', $(document).height());
    //size(0);
    localStorage.setItem("history", null);
    if(id_cotizacion > 0){
        getCotizacion(id_cotizacion, 0);
    }
   
});
function in_array(arr, id){
    var len = arr.length;
    if (len > 0){
        for(var i=0; i<len; i++){
            if(arr[i].id == id){
                return false;
            }
        }
    }else{
        return true;
    }
    return true;
}
function in_array2(arr, id){
    var len = arr.length;
    if (len > 0){
        for(var i=0; i<len; i++){
            if(arr[i] == id){
                return false;
            }
        }
    }else{
        return true;
    }
    return true;
}
$(window).resize(function() {
    //size(1);
});
function showFrom(num, that){
    var el = that.parentElement.children;
    for (var i=0; i<el.length; i++){
        if (i == num) {
            $(".Prop-"+num).show();
            el[i].classList.add("selected")
        }else{
            el[i].classList.remove("selected");
            $(".Prop-"+i).hide();
        }
    }    
}
function show_cart(){
    if (cartvisible == 0){
        cartvisible = 1;
        cartmostrar();
    }else{
        cartvisible = 0;
        cartesconder();
    }
}
function ToogleInfo(tipo, id){

    if(tipo == 0){
        if($(".lp-"+id).is(":visible")){
            $(".lp-"+id).hide(500);
            $(".lpb-"+id).addClass("mas");
            $(".lpb-"+id).removeClass("menos");
        }else{
            $(".lp-"+id).show(500);
            $(".lpb-"+id).addClass("menos");
            $(".lpb-"+id).removeClass("mas");
        }
    }
    if(tipo == 1){
        if($(".la-"+id).is(":visible")){
            $(".la-"+id).hide(500);
            $(".lab-"+id).addClass("mas");
            $(".lab-"+id).removeClass("menos");
        }else{
            $(".la-"+id).show(500);
            $(".lab-"+id).addClass("menos");
            $(".lab-"+id).removeClass("mas");
        }
    }

}
function renderCotizacion(data){

    $(".list_cart").html("");
    if(data.Op == 1){
        $(".t5").html(data.Lista.length);
        for(var i=0; i<data.Lista.length; i++){
            $(".list_cart").append("<div class='cart_item clearfix'><div class='cartnombre'><div class='cn1'>"+data.Lista[i].Propiedad+"</div><div class='cn2'>"+data.Lista[i].NombreAle+"</div></div><div onclick='rmCart("+data.Lista[i].IdAle+","+data.Lista[i].IdPro+","+data.IdCot+")' class='cartaccion'></div></div>"); 
        }
    }else{
        $(".t5").html("0");
    }

}
function getZoom(diff){

    if(diff < 0.03775){
        return 14;
    }else if(diff < 0.074658){
        return 13;
    }else if(diff < 0.149502){
        return 12;
    }else if(diff < 0.295757){
        return 11;
    }else if(diff < 0.592388){
        return 10;
    }else if(diff < 1.166424){
        return 9;
    }else if(diff < 2.385906){
        return 8;
    }else if(diff < 4.780926){
        return 7;
    }else if(diff < 9.52702){
        return 6;
    }else if(diff < 18.887371){
        return 5;
    }else if(diff < 37.695965){
        return 4;
    }else if(diff < 75.664715){
        return 3;
    }else if(diff < 153.71159){
        return 2;
    }else{
        return 1;
    }
}
function getCotizacion(id_cot, n){
    var send = {id_cot: id_cot, accion: "get"};
    $.ajax({
        url: "cart/",
        type: "POST",
        data: send,
        success: function(data){
            if(n == 1){
                cartvisible = 1;
                cartmostrar();
            }
            renderCotizacion(data)
        }, error: function(e){
            console.log(e);
        }
    });
}
function rmCart(id_ale, id_pro, id_cot){
    var send = {id_ale: id_ale, id_pro: id_pro, id_cot: id_cot, accion: "rm"};
    $.ajax({
        url: "cart/",
        type: "POST",
        data: send,
        success: function(data){
            renderCotizacion(data)
        }, error: function(e){
            console.log(e);
        }
    });
}
function SendCart(id_ale, id_pro, id_cot){
    var send = {id_ale: id_ale, id_pro: id_pro, id_cot: id_cot, accion: "add"};
    $.ajax({
        url: "cart/",
        type: "POST",
        data: send,
        success: function(data){
            cartvisible = 1;
            cartmostrar();
            renderCotizacion(data)
        }, error: function(e){
            console.log(e);
        }
    });
}
function AddCart(id_ale, id_pro){

    if(id_cotizacion > 0){
        swal({   
            title: "Cotizacion",   
            text: "Desea crear una nueva cotizacion o ocupar la existente",   
            type: "warning",   
            showCancelButton: true,   
            confirmButtonColor: "#DD6B55",   
            confirmButtonText: "Nueva",   
            closeOnConfirm: true,
            showLoaderOnConfirm: false
        }, function(isConfirm){
            if(isConfirm){ 
                SendCart(id_ale, id_pro, 0);
            }else{ 
                SendCart(id_ale, id_pro, id_cotizacion);
            }
        });
    }else{
        SendCart(id_ale, id_pro, 0);
    }

}
function cartmostrar(){
    $('.cart').animate({ right: 0 }, 500);
}
function cartesconder(){
    $('.cart').animate({ right: -301 }, 500);
}
function topscroll(){
    $('html, body').animate({ scrollTop: 0 }, 500);
}
function backurl(){
    
    var history = JSON.parse(window.localStorage.getItem("history"));
    var len = history.length;
    var i = 1;
    if(len > 1){
        history.pop();
        i++;
    }
    navlinks(history[len - i]);
    localStorage.setItem("history", JSON.stringify(history));
    
}
function addhistorylink(url){
    
    var historyList = JSON.parse(window.localStorage.getItem("history"));
    if(historyList == null){
        historyList = new Array();
    }
    historyList.push(url);
    history.pushState({ url: url }, "Redigo", null);
    localStorage.setItem("history", JSON.stringify(historyList));
    
}
function navlink(href){
    addhistorylink(href);
    topscroll();
    $.ajax({
        url: href,
        type: "GET",
        data: "",
        beforeSend: function(){
            //$(".loading").show();
            //$(".error").hide();
        },
        success: function(data, status){
            $(".h").html(data);
            //$(".loading").hide();
        },
        error: function(){
            //$(".error").show();
            //$(".loading").hide();
        },
        complete: function(){
        }
    });
    return false;
}
function navlinks(href){
    topscroll();
    $.ajax({
        url: href,
        type: "GET",
        data: "",
        beforeSend: function(){
            //$(".loading").show();
            //$(".error").hide();
        },
        success: function(data, status){
            $(".h").html(data);
            //$(".loading").hide();
        },
        error: function(){
            //$(".error").show();
            //$(".loading").hide();
        },
        complete: function(){
        }
    });
    return false;
}
function eliminar(accion, id, tipo, name){

    var msg = {
        title: "Eliminar "+tipo, 
        text: "Esta seguro que desea eliminar a "+name, 
        confirm: "Si, deseo eliminarlo",
        name: name,
        accion: accion,
        id: id,
    };

    confirm(msg);
        
}
function confirm(message){
    
    swal({   
        title: message['title'],   
        text: message['text'],   
        type: "error",   
        showCancelButton: true,   
        confirmButtonColor: "#DD6B55",   
        confirmButtonText: message['confirm'],   
        closeOnConfirm: false,
        showLoaderOnConfirm: true
    }, function(isConfirm){

        if(isConfirm){
            
            var send = {accion: message['accion'], id: message['id'], nombre: message['name']};
            $.ajax({
                url: "delete/",
                type: "POST",
                data: send,
                success: function(data){
                    
                    setTimeout(function(){  
                        swal({
                            title: data.Titulo,
                            text: data.Texto,
                            type: data.Tipo,
                            timer: 2000,
                            showConfirmButton: false
                        });
                        if(data.Reload)
                            navlinks('pages/'+data.Page);
                    }, 10);

                }, error: function(e){
                    console.log(e);
                }
            });
 
        }
        
    });
    
}
function openwn(url, w, h){
    var myWindow = window.open(url, "_blank", "width="+w+",height="+h);
}
function download(that){

    var send = true;
    var data = new FormData();

    var nombre = $("#nombre").val()

    data.append("id", $("#id").val())
    data.append("accion", "descargar_custom_pdf")
    $(that).parents('form').find('input').each(function(){
        if($(this).attr('type') == "checkbox" && $(this).is(':checked')){
            data.append($(this).attr('id'), "1");
        }
        if($(this).attr('type') == "checkbox" && !$(this).is(':checked')){
            data.append($(this).attr('id'), "0");
        }
    });
    if(send){
        $('.loading').show();
        $.ajax({
            url: "save/",
            type: "POST",
            contentType: false,
            data: data,
            dataType: 'binary',
            processData: false,
            cache: false,
            xhrFields: {
                'responseType': 'blob'
            },
            success: function(data){
                if(data != null){
                    backurl()
                    mensaje(1, "Pdf "+nombre+" Creado");
                    var link = document.createElement('a'),
                    filename = nombre+".pdf";
                    link.href = URL.createObjectURL(data);
                    link.download = filename;
                    link.click();
                }
            },
            error: function(){}
        });
    }
    return false;

}
function downloadexcel(lista, busqueda){

    var send = true;
    var data = new FormData();

    data.append("lista", JSON.stringify(lista))
    data.append("busqueda", JSON.stringify(busqueda))
    data.append("accion", "descargar_lista_propiedades")
    
    if(send){
        $('.loading').show();
        $.ajax({
            url: "save/",
            type: "POST",
            contentType: false,
            data: data,
            dataType: 'binary',
            processData: false,
            cache: false,
            xhrFields: {
                'responseType': 'blob'
            },
            success: function(data){
                if(data != null){
                    backurl()
                    mensaje(1, "Excel Propiedades Creado");
                    var link = document.createElement('a'),
                    filename = "propiedades.xlsx";
                    link.href = URL.createObjectURL(data);
                    link.download = filename;
                    link.click();
                }
            },
            error: function(){}
        });
    }
    return false;

}
function fm(that){
    
    var inputs = new Array();
    var selects = new Array();
    var textareas = new Array();
    var data = new FormData();
    var require = "";
    var func = "";
    var send = true;
    
    $(that).parents('form').find('input').each(function(){
        
        if($(this).attr('require')){
            require = $(this).attr('require').split(" ");
            for(var i=0; i<require.length; i++){

                func = require[i].split("-");
                if(func[0] == "email"){
                    if(!email($(this).val())){
                        send = false;
                        $(this).parent('label').find('.mensaje').html("No es un correo electronico");
                    }else{
                        $(this).parent('label').find('.mensaje').html("");
                    }
                }
                if(func[0] == "distnada"){
                    if(!distnada($(this).val())){
                        send = false;
                        $(this).parent('label').find('.mensaje').html("Debe completar este campo");
                    }else{
                        $(this).parent('label').find('.mensaje').html("");
                    }
                }
                if(func[0] == "distzero"){
                    if(!distzero($(this).val())){
                        send = false;
                        $(this).parent('label').find('.mensaje').html("Debe seleccionar una opcion");
                    }else{
                        $(this).parent('label').find('.mensaje').html("");
                    }
                }
                if(func[0] == "textma"){
                    if(!textma($(this).val(), func[1])){
                        send = false;
                        $(this).parent('label').find('.mensaje').html("Debe tener a lo menos "+func[1]+" caracteres");
                    }else{
                        $(this).parent('label').find('.mensaje').html("");
                    }
                }
                if(func[0] == "textme"){
                    if(!textme($(this).val(), func[1])){
                        send = false;
                        $(this).parent('label').find('.mensaje').html("Debe tener a lo mas "+func[1]+" caracteres");
                    }else{
                        $(this).parent('label').find('.mensaje').html("");
                    }
                }
            }
        }
        
        if($(this).attr('type') == "password"){
            data.append($(this).attr('id'), $(this).val());
            //inputs.push($(this));
        }
        if($(this).attr('type') == "text"){
            data.append($(this).attr('id'), $(this).val());
            //inputs.push($(this));
        }
        if($(this).attr('type') == "date"){
            data.append($(this).attr('id'), $(this).val());
            //inputs.push($(this));
        }
        if($(this).attr('type') == "hidden"){
            data.append($(this).attr('id'), $(this).val());
            //inputs.push($(this));
        }
        if($(this).attr('type') == "checkbox" && $(this).is(':checked')){
            data.append($(this).attr('id'), "1");
            //inputs.push($(this));
        }
        if($(this).attr('type') == "checkbox" && !$(this).is(':checked')){
            data.append($(this).attr('id'), "0");
        }
        if($(this).attr('type') == "radio" && $(this).is(':checked')){
            data.append($(this).attr('id'), $(this).val());
            //inputs.push($(this));
        }
        if($(this).attr('type') == "file"){
            var inputFileImage = document.getElementById($(this).attr('id'));
            for(var i=0; i<inputFileImage.files.length; i++){
                var file = inputFileImage.files[i];
                data.append($(this).attr('id'), file);
            }
        }
    });
    $(that).parents('form').find('select').each(function(){
        data.append($(this).attr('id'), $(this).val());
        //selects.push($(this));
    });
    $(that).parents('form').find('textarea').each(function(){
        data.append($(this).attr('id'), $(this).val());
        //textareas.push($(this));
    });
    
    console.log(data)

    if(send){
        $('.loading').show();
        $.ajax({
            url: "save/",
            type: "POST",
            contentType: false,
            data: data,
            dataType: 'json',
            processData: false,
            cache: false,
            success: function(data){
                if(data != null){
                    if(data.Reload == 1)
                        navlinks('pages/'+data.Page);
                    if(data.Op != null)
                        mensaje(data.Op, data.Msg);
                }
            },
            error: function(){}
        });
    }
    return false;

}
function mensaje(op, mens){
    
    if(op == 1){
        var type = "success";
        var timer = 3000;
    }
    if(op == 2){
        var type = "error";
        var timer = 6000;
    }
    if(op == 3){
        var type = "warning";
        var timer = 6000;
    }
    swal({
        title: "",
        text: mens,
        html: true,
        timer: timer,
        type: type
    });
    
}
function salir(){
    
    var send = {accion: "salir"};
    $.ajax({
        url: "ajax/index.php",
        type: "POST",
        data: send,
        success: function(data){
            
            location.reload();
            
        }
    });
    return false;
}

window.addEventListener('popstate', function(e){
    if(e.state !== null){
        navlinks(e.state.url)
    }
}, false);