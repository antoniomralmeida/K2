var div_id = document.getElementById('form_di_37232723')

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


function SubmitDataInput(form) {
    var formData = {};
    $(form).find("input[name]").each(function (index, node) {
        formData[node.name] = node.value;
    });
    
    $(form).find("select[name]").each(function (index, node) {
        formData[node.name] = node.value;
    });

    $.post(apikernel +'/attributes' ,formData).done(function(data){
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
    background.onload = function() {
        canvas.getContext('2d').drawImage(background, 0,0);
    }
    background.width = '100%';
    background.height = '100%';
    background.src = img;
}

$(function () {
    LoadDataInput();
});

