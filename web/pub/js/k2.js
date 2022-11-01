

var div_id = document.getElementById('form_di_37232723')
$.ajax({
    url: window.location.href + 'api-datainput',
    type: 'GET',
    dataType: 'json',
    success: function (data) {
        htmltext = '';
        fields = '';
        for (i = 0; i < data.length; i++) {
            //id = data[i]["id"];
            name = data[i]["name"];
            atype = data[i]["atype"];
            options = data[i]["options"];
            fields = fields + '|' + name; 
            switch (atype) {
                case 'String':
                    htmltext = htmltext +
                        '<div class="mb-3"><label for="exampleFormControlInput1">' + name + '</label><input class="form-control" name="' + name + '" type="text" ></div>'
                    break;
                case 'Number':
                    htmltext = htmltext +
                        '<div class="mb-3"><label for="exampleFormControlInput1">' + name + '</label><input class="form-control" name="' + name + '" type="number" ></div>'
                    break;
                case 'Date':
                    htmltext = htmltext + '<div class="mb-3"><label for="exampleFormControlInput1">' + name + '</label>' +
                        '<input class="form-control ps-0" id="datepicker" type="text" name="' + name + '" /></div>';
                    break;
                case 'List':
                    htmltext = htmltext + '<div class="mb-3"><label for="exampleFormControlSelect1">' + name + '</label><select class="form-control form-control-solid" name="' + name + '">';
                    for (var i = 0; i < options.length; i++) {
                        htmltext = htmltext + '<option>' + options[i] + '</option>';
                    }
                    htmltext = htmltext + '</select></div>';
                    break;
            }
        }
        htmltext = htmltext + '<input type="hidden" name="fields" value="' + fields + '">'        
        if (fields.length > 0) {
            htmltext = htmltext + '<input type="submit" value="Submit">';
        }
        div_id.innerHTML = htmltext;
    }
});


$(function () {
    $("#datepicker").datepicker({
        dateFormat: "dd/mm/yy"
    });
});
