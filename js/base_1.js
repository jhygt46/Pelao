$(document).ready(function(){

    $('.conthtml').css('min-height', $(document).height());
    //size(0);
    localStorage.setItem("history", null);
   
});

$(window).resize(function() {
    //size(1);
});

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
    
    var history = JSON.parse(window.localStorage.getItem("history"));
    if(history == null){
        history = new Array();
    }
    history.push(url);
    localStorage.setItem("history", JSON.stringify(history));
    
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
        type: "POST",
        data: "",
        beforeSend: function(){
            $(".loading").show();
            $(".error").hide();
        },
        success: function(data, status){
            $(".conthtml").html(data);
            $(".loading").hide();
        },
        error: function(){
            $(".error").show();
            $(".loading").hide();
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
                            navlink('pages/'+data.Page);
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
                console.log(data);
                if(data != null){
                    if(data.Reload == 1)
                        navlink('pages/'+data.Page);
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

