var div_id = document.getElementById('datainput')
var div_face_id = document.getElementById('face')
const params = new URL(location.href).searchParams;
const lang = params.get('lang');


function SendMessage(text) {
    $.ajax({
        url: apikernel + '/chats',
        data: { msg: text },
        type: 'GET',
        dataType: 'text',
        error: function (jqXhr, Status) {
            div_face_id.innerHTML = Status;
        },
        success: function (data) {
            Speak(data);
        }
    });
}

function LoadDataInput() {
    $.ajax({
        url: apikernel + '/attributes',
        type: 'GET',
        dataType: 'json',
        error: function (jqXhr, Status) {
            div_id.innerHTML = Status;
        },
        success: function (data) {
            htmltext = '';
            for (i = 0; i < data.length; i++) {
                namefield = data[i]["name"];
                atype = data[i]["atype"];
                options = data[i]["options"];

                switch (atype) {
                    case 'String':
                        htmltext = htmltext +
                            '<div class="mb-3"><label for="exampleFormControlInput1">' + namefield + '</label><input class="form-control" name="' + namefield + '" type="text" ></div>'
                        break;
                    case 'Number':
                        htmltext = htmltext +
                            '<div class="mb-3"><label for="exampleFormControlInput1">' + namefield + '</label><input class="form-control" name="' + namefield + '" type="number" ></div>'
                        break;
                    case 'Date':
                        htmltext = htmltext + '<div class="mb-3"><label for="exampleFormControlInput1">' + namefield + '</label>' +
                            '<input class="form-control ps-0" id="datepicker" type="date" name="' + namefield + '" /></div>';
                        break;
                    case 'List':
                        htmltext = htmltext + '<div class="mb-3"><label for="exampleFormControlSelect1">' + namefield + '</label><select class="form-control form-control-solid" name="' + namefield + '">';
                        for (var j = 0; j < options.length; j++) {
                            htmltext = htmltext + '<option>' + options[j] + '</option>';
                        }
                        htmltext = htmltext + '</select></div>';
                        break;
                }
            }
            if (htmltext != '') {
                htmltext = htmltext + '<input type="button" onclick="return SubmitDataInput(this.form)" value="Enviar">';
            } else {
                htmltext = '<p>vazio</p>'
            }
            div_id.innerHTML = htmltext;
        }
    });
}

//TODO: Submit LoadDataInput in home /

function SubmitDataInput(form) {
    var formData = {};
    $(form).find("input[name]").each(function (index, node) {
        formData[node.name] = node.value;
    });

    $(form).find("select[name]").each(function (index, node) {
        formData[node.name] = node.value;
    });

    $.post(apikernel + '/attributes', formData).done(function (data) {
        LoadDataInput();
    }).error(function (error) {
        alert(error);
    });

    return true;
}

function LoadWorkspace(name, img) {
    var ctx = document.getElementById("worktitle");
    ctx.innerHTML = 'Workspace - ' + name;
    var canvas = document.getElementById("workspace");
    var background = new Image();
    background.onload = function () {
        canvas.getContext('2d').drawImage(background, 0, 0);
    }
    background.width = '100%';
    background.height = '100%';
    background.src = img;
}

function getCookie(cookieName) {
    let cookie = {};
    document.cookie.split(';').forEach(function (el) {
        let [key, value] = el.split('=');
        cookie[key.trim()] = value;
    })
    return cookie[cookieName];
}

