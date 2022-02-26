<div class="mt">
    <h1>{{.Titulo}}</h1>
    <ul class="clearfix valign">
        <li class="back" onclick="backurl()"></li>
    </ul>
</div>
<hr>
<div class="info">
    <div class="fc" id="info-0">
        <div class="minimizar m1"></div>
        <div class="close"></div>
        <div class="name">{{.SubTitulo}}</div>
        <div class="name2">{{.SubTitulo2}}</div>
        <div class="message"></div>
        <div class="sucont">

            <form action="" method="post" class="basic-grey">
                <fieldset>
                    <input id="id" type="hidden" value="{{.FormId}}" />
                    <input id="accion" type="hidden" value="{{.FormAccion}}" />
                    <label class="nboleta">
                        <span>Nombre:</span>
                        <input id="nombre" type="text" value="{{.FormNombre}}" />
                        <div class="mensaje"></div>
                    </label>
                    <label>
                        <span>Foto:</span>
                        <input id="file_image" type="file" />
                        <div class="mensaje"></div>
                    </label>
                    <label>
                        <span>Descripcion:</span>
                        <textarea id="desc" >{{.FormDescripcion}}</textarea>
                        <div class="mensaje"></div>
                    </label>
                    <label style='margin-top:20px'>
                        <span>&nbsp;</span>
                        <a id='button' onclick="form()">Enviar</a>
                    </label>
                </fieldset>
            </form>
            
        </div>
    </div>
</div>

{{if .Lista}}
<div class="info">
    <div class="fc" id="info-0">
        <div class="minimizar m1"></div>
        <div class="close"></div>
        <div class="name">{{.TituloLista}}</div>
        <ul class="options sort clearfix">
            <li class="opt">
                <div onclick="navlinks('{{.PageMod}}?sortable=1')" class="order"></div>
            </li>
        </ul>
        <div class="message"></div>
        <div class="sucont">
            <ul class='listUser'>

                {{range .Lista}}
                <li class="user" rel="{{.Idemp}}">
                    <ul class="clearfix">
                        <li class="nombre">{{.Nombre}}</li>
                        <a title="Eliminar" class="icn borrar" onclick="eliminar('{{.Idemp}}', '{{.Nombre}}')"></a>
                        <a title="Modificar" class="icn modificar" onclick="navlink('/pages/categoria?id_emp={{.Idemp}}')"></a>
                    </ul>
                </li>
                {{end}}
            </ul>
        </div>
    </div>
</div>
{{end}}
<br />
<br />