function beep() {
    var snd = new Audio("data:audio/wav;base64,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qJrx4KOQjahZxWbcZgztj2c49nKmkId44S71j0c8eV9yDK6uPRzx5X18eDvjvQ6yKo9ZSS6l//8elePK/Lf//IInrOF/FvDoADYAGBMGb7FtErm5MXMlmPAJQVgWta7Zx2go+8xJ0UiCb8LHHdftWyLJE0QIAIsI+UbXu67dZMjmgDGCGl1H+vpF4NSDckSIkk7Vd+sxEhBQMRU8j/12UIRhzSaUdQ+rQU5kGeFxm+hb1oh6pWWmv3uvmReDl0UnvtapVaIzo1jZbf/pD6ElLqSX+rUmOQNpJFa/r+sa4e/pBlAABoAAAAA3CUgShLdGIxsY7AUABPRrgCABdDuQ5GC7DqPQCgbbJUAoRSUj+NIEig0YfyWUho1VBBBA//uQZB4ABZx5zfMakeAAAAmwAAAAF5F3P0w9GtAAACfAAAAAwLhMDmAYWMgVEG1U0FIGCBgXBXAtfMH10000EEEEEECUBYln03TTTdNBDZopopYvrTTdNa325mImNg3TTPV9q3pmY0xoO6bv3r00y+IDGid/9aaaZTGMuj9mpu9Mpio1dXrr5HERTZSmqU36A3CumzN/9Robv/Xx4v9ijkSRSNLQhAWumap82WRSBUqXStV/YcS+XVLnSS+WLDroqArFkMEsAS+eWmrUzrO0oEmE40RlMZ5+ODIkAyKAGUwZ3mVKmcamcJnMW26MRPgUw6j+LkhyHGVGYjSUUKNpuJUQoOIAyDvEyG8S5yfK6dhZc0Tx1KI/gviKL6qvvFs1+bWtaz58uUNnryq6kt5RzOCkPWlVqVX2a/EEBUdU1KrXLf40GoiiFXK///qpoiDXrOgqDR38JB0bw7SoL+ZB9o1RCkQjQ2CBYZKd/+VJxZRRZlqSkKiws0WFxUyCwsKiMy7hUVFhIaCrNQsKkTIsLivwKKigsj8XYlwt/WKi2N4d//uQRCSAAjURNIHpMZBGYiaQPSYyAAABLAAAAAAAACWAAAAApUF/Mg+0aohSIRobBAsMlO//Kk4soosy1JSFRYWaLC4qZBYWFRGZdwqKiwkNBVmoWFSJkWFxX4FFRQWR+LsS4W/rFRb/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////VEFHAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAU291bmRib3kuZGUAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMjAwNGh0dHA6Ly93d3cuc291bmRib3kuZGUAAAAAAAAAACU=");
    console.log(snd);  
    snd.play();
}

function PostLogout() {
    var errlabel = document.getElementById('errlabel')
    const email = document.getElementById("email");
    const pwd = document.getElementById("password");
    const params = new URLSearchParams(window.location.search);
    const lang = params.get("lang");

    $.ajax({
        url: location.url,
        type: 'POST',
        data: { "email": email.value, "password": pwd.value, "csrf" :getCookie("csrf_") },
        error: function (xmlHttpRequest, textStatus, errorThrown) {
            beep();
            errlabel.innerHTML = xmlHttpRequest.responseText;
        },
        success: function (data) {
            errlabel.innerHTML = "";
            if (lang == "") {
                window.location.href = "/home";
            } else {
                window.location.href = "/home?lang=" + lang
            }
        }
    });
    return false;
}

function PostSignup() {
    const params = new URLSearchParams(window.location.search);
    const lang = params.get("lang");
    data = new FormData();
    data.append( "csrf" ,getCookie("csrf_"))
    $.ajax({
        url: location.url,
        type: 'POST',
        contentType: 'multipart/form-data',
        data: data,
        processData: false,
        contentType: false,
        error: function (xmlHttpRequest, textStatus, errorThrown) {
            beep();
        },
        success: function (data) {
            if (lang == "") {
                window.location.href = "/login";
            } else {
                window.location.href = "/login?lang=" + lang
            }
        }
    });
    return false;
}


function PostLogin() {
    var errlabel = document.getElementById('errlabel')
    const email = document.getElementById("email");
    const pwd = document.getElementById("password");
    const params = new URLSearchParams(window.location.search);
    const lang = params.get("lang");

    $.ajax({
        url: location.url,
        type: 'POST',
        data: { "email": email.value, "password": pwd.value, "csrf" :getCookie("csrf_") },
        error: function (xmlHttpRequest, textStatus, errorThrown) {
            beep();
            errlabel.innerHTML = xmlHttpRequest.responseText;
        },
        success: function (data) {
            errlabel.innerHTML = "";
            if (lang == "") {
                window.location.href = "/home";
            } else {
                window.location.href = "/home?lang=" + lang
            }
        }
    });
    return false;
}

function PostSignup() {
    var errlabel = document.getElementById('errlabel')
    //const faceimage = document.getElementById("faceimage");
    const form = document.getElementById("signup");
    const params = new URLSearchParams(window.location.search);
    const lang = params.get("lang");
    data = new FormData( form );
    data.append( "csrf" ,getCookie("csrf_"))
    $.ajax({
        url: location.url,
        type: 'POST',
        contentType: 'multipart/form-data',
        data: data,
        processData: false,
        contentType: false,
        error: function (xmlHttpRequest, textStatus, errorThrown) {
            beep();
            errlabel.innerHTML = xmlHttpRequest.responseText;
        },
        success: function (data) {
            errlabel.innerHTML = "";
            if (lang == "") {
                window.location.href = "/login";
            } else {
                window.location.href = "/login?lang=" + lang
            }
        }
    });
    return false;
}